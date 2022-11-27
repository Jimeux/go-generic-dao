package like

import (
	"context"

	"github.com/Jimeux/go-generic-dao/db"
)

type DAO struct{}

const createStmt = "INSERT INTO " + Table + " (`user_id`, `partner_id`, `created_at`) VALUES (?, ?, ?);"

func (DAO) Create(ctx context.Context, l Like) (Like, error) {
	id, err := db.Create(ctx, createStmt, &l.UserID, &l.PartnerID, &l.CreatedAt)
	l.ID = id
	return l, err
}

const getByPairQuery = "SELECT " + Columns + " FROM " + Table + " WHERE `user_id` = ? AND `partner_id` = ? LIMIT 1;"

func (DAO) GetByPair(ctx context.Context, userID, partnerID int64) (Like, error) {
	return db.GetRow[Like](ctx, getByPairQuery, userID, partnerID)
}

const countQuery = "SELECT COUNT(*) FROM " + Table + ";"

func (DAO) Count(ctx context.Context) (int64, error) {
	return db.GetColumn[int64](ctx, countQuery)
}

const findByUserQuery = "SELECT " + Columns + " FROM " + Table + " WHERE `user_id` = ? ORDER BY `id`;"

func (DAO) FindByUser(ctx context.Context, userID int64) ([]Like, error) {
	return db.FindRows[Like](ctx, findByUserQuery, userID)
}
