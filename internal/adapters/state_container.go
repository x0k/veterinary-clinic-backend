package adapters

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type StateId string

type StateSaver[S any] interface {
	Save(state S) StateId
}

type ManualStateSaver[S any] interface {
	SaveByKey(key StateId, value S)
}

type StateLoader[S any] interface {
	Load(stateId StateId) (S, bool)
}

type TelegramDatePickerState struct {
	ServiceId entity.ServiceId
	Date      time.Time
}
