package enrichers

import (
	"net"

	"github.com/sheacloud/goflow-addons/utils"
)

type FlowDirection struct {
	Direction  string `json:",omitempty"` // inbound, outbound, local, external
	IsResponse bool
}

type FlowDirectionEnricher struct {
	LocalNetworks []net.IPNet
}

func (e *FlowDirectionEnricher) isIPLocal(ip net.IP) bool {
	for _, net := range e.LocalNetworks {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}

func (e *FlowDirectionEnricher) determineFlowDirection(msg *utils.ExtendedFlowMessage) string {
	srcIsLocal := e.isIPLocal(net.IP(msg.SrcAddr))
	dstIsLocal := e.isIPLocal(net.IP(msg.DstAddr))

	if srcIsLocal && dstIsLocal {
		return "local"
	} else if srcIsLocal {
		return "outbound"
	} else if dstIsLocal {
		return "inbound"
	}

	return "external"

}

func (e *FlowDirectionEnricher) determineIsResponse(msg *utils.ExtendedFlowMessage) bool {
	if msg.Proto == 6 || msg.Proto == 17 {
		if msg.DstPort > msg.SrcPort {
			return true
		}
	}
	return false
}

func (e *FlowDirectionEnricher) Enrich(msgs []*utils.ExtendedFlowMessage) {
	for _, msg := range msgs {
		flowDirection := FlowDirection{
			Direction:  e.determineFlowDirection(msg),
			IsResponse: e.determineIsResponse(msg),
		}
		msg.Metadata["FlowDirection"] = flowDirection
	}
}
