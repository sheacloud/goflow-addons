package transport

import (
	"fmt"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
	"github.com/sheacloud/goflow-addons/utils"
)

type ExtendedWrapperState struct {
	ExtendedTransports []utils.ExtendedTransport
	Enrichers          []utils.Enricher
}

func (s ExtendedWrapperState) Publish(msgs []*flowmessage.FlowMessage) {
	fmt.Printf("received message with %v flows\n", len(msgs))
	extendedMsgs := utils.ConvertMessages(msgs)
	for _, enricher := range s.Enrichers {
		enricher.Enrich(extendedMsgs)
	}
	for _, transport := range s.ExtendedTransports {
		transport.Publish(extendedMsgs)
	}
}
