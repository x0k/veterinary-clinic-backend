package appointment

import (
	"reflect"
	"testing"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

func TestCalculateSchedulePeriods(t *testing.T) {
	type args struct {
		freePeriods FreeTimeSlots
		busyPeriods BusyPeriods
		workBreaks  DayWorkBreaks
	}
	tests := []struct {
		name string
		args args
		want ScheduleEntries
	}{
		{
			name: "Vacation",
			args: args{
				freePeriods: FreeTimeSlots{
					{
						Start: shared.Time{
							Hours:   9,
							Minutes: 30,
						},
						End: shared.Time{
							Hours:   17,
							Minutes: 0,
						},
					},
				},
				busyPeriods: BusyPeriods{},
				workBreaks: DayWorkBreaks{
					{
						Id:              "lunch",
						MatchExpression: `^[1-5]`,
						Title:           "Перерыв на обед",
						Period: shared.TimePeriod{
							Start: shared.Time{
								Hours:   12,
								Minutes: 30,
							},
							End: shared.Time{
								Hours:   13,
								Minutes: 30,
							},
						},
					},
					{
						Id:              "vacation",
						MatchExpression: `^\d 2024-03-(2[6-9]|30|31)`,
						Title:           "Отпуск",
						Period: shared.TimePeriod{
							Start: shared.Time{
								Hours:   0,
								Minutes: 0,
							},
							End: shared.Time{
								Hours:   23,
								Minutes: 59,
							},
						},
					},
				},
			},
			want: ScheduleEntries{
				{
					DateTimePeriod: shared.DateTimePeriod{
						Start: shared.DateTime{
							Time: shared.Time{
								Hours:   0,
								Minutes: 0,
							},
						},
						End: shared.DateTime{
							Time: shared.Time{
								Hours:   23,
								Minutes: 59,
							},
						},
					},
					Type:  BusyPeriod,
					Title: "Отпуск",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newScheduleEntries(shared.Date{}, tt.args.freePeriods, tt.args.busyPeriods, tt.args.workBreaks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateSchedulePeriods() = %v, want %v", got, tt.want)
			}
		})
	}
}
