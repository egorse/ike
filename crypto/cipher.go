package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/dgryski/go-camellia"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/msgboxio/ike/protocol"
)

// Must returm an interface
// Interface can be either cipher.BlockMode or cipher.Stream
type cipherFunc func(key, iv []byte, isRead bool) interface{}

func (cipherFunc) MarshalJSON() ([]byte, error) { return []byte("{}"), nil }

// TODO - check if the parameters are valid
func cipherTransform(cipherId uint16, keyLen int, cipher *simpleCipher) (*simpleCipher, bool) {
	blockSize, cipherFunc, ok := _cipherTransform(cipherId)
	if !ok {
		return nil, false
	}
	if cipher == nil {
		cipher = &simpleCipher{}
	}
	cipher.keyLen = keyLen
	cipher.blockLen = blockSize
	cipher.ivLen = blockSize
	cipher.cipherFunc = cipherFunc
	cipher.EncrTransformId = protocol.EncrTransformId(cipherId)
	return cipher, true
}

func _cipherTransform(cipherId uint16) (int, cipherFunc, bool) {
	switch protocol.EncrTransformId(cipherId) {
	case protocol.ENCR_CAMELLIA_CBC:
		return camellia.BlockSize, cipherCamellia, true
	case protocol.ENCR_AES_CBC:
		return aes.BlockSize, cipherAES, true
	case protocol.ENCR_NULL:
		return 0, cipherNull, true
	default:
		return 0, nil, false
	}
}

// Cipher interface implementation

type simpleCipher struct {
	macTruncLen, macLen int
	macFunc

	keyLen, ivLen, blockLen int
	cipherFunc

	protocol.EncrTransformId
	protocol.AuthTransformId
}

func (cs *simpleCipher) String() string {
	return cs.EncrTransformId.String() + "+" + cs.AuthTransformId.String()
}

func (cs *simpleCipher) Overhead(clear []byte) int {
	return cs.blockLen - len(clear)%cs.blockLen + cs.macLen + cs.ivLen
}
func (cs *simpleCipher) VerifyDecrypt(ike, skA, skE []byte, log log.Logger) (dec []byte, err error) {
	level.Debug(log).Log(
		"msg", "simple verify&decrypt",
		"Clear", hex.Dump(ike), "SkA", hex.Dump(skA), "SkE", hex.Dump(skE))
	// MAC-then-decrypt
	if err = verifyMac(skA, ike, cs.macLen, cs.macFunc); err != nil {
		return
	}
	b := ike[protocol.IKE_HEADER_LEN:]
	dec, err = decrypt(b[protocol.PAYLOAD_HEADER_LENGTH:len(b)-cs.macLen], skE, cs.ivLen, cs.cipherFunc, log)
	return
}

func (cs *simpleCipher) EncryptMac(headers, payload, skA, skE []byte, log log.Logger) (b []byte, err error) {
	// encrypt-then-MAC
	encr, err := encrypt(payload, skE, cs.ivLen, cs.cipherFunc, log)
	if err != nil {
		return
	}
	data := append(headers, encr...)
	mac := cs.macFunc(skA, data)
	b = append(data, mac...)
	level.Debug(log).Log(
		"msg", "simple encrypt&mac",
		"Mac", hex.Dump(mac), "SkA", hex.Dump(skA), "SkE", hex.Dump(skE))
	return
}

// cipherFunc Implementations

func cipherAES(key, iv []byte, isRead bool) interface{} {
	block, _ := aes.NewCipher(key)
	if isRead {
		return cipher.NewCBCDecrypter(block, iv)
	}
	return cipher.NewCBCEncrypter(block, iv)
}

func cipherCamellia(key, iv []byte, isRead bool) interface{} {
	block, _ := camellia.New(key)
	if isRead {
		return cipher.NewCBCDecrypter(block, iv)
	}
	return cipher.NewCBCEncrypter(block, iv)
}

// TODO - this needs a proper do nothing implementation
func cipherNull([]byte, []byte, bool) interface{} { return nil }

// decryption & encryption routines

func decrypt(b, key []byte, ivLen int, cipherFn cipherFunc, log log.Logger) (dec []byte, err error) {
	iv := b[0:ivLen]
	ciphertext := b[ivLen:]
	// block ciphers only yet
	mode := cipherFn(key, iv, true)
	if mode == nil {
		// null transform
		return b, nil
	}
	block := mode.(cipher.BlockMode)
	// CBC mode always works in whole blocks.
	if len(ciphertext)%block.BlockSize() != 0 {
		err = errors.New("ciphertext is not a multiple of the block size")
		return
	}
	clear := make([]byte, len(ciphertext))
	block.CryptBlocks(clear, ciphertext)
	padlen := clear[len(clear)-1] + 1 // padlen byte itself
	if int(padlen) > block.BlockSize() {
		err = errors.New("pad length is larger than block size")
		return
	}
	dec = clear[:len(clear)-int(padlen)]
	level.Debug(log).Log(
		"Pad ", padlen,
		"Clear", hex.Dump(clear),
		"Cyp", hex.Dump(ciphertext),
		"IV", hex.Dump(iv))
	return
}

func encrypt(clear, key []byte, ivLen int, cipherFn cipherFunc, log log.Logger) (b []byte, err error) {
	iv, err := rand.Prime(rand.Reader, ivLen*8) // bits
	if err != nil {
		return
	}
	mode := cipherFn(key, iv.Bytes(), false)
	if mode == nil {
		// null transform
		return clear, nil
	}
	// TODO - block mode supported only
	block := mode.(cipher.BlockMode)
	// CBC mode always works in whole blocks.
	// (b - (length % b)) % b
	// pl := (block.BlockSize() - (len(clear) % block.BlockSize())) % block.BlockSize()
	padlen := block.BlockSize() - len(clear)%block.BlockSize()
	if padlen != 0 {
		pad := make([]byte, padlen)
		pad[padlen-1] = byte(padlen - 1)
		clear = append(clear, pad...)
	}
	ciphertext := make([]byte, len(clear))
	block.CryptBlocks(ciphertext, clear)
	b = append(iv.Bytes(), ciphertext...)
	level.Debug(log).Log(
		"Pad ", padlen,
		"Clear", hex.Dump(clear),
		"Cyp", hex.Dump(ciphertext),
		"IV", hex.Dump(iv.Bytes()))
	return
}
