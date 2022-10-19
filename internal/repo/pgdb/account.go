package pgdb

import (
	"account-management-service/internal/entity"
	"account-management-service/internal/repo/repoerrs"
	"account-management-service/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type AccountRepo struct {
	*postgres.Postgres
}

func NewAccountRepo(pg *postgres.Postgres) *AccountRepo {
	return &AccountRepo{pg}
}

func (r *AccountRepo) CreateAccount(ctx context.Context) (int, error) {
	sql, args, err := r.Builder.
		Insert("account").
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("AccountRepo.CreateAccount - r.Builder: %v", err)
	}

	var id int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrs.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("AccountRepo.CreateAccount - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *AccountRepo) GetAccountById(ctx context.Context, id int) (entity.Account, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("account").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return entity.Account{}, fmt.Errorf("AccountRepo.GetAccountById - r.Builder: %v", err)
	}

	var account entity.Account
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&account.Id,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Account{}, repoerrs.ErrNotFound
		}
		return entity.Account{}, fmt.Errorf("AccountRepo.GetAccountById - r.Pool.QueryRow: %v", err)
	}

	return account, nil
}

func (r *AccountRepo) Deposit(ctx context.Context, id, amount int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - r.Pool.Begin: %v", err)
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - tx.Exec: %v", err)
	}

	sql, args, err = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type").
		Values(id, amount, entity.OperationTypeDeposit).
		ToSql()

	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - tx.Commit: %v", err)
	}

	return nil
}

func (r *AccountRepo) Withdraw(ctx context.Context, id, amount int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - r.Pool.Begin: %v", err)
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance - ?", amount)).
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - tx.Exec: %v", err)
	}

	sql, args, err = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type").
		Values(id, amount, entity.OperationTypeWithdraw).
		ToSql()

	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - tx.Commit: %v", err)
	}

	return nil
}

func (r *AccountRepo) Transfer(ctx context.Context, from, to, amount int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - r.Pool.Begin: %v", err)
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance - ?", amount)).
		Where("id = ?", from).
		ToSql()

	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Exec: %v", err)
	}

	sql, args, err = r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where("id = ?", to).
		ToSql()

	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Exec: %v", err)
	}

	sql, args, err = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type").
		Values(from, amount, entity.OperationTypeTransferFrom).
		ToSql()

	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Exec: %v", err)
	}

	sql, args, err = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type").
		Values(to, amount, entity.OperationTypeTransferTo).
		ToSql()

	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Commit: %v", err)
	}

	return nil
}
