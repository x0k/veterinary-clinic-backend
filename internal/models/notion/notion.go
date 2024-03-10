package notion

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
	Awaits    ClinicRecordStatus = "Ожидает"
	InWork    ClinicRecordStatus = "В работе"
	Done      ClinicRecordStatus = "Выполнено"
	NotAppear ClinicRecordStatus = "Не пришел"
)
