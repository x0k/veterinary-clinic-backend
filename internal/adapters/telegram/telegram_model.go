package adapters_telegram

type Token string

func (t Token) String() string {
	return string(t)
}
