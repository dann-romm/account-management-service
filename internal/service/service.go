package service

import (
	"account-management-service/internal/entity"
	"account-management-service/internal/repo"
	"account-management-service/pkg/hasher"
	"context"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go -package=mocks

type AuthCreateUserInput struct {
	Username string
	Password string
}

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (int, error)
}

type AccountDepositInput struct {
	Id     int
	Amount int
}

type AccountWithdrawInput struct {
	Id     int
	Amount int
}

type AccountTransferInput struct {
	From   int
	To     int
	Amount int
}

type Account interface {
	CreateAccount(ctx context.Context) (int, error)
	Deposit(ctx context.Context, input AccountDepositInput) error
	Withdraw(ctx context.Context, input AccountWithdrawInput) error
	Transfer(ctx context.Context, input AccountTransferInput) error
}

type Product interface {
	CreateProduct(ctx context.Context, name string) (int, error)
	GetProductById(ctx context.Context, id int) (entity.Product, error)
}

type ReservationCreateInput struct {
	AccountId int
	ProductId int
	OrderId   int
	Amount    int
}

type Reservation interface {
	CreateReservation(ctx context.Context, input ReservationCreateInput) (int, error)
	RefundReservationByOrderId(ctx context.Context, orderId int) error
	RevenueReservationByOrderId(ctx context.Context, orderId int) error
}

type Services struct {
	Auth        Auth
	Account     Account
	Product     Product
	Reservation Reservation
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth:        NewAuthService(deps.Repos, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Account:     NewAccountService(deps.Repos),
		Product:     NewProductService(deps.Repos),
		Reservation: NewReservationService(deps.Repos),
	}
}
