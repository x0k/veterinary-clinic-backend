package entity

import "maps"

type RecordChangeType int

const (
	RecordCreated RecordChangeType = iota
	RecordStatusChanged
	RecordDateTimeChanged
	RecordRemoved
)

type RecordChange struct {
	Type   RecordChangeType
	Record Record
}

type ActualRecordsState struct {
	records map[RecordId]Record
}

func NewActualRecordsState(records map[RecordId]Record) ActualRecordsState {
	return ActualRecordsState{
		records: records,
	}
}

func (r ActualRecordsState) Records() map[RecordId]Record {
	return r.records
}

func (r *ActualRecordsState) Update(actualRecords []Record) []RecordChange {
	stateCopy := maps.Clone(r.records)
	changes := make([]RecordChange, 0, len(actualRecords))
	for _, actualRecord := range actualRecords {
		r.records[actualRecord.Id] = actualRecord
		oldRecord, ok := stateCopy[actualRecord.Id]
		// created
		if !ok {
			changes = append(changes, RecordChange{
				Type:   RecordCreated,
				Record: actualRecord,
			})
			continue
		}
		if oldRecord.Status != actualRecord.Status {
			changes = append(changes, RecordChange{
				Type:   RecordStatusChanged,
				Record: actualRecord,
			})
		} else if DateTimePeriodApi.ComparePeriods(oldRecord.DateTimePeriod, actualRecord.DateTimePeriod) != 0 {
			changes = append(changes, RecordChange{
				Type:   RecordDateTimeChanged,
				Record: actualRecord,
			})
		}
		delete(stateCopy, actualRecord.Id)
	}
	for _, record := range stateCopy {
		delete(r.records, record.Id)
		changes = append(changes, RecordChange{
			Type:   RecordRemoved,
			Record: record,
		})
	}
	return changes
}
