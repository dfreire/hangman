package deltas

type Commit struct {
	Number int
	Deltas []Delta
}

type Delta struct {
	Id            string
	Operation     Operation
	RecordType    string
	RecordVersion int
	Record        interface{}
}

type Operation string

const (
	CREATE Operation = "CREATE"
	UPDATE Operation = "UPDATE"
	UPSERT Operation = "UPSERT"
	REMOVE Operation = "REMOVE"
)

type DeltaHandler interface {
	OnDelta(delta Delta)
}

