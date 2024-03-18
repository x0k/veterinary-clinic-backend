package adapters

type StateId string

type StateEncoder[S comparable] interface {
	Encode(state S) StateId
}

type StateDecoder[S comparable] interface {
	Decode(stateId StateId) (S, bool)
}
