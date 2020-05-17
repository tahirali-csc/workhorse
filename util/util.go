package util

import (
	"bytes"
	"encoding/gob"
	"workhorse/api"
)

func ConvertToWorkflowObject(data []byte) *api.WorkflowTransferObject {
	dec := gob.NewDecoder(bytes.NewReader(data))

	j := &api.WorkflowTransferObject{}
	dec.Decode(j)
	return j
}
