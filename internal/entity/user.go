package entity

type User struct {
	ID                         int
	TelegramID                 int64
	IsBlocked                  bool
	IsSubscribedToBroadcasting bool
}
