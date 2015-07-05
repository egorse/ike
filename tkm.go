package ike

import (
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"errors"
	"math/big"
)

// ike-seperation.pdf

// 2.1.2 IKE_SA_INIT
// tkm creates KEi, Ni

// get SKEYSEED
// derive SK_e (encryption) and SK_a (authentication)

// 2.1.3 IKE_AUTH
// tkm creates SK, AUTH

// 2.1.4 CREATE_CHILD_SA
// tkm creates SK, Ni, [KEi]

type Tkm struct {
	suite       *cipherSuite
	isInitiator bool

	ids Identities

	Nr, Ni *big.Int

	DhPrivate, DhPublic *big.Int
	DhShared            *big.Int

	// for debug
	SKEYSEED, KEYMAT []byte

	skD        []byte // further keying material for child sa
	skPi, skPr []byte
	skAi, skAr []byte // integrity protection keys
	skEi, skEr []byte // encryption keys
}

func NewTkmInitiator(suite *cipherSuite, ids Identities) (tkm *Tkm, err error) {
	tkm = &Tkm{
		suite:       suite,
		isInitiator: true,
		ids:         ids,
	}
	// standard says nonce shwould be at least half of size of negotiated prf
	if err := tkm.NcCreate(suite.prfLen * 8); err != nil {
		return nil, err
	}
	// for sending public key
	if _, err := tkm.DhCreate(); err != nil {
		return nil, err
	}
	return tkm, nil
}

func NewTkmResponder(suite *cipherSuite, theirPublic, no *big.Int, ids Identities) (tkm *Tkm, err error) {
	tkm = &Tkm{
		suite: suite,
		Ni:    no,
		ids:   ids,
	}
	// at least 128 bits & at least half the key size of the negotiated prf
	if err := tkm.NcCreate(no.BitLen()); err != nil {
		return nil, err
	}
	if _, err := tkm.DhCreate(); err != nil {
		return nil, err
	}
	if err := tkm.DhGenerateKey(theirPublic); err != nil {
		return nil, err
	}
	return tkm, nil
}

// 4.1.2 creation of ike sa

// The client gets the Nr
func (t *Tkm) NcCreate(bits int) (err error) {
	no, err := rand.Prime(rand.Reader, bits)
	if t.isInitiator {
		t.Ni = no
	} else {
		t.Nr = no
	}
	return
}

// the client get the dh public value
func (t *Tkm) DhCreate() (n *big.Int, err error) {
	t.DhPrivate, err = t.suite.dhGroup.private(rand.Reader)
	if err != nil {
		return nil, err
	}
	t.DhPublic = t.suite.dhGroup.public(t.DhPrivate)
	return t.DhPublic, nil
}

// upon receipt of peers resp, a dh shared secret can be calculated
// client creates & stores the dh key
func (t *Tkm) DhGenerateKey(theirPublic *big.Int) (err error) {
	t.DhShared, err = t.suite.dhGroup.diffieHellman(theirPublic, t.DhPrivate)
	return
}

func (t *Tkm) prfplus(key, data []byte, bits int) []byte {
	var ret, prev []byte
	var round int = 1
	for len(ret) < bits {
		prev = t.suite.prf(key, append(append(prev, data...), byte(round)))
		ret = append(ret, prev...)
		round += 1
	}
	return ret[:bits]
}

// create ike sa
func (t *Tkm) IsaCreate(spiI, spiR []byte) {
	// SKEYSEED = prf(Ni | Nr, g^ir)
	SKEYSEED := t.suite.prf(append(t.Ni.Bytes(), t.Nr.Bytes()...), t.DhShared.Bytes())
	// keymat =  = prf+ (SKEYSEED, Ni | Nr | SPIi | SPIr)
	kmLen := 3*t.suite.prfLen + 2*t.suite.keyLen + 2*t.suite.macKeyLen
	// TODO - supports 32B keys only
	KEYMAT := t.prfplus(SKEYSEED,
		append(append(t.Ni.Bytes(), t.Nr.Bytes()...), append(spiI, spiR...)...),
		kmLen)

	// SK_d, SK_pi, and SK_pr MUST be prfLength
	offset := t.suite.prfLen
	t.skD = KEYMAT[0:offset]
	t.skAi = KEYMAT[offset : offset+t.suite.macKeyLen]
	offset += t.suite.macKeyLen
	t.skAr = KEYMAT[offset : offset+t.suite.macKeyLen]
	offset += t.suite.macKeyLen
	t.skEi = KEYMAT[offset : offset+t.suite.keyLen]
	offset += t.suite.keyLen
	t.skEr = KEYMAT[offset : offset+t.suite.keyLen]
	offset += t.suite.keyLen
	t.skPi = KEYMAT[offset : offset+t.suite.prfLen]
	offset += t.suite.prfLen
	t.skPr = KEYMAT[offset : offset+t.suite.prfLen]

	// for test
	t.KEYMAT = KEYMAT
	t.SKEYSEED = SKEYSEED
	// fmt.Printf("keymat length %d\n", len(KEYMAT))
	// fmt.Printf("\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
	// 	hex.Dump(t.skD),
	// 	hex.Dump(t.skAi),
	// 	hex.Dump(t.skAr),
	// 	hex.Dump(t.skEi),
	// 	hex.Dump(t.skEr),
	// 	hex.Dump(t.skPi),
	// 	hex.Dump(t.skPr))
}

// verify message appended with mac
func (t *Tkm) VerifyMac(b []byte) (err error) {
	key := t.skAi
	if t.isInitiator {
		key = t.skAr
	}
	l := len(b)
	msg := b[:l-t.suite.macLen]
	msgMAC := b[l-t.suite.macLen:]
	expectedMAC := t.suite.mac(key, msg)[:t.suite.macLen]
	if !hmac.Equal(msgMAC, expectedMAC) {
		return errors.New("HMAC verify failed: ")
	}
	return
}

func (t *Tkm) Decrypt(b []byte) (dec []byte, err error) {
	key := t.skEi
	if t.isInitiator {
		key = t.skEr
	}
	iv := b[0:t.suite.ivLen]
	ciphertext := b[t.suite.ivLen:]
	// block ciphers only yet
	mode := t.suite.cipher(key, iv, true)
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
	block.CryptBlocks(ciphertext, ciphertext)
	padlen := ciphertext[len(ciphertext)-1] + 1 // padlen byte itself
	dec = ciphertext[:len(ciphertext)-int(padlen)]
	return
}

func (t *Tkm) VerifyDecrypt(ike []byte) (dec []byte, err error) {
	b := ike[IKE_HEADER_LEN:]
	if err = t.VerifyMac(ike); err != nil {
		return
	}
	dec, err = t.Decrypt(b[PAYLOAD_HEADER_LENGTH : len(b)-t.suite.macLen])
	return
}

func (t *Tkm) Encrypt(clear []byte) (b []byte, err error) {
	key := t.skEr
	if t.isInitiator {
		key = t.skEi
	}
	iv, err := rand.Prime(rand.Reader, t.suite.ivLen*8) // bits
	if err != nil {
		return
	}
	mode := t.suite.cipher(key, iv.Bytes(), false)
	if mode == nil {
		// null transform
		return clear, nil
	}
	block := mode.(cipher.BlockMode)
	// CBC mode always works in whole blocks.
	if padlen := block.BlockSize() - len(clear)%block.BlockSize(); padlen != 0 {
		pad := make([]byte, padlen)
		pad[padlen-1] = byte(padlen - 1)
		clear = append(clear, pad...)
	}
	block.CryptBlocks(clear, clear)
	b = append(iv.Bytes(), clear...)
	return
}

func (t *Tkm) mac(b []byte) []byte {
	macKey := t.skAr
	if t.isInitiator {
		macKey = t.skAi
	}
	return t.suite.mac(macKey, b)
}

func (t *Tkm) EncryptMac(s *Message) (b []byte, err error) {
	// encrypt the remaining payloads
	encr, err := t.Encrypt(encodePayloads(s.Payloads))
	if err != nil {
		return
	}
	firstPayload := s.Payloads.Array[0].Type()
	// append to new secure payload
	b = append(encodePayloadHeader(firstPayload, uint16(len(encr)+t.suite.macLen)), encr...)
	// prepare proper ike header
	s.IkeHeader.MsgLength = uint32(len(b) + IKE_HEADER_LEN + t.suite.macLen)
	// encode and append ike header
	b = append(s.IkeHeader.Encode(), b...)
	// finally attach mac
	b = append(b, t.mac(b)...)
	return
}

func (t *Tkm) AuthId(idType IdType) []byte {
	return t.ids.ForAuthentication(idType)
}

// signed
//  intiator:  signed1 | prf(sk_pi | IDi )
//  responder: signed1 | prf(sk_pr | IDr )
// AUTH = prf( prf(Shared Secret, "Key Pad for IKEv2"), signed)
func (t *Tkm) Auth(signed1 []byte, id *IdPayload, method AuthMethod, flag IkeFlags) []byte {
	key := t.skPr
	if flag.IsInitiator() {
		key = t.skPi
	}
	signed := append(signed1, t.suite.prf(key, id.Encode())...)
	secret := t.ids.AuthData(id.Data, method)
	secret = t.suite.prf(secret, []byte("Key Pad for IKEv2"))
	return t.suite.prf(secret, signed)[:t.suite.prfLen]
}

func (t *Tkm) IpsecSaCreate(spiI, spiR []byte) (espEi, espAi, espEr, espAr []byte) {
	kmLen := 2*t.suite.keyLen + 2*t.suite.macKeyLen
	// KEYMAT = prf+(SK_d, Ni | Nr)
	KEYMAT := t.prfplus(t.skD, append(t.Ni.Bytes(), t.Nr.Bytes()...),
		kmLen)

	offset := t.suite.keyLen
	espEi = KEYMAT[0:offset]
	espAi = KEYMAT[offset : offset+t.suite.macKeyLen]
	offset += t.suite.macKeyLen
	espEr = KEYMAT[offset : offset+t.suite.keyLen]
	offset += t.suite.keyLen
	espAr = KEYMAT[offset : offset+t.suite.macKeyLen]
	// fmt.Printf("ESP keys : \n%s\n%s\n%s\n%s\n",
	// 	hex.Dump(espEi),
	// 	hex.Dump(espAi),
	// 	hex.Dump(espEr),
	// 	hex.Dump(espAr))
	return
}

// request signed data from tkm
func isa_sign(isa_id, lc_id, init_message []byte) (AUTH_loc []byte) { return }

// cert validation
// start vaildating cert chain
func cc_set_user_certficate(cc_id, ri_id, au_tha_id, CERT []byte) {}

// add remianing certs in chain
func cc_add_certificate(cc_id, autha_id, CERT []byte) {}

// validate
func cc_check_ca(cc_id, ca_id []byte) {}

// after cert validtaion, authenticate peer
func isa_auth(isa_id, cc_id, init_message, AUTH_rem []byte) {}

// create first child sa
func esa_create_first(esa_id, isa_id, sp_id, ea_id, esp_spi_loc, esp_spi_rem []byte) {}
