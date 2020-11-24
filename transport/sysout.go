package transport

import (
	"fmt"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type SysoutState struct {
}

func (s SysoutState) Publish(msgs []*flowmessage.FlowMessage) {
	for _, msg := range msgs {
		str, _ := HumanReadableJSONMarshalIndent(msg, "", "  ")
		fmt.Println(string(str))
	}
}
