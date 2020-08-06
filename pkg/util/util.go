package util

import (
	"bytes"
	"encoding/gob"
	"math/rand"
	"workhorse/pkg/api"
)

func ConvertToWorkflowObject(data []byte) *api.Workflow {
	dec := gob.NewDecoder(bytes.NewReader(data))

	j := &api.Workflow{}
	dec.Decode(j)
	return j
}

func ConvertToJobObject(data []byte) *api.WorkflowJob {
	dec := gob.NewDecoder(bytes.NewReader(data))

	j := &api.WorkflowJob{}
	dec.Decode(j)
	return j
}

func ConvertToByteArray(obj interface{}) ([]byte, error) {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return network.Bytes(), nil
}

func RandomBetween(min int, max int) int {
	return rand.Intn(max-min) + min
}
