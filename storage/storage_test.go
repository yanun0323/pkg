package storage

import (
	"context"
	"math/big"
	"testing"

	"github.com/pkg/errors"
	"github.com/yanun0323/pkg/test"
)

func TestNew(t *testing.T) {
	defer func() {
		test.RequireNoError(t, Delete("./test_int.db"))
		test.RequireNoError(t, Delete("./test_string.db"))
	}()

	{
		db, err := New[int]("./test_int.db")
		test.RequireNoError(t, err)
		test.RequireNotNil(t, db)
	}

	{
		db, err := New[string]("./test_string.db")
		test.RequireNoError(t, err)
		test.RequireNotNil(t, db)
		test.RequireNoError(t, db.Close())
	}

	{
		db, err := New[bool]("./test_int.db")
		test.RequireErrorIs(t, ErrTypeMismatch, err)
		test.RequireNil(t, db)
	}
}

func TestLocal_CURD(t *testing.T) {
	defer func() {
		test.RequireNoError(t, Delete("./test_curd_int.db"))
	}()

	ctx := context.Background()

	{
		db, err := New[int]("./test_curd_int.db")
		test.RequireNoError(t, err)
		test.RequireNotNil(t, db)

		{
			ok, err := db.Exists(ctx, "hello")
			test.RequireNoError(t, err)
			test.RequireFalse(t, ok)

			val, err := db.Get(ctx, "hello")
			test.RequireErrorIs(t, ErrNotFound, err)
			test.RequireEqual(t, 0, val)
		}

		{
			test.RequireNoError(t, db.Set(ctx, "hello", 1))
			test.RequireNoError(t, db.Set(ctx, "world", 2))
		}

		{
			val, err := db.Get(ctx, "hello")
			test.RequireNoError(t, err)
			test.RequireEqual(t, 1, val)

			val, err = db.Get(ctx, "world")
			test.RequireNoError(t, err)
			test.RequireEqual(t, 2, val)
		}

		{
			test.RequireNoError(t, db.Delete(ctx, "hello"))

			ok, err := db.Exists(ctx, "hello")
			test.RequireNoError(t, err)
			test.RequireFalse(t, ok)

			val, err := db.Get(ctx, "hello")
			test.RequireErrorIs(t, err, ErrNotFound)
			test.RequireEqual(t, 0, val)
		}
	}

	defer func() {
		test.RequireNoError(t, Delete("./test_curd_object.db"))
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
		test.RequireNoError(t, err)
		test.RequireNotNil(t, db)

		{
			ok, err := db.Exists(ctx, "order")
			test.RequireNoError(t, err)
			test.RequireFalse(t, ok)

			val, err := db.Get(ctx, "order")
			test.RequireErrorIs(t, err, ErrNotFound)
			test.RequireNil(t, val)
		}

		{
			test.RequireNoError(t, db.Set(ctx, "order", order))
		}

		{
			val, err := db.Get(ctx, "order")
			test.RequireNoError(t, err)
			test.RequireEqual(t, 1, val.ID)
			test.RequireEqual(t, 1, len(val.RelativeOrder))
			test.RequireEqual(t, 2, val.RelativeOrder[0].ID)
			test.RequireEqual(t, 2, len(val.FilledAmount))
			test.RequireEqual(t, 0, big.NewInt(100).Cmp(val.FilledAmount["BTC"]))
			test.RequireEqual(t, 0, big.NewInt(200).Cmp(val.FilledAmount["USDT"]))
		}

		{
			test.RequireNoError(t, db.Delete(ctx, "order"))

			ok, err := db.Exists(ctx, "order")
			test.RequireNoError(t, err)
			test.RequireFalse(t, ok)

			val, err := db.Get(ctx, "order")
			test.RequireErrorIs(t, err, ErrNotFound)
			test.RequireNil(t, val)
		}
	}
}

func TestLocal_Find(t *testing.T) {
	defer func() {
		test.RequireNoError(t, Delete("./test_find_int.db"))
	}()

	ctx := context.Background()

	db, err := New[int]("./test_find_int.db")
	test.RequireNoError(t, err)
	test.RequireNotNil(t, db)

	{
		test.RequireNoError(t, db.Clear(ctx))
		test.RequireNoError(t, db.Set(ctx, "hello", 1))
		test.RequireNoError(t, db.Set(ctx, "world", 2))

		vals, err := db.Find(ctx)
		test.RequireNoError(t, err)
		test.RequireEqual(t, 2, len(vals))
		test.RequireEqual(t, 1, vals[0])
		test.RequireEqual(t, 2, vals[1])

		vals, err = db.Find(ctx, "hello", "world")
		test.RequireNoError(t, err)
		test.RequireEqual(t, 2, len(vals))
		test.RequireEqual(t, 1, vals[0])
		test.RequireEqual(t, 2, vals[1])

		vals, err = db.Find(ctx, "hello")
		test.RequireNoError(t, err)
		test.RequireEqual(t, 1, len(vals))
		test.RequireEqual(t, 1, vals[0])

		vals, err = db.Find(ctx, "not_exist")
		test.RequireNoError(t, err)
		test.RequireEqual(t, 0, len(vals))
	}
}

func TestLocal_Atomic(t *testing.T) {
	defer func() {
		test.RequireNoError(t, Delete("./test_atomic_int.db"))
	}()

	ctx := context.Background()

	db, err := New[int]("./test_atomic_int.db")
	test.RequireNoError(t, err)
	test.RequireNotNil(t, db)

	{
		test.RequireNoError(t, db.Clear(ctx))

		err := db.Atomic(ctx, func(tx Local[int]) error {
			if err := tx.Set(ctx, "hello", 1); err != nil {
				return err
			}

			if err := tx.Set(ctx, "world", 2); err != nil {
				return err
			}

			return nil
		})
		test.RequireNoError(t, err)

		ok, err := db.Exists(ctx, "hello")
		test.RequireNoError(t, err)
		test.RequireTrue(t, ok)

		ok, err = db.Exists(ctx, "world")
		test.RequireNoError(t, err)
		test.RequireTrue(t, ok)
	}

	{
		test.RequireNoError(t, db.Clear(ctx))

		err := db.Atomic(ctx, func(tx Local[int]) error {
			if err := tx.Set(ctx, "hello", 1); err != nil {
				return err
			}

			if err := tx.Set(ctx, "world", 2); err != nil {
				return err
			}

			return tx.Close()
		})
		test.RequireNoError(t, err)

		db, err = New[int]("./test_atomic_int.db")
		test.RequireNoError(t, err)
		test.RequireNotNil(t, db)

		ok, err := db.Exists(ctx, "hello")
		test.RequireNoError(t, err)
		test.RequireTrue(t, ok)

		ok, err = db.Exists(ctx, "world")
		test.RequireNoError(t, err)
		test.RequireTrue(t, ok)
	}

	{
		test.RequireNoError(t, db.Clear(ctx))

		err := db.Atomic(ctx, func(tx Local[int]) error {
			if err := tx.Set(ctx, "not_hello", 1); err != nil {
				return err
			}

			if err := tx.Set(ctx, "not_world", 2); err != nil {
				return err
			}

			return errors.New("rollback")
		})
		test.RequireError(t, err)

		ok, err := db.Exists(ctx, "not_hello")
		test.RequireNoError(t, err)
		test.RequireFalse(t, ok)

		ok, err = db.Exists(ctx, "not_world")
		test.RequireNoError(t, err)
		test.RequireFalse(t, ok)
	}
}
