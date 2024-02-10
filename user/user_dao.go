package user

import (
	"context"
	"fmt"
	"iter"

	"github.com/Jimeux/go-generic-dao/db"
)

type DAO struct{}

const createStmt = "INSERT INTO " + Table + " (`nickname`, `created_at`) VALUES (?, ?);"

func (DAO) Create(ctx context.Context, u User) (User, error) {
	id, err := db.Create(ctx, createStmt, &u.Nickname, &u.CreatedAt)
	u.ID = id
	return u, err
}

const getByIDQuery = "SELECT " + Columns + " FROM " + Table + " WHERE `id` = ? LIMIT 1;"

func (DAO) GetByID(ctx context.Context, id int64) (User, error) {
	return db.GetRow[User](ctx, getByIDQuery, id)
}

const countQuery = "SELECT COUNT(*) FROM " + Table + ";"

func (DAO) Count(ctx context.Context) (int64, error) {
	return db.GetColumn[int64](ctx, countQuery)
}

const findByIDsQuery = "SELECT " + Columns + " FROM " + Table + " WHERE `id` IN (%s) ORDER BY `id`;"

func (DAO) FindByIDs(ctx context.Context, ids []int64) iter.Seq2[User, error] {
	placeholders, args := db.InArgs(ids)
	q := fmt.Sprintf(findByIDsQuery, placeholders)
	return db.FindRows[User](ctx, q, args...)
}

const findIDsWithBioQuery = "SELECT `id` FROM " + Table + " WHERE `bio` IS NOT NULL;"

func (DAO) FindIDsWithBio(ctx context.Context) iter.Seq2[int64, error] {
	return db.FindColumns[int64](ctx, findIDsWithBioQuery)
}
