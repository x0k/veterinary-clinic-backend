package appointment_presenter

import (
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

func recordStatus(status appointment.RecordStatus) (string, error) {
	switch status {
	case appointment.RecordAwaits:
		return "ожидает", nil
	case appointment.RecordDone:
		return "выполнено", nil
	case appointment.RecordNotAppear:
		return "не пришел", nil
	default:
		return "", fmt.Errorf("unknown status: %s", status)
	}
}

func RecordState(status appointment.RecordStatus, isArchived bool) (string, error) {
	st, err := recordStatus(status)
	if err != nil {
		return "", err
	}
	if isArchived {
		return fmt.Sprintf("%s (архив)", st), nil
	}
	return st, nil
}
