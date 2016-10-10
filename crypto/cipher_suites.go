package crypto

import (
	"fmt"

	"github.com/msgboxio/ike/protocol"
	"github.com/msgboxio/log"
)

// Cipher interface provides Encryption & Integrity Protection
type Cipher interface {
	Overhead(clear []byte) int
	VerifyDecrypt(ike, skA, skE []byte) (dec []byte, err error)
	EncryptMac(headers, payload, skA, skE []byte) (b []byte, err error)
}

type CipherSuite struct {
	Cipher  // aead or nonAead
	Prf     *Prf
	DhGroup dhGroup

	// Lengths, in bytes, of the key material needed for each component.
	KeyLen, MacKeyLen int
}

// Build a CipherSuite from the given transfom
// TODO - check that the entire suite makes sense
func NewCipherSuite(trs protocol.Transforms) (*CipherSuite, error) {
	cs := &CipherSuite{}
	// empty variables, filled in later
	var aead *aeadCipher
	var cipher *simpleCipher

	for _, tr := range trs {
		switch tr.Transform.Type {
		case protocol.TRANSFORM_TYPE_DH:
			dh, ok := kexAlgoMap[protocol.DhTransformId(tr.Transform.TransformId)]
			if !ok {
				return nil, fmt.Errorf("Unsupported dh transfom %d", tr.Transform.TransformId)
			}
			cs.DhGroup = dh
		case protocol.TRANSFORM_TYPE_PRF:
			// for hmac based Prf, preferred key size is size of output
			prf, err := prfTranform(tr.Transform.TransformId)
			if err != nil {
				return nil, err
			}
			cs.Prf = prf
		case protocol.TRANSFORM_TYPE_ENCR:
			keyLen := int(tr.KeyLength) / 8 // from attribute; in bits
			var ok bool
			if cipher, ok = cipherTransform(tr.Transform.TransformId, keyLen, cipher); !ok {
				if aead, keyLen, ok = aeadTransform(tr.Transform.TransformId, keyLen, aead); !ok {
					return nil, fmt.Errorf("Unsupported cipher transfom %d", tr.Transform.TransformId)
				}
			}
			cs.KeyLen = keyLen // TODO - 2 places
		case protocol.TRANSFORM_TYPE_INTEG:
			var ok bool
			if cipher, ok = integrityTransform(tr.Transform.TransformId, cipher); !ok {
				return nil, fmt.Errorf("Unsupported mac transfom %d", tr.Transform.TransformId)
			}
			cs.MacKeyLen = cipher.macKeyLen // TODO - 2 places
		case protocol.TRANSFORM_TYPE_ESN:
		// nothing
		default:
			return nil, fmt.Errorf("Unsupported transfom type %d", tr.Transform.Type)
		} // end switch
	} // end loop
	if cipher == nil && aead == nil {
		return nil, fmt.Errorf("cipher transfoms were not set")
	}
	if cipher != nil && aead != nil {
		return nil, fmt.Errorf("invalid cipher transfoms combination")
	}
	if cipher != nil {
		cs.Cipher = cipher
	}
	if aead != nil {
		cs.Cipher = aead
	}
	return cs, nil
}

func (cs *CipherSuite) CheckIkeTransforms() error {
	if cs.DhGroup == nil || cs.Prf == nil {
		return fmt.Errorf("invalid cipher transfoms combination")
	}
	if log.V(2) {
		log.Infof("IKE CipherSuite: %+v", *cs)
	}
	return nil
}

func (cs *CipherSuite) CheckEspTransforms() error {
	if log.V(2) {
		log.Infof("ESP CipherSuite: %+v", *cs)
	}
	return nil
}
