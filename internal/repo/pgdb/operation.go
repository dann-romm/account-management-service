package pgdb

import (
	"account-management-service/internal/entity"
	"account-management-service/pkg/postgres"
	"context"
	"fmt"
)

const (
	maxPaginationLimit     = 10
	defaultPaginationLimit = 10

	DateSortType   string = "date"
	AmountSortType string = "amount"
)

type OperationRepo struct {
	*postgres.Postgres
}

func NewOperationRepo(pg *postgres.Postgres) *OperationRepo {
	return &OperationRepo{pg}
}

func (r *OperationRepo) GetAllRevenueOperationsGroupedByProductId(ctx context.Context) ([]entity.Operation, error) {
	sql, args, _ := r.Builder.
		Select("product_id", "sum(amount)").
		From("operations").
		Where("operation_type = ?", entity.OperationTypeRevenue).
		GroupBy("product_id").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("OperationRepo.GetAllRevenueOperationsGroupedByProductId - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var operations []entity.Operation
	for rows.Next() {
		var operation entity.Operation
		err = rows.Scan(&operation.ProductId, &operation.Amount)
		if err != nil {
			return nil, fmt.Errorf("OperationRepo.GetAllRevenueOperationsGroupedByProductId - rows.Scan: %v", err)
		}
		operations = append(operations, operation)
	}

	return operations, nil
}

func (r *OperationRepo) OperationsPagination(ctx context.Context, accountId int, sortType string, offset int, limit int) ([]entity.Operation, []string, error) {
	if limit > maxPaginationLimit {
		limit = maxPaginationLimit
	}
	if limit == 0 {
		limit = defaultPaginationLimit
	}

	var orderBySql string
	switch sortType {
	case "":
		orderBySql = "created_at DESC"
	case DateSortType:
		orderBySql = "created_at DESC"
	case AmountSortType:
		orderBySql = "amount DESC"
	default:
		return nil, nil, fmt.Errorf("OperationRepo.PaginationOperations: unknown sort type - %s", sortType)
	}

	sqlQuery, args, _ := r.Builder.
		Select("operations.id", "account_id", "amount", "operation_type", "created_at", "COALESCE((case when operations.product_id is null then null else products.name end), '') as product_name", "order_id", "COALESCE(description, '')").
		From("operations").
		InnerJoin("products on operations.product_id = products.id or operations.product_id is null").
		Where("account_id = ?", accountId).
		OrderBy(orderBySql).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	rows, err := r.Pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("OperationRepo.paginationOperationsByDate - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var operations []entity.Operation
	var productNames []string
	for rows.Next() {
		var operation entity.Operation
		var productName string
		err = rows.Scan(&operation.Id, &operation.AccountId, &operation.Amount, &operation.OperationType, &operation.CreatedAt, &productName, &operation.OrderId, &operation.Description)
		if err != nil {
			return nil, nil, fmt.Errorf("OperationRepo.paginationOperationsByDate - rows.Scan: %v", err)
		}
		operations = append(operations, operation)
		productNames = append(productNames, productName)
	}

	return operations, productNames, nil
}
