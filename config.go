package ike

import (
	"errors"
	"net"

	"github.com/msgboxio/ike/protocol"
	"github.com/msgboxio/log"
)

type Config struct {
	ProposalIke, ProposalEsp protocol.Transforms

	TsI, TsR []*protocol.Selector

	IsTransportMode bool
}

func DefaultConfig() *Config {
	return &Config{
		ProposalIke: protocol.IKE_AES_CBC_SHA256_MODP2048,
		// ProposalIke: protocol.IKE_AES_CBC_SHA256_MODP3072,
		// ProposalIke: protocol.IKE_AES_GCM_16_MODP3072,
		ProposalEsp: protocol.ESP_AES_CBC_SHA2_256,
		// ProposalEsp: protocol.ESP_AES_GCM_16,
	}
}

// CheckProposals checks if incoming proposals include our configuration
func (cfg *Config) CheckProposals(prot protocol.ProtocolId, proposals protocol.Proposals) error {
	for _, prop := range proposals {
		if prop.ProtocolId != prot {
			continue
		}
		// select first acceptable one from the list
		switch prot {
		case protocol.IKE:
			if cfg.ProposalIke.Within(prop.SaTransforms) {
				return nil
			}
		case protocol.ESP:
			if cfg.ProposalEsp.Within(prop.SaTransforms) {
				return nil
			}
		}
	}
	return errors.New("acceptable proposals are missing")
}

// AddSelector builds selector from address & mask
func (cfg *Config) AddSelector(initiator, responder *net.IPNet) (err error) {
	first, last, err := IPNetToFirstLastAddress(initiator)
	if err != nil {
		return
	}
	cfg.TsI = []*protocol.Selector{&protocol.Selector{
		Type:         protocol.TS_IPV4_ADDR_RANGE,
		IpProtocolId: 0,
		StartPort:    0,
		Endport:      65535,
		StartAddress: first,
		EndAddress:   last,
	}}
	first, last, err = IPNetToFirstLastAddress(responder)
	if err != nil {
		return
	}
	cfg.TsR = []*protocol.Selector{&protocol.Selector{
		Type:         protocol.TS_IPV4_ADDR_RANGE,
		IpProtocolId: 0,
		StartPort:    0,
		Endport:      65535,
		StartAddress: first,
		EndAddress:   last,
	}}
	return
}

// CheckFromInit takes an IkeSaInit message and checks
// if acceptable IKE proposal is available
func (cfg *Config) CheckFromInit(initI *Message) error {
	// get SA payload
	ikeSa := initI.Payloads.Get(protocol.PayloadTypeSA).(*protocol.SaPayload)
	return cfg.CheckProposals(protocol.IKE, ikeSa.Proposals)
}

// CheckromAuth checks esp proposal & selector
func (cfg *Config) CheckromAuth(authI *Message) error {
	espSa := authI.Payloads.Get(protocol.PayloadTypeSA).(*protocol.SaPayload)
	if err := cfg.CheckProposals(protocol.ESP, espSa.Proposals); err != nil {
		return err
	}
	// get selectors
	tsI := authI.Payloads.Get(protocol.PayloadTypeTSi).(*protocol.TrafficSelectorPayload).Selectors
	tsR := authI.Payloads.Get(protocol.PayloadTypeTSr).(*protocol.TrafficSelectorPayload).Selectors
	if len(tsI) == 0 || len(tsR) == 0 {
		return errors.New("acceptable traffic selectors are missing")
	}
	log.Infof("Configured selectors: [INI]%s<=>%s[RES]", cfg.TsI, cfg.TsR)
	log.Infof("Offered selectors: [INI]%s<=>%s[RES]", tsI, tsR)
	// TODO - check selectors
	return nil
}

func ProposalFromTransform(prot protocol.ProtocolId, trs protocol.Transforms, spi []byte) []*protocol.SaProposal {
	return []*protocol.SaProposal{
		&protocol.SaProposal{
			IsLast:       true,
			Number:       1,
			ProtocolId:   prot,
			Spi:          append([]byte{}, spi...),
			SaTransforms: trs.AsList(),
		},
	}
}
