package appointment

import (
	"reflect"
	"testing"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

func TestCalculateSchedulePeriods(t *testing.T) {
	type args struct {
		freePeriods FreePeriods
		busyPeriods BusyPeriods
		workBreaks  CalculatedWorkBreaks
	}
	tests := []struct {
		name string
		args args
		want SchedulePeriods
	}{
		{
			name: "Vacation",
			args: args{
				freePeriods: FreePeriods{
					{
						Start: entity.Time{
							Hours:   9,
							Minutes: 30,
						},
						End: entity.Time{
							Hours:   17,
							Minutes: 0,
						},
					},
				},
				busyPeriods: BusyPeriods{},
				workBreaks: CalculatedWorkBreaks{
					{
						Id:              "lunch",
						MatchExpression: `^[1-5]`,
						Title:           "Перерыв на обед",
						Period: entity.TimePeriod{
							Start: entity.Time{
								Hours:   12,
								Minutes: 30,
							},
							End: entity.Time{
								Hours:   13,
								Minutes: 30,
							},
						},
					},
					{
						Id:              "vacation",
						MatchExpression: `^\d 2024-03-(2[6-9]|30|31)`,
						Title:           "Отпуск",
						Period: entity.TimePeriod{
							Start: entity.Time{
								Hours:   0,
								Minutes: 0,
							},
							End: entity.Time{
								Hours:   23,
								Minutes: 59,
							},
						},
					},
				},
			},
			want: SchedulePeriods{
				{
					TimePeriod: entity.TimePeriod{
						Start: entity.Time{
							Hours:   0,
							Minutes: 0,
						},
						End: entity.Time{
							Hours:   23,
							Minutes: 59,
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
			if got := NewSchedulePeriods(tt.args.freePeriods, tt.args.busyPeriods, tt.args.workBreaks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateSchedulePeriods() = %v, want %v", got, tt.want)
			}
		})
	}
}
