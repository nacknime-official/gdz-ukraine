package sqlc

import "github.com/nacknime-official/gdz-ukraine/internal/entity"

func (u *User) ToEntity() *entity.User {
	return &entity.User{
		ID:                         int(u.ID),
		IsBlocked:                  u.IsBlocked,
		IsSubscribedToBroadcasting: u.IsSubscribedToBroadcasting,
	}
}
