package like

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jimeux/go-generic-dao/db"
)

type DAO struct{}

const createStmt = "INSERT INTO " + Table + " (`user_id`, `partner_id`, `created_at`) VALUES (?, ?, ?);"

func (DAO) Create(ctx context.Context, l Like) (Like, error) {
	res, err := db.DB().ExecContext(ctx, createStmt, &l.UserID, &l.PartnerID, &l.CreatedAt)
	if err != nil {
		return l, fmt.Errorf("like.DAO#Create ExecContext error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return l, fmt.Errorf("like.DAO#Create LastInsertId error: %w", err)
	}
	l.ID = id
	return l, nil
}

const getByPairQuery = "SELECT " + Columns + " FROM " + Table + " WHERE `user_id` = ? AND `partner_id` = ? LIMIT 1;"

func (DAO) GetByPair(ctx context.Context, userID, partnerID int64) (Like, error) {
	row := db.DB().QueryRowContext(ctx, getByPairQuery, userID, partnerID)
	var l Like
	if err := row.Scan(l.PtrFields()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return l, db.ErrNotFound
		}
		return l, fmt.Errorf("like.DAO#Get row.Scan error: %w", err)
	}
	return l, nil
}

const countQuery = "SELECT COUNT(*) FROM " + Table + ";"

func (DAO) Count(ctx context.Context) (int64, error) {
	row := db.DB().QueryRowContext(ctx, countQuery)
	var count int64
	if err := row.Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("like.DAO#Count row.Scan error: %w", err)
	}
	return count, nil
}

const findByUserQuery = "SELECT " + Columns + " FROM " + Table + " WHERE `user_id` = ? ORDER BY `id`;"

func (DAO) FindByUser(ctx context.Context, userID int64) ([]Like, error) {
	rows, err := db.DB().QueryContext(ctx, findByUserQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("like.DAO#Find failed\n%s: %w", findByUserQuery, err)
	}
	defer func() { _ = rows.Close() }()

	var result []Like
	for rows.Next() {
		var l Like
		if err := rows.Scan(l.PtrFields()...); err != nil {
			return nil, fmt.Errorf("like.DAO#Find scan error: %w", err)
		}
		result = append(result, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("like.DAO#Find rows.Err(): %w", err)
	}
	return result, nil
}
