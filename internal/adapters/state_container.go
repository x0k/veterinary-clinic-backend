package adapters

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type StateId string

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

type TelegramDatePickerState struct {
	ServiceId entity.ServiceId
	Date      time.Time
}
