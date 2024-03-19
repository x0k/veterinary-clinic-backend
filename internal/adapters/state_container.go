package adapters

type StateId string

type StateSaver[S any] interface {
	Save(state S) StateId
}

type StateLoader[S any] interface {
	Load(stateId StateId) (S, bool)
}
