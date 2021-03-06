package ike

import (
	"crypto/rand"
	"math/big"

	"github.com/msgboxio/ike/crypto"
	"github.com/pkg/errors"
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
	suite    *crypto.CipherSuite
	espSuite *crypto.CipherSuite

	Nr, Ni *big.Int

	dhPrivate, DhPublic *big.Int
	DhShared            *big.Int

	skD        []byte // further keying material for child sa
	skPi, skPr []byte // used when generating an AUTH
	skAi, skAr []byte // integrity protection keys
	skEi, skEr []byte // encryption keys
}

var errMissingCryptoKeys = errors.New("Missing crypto keys")

func NewTkm(cfg *Config, ni *big.Int) (*Tkm, error) {
	suite, err := crypto.NewCipherSuite(cfg.ProposalIke)
	if err != nil {
		return nil, err
	}
	espSuite, err := crypto.NewCipherSuite(cfg.ProposalEsp)
	if err != nil {
		return nil, err
	}
	if ni != nil {
		return newTkmResponder(suite, espSuite, ni)
	}
	return newTkmInitiator(suite, espSuite)
}

func newTkmInitiator(suite, espSuite *crypto.CipherSuite) (tkm *Tkm, err error) {
	if err = suite.CheckIkeTransforms(); err != nil {
		return
	}
	if err = espSuite.CheckEspTransforms(); err != nil {
		return
	}
	// standard says nonce should be at least half of size of negotiated prf
	ni, err := createNonce(suite.Prf.Length * 8)
	if err != nil {
		return
	}
	tkm = &Tkm{
		suite:    suite,
		espSuite: espSuite,
		Ni:       ni,
	}
	// for sending public key
	err = tkm.dhCreate()
	return
}

func newTkmResponder(suite, espSuite *crypto.CipherSuite, ni *big.Int) (tkm *Tkm, err error) {
	if err = suite.CheckIkeTransforms(); err != nil {
		return
	}
	if err = espSuite.CheckEspTransforms(); err != nil {
		return
	}
	// at least 128 bits & at least half the key size of the negotiated prf
	bitLen := ni.BitLen()
	if bitLen < 128 || bitLen < (suite.Prf.Length*8)/2 {
		err = errors.New("Proposed nonce is too small")
		return
	}
	nr, err := createNonce(bitLen)
	if err != nil {
		return
	}
	tkm = &Tkm{
		suite:    suite,
		espSuite: espSuite,
		Ni:       ni,
		Nr:       nr,
	}
	err = tkm.dhCreate()
	return
}

// 4.1.2 creation of ike sa

func createNonce(bits int) (no *big.Int, err error) {
	return rand.Prime(rand.Reader, bits)
}

func (t *Tkm) dhCreate() (err error) {
	t.dhPrivate, t.DhPublic, err = t.suite.DhGroup.Generate(rand.Reader)
	return
}

// DhGenerateKey creates & stores the dh key
// upon receipt of peers resp, a dh shared secret can be calculated
func (t *Tkm) DhGenerateKey(theirPublic *big.Int) (err error) {
	t.DhShared, err = t.suite.DhGroup.DiffieHellman(theirPublic, t.dhPrivate)
	return
}

func (t *Tkm) prfplus(key, data []byte, bits int) []byte {
	var ret, prev []byte
	var round = 1
	for len(ret) < bits {
		prev = t.suite.Prf.Apply(key, append(append(prev, data...), byte(round)))
		ret = append(ret, prev...)
		round++
	}
	return ret[:bits]
}

func (t *Tkm) skeySeedInitial() []byte {
	// SKEYSEED = prf(Ni | Nr, g^ir)
	return t.suite.Prf.Apply(append(t.Ni.Bytes(), t.Nr.Bytes()...), t.DhShared.Bytes())
}

func (t *Tkm) skeySeedRekey(old_SK_D []byte) []byte {
	// SKEYSEED = prf(SK_d (old), g^ir (new) | Ni | Nr)
	return t.suite.Prf.Apply(old_SK_D, append(t.DhShared.Bytes(), append(t.Ni.Bytes(), t.Nr.Bytes()...)...))
}

// IkeSaKeys creates ike sa keys
func (t *Tkm) IkeSaKeys(spiI, spiR []byte, old_skD []byte) {
	// fmt.Printf("key inputs: \nni:\n%snr:\n%sshared:\n%sspii:\n%sspir:\n%s",
	// 	hex.Dump(t.Ni.Bytes()), hex.Dump(t.Nr.Bytes()), hex.Dump(t.DhShared.Bytes()),
	// 	hex.Dump(spiI), hex.Dump(spiR))
	SKEYSEED := []byte{}
	if len(old_skD) == 0 {
		SKEYSEED = t.skeySeedInitial()
	} else {
		SKEYSEED = t.skeySeedRekey(old_skD)
	}
	kmLen := 3*t.suite.Prf.Length + 2*t.suite.KeyLen + 2*t.suite.MacTruncLen
	// KEYMAT =  = prf+ (SKEYSEED, Ni | Nr | SPIi | SPIr)
	KEYMAT := t.prfplus(SKEYSEED,
		append(append(t.Ni.Bytes(), t.Nr.Bytes()...), append(spiI, spiR...)...),
		kmLen)

	// SK_d, SK_pi, and SK_pr MUST be prfLength
	offset := t.suite.Prf.Length
	t.skD = append([]byte{}, KEYMAT[0:offset]...)
	t.skAi = append([]byte{}, KEYMAT[offset:offset+t.suite.MacTruncLen]...)
	offset += t.suite.MacTruncLen
	t.skAr = append([]byte{}, KEYMAT[offset:offset+t.suite.MacTruncLen]...)
	offset += t.suite.MacTruncLen
	t.skEi = append([]byte{}, KEYMAT[offset:offset+t.suite.KeyLen]...)
	offset += t.suite.KeyLen
	t.skEr = append([]byte{}, KEYMAT[offset:offset+t.suite.KeyLen]...)
	offset += t.suite.KeyLen
	t.skPi = append([]byte{}, KEYMAT[offset:offset+t.suite.Prf.Length]...)
	offset += t.suite.Prf.Length
	t.skPr = append([]byte{}, KEYMAT[offset:offset+t.suite.Prf.Length]...)

	// fmt.Printf("keymat length %d\n", len(KEYMAT))
	// fmt.Printf("skD:\n%sskAi:\n%sskAr:\n%sskEi:\n%sskEr:\n%sskPi:\n%sskPr:\n%s",
	// 	hex.Dump(t.skD),
	// 	hex.Dump(t.skAi),
	// 	hex.Dump(t.skAr),
	// 	hex.Dump(t.skEi),
	// 	hex.Dump(t.skEr),
	// 	hex.Dump(t.skPi),
	// 	hex.Dump(t.skPr))
}

func (t *Tkm) CryptoOverhead(b []byte) int {
	return t.suite.Overhead(b)
}

// MAC-then-decrypt
func (t *Tkm) VerifyDecrypt(ike []byte, forInitiator bool) (dec []byte, err error) {
	skA, skE := t.skAi, t.skEi
	if forInitiator {
		skA, skE = t.skAr, t.skEr
	}
	if skA == nil || skE == nil {
		return nil, errors.Wrap(errMissingCryptoKeys, "Decrypting")
	}
	dec, err = t.suite.VerifyDecrypt(ike, skA, skE)
	return
}

// encrypt-then-MAC
func (t *Tkm) EncryptMac(ike []byte, forInitiator bool) (b []byte, err error) {
	skA, skE := t.skAr, t.skEr
	if forInitiator {
		skA, skE = t.skAi, t.skEi
	}
	if skA == nil || skE == nil {
		return nil, errors.Wrap(errMissingCryptoKeys, "Encrypting")
	}
	b, err = t.suite.EncryptMac(ike, skA, skE)
	return
}

// IpsecSaKeys generates & returns Ipsec Sa keys
func (t *Tkm) IpsecSaKeys(ni, nr, dhShared *big.Int) (espEi, espAi, espEr, espAr []byte) {
	kmLen := 2*t.espSuite.KeyLen + 2*t.espSuite.MacTruncLen
	// KEYMAT = prf+(SK_d, Ni | Nr)
	KEYMAT := t.prfplus(t.skD, append(ni.Bytes(), nr.Bytes()...), kmLen)
	// KEYMAT = prf+(SK_d, g^ir (new) | Ni | Nr)
	if dhShared != nil {
		KEYMAT = t.prfplus(t.skD,
			append(dhShared.Bytes(), append(ni.Bytes(), nr.Bytes()...)...), kmLen)
	}
	offset := t.espSuite.KeyLen
	espEi = append([]byte{}, KEYMAT[0:offset]...)
	espAi = append([]byte{}, KEYMAT[offset:offset+t.espSuite.MacTruncLen]...)
	offset += t.espSuite.MacTruncLen
	espEr = append([]byte{}, KEYMAT[offset:offset+t.espSuite.KeyLen]...)
	offset += t.espSuite.KeyLen
	espAr = append([]byte{}, KEYMAT[offset:offset+t.espSuite.MacTruncLen]...)
	// fmt.Printf("ESP keys :\nEi:\n%sAi:\n%sEr:\n%sAr\n%s",
	// 	hex.Dump(espEi),
	// 	hex.Dump(espAi),
	// 	hex.Dump(espEr),
	// 	hex.Dump(espAr))
	return
}

// SignB gets signed data from tkm
// section 2.15
// For the responder, the octets to be signed
// start with the first octet of the first SPI in the
// header of the second message (IKE_SA_INIT response) and end with the
// last octet of the last payload in the second message.  => initIRB
// Appended to this (for the purposes of computing the signature) are the
// initiator's nonce Ni (just the value, not the payload containing it),
// and the value prf(SK_pr, IDr')
// so signB :=
// responder: initRB | Ni | prf(SK_pr, IDr')
// initiator: initIB | Nr | prf(SK_pi, IDi')
// this method can be used by signer & verifier
func (t *Tkm) SignB(initB []byte, id []byte, forInitiator bool) []byte {
	// ResponderSignedOctets = RealMessage2 | NonceIData | MACedIDForR
	// InitiatorSignedOctets = RealMessage1 | NonceRData | MACedIDForI
	key := t.skPr
	nonce := t.Ni
	if forInitiator {
		key = t.skPi
		nonce = t.Nr
	}
	macedID := t.suite.Prf.Apply(key, id)
	signB := append(append(initB, nonce.Bytes()...), macedID...)
	return signB
}
