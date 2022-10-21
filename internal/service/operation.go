package service

import (
	"account-management-service/internal/repo"
	"context"
)

type OperationService struct {
	operationRepo repo.Operation
	productRepo   repo.Product
}

func NewOperationService(operationRepo repo.Operation, productRepo repo.Product) *OperationService {
	return &OperationService{
		operationRepo: operationRepo,
		productRepo:   productRepo,
	}
}

func (s *OperationService) OperationHistory(ctx context.Context, input OperationHistoryInput) ([]OperationHistoryOutput, error) {
	operations, productNames, err := s.operationRepo.OperationsPagination(ctx, input.AccountId, input.SortType, input.Offset, input.Limit)
	if err != nil {
		return nil, err
	}

	output := make([]OperationHistoryOutput, 0, len(operations))
	for i, operation := range operations {

		output = append(output, OperationHistoryOutput{
			Amount:      operation.Amount,
			Operation:   operation.OperationType,
			Time:        operation.CreatedAt,
			Product:     productNames[i],
			Order:       operation.OrderId,
			Description: operation.Description,
		})
	}
	return output, nil
}
