package models

type ClinicServiceProperty string

const (
	ServiceTitle       ClinicServiceProperty = "Наименование"
	ServiceDuration    ClinicServiceProperty = "Продолжительность в минутах"
	ServiceDescription ClinicServiceProperty = "Описание"
	ServiceCost        ClinicServiceProperty = "Стоимость"
)

type ClinicRecordProperty string

const (
	RecordTitle          ClinicRecordProperty = "ФИО"
	RecordService        ClinicRecordProperty = "Услуга"
	RecordPhoneNumber    ClinicRecordProperty = "Телефон"
	RecordEmail          ClinicRecordProperty = "Почта"
	RecordDateTimePeriod ClinicRecordProperty = "Время записи"
	RecordState          ClinicRecordProperty = "Статус"
	RecordUserId         ClinicRecordProperty = "identity"
)

type ClinicRecordStatus string

const (
	ClinicRecordAwaits    ClinicRecordStatus = "Ожидает"
	ClinicRecordInWork    ClinicRecordStatus = "В работе"
	ClinicRecordDone      ClinicRecordStatus = "Выполнено"
	ClinicRecordNotAppear ClinicRecordStatus = "Не пришел"
)
