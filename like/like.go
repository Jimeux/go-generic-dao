package like

import "time"

const (
	Table   = "`like`"
	Columns = "`id`, `user_id`, `partner_id`, `created_at`"
)

type Like struct {
	ID        int64
	UserID    int64
	PartnerID int64
	CreatedAt time.Time
}

func (l *Like) PtrFields() []any {
	return []any{&l.ID, &l.UserID, &l.PartnerID, &l.CreatedAt}
}
