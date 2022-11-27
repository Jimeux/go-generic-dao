package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/Jimeux/go-generic-dao/db"
	"github.com/Jimeux/go-generic-dao/test"
)

func defaultUser() User {
	return User{
		ID:        1,
		Nickname:  "Socks",
		Bio:       sql.NullString{},
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}
}

func TestMain(m *testing.M) {
	test.InitConfig()
	db.Init()
	m.Run()
	db.Close()
}

func TestUserDAO_Create(t *testing.T) {
	ctx := context.Background()
	dao := DAO{}

	t.Run("Create", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		want := defaultUser()
		got, err := dao.Create(ctx, want)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, want) {
			t.Fatal(cmp.Diff(want, got))
		}
	})
}

func TestUserDAO_GetByID(t *testing.T) {
	ctx := context.Background()
	dao := DAO{}

	t.Run("found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		want := defaultUser()
		want, err := dao.Create(ctx, want)
		if err != nil {
			t.Fatal(err)
		}
		got, err := dao.GetByID(ctx, want.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, want) {
			t.Fatal(cmp.Diff(want, got))
		}
	})
	t.Run("not found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		got, err := dao.GetByID(ctx, 1000)
		if !errors.Is(err, db.ErrNotFound) {
			t.Fatal(err)
		}
		if got != (User{}) {
			t.Fatalf("expected default value but got %+v", got)
		}
	})
}

func TestUserDAO_Count(t *testing.T) {
	ctx := context.Background()
	dao := DAO{}

	t.Run("records found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		if _, err := dao.Create(ctx, defaultUser()); err != nil {
			t.Fatal(err)
		}
		got, err := dao.Count(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if got != 1 {
			t.Fatalf("got %d but want 1", got)
		}
	})
	t.Run("no records", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		got, err := dao.Count(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if got != 0 {
			t.Fatalf("got %d but want 0", got)
		}
	})
}

func TestUserDAO_FindByIDs(t *testing.T) {
	ctx := context.Background()
	dao := DAO{}

	t.Run("found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		user1 := defaultUser()
		user1, err := dao.Create(ctx, user1)
		if err != nil {
			t.Fatal(err)
		}
		user2 := defaultUser()
		user2, err = dao.Create(ctx, user2)
		if err != nil {
			t.Fatal(err)
		}

		want := []User{user1, user2}
		got, err := dao.FindByIDs(ctx, []int64{user1.ID, user2.ID})
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, want) {
			t.Fatalf(cmp.Diff(want, got))
		}
	})
	t.Run("not found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		got, err := dao.FindByIDs(ctx, []int64{1000})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("expected got to be empty but was len=%d", len(got))
		}
	})
}

func TestUserDAO_FindIDsWithBio(t *testing.T) {
	ctx := context.Background()
	dao := DAO{}

	t.Run("not found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		got, err := dao.FindIDsWithBio(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("expected got to be empty but was len=%d", len(got))
		}
	})
}
