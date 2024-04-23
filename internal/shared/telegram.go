package shared

// var ErrNotTelegramUser = errors.New("not a telegram user")

type TelegramUserId int64

func NewTelegramUserId(id int64) TelegramUserId {
	return TelegramUserId(id)
}

func (id TelegramUserId) Int() int64 {
	return int64(id)
}
