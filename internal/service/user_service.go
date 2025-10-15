// internal/service/user_service.go
package service

import (
	"context"
	"errors"

	"github.com/knyazushka/verbose-engine/internal/domain"
	"github.com/knyazushka/verbose-engine/internal/repository"
)

var (
	ErrUserAlreadyExists   = errors.New("user with this email already exists")
	ErrInvaliedCredentials = errors.New("invalid email or password")
)

type UserService struct {
	userRepo   *repository.UserRepository
	jwtService *JWTService
}

func NewUserService(userRepo *repository.UserRepository, jwtService *JWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *UserService) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrUserAlreadyExists
	}

	user := &domain.User{
		Email:    req.Email,
		Username: req.Username,
		IsActive: true,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *UserService) Login(ctx context.Context, req *domain.LoginRequest) (string, *domain.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, ErrInvaliedCredentials
	}

	if !user.CheckPassword(req.Password) {
		return "", nil, ErrInvaliedCredentials
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return "", nil, err
	}

	return token, &domain.UserResponse{}, nil
}

func (s *UserService) Profile(ctx context.Context, userId string) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}
