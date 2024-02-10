package like

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/Jimeux/go-generic-dao/db"
	"github.com/Jimeux/go-generic-dao/test"
)

func defaultLike() Like {
	return Like{
		ID:        1,
		UserID:    1,
		PartnerID: 2,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}
}

func TestMain(m *testing.M) {
	test.InitConfig()
	db.Init()
	m.Run()
	db.Close()
}

func TestLikeDAO_Create(t *testing.T) {
	ctx := context.Background()
	dao := DAO{}

	t.Run("Create", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		want := defaultLike()
		got, err := dao.Create(ctx, want)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, want) {
			t.Fatal(cmp.Diff(want, got))
		}
	})
}

func TestLikeDAO_GetByPair(t *testing.T) {
	ctx := context.Background()
	dao := DAO{}

	t.Run("found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		want := defaultLike()
		want, err := dao.Create(ctx, want)
		if err != nil {
			t.Fatal(err)
		}
		got, err := dao.GetByPair(ctx, want.UserID, want.PartnerID)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, want) {
			t.Fatal(cmp.Diff(want, got))
		}
	})
	t.Run("not found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		got, err := dao.GetByPair(ctx, 1000, 1001)
		if !errors.Is(err, db.ErrNotFound) {
			t.Fatal(err)
		}
		if got != (Like{}) {
			t.Fatalf("expected default value but got %+v", got)
		}
	})

}

func TestLikeDAO_Count(t *testing.T) {
	ctx := context.Background()
	dao := DAO{}

	t.Run("records counted", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		if _, err := dao.Create(ctx, defaultLike()); err != nil {
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

func TestLikeDAO_FindByUser(t *testing.T) {
	ctx := context.Background()
	dao := DAO{}

	t.Run("found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		like1 := defaultLike()
		like1.PartnerID = 10
		like1, err := dao.Create(ctx, like1)
		if err != nil {
			t.Fatal(err)
		}
		like2 := defaultLike()
		like2.PartnerID = 11
		like2, err = dao.Create(ctx, like2)
		if err != nil {
			t.Fatal(err)
		}

		want := []Like{like1, like2}
		var got []Like
		for l, err := range dao.FindByUser(ctx, like2.UserID) {
			if err != nil {
				t.Fatal(err)
			}
			got = append(got, l)
		}
		if !cmp.Equal(got, want) {
			t.Fatalf(cmp.Diff(want, got))
		}
	})
	t.Run("not found", func(t *testing.T) {
		t.Cleanup(test.Truncate(t, Table))
		for got, err := range dao.FindByUser(ctx, 1000) {
			t.Fatalf("want empty iterator, got %v, %v", got, err)
		}
	})
}
