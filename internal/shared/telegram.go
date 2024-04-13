package shared

// var ErrNotTelegramUser = errors.New("not a telegram user")

type TelegramUserId int64

func NewTelegramUserId(id int64) TelegramUserId {
	return TelegramUserId(id)
}

func (id TelegramUserId) Int() int64 {
	return int64(id)
}

// func NewTelegramUser(
// 	id TelegramUserId,
// 	username string,
// 	firstName string,
// 	lastName string,
// 	phoneNumber string,
// ) User {
// 	return User{
// 		Id:          TelegramUserIdToUserId(id),
// 		Name:        fmt.Sprintf("%s %s", firstName, lastName),
// 		PhoneNumber: phoneNumber,
// 		Email:       fmt.Sprintf("https://t.me/%s", username),
// 	}
// }

// func IsTelegramUserId(userId UserId) bool {
// 	return strings.HasPrefix(string(userId), "tg-")
// }

// func UserIdToTelegramUserId(userId UserId) (TelegramUserId, error) {
// 	if !IsTelegramUserId(userId) {
// 		fmt.Println(string(userId))
// 		return 0, ErrNotTelegramUser
// 	}
// 	id, err := strconv.ParseInt(string(userId)[3:], 10, 64)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return TelegramUserId(id), nil
// }
