package transport

import (
	"fmt"

	"github.com/sheacloud/goflow-addons/utils"
)

type SysoutState struct {
}

func (s SysoutState) Publish(msgs []*utils.ExtendedFlowMessage) {
	for _, msg := range msgs {
		str, _ := utils.HumanReadableJSONMarshalIndent(msg, "", "  ")
		fmt.Println(string(str))
	}
}
