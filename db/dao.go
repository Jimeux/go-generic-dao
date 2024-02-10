package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"
	"strings"
	"time"
)

// Row is a type constraint for types representing
// a single database row.
type Row[T any] interface {
	// PtrFields returns all fields of a struct for use with row.Scan.
	// It must be implemented with a pointer receiver type, and all elements
	// in the returned slice must also be pointers.
	PtrFields() []any
	// *T allows creating instances of type T as pointers.
	// See: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#pointer-method-example.
	*T
}

// Column is a type constraint for types representing
// a single database column.
type Column interface {
	~byte | ~int16 | ~int32 | ~int64 | ~float64 |
		~string | ~bool | time.Time
}

// Create executes the insert statement found in q.
// It returns the last inserted ID if any.
func Create(ctx context.Context, q string, args ...any) (int64, error) {
	res, err := database.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("DAO#Create failed\n%s: %w", q, err)
	}
	return res.LastInsertId()
}

// GetRow returns a single row from the database as type T.
// If no row is found, ErrNotFound is returned with a default T value.
func GetRow[T any, PT Row[T]](ctx context.Context, q string, args ...any) (T, error) {
	row := database.QueryRowContext(ctx, q, args...)
	var t T
	ptr := PT(&t)
	if err := row.Scan(ptr.PtrFields()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, ErrNotFound
		}
		return t, fmt.Errorf("DAO#GetRow row.Scan error\n%s: %w", q, err)
	}
	return t, nil
}

// GetColumn returns a single column value from the database as type T.
// If no column is found, a default T value is returned with no error.
func GetColumn[T Column](ctx context.Context, q string, args ...any) (T, error) {
	row := database.QueryRowContext(ctx, q, args...)
	var t T
	if err := row.Scan(&t); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, ErrNotFound
		}
		return t, fmt.Errorf("DAO#GetColumn row.Scan error\n%s: %w", q, err)
	}
	return t, nil
}

// FindRows executes q as a query, and returns an iterator of type (T, error).
// T must implement a ToPtrArgs method with a pointer receiver type.
func FindRows[T any, PT Row[T]](ctx context.Context, q string, args ...any) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		rows, err := database.QueryContext(ctx, q, args...)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return // empty iterator
			}
			var t T
			yield(t, fmt.Errorf("DAO#FindRows QueryContext failed\n%s: %w", q, err))
			return
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var t T
			ptr := PT(&t)
			if err := rows.Scan(ptr.PtrFields()...); err != nil {
				yield(t, fmt.Errorf("DAO#FindRows row.Scan error\n%s: %w", q, err))
				return
			}
			if !yield(t, nil) {
				break
			}
		}
		if err := rows.Err(); err != nil {
			var t T
			yield(t, fmt.Errorf("DAO#FindRows rows.Err()\n%s: %w", q, err))
		}
	}
}

// FindColumns executes q as a query, and returns an iterator of type (T, error).
// T must satisfy the Column constraint.
func FindColumns[T Column](ctx context.Context, q string, args ...any) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		rows, err := database.QueryContext(ctx, q, args...)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return // empty iterator
			}
			var t T
			yield(t, fmt.Errorf("DAO#FindColumns QueryContext failed\n%s: %w", q, err))
			return
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var t T
			if err := rows.Scan(&t); err != nil {
				yield(t, fmt.Errorf("DAO#FindColumns row.Scan error\n%s: %w", q, err))
				return
			}
			if !yield(t, nil) {
				break
			}
		}
		if err := rows.Err(); err != nil {
			var t T
			yield(t, fmt.Errorf("DAO#FindColumns rows.Err()\n%s: %w", q, err))
		}
	}
}

// InArgs returns placeholders and args formatted for a WHERE IN clause.
// Calling InArgs([]int{1,2,3}) will return ("?,?,?", []any{1,2,3}).
func InArgs[T Column](tt []T) (string, []any) {
	args := make([]any, len(tt))
	for i, t := range tt {
		args[i] = t
	}
	return strings.Repeat("?,", len(args)-1) + "?", args
}
