package transport

import (
	"github.com/sheacloud/goflow-addons/utils"
)

type NullState struct {
}

func (s NullState) Publish(msgs []*utils.ExtendedFlowMessage) {
	return
}
