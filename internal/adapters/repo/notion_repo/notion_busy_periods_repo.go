package notion_repo

import (
	"context"
	"errors"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

var ErrLoadingBusyPeriodsFailed = errors.New("failed to load busy periods")

type BusyPeriods struct {
	log               *logger.Logger
	client            *notionapi.Client
	recordsDatabaseId notionapi.DatabaseID
}

func NewBusyPeriods(
	log *logger.Logger,
	client *notionapi.Client,
	recordsDatabaseId notionapi.DatabaseID,
) *BusyPeriods {
	return &BusyPeriods{
		log:               log,
		client:            client,
		recordsDatabaseId: recordsDatabaseId,
	}
}

func (s *BusyPeriods) BusyPeriods(ctx context.Context, t time.Time) ([]entity.TimePeriod, error) {
	after := t.Add(
		-(time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second),
	)
	afterDate := notionapi.Date(after)
	beforeDate := notionapi.Date(after.AddDate(0, 0, 1))
	r, err := s.client.Database.Query(ctx, s.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: RecordDateTimePeriod,
				Date: &notionapi.DateFilterCondition{
					After:  &afterDate,
					Before: &beforeDate,
				},
			},
			notionapi.OrCompoundFilter{
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: ClinicRecordInWork,
					},
				},
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: ClinicRecordAwaits,
					},
				},
			},
		},
		Sorts: []notionapi.SortObject{
			{
				Property:  RecordDateTimePeriod,
				Direction: notionapi.SortOrderASC,
			},
		},
	})
	if err != nil {
		s.log.Error(ctx, "failed to load busy periods", sl.Err(err))
		return nil, ErrLoadingBusyPeriodsFailed
	}
	periods := make([]entity.TimePeriod, 0, len(r.Results))
	for _, page := range r.Results {
		if period := DateTimePeriodFromRecord(page.Properties); period != nil {
			periods = append(periods, entity.TimePeriod{
				Start: period.Start.Time,
				End:   period.End.Time,
			})
		}
	}
	return periods, nil
}
