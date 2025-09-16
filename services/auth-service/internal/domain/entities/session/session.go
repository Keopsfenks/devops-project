package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id" bson:"id"`
	UserID    uuid.UUID `json:"user_id" bson:"user_id"`
	Token     string    `json:"token" bson:"token"`
	IPAddress string    `json:"ip_address" bson:"ip_address"`
	UserAgent string    `json:"user_agent" bson:"user_agent"`

	IsDeleted bool      `json:"is_active" bson:"is_active"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`

	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func NewSession(userID uuid.UUID, token, ipAddress, userAgent string, expiresAt time.Time) *Session {
	return &Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		IsDeleted: false,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Invalidate() {
	s.IsDeleted = false
	now := time.Now()
	s.UpdatedAt = &now
}

func (s *Session) ExtendExpiry(duration time.Duration) {
	s.ExpiresAt = time.Now().Add(duration)
	now := time.Now()
	s.UpdatedAt = &now
}
