package protocol

import "fmt"

type IkeErrorCode uint16

type IkeError struct {
	IkeErrorCode
	Message string
}

func ErrF(e IkeErrorCode, format string, a ...interface{}) IkeError {
	return IkeError{e, fmt.Sprintf(format, a...)}
}

func (e IkeError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s:%s", e.IkeErrorCode, e.Message)
	}
	return fmt.Sprintf("%s", e.IkeErrorCode)
}

const (
	ERR_UNSUPPORTED_CRITICAL_PAYLOAD IkeErrorCode = 1
	ERR_INVALID_IKE_SPI              IkeErrorCode = 4
	ERR_INVALID_MAJOR_VERSION        IkeErrorCode = 5
	ERR_INVALID_SYNTAX               IkeErrorCode = 7
	ERR_INVALID_MESSAGE_ID           IkeErrorCode = 9
	ERR_INVALID_SPI                  IkeErrorCode = 11
	ERR_NO_PROPOSAL_CHOSEN           IkeErrorCode = 14
	ERR_INVALID_KE_PAYLOAD           IkeErrorCode = 17
	ERR_AUTHENTICATION_FAILED        IkeErrorCode = 24
	ERR_SINGLE_PAIR_REQUIRED         IkeErrorCode = 34
	ERR_NO_ADDITIONAL_SAS            IkeErrorCode = 35
	ERR_INTERNAL_ADDRESS_FAILURE     IkeErrorCode = 36
	ERR_FAILED_CP_REQUIRED           IkeErrorCode = 37
	ERR_TS_UNACCEPTABLE              IkeErrorCode = 38
	ERR_INVALID_SELECTORS            IkeErrorCode = 39
	ERR_TEMPORARY_FAILURE            IkeErrorCode = 43
	ERR_CHILD_SA_NOT_FOUND           IkeErrorCode = 44
)

func GetIkeErrorCode(nt NotificationType) (IkeErrorCode, bool) {
	switch nt {
	case UNSUPPORTED_CRITICAL_PAYLOAD:
		return ERR_UNSUPPORTED_CRITICAL_PAYLOAD, true
	case INVALID_IKE_SPI:
		return ERR_INVALID_IKE_SPI, true
	case INVALID_MAJOR_VERSION:
		return ERR_INVALID_MAJOR_VERSION, true
	case INVALID_SYNTAX:
		return ERR_INVALID_SYNTAX, true
	case INVALID_MESSAGE_ID:
		return ERR_INVALID_MESSAGE_ID, true
	case INVALID_SPI:
		return ERR_INVALID_SPI, true
	case NO_PROPOSAL_CHOSEN:
		return ERR_NO_PROPOSAL_CHOSEN, true
	case INVALID_KE_PAYLOAD:
		return ERR_INVALID_KE_PAYLOAD, true
	case AUTHENTICATION_FAILED:
		return ERR_AUTHENTICATION_FAILED, true
	case SINGLE_PAIR_REQUIRED:
		return ERR_SINGLE_PAIR_REQUIRED, true
	case NO_ADDITIONAL_SAS:
		return ERR_NO_ADDITIONAL_SAS, true
	case INTERNAL_ADDRESS_FAILURE:
		return ERR_INTERNAL_ADDRESS_FAILURE, true
	case FAILED_CP_REQUIRED:
		return ERR_FAILED_CP_REQUIRED, true
	case TS_UNACCEPTABLE:
		return ERR_TS_UNACCEPTABLE, true
	case INVALID_SELECTORS:
		return ERR_INVALID_SELECTORS, true
	case TEMPORARY_FAILURE:
		return ERR_TEMPORARY_FAILURE, true
	case CHILD_SA_NOT_FOUND:
		return ERR_CHILD_SA_NOT_FOUND, true
	default:
		return 0, false
	}
}

func (e IkeErrorCode) Error() string {
	switch e {
	case ERR_UNSUPPORTED_CRITICAL_PAYLOAD:
		return "UNSUPPORTED_CRITICAL_PAYLOAD"
	case ERR_INVALID_IKE_SPI:
		return "INVALID_IKE_SPI"
	case ERR_INVALID_MAJOR_VERSION:
		return "INVALID_MAJOR_VERSION"
	case ERR_INVALID_SYNTAX:
		return "INVALID_SYNTAX"
	case ERR_INVALID_MESSAGE_ID:
		return "INVALID_MESSAGE_ID"
	case ERR_INVALID_SPI:
		return "INVALID_SPI"
	case ERR_NO_PROPOSAL_CHOSEN:
		return "NO_PROPOSAL_CHOSEN"
	case ERR_INVALID_KE_PAYLOAD:
		return "INVALID_KE_PAYLOAD"
	case ERR_AUTHENTICATION_FAILED:
		return "AUTHENTICATION_FAILED"
	case ERR_SINGLE_PAIR_REQUIRED:
		return "SINGLE_PAIR_REQUIRED"
	case ERR_NO_ADDITIONAL_SAS:
		return "NO_ADDITIONAL_SAS"
	case ERR_INTERNAL_ADDRESS_FAILURE:
		return "INTERNAL_ADDRESS_FAILURE"
	case ERR_FAILED_CP_REQUIRED:
		return "FAILED_CP_REQUIRED"
	case ERR_TS_UNACCEPTABLE:
		return "TS_UNACCEPTABLE"
	case ERR_INVALID_SELECTORS:
		return "INVALID_SELECTORS"
	case ERR_TEMPORARY_FAILURE:
		return "TEMPORARY_FAILURE"
	case ERR_CHILD_SA_NOT_FOUND:
		return "CHILD_SA_NOT_FOUND"
	default:
		return "Unknown Error"
	}
}
