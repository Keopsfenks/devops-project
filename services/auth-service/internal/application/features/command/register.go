package command

import (
	"context"
	"net/http"

	"auth-service/internal/application/dto"
	"auth-service/internal/domain/entities/user"
	"auth-service/internal/interfaces"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	FirstName string `json:"first_name" validate:"required,alpha,max=50"`
	LastName  string `json:"last_name" validate:"required,alpha,max=50"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=8,password_complex"`
}

type RegisterHandler struct {
	userRepository interfaces.IRepository[user.User]
}

func NewRegisterHandler(userRepo interfaces.IRepository[user.User]) *RegisterHandler {
	return &RegisterHandler{
		userRepository: userRepo,
	}
}

var validate *validator.Validate

func passwordComplex(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, c := range password {
		switch {
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case 'a' <= c && c <= 'z':
			hasLower = true
		case '0' <= c && c <= '9':
			hasNumber = true
		case (c >= 33 && c <= 47) || (c >= 58 && c <= 64) || (c >= 91 && c <= 96) || (c >= 123 && c <= 126):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("password_complex", passwordComplex)
}

func (r *Request) Handle(ctx context.Context, handler *RegisterHandler) *dto.Result[*user.User] {
	// Validate request
	if err := validate.Struct(r); err != nil {
		return dto.NewResult[*user.User](nil, false, err.Error(), http.StatusBadRequest)
	}

	// Check if user already exists
	exists, err := handler.userRepository.Exists(
		ctx, map[string]interface{}{
			"email": r.Email,
		},
	)
	if err != nil {
		return dto.NewResult[*user.User](
			nil, false, "Failed to check existing user", http.StatusInternalServerError,
		)
	}
	if exists {
		return dto.NewResult[*user.User](
			nil, false, "User with this email already exists", http.StatusConflict,
		)
	}
	// Create new user
	newUser, err := user.NewUser(r.FirstName, r.LastName, r.Email, r.Password)
	if err != nil {
		return dto.NewResult[*user.User](nil, false, "Failed to create user", http.StatusInternalServerError)
	}

	// Save user to database
	createdUser, err := handler.userRepository.InsertOne(ctx, newUser)
	if err != nil {
		return dto.NewResult[*user.User](nil, false, "Failed to save user", http.StatusInternalServerError)
	}

	return dto.NewResult[*user.User](&createdUser, true, "", http.StatusCreated)
}
