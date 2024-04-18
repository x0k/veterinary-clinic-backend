package adapters

type StateId string

func (s StateId) String() string {
	return string(s)
}

func NewStateId(str string) StateId {
	return StateId(str)
}

type StateSaver[S any] interface {
	Save(state S) StateId
}

type StateByKeySaver[S any] interface {
	SaveByKey(key StateId, state S)
}

type StateLoader[S any] interface {
	Load(stateId StateId) (S, bool)
}

type StatePopper[S any] interface {
	Pop(stateId StateId) (S, bool)
}
