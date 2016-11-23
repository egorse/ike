// Code generated by "stringer -type AuthTransformId,DhTransformId,EncrTransformId,HashAlgorithmId,IdType,IkeExchangeType,NotificationType,PrfTransformId,AuthMethod -output protocol_strings.go"; DO NOT EDIT

package protocol

import "fmt"

const _AuthTransformId_name = "AUTH_NONEAUTH_HMAC_MD5_96AUTH_HMAC_SHA1_96AUTH_DES_MACAUTH_KPDK_MD5AUTH_AES_XCBC_96AUTH_HMAC_MD5_128AUTH_HMAC_SHA1_160AUTH_AES_CMAC_96AUTH_AES_128_GMACAUTH_AES_192_GMACAUTH_AES_256_GMACAUTH_HMAC_SHA2_256_128AUTH_HMAC_SHA2_384_192AUTH_HMAC_SHA2_512_256"

var _AuthTransformId_index = [...]uint8{0, 9, 25, 42, 54, 67, 83, 100, 118, 134, 151, 168, 185, 207, 229, 251}

func (i AuthTransformId) String() string {
	if i >= AuthTransformId(len(_AuthTransformId_index)-1) {
		return fmt.Sprintf("AuthTransformId(%d)", i)
	}
	return _AuthTransformId_name[_AuthTransformId_index[i]:_AuthTransformId_index[i+1]]
}

const (
	_DhTransformId_name_0 = "MODP_NONEMODP_768MODP_1024"
	_DhTransformId_name_1 = "MODP_1536"
	_DhTransformId_name_2 = "MODP_2048MODP_3072MODP_4096MODP_6144MODP_8192ECP_256ECP_384ECP_521MODP_1024_PRIME_160MODP_2048_PRIME_224MODP_2048_PRIME_256ECP_192ECP_224BRAINPOOLP224R1BRAINPOOLP256R1BRAINPOOLP384R1BRAINPOOLP512R1"
)

var (
	_DhTransformId_index_0 = [...]uint8{0, 9, 17, 26}
	_DhTransformId_index_1 = [...]uint8{0, 9}
	_DhTransformId_index_2 = [...]uint8{0, 9, 18, 27, 36, 45, 52, 59, 66, 85, 104, 123, 130, 137, 152, 167, 182, 197}
)

func (i DhTransformId) String() string {
	switch {
	case 0 <= i && i <= 2:
		return _DhTransformId_name_0[_DhTransformId_index_0[i]:_DhTransformId_index_0[i+1]]
	case i == 5:
		return _DhTransformId_name_1
	case 14 <= i && i <= 30:
		i -= 14
		return _DhTransformId_name_2[_DhTransformId_index_2[i]:_DhTransformId_index_2[i+1]]
	default:
		return fmt.Sprintf("DhTransformId(%d)", i)
	}
}

const (
	_EncrTransformId_name_0 = "ENCR_DES_IV64ENCR_DESENCR_3DESENCR_RC5ENCR_IDEAENCR_CASTENCR_BLOWFISHENCR_3IDEAENCR_DES_IV32"
	_EncrTransformId_name_1 = "ENCR_NULLENCR_AES_CBCENCR_AES_CTRAEAD_AES_CCM_SHORT_8AEAD_AES_CCM_SHORT_12AEAD_AES_CCM_SHORT_16"
	_EncrTransformId_name_2 = "AEAD_AES_GCM_8AEAD_AES_GCM_12AEAD_AES_GCM_16ENCR_NULL_AUTH_AES_GMAC"
	_EncrTransformId_name_3 = "ENCR_CAMELLIA_CBCENCR_CAMELLIA_CTRENCR_CAMELLIA_CCM_8_ICVENCR_CAMELLIA_CCM_12_ICVENCR_CAMELLIA_CCM_16_ICVAEAD_CHACHA20_POLY1305"
)

var (
	_EncrTransformId_index_0 = [...]uint8{0, 13, 21, 30, 38, 47, 56, 69, 79, 92}
	_EncrTransformId_index_1 = [...]uint8{0, 9, 21, 33, 53, 74, 95}
	_EncrTransformId_index_2 = [...]uint8{0, 14, 29, 44, 67}
	_EncrTransformId_index_3 = [...]uint8{0, 17, 34, 57, 81, 105, 127}
)

func (i EncrTransformId) String() string {
	switch {
	case 1 <= i && i <= 9:
		i -= 1
		return _EncrTransformId_name_0[_EncrTransformId_index_0[i]:_EncrTransformId_index_0[i+1]]
	case 11 <= i && i <= 16:
		i -= 11
		return _EncrTransformId_name_1[_EncrTransformId_index_1[i]:_EncrTransformId_index_1[i+1]]
	case 18 <= i && i <= 21:
		i -= 18
		return _EncrTransformId_name_2[_EncrTransformId_index_2[i]:_EncrTransformId_index_2[i+1]]
	case 23 <= i && i <= 28:
		i -= 23
		return _EncrTransformId_name_3[_EncrTransformId_index_3[i]:_EncrTransformId_index_3[i+1]]
	default:
		return fmt.Sprintf("EncrTransformId(%d)", i)
	}
}

const _HashAlgorithmId_name = "HASH_RESERVEDHASH_SHA1HASH_SHA2_256HASH_SHA2_384HASH_SHA2_512"

var _HashAlgorithmId_index = [...]uint8{0, 13, 22, 35, 48, 61}

func (i HashAlgorithmId) String() string {
	if i >= HashAlgorithmId(len(_HashAlgorithmId_index)-1) {
		return fmt.Sprintf("HashAlgorithmId(%d)", i)
	}
	return _HashAlgorithmId_name[_HashAlgorithmId_index[i]:_HashAlgorithmId_index[i+1]]
}

const (
	_IdType_name_0 = "ID_IPV4_ADDRID_FQDNID_RFC822_ADDR"
	_IdType_name_1 = "ID_IPV6_ADDR"
	_IdType_name_2 = "ID_DER_ASN1_DNID_DER_ASN1_GNID_KEY_ID"
)

var (
	_IdType_index_0 = [...]uint8{0, 12, 19, 33}
	_IdType_index_1 = [...]uint8{0, 12}
	_IdType_index_2 = [...]uint8{0, 14, 28, 37}
)

func (i IdType) String() string {
	switch {
	case 1 <= i && i <= 3:
		i -= 1
		return _IdType_name_0[_IdType_index_0[i]:_IdType_index_0[i+1]]
	case i == 5:
		return _IdType_name_1
	case 9 <= i && i <= 11:
		i -= 9
		return _IdType_name_2[_IdType_index_2[i]:_IdType_index_2[i+1]]
	default:
		return fmt.Sprintf("IdType(%d)", i)
	}
}

const _IkeExchangeType_name = "IKE_SA_INITIKE_AUTHCREATE_CHILD_SAINFORMATIONALIKE_SESSION_RESUMEGSA_AUTHGSA_REGISTRATIONGSA_REKEY"

var _IkeExchangeType_index = [...]uint8{0, 11, 19, 34, 47, 65, 73, 89, 98}

func (i IkeExchangeType) String() string {
	i -= 34
	if i >= IkeExchangeType(len(_IkeExchangeType_index)-1) {
		return fmt.Sprintf("IkeExchangeType(%d)", i+34)
	}
	return _IkeExchangeType_name[_IkeExchangeType_index[i]:_IkeExchangeType_index[i+1]]
}

const _NotificationType_name = "UNSUPPORTED_CRITICAL_PAYLOADINVALID_IKE_SPIINVALID_MAJOR_VERSIONINVALID_SYNTAXINVALID_MESSAGE_IDINVALID_SPINO_PROPOSAL_CHOSENINVALID_KE_PAYLOADAUTHENTICATION_FAILEDSINGLE_PAIR_REQUIREDNO_ADDITIONAL_SASINTERNAL_ADDRESS_FAILUREFAILED_CP_REQUIREDTS_UNACCEPTABLEINVALID_SELECTORSTEMPORARY_FAILURECHILD_SA_NOT_FOUNDINITIAL_CONTACTSET_WINDOW_SIZEADDITIONAL_TS_POSSIBLEIPCOMP_SUPPORTEDNAT_DETECTION_SOURCE_IPNAT_DETECTION_DESTINATION_IPCOOKIEUSE_TRANSPORT_MODEHTTP_CERT_LOOKUP_SUPPORTEDREKEY_SAESP_TFC_PADDING_NOT_SUPPORTEDNON_FIRST_FRAGMENTS_ALSOMOBIKE_SUPPORTEDADDITIONAL_IP4_ADDRESSADDITIONAL_IP6_ADDRESSNO_ADDITIONAL_ADDRESSESUPDATE_SA_ADDRESSESCOOKIE2NO_NATS_ALLOWEDAUTH_LIFETIMEMULTIPLE_AUTH_SUPPORTEDANOTHER_AUTH_FOLLOWSREDIRECT_SUPPORTEDREDIRECTREDIRECTED_FROMTICKET_LT_OPAQUETICKET_REQUESTTICKET_ACKTICKET_NACKTICKET_OPAQUELINK_IDUSE_WESP_MODEROHC_SUPPORTEDEAP_ONLY_AUTHENTICATIONCHILDLESS_IKEV2_SUPPORTEDQUICK_CRASH_DETECTIONIKEV2_MESSAGE_ID_SYNC_SUPPORTEDIPSEC_REPLAY_COUNTER_SYNC_SUPPORTEDIKEV2_MESSAGE_ID_SYNCIPSEC_REPLAY_COUNTER_SYNCSECURE_PASSWORD_METHODSPSK_PERSISTPSK_CONFIRMERX_SUPPORTEDIFOM_CAPABILITYSENDER_REQUEST_IDIKEV2_FRAGMENTATION_SUPPORTEDSIGNATURE_HASH_ALGORITHMS"

var _NotificationType_map = map[NotificationType]string{
	1:     _NotificationType_name[0:28],
	4:     _NotificationType_name[28:43],
	5:     _NotificationType_name[43:64],
	7:     _NotificationType_name[64:78],
	9:     _NotificationType_name[78:96],
	11:    _NotificationType_name[96:107],
	14:    _NotificationType_name[107:125],
	17:    _NotificationType_name[125:143],
	24:    _NotificationType_name[143:164],
	34:    _NotificationType_name[164:184],
	35:    _NotificationType_name[184:201],
	36:    _NotificationType_name[201:225],
	37:    _NotificationType_name[225:243],
	38:    _NotificationType_name[243:258],
	39:    _NotificationType_name[258:275],
	43:    _NotificationType_name[275:292],
	44:    _NotificationType_name[292:310],
	16384: _NotificationType_name[310:325],
	16385: _NotificationType_name[325:340],
	16386: _NotificationType_name[340:362],
	16387: _NotificationType_name[362:378],
	16388: _NotificationType_name[378:401],
	16389: _NotificationType_name[401:429],
	16390: _NotificationType_name[429:435],
	16391: _NotificationType_name[435:453],
	16392: _NotificationType_name[453:479],
	16393: _NotificationType_name[479:487],
	16394: _NotificationType_name[487:516],
	16395: _NotificationType_name[516:540],
	16396: _NotificationType_name[540:556],
	16397: _NotificationType_name[556:578],
	16398: _NotificationType_name[578:600],
	16399: _NotificationType_name[600:623],
	16400: _NotificationType_name[623:642],
	16401: _NotificationType_name[642:649],
	16402: _NotificationType_name[649:664],
	16403: _NotificationType_name[664:677],
	16404: _NotificationType_name[677:700],
	16405: _NotificationType_name[700:720],
	16406: _NotificationType_name[720:738],
	16407: _NotificationType_name[738:746],
	16408: _NotificationType_name[746:761],
	16409: _NotificationType_name[761:777],
	16410: _NotificationType_name[777:791],
	16411: _NotificationType_name[791:801],
	16412: _NotificationType_name[801:812],
	16413: _NotificationType_name[812:825],
	16414: _NotificationType_name[825:832],
	16415: _NotificationType_name[832:845],
	16416: _NotificationType_name[845:859],
	16417: _NotificationType_name[859:882],
	16418: _NotificationType_name[882:907],
	16419: _NotificationType_name[907:928],
	16420: _NotificationType_name[928:959],
	16421: _NotificationType_name[959:994],
	16422: _NotificationType_name[994:1015],
	16423: _NotificationType_name[1015:1040],
	16424: _NotificationType_name[1040:1063],
	16425: _NotificationType_name[1063:1074],
	16426: _NotificationType_name[1074:1085],
	16427: _NotificationType_name[1085:1098],
	16428: _NotificationType_name[1098:1113],
	16429: _NotificationType_name[1113:1130],
	16430: _NotificationType_name[1130:1159],
	16431: _NotificationType_name[1159:1184],
}

func (i NotificationType) String() string {
	if str, ok := _NotificationType_map[i]; ok {
		return str
	}
	return fmt.Sprintf("NotificationType(%d)", i)
}

const _PrfTransformId_name = "PRF_HMAC_MD5PRF_HMAC_SHA1PRF_HMAC_TIGERPRF_AES128_XCBCPRF_HMAC_SHA2_256PRF_HMAC_SHA2_384PRF_HMAC_SHA2_512PRF_AES128_CMAC"

var _PrfTransformId_index = [...]uint8{0, 12, 25, 39, 54, 71, 88, 105, 120}

func (i PrfTransformId) String() string {
	i -= 1
	if i >= PrfTransformId(len(_PrfTransformId_index)-1) {
		return fmt.Sprintf("PrfTransformId(%d)", i+1)
	}
	return _PrfTransformId_name[_PrfTransformId_index[i]:_PrfTransformId_index[i+1]]
}

const (
	_AuthMethod_name_0 = "AUTH_RSA_DIGITAL_SIGNATUREAUTH_SHARED_KEY_MESSAGE_INTEGRITY_CODEAUTH_DSS_DIGITAL_SIGNATURE"
	_AuthMethod_name_1 = "AUTH_ECDSA_256AUTH_ECDSA_384AUTH_ECDSA_521"
	_AuthMethod_name_2 = "AUTH_DIGITAL_SIGNATURE"
)

var (
	_AuthMethod_index_0 = [...]uint8{0, 26, 64, 90}
	_AuthMethod_index_1 = [...]uint8{0, 14, 28, 42}
	_AuthMethod_index_2 = [...]uint8{0, 22}
)

func (i AuthMethod) String() string {
	switch {
	case 1 <= i && i <= 3:
		i -= 1
		return _AuthMethod_name_0[_AuthMethod_index_0[i]:_AuthMethod_index_0[i+1]]
	case 9 <= i && i <= 11:
		i -= 9
		return _AuthMethod_name_1[_AuthMethod_index_1[i]:_AuthMethod_index_1[i+1]]
	case i == 14:
		return _AuthMethod_name_2
	default:
		return fmt.Sprintf("AuthMethod(%d)", i)
	}
}
