package deltas

import (
	"github.com/satori/go.uuid"
)

type Delta struct {
	Id         string
	Operations []Operation
}

type Operation struct {
	Type   OperationType
	Record Record
}

type OperationType string

const (
	CREATE OperationType = "CREATE"
	UPDATE OperationType = "UPDATE"
	UPSERT OperationType = "UPSERT"
	REMOVE OperationType = "REMOVE"
)

type Record struct {
	Type    string
	Id      string
	Version int
	Value   interface{}
}

type DeltaHandler interface {
	OnDelta(delta Delta)
}

func New(operations []Operation) (delta Delta) {
	delta.Id = uuid.NewV1().String()
	delta.Operations = operations
	return
}
