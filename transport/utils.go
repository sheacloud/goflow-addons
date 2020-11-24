package transport

import (
	"encoding/json"
	"net"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type FlowMessageAlias struct {
	SrcAddr        string `json:"SrcAddr,omitempty"`
	DstAddr        string `json:"DstAddr,omitempty"`
	NextHop        string `json:"NextHop,omitempty"`
	SamplerAddress string `json:"SamplerAddress,omitempty"`
	*flowmessage.FlowMessage
}

func convertFlowMessageToHumanFormat(msg *flowmessage.FlowMessage) interface{} {
	//TODO figure out how to check if srcaddr/dstaddr/etc are set and conditionally add them to new struct
	newStruct := FlowMessageAlias{
		SrcAddr:        net.IP(msg.SrcAddr).String(),
		DstAddr:        net.IP(msg.DstAddr).String(),
		NextHop:        net.IP(msg.NextHop).String(),
		SamplerAddress: net.IP(msg.SamplerAddress).String(),
		FlowMessage:    (msg),
	}

	return newStruct
}

func HumanReadableJSONMarshal(msg *flowmessage.FlowMessage) ([]byte, error) {
	newStruct := convertFlowMessageToHumanFormat(msg)

	return json.Marshal(newStruct)
}

func HumanReadableJSONMarshalIndent(msg *flowmessage.FlowMessage, prefix, indent string) ([]byte, error) {
	newStruct := convertFlowMessageToHumanFormat(msg)

	return json.MarshalIndent(newStruct, prefix, indent)
}
