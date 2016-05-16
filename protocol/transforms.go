package protocol

var (
	_ENCR_AES_CTR      = Transform{Type: TRANSFORM_TYPE_ENCR, TransformId: uint16(ENCR_AES_CTR)}
	_ENCR_AES_CBC      = Transform{Type: TRANSFORM_TYPE_ENCR, TransformId: uint16(ENCR_AES_CBC)}
	_ENCR_CAMELLIA_CBC = Transform{Type: TRANSFORM_TYPE_ENCR, TransformId: uint16(ENCR_CAMELLIA_CBC)}
	_ENCR_CAMELLIA_CTR = Transform{Type: TRANSFORM_TYPE_ENCR, TransformId: uint16(ENCR_CAMELLIA_CTR)}
	_ENCR_NULL         = Transform{Type: TRANSFORM_TYPE_ENCR, TransformId: uint16(ENCR_NULL)}

	_AEAD_AES_GCM_16 = Transform{Type: TRANSFORM_TYPE_ENCR, TransformId: uint16(AEAD_AES_GCM_16)}

	_PRF_AES128_XCBC   = Transform{Type: TRANSFORM_TYPE_PRF, TransformId: uint16(PRF_AES128_XCBC)}
	_PRF_HMAC_SHA1     = Transform{Type: TRANSFORM_TYPE_PRF, TransformId: uint16(PRF_HMAC_SHA1)}
	_PRF_HMAC_SHA2_256 = Transform{Type: TRANSFORM_TYPE_PRF, TransformId: uint16(PRF_HMAC_SHA2_256)}
	_PRF_HMAC_SHA2_384 = Transform{Type: TRANSFORM_TYPE_PRF, TransformId: uint16(PRF_HMAC_SHA2_384)}
	_PRF_HMAC_SHA2_512 = Transform{Type: TRANSFORM_TYPE_PRF, TransformId: uint16(PRF_HMAC_SHA2_512)}

	_AUTH_AES_XCBC_96       = Transform{Type: TRANSFORM_TYPE_INTEG, TransformId: uint16(AUTH_AES_XCBC_96)}
	_AUTH_HMAC_SHA1_96      = Transform{Type: TRANSFORM_TYPE_INTEG, TransformId: uint16(AUTH_HMAC_SHA1_96)}
	_AUTH_HMAC_SHA2_256_128 = Transform{Type: TRANSFORM_TYPE_INTEG, TransformId: uint16(AUTH_HMAC_SHA2_256_128)}
	_AUTH_HMAC_SHA2_384_192 = Transform{Type: TRANSFORM_TYPE_INTEG, TransformId: uint16(AUTH_HMAC_SHA2_384_192)}
	_AUTH_HMAC_SHA2_512_256 = Transform{Type: TRANSFORM_TYPE_INTEG, TransformId: uint16(AUTH_HMAC_SHA2_512_256)}

	_MODP_768  = Transform{Type: TRANSFORM_TYPE_DH, TransformId: uint16(MODP_768)}
	_MODP_1024 = Transform{Type: TRANSFORM_TYPE_DH, TransformId: uint16(MODP_1024)}
	_MODP_1536 = Transform{Type: TRANSFORM_TYPE_DH, TransformId: uint16(MODP_1536)}
	_MODP_2048 = Transform{Type: TRANSFORM_TYPE_DH, TransformId: uint16(MODP_2048)}
	_MODP_3072 = Transform{Type: TRANSFORM_TYPE_DH, TransformId: uint16(MODP_3072)}
	_MODP_4096 = Transform{Type: TRANSFORM_TYPE_DH, TransformId: uint16(MODP_4096)}
	_MODP_6144 = Transform{Type: TRANSFORM_TYPE_DH, TransformId: uint16(MODP_6144)}
	_MODP_8192 = Transform{Type: TRANSFORM_TYPE_DH, TransformId: uint16(MODP_8192)}

	_ESN    = Transform{Type: TRANSFORM_TYPE_ESN, TransformId: uint16(ESN)}
	_NO_ESN = Transform{Type: TRANSFORM_TYPE_ESN, TransformId: uint16(ESN_NONE)}
)

var transforms = map[Transform]string{
	_ENCR_AES_CBC:      "ENCR_AES_CBC",
	_ENCR_AES_CTR:      "ENCR_AES_CTR",
	_ENCR_CAMELLIA_CBC: "ENCR_CAMELLIA_CBC",
	_ENCR_CAMELLIA_CTR: "ENCR_CAMELLIA_CTR",
	_ENCR_NULL:         "ENCR_NULL",

	_AEAD_AES_GCM_16: "AEAD_AES_GCM_16",

	_PRF_AES128_XCBC:   "PRF_AES128_XCBC",
	_PRF_HMAC_SHA1:     "PRF_HMAC_SHA1",
	_PRF_HMAC_SHA2_256: "PRF_HMAC_SHA2_256",
	_PRF_HMAC_SHA2_384: "PRF_HMAC_SHA2_384",
	_PRF_HMAC_SHA2_512: "PRF_HMAC_SHA2_512",

	_AUTH_AES_XCBC_96:       "AUTH_AES_XCBC_96",
	_AUTH_HMAC_SHA1_96:      "AUTH_HMAC_SHA1_96",
	_AUTH_HMAC_SHA2_256_128: "AUTH_HMAC_SHA2_256_128",
	_AUTH_HMAC_SHA2_384_192: "AUTH_HMAC_SHA2_384_192",
	_AUTH_HMAC_SHA2_512_256: "AUTH_HMAC_SHA2_512_256",

	_MODP_768:  "MODP_768",
	_MODP_1024: "MODP_1024",
	_MODP_1536: "MODP_1536",
	_MODP_2048: "MODP_2048",
	_MODP_3072: "MODP_3072",
	_MODP_4096: "MODP_4096",
	_MODP_6144: "MODP_6144",
	_MODP_8192: "MODP_8192",

	_ESN:    "ESN",
	_NO_ESN: "NO_ESN",
}

type Transforms map[TransformType]*SaTransform

var (
	IKE_AES_CBC_SHA1_96_DH_1024 = Transforms{
		TRANSFORM_TYPE_ENCR:  &SaTransform{Transform: _ENCR_AES_CBC, KeyLength: 128},
		TRANSFORM_TYPE_PRF:   &SaTransform{Transform: _PRF_HMAC_SHA1},
		TRANSFORM_TYPE_INTEG: &SaTransform{Transform: _AUTH_HMAC_SHA1_96},
		TRANSFORM_TYPE_DH:    &SaTransform{Transform: _MODP_1024, IsLast: true},
	}

	// key length is set to 128b
	// 16B icv
	IKE_AES_GCM_16_DH_1024 = Transforms{
		TRANSFORM_TYPE_ENCR:  &SaTransform{Transform: _AEAD_AES_GCM_16, KeyLength: 128}, // AEAD_AES_128_GCM
		TRANSFORM_TYPE_INTEG: &SaTransform{Transform: _PRF_HMAC_SHA1},
		TRANSFORM_TYPE_DH:    &SaTransform{Transform: _MODP_1024, IsLast: true},
	}

	IKE_AES_GCM_16_DH_2048 = Transforms{
		TRANSFORM_TYPE_ENCR:  &SaTransform{Transform: _AEAD_AES_GCM_16, KeyLength: 128}, // AEAD_AES_128_GCM
		TRANSFORM_TYPE_INTEG: &SaTransform{Transform: _PRF_HMAC_SHA1},
		TRANSFORM_TYPE_DH:    &SaTransform{Transform: _MODP_2048, IsLast: true},
	}

	IKE_CAMELLIA_CBC_SHA2_256_128_DH_2048 = Transforms{
		TRANSFORM_TYPE_ENCR:  &SaTransform{Transform: _ENCR_CAMELLIA_CBC, KeyLength: 128},
		TRANSFORM_TYPE_PRF:   &SaTransform{Transform: _PRF_HMAC_SHA2_256},
		TRANSFORM_TYPE_INTEG: &SaTransform{Transform: _AUTH_HMAC_SHA2_256_128},
		TRANSFORM_TYPE_DH:    &SaTransform{Transform: _MODP_2048, IsLast: true},
	}

	ESP_AES_CBC_SHA1_96 = Transforms{
		TRANSFORM_TYPE_ENCR:  &SaTransform{Transform: _ENCR_AES_CBC, KeyLength: 128},
		TRANSFORM_TYPE_INTEG: &SaTransform{Transform: _AUTH_HMAC_SHA1_96},
		TRANSFORM_TYPE_ESN:   &SaTransform{Transform: _NO_ESN, IsLast: true},
	}

	// key length is set to 128
	// This is due to lack of support for 256b keys in older kernels (Jessie)
	ESP_AES_GCM_16 = Transforms{
		TRANSFORM_TYPE_ENCR: &SaTransform{Transform: _AEAD_AES_GCM_16, KeyLength: 128},
		TRANSFORM_TYPE_ESN:  &SaTransform{Transform: _NO_ESN, IsLast: true},
	}

	ESP_NULL_SHA1_96 = Transforms{
		TRANSFORM_TYPE_ENCR:  &SaTransform{Transform: _ENCR_NULL},
		TRANSFORM_TYPE_INTEG: &SaTransform{Transform: _AUTH_HMAC_SHA1_96},
		TRANSFORM_TYPE_ESN:   &SaTransform{Transform: _NO_ESN, IsLast: true},
	}

	ESP_CAMELLIA_CBC_SHA2_256_128 = Transforms{
		TRANSFORM_TYPE_ENCR:  &SaTransform{Transform: _ENCR_CAMELLIA_CBC, KeyLength: 128},
		TRANSFORM_TYPE_INTEG: &SaTransform{Transform: _AUTH_HMAC_SHA2_256_128},
		TRANSFORM_TYPE_ESN:   &SaTransform{Transform: _ESN, IsLast: true},
	}
)

func listHas(trsList []*SaTransform, trs *SaTransform) bool {
	for _, tr := range trsList {
		if trs.IsEqual(tr) {
			return true
		}
	}
	return false
}

func (configured Transforms) AsList() (trs []*SaTransform) {
	for _, trsVal := range configured {
		trs = append(trs, trsVal)
	}
	return
}

// checks if the configured set of transforms occurs within list of porposed transforms
func (configured Transforms) Within(trs []*SaTransform) bool {
	for _, trsVal := range configured {
		if !listHas(trs, trsVal) {
			return false
		}
	}
	return true
}

func (tr *SaTransform) IsEqual(other *SaTransform) bool {
	if tr == nil || other == nil {
		return false
	}
	if tr.KeyLength != other.KeyLength {
		return false
	}
	if tr.Transform.Type != other.Transform.Type {
		return false
	}
	if tr.Transform.TransformId != other.Transform.TransformId {
		return false
	}
	return true
}

// mutualTransform returns a cipherSuite given
// a list requested by the peer.
func mutualTransform(want [][]*SaTransform) bool {
	for _, w := range want {
	next:
		for {
			for _, t := range w {
				if _, ok := transforms[t.Transform]; !ok {
					break next
				}
			}
			// have all
			return true
		}
	}
	return false
}
