package protocol

import "github.com/msgboxio/packets"

func (s *IdPayload) Type() PayloadType {
	return s.IdPayloadType
}

func (s *IdPayload) Encode() (b []byte) {
	b = []byte{uint8(s.IdType), 0, 0, 0}
	return append(b, s.Data...)
}

func (s *IdPayload) Decode(b []byte) (err error) {
	if len(b) < 4 {
		err = ErrF(ERR_INVALID_SYNTAX, "id too small %d < %d", len(b), 4)
		return
	}
	// Header has already been decoded
	Idt, _ := packets.ReadB8(b, 0)
	s.IdType = IdType(Idt)
	s.Data = append([]byte{}, b[4:]...)
	return
}
