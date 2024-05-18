package adapters

type StateId string

func (s StateId) String() string {
	return string(s)
}

func NewStateId(str string) StateId {
	return StateId(str)
}

type StateSaver[S any] func(state S) StateId

type StateByKeySaver[S any] func(key StateId, state S)

type StateLoader[S any] func(stateId StateId) (S, bool)

type StatePopper[S any] func(stateId StateId) (S, bool)
