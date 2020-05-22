package util

import (
	"bytes"
	"encoding/gob"
	"log"
	"workhorse/api"
)

func ConvertToWorkflowObject(data []byte) *api.WorkflowTransferObject {
	dec := gob.NewDecoder(bytes.NewReader(data))

	j := &api.WorkflowTransferObject{}
	dec.Decode(j)
	return j
}

func ConvertToJobObject(data []byte) *api.JobTransferObject {
	dec := gob.NewDecoder(bytes.NewReader(data))

	j := &api.JobTransferObject{}
	dec.Decode(j)
	return j
}

func ConvertToByteArray(workflow interface{}) []byte {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	err := enc.Encode(workflow)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return network.Bytes()
}
