package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"id" bson:"id"`
	FirstName    string    `json:"first_name" bson:"first_name"`
	LastName     string    `json:"last_name" bson:"last_name"`
	Email        string    `json:"email" bson:"email"`
	PasswordHash string    `json:"password_hash" bson:"password_hash"`

	IsDeleted bool `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func NewUser(firstName, lastName, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           uuid.New(),
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
		DeletedAt:    nil,
		IsDeleted:    false,
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) UpdateFirstName(firstName string) {
	u.FirstName = firstName
	now := time.Now()
	u.UpdatedAt = &now
}

func (u *User) UpdateLastName(lastName string) {
	u.LastName = lastName
	now := time.Now()
	u.UpdatedAt = &now
}

func (u *User) UpdateEmail(email string) {
	u.Email = email
	now := time.Now()
	u.UpdatedAt = &now
}

func (u *User) UpdatePassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	now := time.Now()
	u.UpdatedAt = &now
	return nil
}

func (u *User) SoftDelete() {
	u.IsDeleted = true
	now := time.Now()
	u.DeletedAt = &now
}
