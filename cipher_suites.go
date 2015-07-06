package ike

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"hash"

	"github.com/dgryski/go-camellia"

	"msgbox.io/log"
)

type macFunc func(key, data []byte) []byte
type prfFunc func(key, data []byte) []byte
type cipherFunc func(key, iv []byte, isRead bool) interface{}

type cipherSuite struct {
	prfLen int
	prf    prfFunc

	dhGroup *dhGroup

	// the lengths, in bytes, of the key material needed for each component.
	keyLen            int
	macLen, macKeyLen int
	ivLen             int

	cipher cipherFunc
	mac    macFunc
	// aead   func(key, fixedNonce []byte) cipher.AEAD
}

// assume that transforms are supported
// TODO - check that the entire sute makes sense
func NewCipherSuite(trs []*SaTransform) *cipherSuite {
	cs := &cipherSuite{}
	for _, tr := range trs {
		switch tr.Transform.Type {
		case TRANSFORM_TYPE_DH:
			dhGroup, ok := kexAlgoMap[DhTransformId(tr.Transform.TransformId)]
			if !ok {
				log.Errorf("Missing dh transfom %s", tr.Transform.TransformId)
				return nil
			}
			cs.dhGroup = dhGroup
		case TRANSFORM_TYPE_PRF:
			// for hmac based prf, preferred key size is size of output
			cs.prfLen, cs.prf = prfTranform(tr.Transform.TransformId)
		case TRANSFORM_TYPE_ENCR:
			// for block mode ciphers, equal to block length
			cs.ivLen, cs.cipher = cipherTransform(tr.Transform.TransformId)
			cs.keyLen = int(tr.KeyLength) / 8 // from attribute; in bits; 256
		case TRANSFORM_TYPE_INTEG:
			cs.macLen, cs.macKeyLen, cs.mac = integrityTransform(tr.Transform.TransformId)
		}
	}
	return cs
}

func prfTranform(prfId uint16) (prfLen int, prfFunc prfFunc) {
	switch PrfTransformId(prfId) {
	case PRF_HMAC_SHA2_256:
		return sha256.Size, macPrf(sha256.New)
	case PRF_HMAC_SHA1:
		return sha1.Size, macPrf(sha1.New)
	default:
		panic("unsupported")
	}
}
func cipherTransform(cipherId uint16) (ivLen int, ciperFunc cipherFunc) {
	switch EncrTransformId(cipherId) {
	case ENCR_CAMELLIA_CBC:
		return camellia.BlockSize, cipherCamellia
	case ENCR_AES_CBC:
		return aes.BlockSize, cipherAES
	case ENCR_NULL:
		return 0, cipherNull
	default:
		panic("unsupported")
	}
}
func integrityTransform(trfId uint16) (macLen, macKeyLength int, macFunc macFunc) {
	switch AuthTransformId(trfId) {
	case AUTH_HMAC_SHA2_256_128:
		return 16 /* truncated */, sha256.Size, hashMac(sha256.New, 16)
	case AUTH_HMAC_SHA1_96:
		return 12 /* truncated */, sha1.Size, hashMac(sha1.New, 12)
	default:
		panic("unsupported")
	}
}

func macPrf(h func() hash.Hash) prfFunc {
	return func(key, data []byte) []byte {
		mac := hmac.New(h, key)
		mac.Write(data)
		return mac.Sum(nil)
	}
}

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

func cipherNull([]byte, []byte, bool) interface{} { return nil }

func hashMac(h func() hash.Hash, macLen int) macFunc {
	return func(key, data []byte) []byte {
		mac := hmac.New(h, key)
		mac.Write(data)
		return mac.Sum(nil)[:macLen]
	}
}