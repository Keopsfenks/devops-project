package session

import (
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/infrastructure/repository"
	"auth-service/internal/interfaces"
)

type SessionRepository struct {
	*repository.Repository[Session]
}

func NewSessionRepository(db *database.MongoDB) interfaces.IRepository[Session] {
	return &SessionRepository{
		Repository: repository.NewRepository[Session](db.Database.Collection("sessions")),
	}
}
