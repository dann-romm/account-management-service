package pgdb

import (
	"account-management-service/internal/entity"
	"account-management-service/pkg/postgres"
	"context"
	"fmt"
)

type ProductRepo struct {
	*postgres.Postgres
}

func NewProductRepo(pg *postgres.Postgres) *ProductRepo {
	return &ProductRepo{pg}
}

func (p ProductRepo) CreateProduct(ctx context.Context, name string) (int, error) {
	sql, args, _ := p.Builder.
		Insert("products").
		Columns("name").
		Values(name).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := p.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ProductRepo.CreateProduct - p.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (p ProductRepo) GetProductById(ctx context.Context, id int) (entity.Product, error) {
	sql, args, _ := p.Builder.
		Select("*").
		From("products").
		Where("id = ?", id).
		ToSql()

	var product entity.Product
	err := p.Pool.QueryRow(ctx, sql, args...).Scan(
		&product.Id,
		&product.Name,
	)
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductRepo.GetProductById - p.Pool.QueryRow: %v", err)
	}

	return product, nil
}
