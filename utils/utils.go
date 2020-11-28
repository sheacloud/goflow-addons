package utils

import (
	"encoding/json"
	"net"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type ExtendedTransport interface {
	Publish([]*ExtendedFlowMessage)
}

type Enricher interface {
	Enrich([]*ExtendedFlowMessage)
}

type FlowMessageAlias struct {
	SrcAddr        string `json:"SrcAddr,omitempty"`
	DstAddr        string `json:"DstAddr,omitempty"`
	NextHop        string `json:"NextHop,omitempty"`
	SamplerAddress string `json:"SamplerAddress,omitempty"`
	*ExtendedFlowMessage
}

type ExtendedFlowMessage struct {
	Metadata map[string]interface{} `json:"Metadata,omitempty"`
	*flowmessage.FlowMessage
}

func convertToExtendedMessage(msg *flowmessage.FlowMessage) *ExtendedFlowMessage {
	return &ExtendedFlowMessage{
		Metadata:    make(map[string]interface{}),
		FlowMessage: (msg),
	}
}

func ConvertMessages(msgs []*flowmessage.FlowMessage) []*ExtendedFlowMessage {
	extendedMsgs := []*ExtendedFlowMessage{}
	for _, msg := range msgs {
		if msg == nil {
			continue
		}
		extendedMsg := convertToExtendedMessage(msg)
		extendedMsgs = append(extendedMsgs, extendedMsg)
	}
	return extendedMsgs
}

func convertFlowMessageToHumanFormat(msg *ExtendedFlowMessage) interface{} {
	// TODO figure out how to check if srcaddr/dstaddr/etc are set and conditionally add them to new struct
	newStruct := FlowMessageAlias{
		SrcAddr:             net.IP(msg.SrcAddr).String(),
		DstAddr:             net.IP(msg.DstAddr).String(),
		NextHop:             net.IP(msg.NextHop).String(),
		SamplerAddress:      net.IP(msg.SamplerAddress).String(),
		ExtendedFlowMessage: (msg),
	}

	return newStruct
}

func HumanReadableJSONMarshal(msg *ExtendedFlowMessage) ([]byte, error) {
	newStruct := convertFlowMessageToHumanFormat(msg)

	return json.Marshal(newStruct)
}

func HumanReadableJSONMarshalIndent(msg *ExtendedFlowMessage, prefix, indent string) ([]byte, error) {
	newStruct := convertFlowMessageToHumanFormat(msg)

	return json.MarshalIndent(newStruct, prefix, indent)
}
