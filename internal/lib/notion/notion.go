package notion

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/period"
)

var ErrInvalidPeriodDate = errors.New("invalid period date")

func RichTextValue(richText []notionapi.RichText) string {
	strs := make([]string, 0, len(richText))
	for _, r := range richText {
		if r.Type != notionapi.ObjectTypeText {
			continue
		}
		strs = append(strs, r.Text.Content)
	}
	return strings.Join(strs, "")
}

func Title(properties notionapi.Properties, titleKey string) string {
	return RichTextValue(properties[titleKey].(*notionapi.TitleProperty).Title)
}

func Number(properties notionapi.Properties, numberKey string) float64 {
	return properties[numberKey].(*notionapi.NumberProperty).Number
}

func Text(properties notionapi.Properties, stringKey string) string {
	return RichTextValue(properties[stringKey].(*notionapi.RichTextProperty).RichText)
}

func Date(properties notionapi.Properties, dateKey string) *notionapi.DateObject {
	return properties[dateKey].(*notionapi.DateProperty).Date
}

func DatePeriod(properties notionapi.Properties, periodKey string) (period.Period[time.Time], error) {
	date := Date(properties, periodKey)
	if date == nil || date.Start == nil || date.End == nil {
		return period.Period[time.Time]{}, fmt.Errorf("%s: %w", periodKey, ErrInvalidPeriodDate)
	}
	return period.Period[time.Time]{
		Start: time.Time(*date.Start),
		End:   time.Time(*date.End),
	}, nil
}

func Phone(properties notionapi.Properties, phoneKey string) string {
	return properties[phoneKey].(*notionapi.PhoneNumberProperty).PhoneNumber
}

func Email(properties notionapi.Properties, emailKey string) string {
	return properties[emailKey].(*notionapi.EmailProperty).Email
}

func Relations(properties notionapi.Properties, relationKey string) []notionapi.Relation {
	return properties[relationKey].(*notionapi.RelationProperty).Relation
}

func Select(properties notionapi.Properties, selectKey string) string {
	return properties[selectKey].(*notionapi.SelectProperty).Select.Name
}

func CreatedTime(properties notionapi.Properties, createdTimeKey string) time.Time {
	return properties[createdTimeKey].(*notionapi.CreatedTimeProperty).CreatedTime
}

func ToRichText(value string) []notionapi.RichText {
	return []notionapi.RichText{
		{
			Type: notionapi.ObjectTypeText,
			Text: &notionapi.Text{Content: value},
		},
	}
}
