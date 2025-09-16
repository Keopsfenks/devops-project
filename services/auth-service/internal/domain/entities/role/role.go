package role

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id" bson:"id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`

	Policies []uuid.UUID `json:"policies" bson:"policies"`

	IsDeleted bool `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func NewRole(name, description string, policies []uuid.UUID) *Role {
	return &Role{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Policies:    policies,
		IsDeleted:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		DeletedAt:   nil,
	}
}

func (r *Role) AddPolicy(policy uuid.UUID) {
	for _, p := range r.Policies {
		if p == policy {
			return
		}
	}
	r.Policies = append(r.Policies, policy)
	now := time.Now()
	r.UpdatedAt = &now
}

func (r *Role) RemovePolicy(policy uuid.UUID) {
	for i, p := range r.Policies {
		if p == policy {
			r.Policies = append(r.Policies[:i], r.Policies[i+1:]...)
			now := time.Now()
			r.UpdatedAt = &now
			return
		}
	}
}
