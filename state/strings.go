// generated by stringer -type IkeSaState,IkeEventId -output strings.go; DO NOT EDIT

package state

import "fmt"

const _IkeSaState_name = "SMI_INITSMI_AUTHSMI_AUTH_WAITSMI_AUTH_PEERSMI_EAPSMI_INSTALLCSA_DLSMI_INSTALLCSASMR_INITSMR_AUTHSMR_AUTH_FINALIZESMR_AUTH_RESPONSE_IDSMR_AUTH_RESPONSESMR_EAP_INITATOR_REQUESTSMR_EAP_AAA_REQUESTSMR_AUTH_DL_PEERSMR_CFG_WAITSM_MATURESM_REKEYSM_CRL_UPDATESM_REAUTHSM_TERMINATESM_DYINGSM_DEAD"

var _IkeSaState_index = [...]uint16{0, 8, 16, 29, 42, 49, 66, 80, 88, 96, 113, 133, 150, 174, 193, 209, 221, 230, 238, 251, 260, 272, 280, 287}

func (i IkeSaState) String() string {
	i -= 1
	if i < 0 || i+1 >= IkeSaState(len(_IkeSaState_index)) {
		return fmt.Sprintf("IkeSaState(%d)", i+1)
	}
	return _IkeSaState_name[_IkeSaState_index[i]:_IkeSaState_index[i+1]]
}

const _IkeEventId_name = "ACQUIRECONNECTREAUTHN_COOKIEN_INVALID_KEN_NO_PROPOSAL_CHOSENIKE_SA_INIT_RESPONSEIKE_AUTH_RESPONSEDELETE_IKE_SA_RESPONSEIKE_SA_INIT_SUCCESSIKE_AUTH_SUCCESSDELETE_IKE_SA_SUCCESSINVALID_KEMSG_IKE_REKEYMSG_IKE_REKEY_RESPONSEMSG_IKE_DPDMSG_IKE_CRL_UPDATEMSG_IKE_REAUTHMSG_IKE_TERMINATEIKE_SA_DELETE_REQUESTIKE_TIMEOUTStateEntryStateExit"

var _IkeEventId_index = [...]uint16{0, 7, 14, 20, 28, 40, 60, 80, 97, 119, 138, 154, 175, 185, 198, 220, 231, 249, 263, 280, 301, 312, 322, 331}

func (i IkeEventId) String() string {
	i -= 1
	if i < 0 || i+1 >= IkeEventId(len(_IkeEventId_index)) {
		return fmt.Sprintf("IkeEventId(%d)", i+1)
	}
	return _IkeEventId_name[_IkeEventId_index[i]:_IkeEventId_index[i+1]]
}
