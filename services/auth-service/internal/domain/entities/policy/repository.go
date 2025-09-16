package policy

import (
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/infrastructure/repository"
	"auth-service/internal/interfaces"
)

type PolicyRepository struct {
	*repository.Repository[Policy]
}

func NewPolicyRepository(db *database.MongoDB) interfaces.IRepository[Policy] {
	return &PolicyRepository{
		Repository: repository.NewRepository[Policy](db.Database.Collection("policies")),
	}
}
