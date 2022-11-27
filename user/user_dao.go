package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Jimeux/go-generic-dao/db"
)

type DAO struct{}

const createStmt = "INSERT INTO " + Table + " (`nickname`, `created_at`) VALUES (?, ?);"

func (DAO) Create(ctx context.Context, u User) (User, error) {
	res, err := db.DB().ExecContext(ctx, createStmt, &u.Nickname, &u.CreatedAt)
	if err != nil {
		return u, fmt.Errorf("user.DAO#Create ExecContext error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return u, fmt.Errorf("user.DAO#Create LastInsertId error: %w", err)
	}
	u.ID = id
	return u, nil
}

const getByIDQuery = "SELECT " + Columns + " FROM " + Table + " WHERE `id` = ? LIMIT 1;"

func (DAO) GetByID(ctx context.Context, id int64) (User, error) {
	row := db.DB().QueryRowContext(ctx, getByIDQuery, id)
	var u User
	if err := row.Scan(u.PtrFields()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, db.ErrNotFound
		}
		return u, fmt.Errorf("user.DAO#GetByID row.Scan error: %w", err)
	}
	return u, nil
}

const countQuery = "SELECT COUNT(*) FROM " + Table + ";"

func (DAO) Count(ctx context.Context) (int64, error) {
	row := db.DB().QueryRowContext(ctx, countQuery)
	var count int64
	if err := row.Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("user.DAO#Count row.Scan error: %w", err)
	}
	return count, nil
}

const findByIDsQuery = "SELECT " + Columns + " FROM " + Table + " WHERE `id` IN (%s) ORDER BY `id`;"

func (DAO) FindByIDs(ctx context.Context, ids []int64) ([]User, error) {
	args := make([]any, len(ids))
	for i, t := range ids {
		args[i] = t
	}
	placeholders := strings.Repeat("?,", len(args)-1) + "?"

	q := fmt.Sprintf(findByIDsQuery, placeholders)
	rows, err := db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("user.DAO#FindByIDs failed\n%s: %w", q, err)
	}
	defer func() { _ = rows.Close() }()

	var result []User
	for rows.Next() {
		var u User
		if err := rows.Scan(u.PtrFields()...); err != nil {
			return nil, fmt.Errorf("user.DAO#FindByIDs scan error: %w", err)
		}
		result = append(result, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("user.DAO#FindByIDs rows.Err(): %w", err)
	}
	return result, nil
}

const findIDsWithBioQuery = "SELECT `id` FROM " + Table + " WHERE `bio` IS NOT NULL;"

func (DAO) FindIDsWithBio(ctx context.Context) ([]int64, error) {
	rows, err := db.DB().QueryContext(ctx, findIDsWithBioQuery)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("user.DAO#FindIDsWithBio failed\n%s: %w", findIDsWithBioQuery, err)
	}
	defer func() { _ = rows.Close() }()

	var result []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("user.DAO#FindIDsWithBio scan error: %w", err)
		}
		result = append(result, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("user.DAO#FindIDsWithBio rows.Err(): %w", err)
	}
	return result, nil
}
