package user

import (
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/infrastructure/repository"
	"auth-service/internal/interfaces"
)

type UserRepository struct {
	*repository.Repository[User]
}

func NewUserRepository(db *database.MongoDB) interfaces.IRepository[User] {
	return &UserRepository{
		Repository: repository.NewRepository[User](db.Database.Collection("users")),
	}
}
