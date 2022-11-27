package user

import (
	"database/sql"
	"time"
)

const (
	Table   = "`user`"
	Columns = "`id`, `nickname`, `bio`, `created_at`"
)

type User struct {
	ID        int64
	Nickname  string
	Bio       sql.NullString
	CreatedAt time.Time
}

func (u *User) PtrFields() []any {
	return []any{&u.ID, &u.Nickname, &u.Bio, &u.CreatedAt}
}
