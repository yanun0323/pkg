package storage

import (
	"context"
	"math/big"
	"testing"

	"github.com/pkg/errors"
	"github.com/yanun0323/pkg/tester"
)

func TestNew(t *testing.T) {
	defer func() {
		tester.RequireNoError(t, Delete("./test_int.db"))
		tester.RequireNoError(t, Delete("./test_string.db"))
	}()

	{
		db, err := New[int]("./test_int.db")
		tester.RequireNoError(t, err)
		tester.RequireNotNil(t, db)
	}

	{
		db, err := New[string]("./test_string.db")
		tester.RequireNoError(t, err)
		tester.RequireNotNil(t, db)
		tester.RequireNoError(t, db.Close())
	}

	{
		db, err := New[bool]("./test_int.db")
		tester.RequireErrorIs(t, ErrTypeMismatch, err)
		tester.RequireNil(t, db)
	}
}

func TestLocal_CURD(t *testing.T) {
	defer func() {
		tester.RequireNoError(t, Delete("./test_curd_int.db"))
	}()

	ctx := context.Background()

	{
		db, err := New[int]("./test_curd_int.db")
		tester.RequireNoError(t, err)
		tester.RequireNotNil(t, db)

		{
			ok, err := db.Exists(ctx, "hello")
			tester.RequireNoError(t, err)
			tester.RequireFalse(t, ok)

			val, err := db.Get(ctx, "hello")
			tester.RequireErrorIs(t, ErrNotFound, err)
			tester.RequireEqual(t, 0, val)
		}

		{
			tester.RequireNoError(t, db.Set(ctx, "hello", 1))
			tester.RequireNoError(t, db.Set(ctx, "world", 2))
		}

		{
			val, err := db.Get(ctx, "hello")
			tester.RequireNoError(t, err)
			tester.RequireEqual(t, 1, val)

			val, err = db.Get(ctx, "world")
			tester.RequireNoError(t, err)
			tester.RequireEqual(t, 2, val)
		}

		{
			tester.RequireNoError(t, db.Delete(ctx, "hello"))

			ok, err := db.Exists(ctx, "hello")
			tester.RequireNoError(t, err)
			tester.RequireFalse(t, ok)

			val, err := db.Get(ctx, "hello")
			tester.RequireErrorIs(t, err, ErrNotFound)
			tester.RequireEqual(t, 0, val)
		}
	}

	defer func() {
		tester.RequireNoError(t, Delete("./test_curd_object.db"))
	}()

	type Order struct {
		ID            int
		RelativeOrder []*Order
		FilledAmount  map[string]*big.Int
	}

	order := &Order{
		ID: 1,
		RelativeOrder: []*Order{
			{ID: 2},
		},
		FilledAmount: map[string]*big.Int{
			"BTC":  big.NewInt(100),
			"USDT": big.NewInt(200),
		},
	}

	{
		db, err := New[*Order]("./test_curd_object.db")
		tester.RequireNoError(t, err)
		tester.RequireNotNil(t, db)

		{
			ok, err := db.Exists(ctx, "order")
			tester.RequireNoError(t, err)
			tester.RequireFalse(t, ok)

			val, err := db.Get(ctx, "order")
			tester.RequireErrorIs(t, err, ErrNotFound)
			tester.RequireNil(t, val)
		}

		{
			tester.RequireNoError(t, db.Set(ctx, "order", order))
		}

		{
			val, err := db.Get(ctx, "order")
			tester.RequireNoError(t, err)
			tester.RequireEqual(t, 1, val.ID)
			tester.RequireEqual(t, 1, len(val.RelativeOrder))
			tester.RequireEqual(t, 2, val.RelativeOrder[0].ID)
			tester.RequireEqual(t, 2, len(val.FilledAmount))
			tester.RequireEqual(t, 0, big.NewInt(100).Cmp(val.FilledAmount["BTC"]))
			tester.RequireEqual(t, 0, big.NewInt(200).Cmp(val.FilledAmount["USDT"]))
		}

		{
			tester.RequireNoError(t, db.Delete(ctx, "order"))

			ok, err := db.Exists(ctx, "order")
			tester.RequireNoError(t, err)
			tester.RequireFalse(t, ok)

			val, err := db.Get(ctx, "order")
			tester.RequireErrorIs(t, err, ErrNotFound)
			tester.RequireNil(t, val)
		}
	}
}

func TestLocal_Find(t *testing.T) {
	defer func() {
		tester.RequireNoError(t, Delete("./test_find_int.db"))
	}()

	ctx := context.Background()

	db, err := New[int]("./test_find_int.db")
	tester.RequireNoError(t, err)
	tester.RequireNotNil(t, db)

	{
		tester.RequireNoError(t, db.Clear(ctx))
		tester.RequireNoError(t, db.Set(ctx, "hello", 1))
		tester.RequireNoError(t, db.Set(ctx, "world", 2))

		vals, err := db.Find(ctx)
		tester.RequireNoError(t, err)
		tester.RequireEqual(t, 2, len(vals))
		tester.RequireEqual(t, 1, vals[0])
		tester.RequireEqual(t, 2, vals[1])

		vals, err = db.Find(ctx, "hello", "world")
		tester.RequireNoError(t, err)
		tester.RequireEqual(t, 2, len(vals))
		tester.RequireEqual(t, 1, vals[0])
		tester.RequireEqual(t, 2, vals[1])

		vals, err = db.Find(ctx, "hello")
		tester.RequireNoError(t, err)
		tester.RequireEqual(t, 1, len(vals))
		tester.RequireEqual(t, 1, vals[0])

		vals, err = db.Find(ctx, "not_exist")
		tester.RequireNoError(t, err)
		tester.RequireEqual(t, 0, len(vals))
	}
}

func TestLocal_Atomic(t *testing.T) {
	defer func() {
		tester.RequireNoError(t, Delete("./test_atomic_int.db"))
	}()

	ctx := context.Background()

	db, err := New[int]("./test_atomic_int.db")
	tester.RequireNoError(t, err)
	tester.RequireNotNil(t, db)

	{
		tester.RequireNoError(t, db.Clear(ctx))

		err := db.Atomic(ctx, func(tx Local[int]) error {
			if err := tx.Set(ctx, "hello", 1); err != nil {
				return err
			}

			if err := tx.Set(ctx, "world", 2); err != nil {
				return err
			}

			return nil
		})
		tester.RequireNoError(t, err)

		ok, err := db.Exists(ctx, "hello")
		tester.RequireNoError(t, err)
		tester.RequireTrue(t, ok)

		ok, err = db.Exists(ctx, "world")
		tester.RequireNoError(t, err)
		tester.RequireTrue(t, ok)
	}

	{
		tester.RequireNoError(t, db.Clear(ctx))

		err := db.Atomic(ctx, func(tx Local[int]) error {
			if err := tx.Set(ctx, "hello", 1); err != nil {
				return err
			}

			if err := tx.Set(ctx, "world", 2); err != nil {
				return err
			}

			return tx.Close()
		})
		tester.RequireNoError(t, err)

		db, err = New[int]("./test_atomic_int.db")
		tester.RequireNoError(t, err)
		tester.RequireNotNil(t, db)

		ok, err := db.Exists(ctx, "hello")
		tester.RequireNoError(t, err)
		tester.RequireTrue(t, ok)

		ok, err = db.Exists(ctx, "world")
		tester.RequireNoError(t, err)
		tester.RequireTrue(t, ok)
	}

	{
		tester.RequireNoError(t, db.Clear(ctx))

		err := db.Atomic(ctx, func(tx Local[int]) error {
			if err := tx.Set(ctx, "not_hello", 1); err != nil {
				return err
			}

			if err := tx.Set(ctx, "not_world", 2); err != nil {
				return err
			}

			return errors.New("rollback")
		})
		tester.RequireError(t, err)

		ok, err := db.Exists(ctx, "not_hello")
		tester.RequireNoError(t, err)
		tester.RequireFalse(t, ok)

		ok, err = db.Exists(ctx, "not_world")
		tester.RequireNoError(t, err)
		tester.RequireFalse(t, ok)
	}
}
