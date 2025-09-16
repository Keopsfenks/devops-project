package policy

import (
	"time"

	"github.com/google/uuid"
)

type PolicyEffect string

const (
	Allow PolicyEffect = "allow"
	Deny  PolicyEffect = "deny"
)

type PolicyCondition struct {
	Field    string      `json:"field" bson:"field"`
	Operator string      `json:"operator" bson:"operator"`
	Value    interface{} `json:"value" bson:"value"`
}

type Policy struct {
	ID          uuid.UUID         `json:"id" bson:"id"`
	Name        string            `json:"name" bson:"name"`
	Description string            `json:"description" bson:"description"`
	Effect      PolicyEffect      `json:"effect" bson:"effect"`
	Resource    string            `json:"resource" bson:"resource"`
	Action      string            `json:"action" bson:"action"`
	Conditions  []PolicyCondition `json:"conditions,omitempty" bson:"conditions,omitempty"`
	Priority    int               `json:"priority" bson:"priority"`

	IsDeleted bool `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func NewPolicy(name, description string, effect PolicyEffect, resource, action string, priority int) *Policy {
	return &Policy{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Effect:      effect,
		Resource:    resource,
		Action:      action,
		Priority:    priority,
		IsDeleted:   false,
		CreatedAt:   time.Now(),
		Conditions:  make([]PolicyCondition, 0),
	}
}

func (p *Policy) AddCondition(field, operator string, value interface{}) {
	condition := PolicyCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
	p.Conditions = append(p.Conditions, condition)
	now := time.Now()
	p.UpdatedAt = &now
}
