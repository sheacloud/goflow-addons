package transport

import (
	"fmt"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
	flowutils "github.com/cloudflare/goflow/v3/utils"
)

type OneToManyState struct {
	Transports []flowutils.Transport
}

func (s OneToManyState) Publish(msgs []*flowmessage.FlowMessage) {
	fmt.Println(len(msgs))
	for _, transport := range s.Transports {
		transport.Publish(msgs)
	}
}
