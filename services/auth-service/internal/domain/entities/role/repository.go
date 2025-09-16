package role

import (
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/infrastructure/repository"
	"auth-service/internal/interfaces"
)

type RoleRepository struct {
	*repository.Repository[Role]
}

func NewRoleRepository(db *database.MongoDB) interfaces.IRepository[Role] {
	return &RoleRepository{
		Repository: repository.NewRepository[Role](db.Database.Collection("roles")),
	}
}
