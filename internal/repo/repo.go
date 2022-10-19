package repo

import (
	"account-management-service/internal/entity"
	"account-management-service/internal/repo/pgdb"
	"account-management-service/pkg/postgres"
	"context"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
	GetUserById(ctx context.Context, id int) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}

type Account interface {
	CreateAccount(ctx context.Context) (int, error)
	GetAccountById(ctx context.Context, id int) (entity.Account, error)
	Deposit(ctx context.Context, id, amount int) error
	Withdraw(ctx context.Context, id, amount int) error
	Transfer(ctx context.Context, from, to, amount int) error
}

type Product interface {
	CreateProduct(ctx context.Context, name string) (int, error)
	GetProductById(ctx context.Context, id int) (entity.Product, error)
}

type Reservation interface {
	CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error)
	GetReservationById(ctx context.Context, id int) (entity.Reservation, error)
	DeleteReservationById(ctx context.Context, id int) error
	DeleteReservationByOrderId(ctx context.Context, orderId int) error
}

type Repositories struct {
	User
	Account
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:    pgdb.NewUserRepo(pg),
		Account: pgdb.NewAccountRepo(pg),
	}
}
