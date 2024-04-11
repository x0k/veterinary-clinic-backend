package telegram_adapters

type Token string

func (t Token) String() string {
	return string(t)
}
