package pgdb

import (
	"account-management-service/internal/entity"
	"account-management-service/pkg/postgres"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

type ReservationRepo struct {
	*postgres.Postgres
}

func NewReservationRepo(pg *postgres.Postgres) *ReservationRepo {
	return &ReservationRepo{pg}
}

func (r *ReservationRepo) CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - r.Pool.Begin: %v", err)
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Insert("reservation").
		Columns("account_id", "product_id", "order_id", "amount").
		Values(
			reservation.AccountId,
			reservation.ProductId,
			reservation.OrderId,
			reservation.Amount,
		).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - r.Builder: %v", err)
	}

	var id int
	err = tx.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - tx.QueryRow: %v", err)
	}

	sql, args, err = r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance - ?", reservation.Amount)).
		Where("id = ?", reservation.AccountId).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - tx.Exec: %v", err)
	}

	sql, args, err = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type", "product_id", "order_id").
		Values(
			reservation.AccountId,
			reservation.Amount,
			entity.OperationTypeReservation,
			reservation.ProductId,
			reservation.OrderId,
		).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - tx.Commit: %v", err)
	}

	return id, nil
}

func (r *ReservationRepo) GetReservationById(ctx context.Context, id int) (entity.Reservation, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("reservation").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return entity.Reservation{}, fmt.Errorf("ReservationRepo.GetReservationById - r.Builder: %v", err)
	}

	var reservation entity.Reservation
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&reservation.Id,
		&reservation.AccountId,
		&reservation.ProductId,
		&reservation.OrderId,
		&reservation.Amount,
		&reservation.CreatedAt,
	)
	if err != nil {
		return entity.Reservation{}, fmt.Errorf("ReservationRepo.GetReservationById - r.Pool.QueryRow: %v", err)
	}

	return reservation, nil
}

func (r *ReservationRepo) RefundReservationById(ctx context.Context, id int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.DeleteReservationById - r.Pool.Begin: %v", err)
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Delete("reservation").
		Where("id = ?", id).
		Suffix("RETURNING account_id, product_id, order_id, amount").
		ToSql()

	if err != nil {
		return fmt.Errorf("ReservationRepo.DeleteReservationById - r.Builder: %v", err)
	}

	var reservation entity.Reservation
	err = tx.QueryRow(ctx, sql, args...).Scan(
		&reservation.AccountId,
		&reservation.ProductId,
		&reservation.OrderId,
		&reservation.Amount,
	)
	if err != nil {
		return fmt.Errorf("ReservationRepo.DeleteReservationById - tx.QueryRow: %v", err)
	}

	sql, args, err = r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance + ?", reservation.Amount)).
		Where("id = ?", reservation.AccountId).
		ToSql()

	if err != nil {
		return fmt.Errorf("ReservationRepo.DeleteReservationById - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReservationRepo.DeleteReservationById - tx.Exec: %v", err)
	}

	sql, args, err = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type", "product_id", "order_id").
		Values(
			reservation.AccountId,
			reservation.Amount,
			entity.OperationTypeRefund,
			reservation.ProductId,
			reservation.OrderId,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("ReservationRepo.DeleteReservationById - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReservationRepo.DeleteReservationById - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.DeleteReservationById - tx.Commit: %v", err)
	}

	return nil
}

func (r *ReservationRepo) RevenueReservationById(ctx context.Context, id int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - r.Pool.Begin: %v", err)
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Delete("reservation").
		Where("id = ?", id).
		Suffix("RETURNING account_id, product_id, order_id, amount").
		ToSql()

	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - r.Builder: %v", err)
	}

	var reservation entity.Reservation
	err = tx.QueryRow(ctx, sql, args...).Scan(
		&reservation.AccountId,
		&reservation.ProductId,
		&reservation.OrderId,
		&reservation.Amount,
	)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - tx.QueryRow: %v", err)
	}

	sql, args, err = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type", "product_id", "order_id").
		Values(
			reservation.AccountId,
			reservation.Amount,
			entity.OperationTypeRevenue,
			reservation.ProductId,
			reservation.OrderId,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - r.Builder: %v", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - tx.Commit: %v", err)
	}

	return nil
}
