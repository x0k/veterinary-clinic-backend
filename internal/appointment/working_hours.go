package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type WorkingHours map[time.Weekday]entity.TimePeriod
