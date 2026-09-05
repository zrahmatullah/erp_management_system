package auth

import (
	"context"
	"errors"

	"cafe-erp-system/backend/internal/delivery/http/middleware"
	"cafe-erp-system/backend/internal/domain"
	"cafe-erp-system/backend/pkg/crypto"
)

type AuthUsecaseImpl struct {
	userRepo domain.UserRepository
}

func NewAuthUsecase(userRepo domain.UserRepository) domain.AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo: userRepo,
	}
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		user, err = u.userRepo.GetByUsername(ctx, req.Email)
		if err != nil {
			return domain.LoginResponse{}, errors.New("invalid credentials")
		}
	}

	if !crypto.CheckPassword(req.Password, user.PasswordHash) {
		return domain.LoginResponse{}, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return domain.LoginResponse{}, errors.New("account is disabled")
	}

	accessToken, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Email, "Super Admin", user.BranchID, []string{})
	if err != nil {
		return domain.LoginResponse{}, err
	}

	return domain.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (u *AuthUsecaseImpl) Register(ctx context.Context, req domain.RegisterRequest) (domain.User, error) {
	_, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return domain.User{}, errors.New("email already in use")
	}

	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        req.Phone,
		IsActive:     true,
	}

	err = u.userRepo.Create(ctx, &user)
	return user, err
}

func (u *AuthUsecaseImpl) RefreshToken(ctx context.Context, token string) (domain.LoginResponse, error) {
	claims, err := middleware.ValidateToken(token)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	accessToken, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Email, "Super Admin", user.BranchID, []string{})
	if err != nil {
		return domain.LoginResponse{}, err
	}

	return domain.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (u *AuthUsecaseImpl) ForgotPassword(ctx context.Context, email string) error {
	return nil
}

func (u *AuthUsecaseImpl) ResetPassword(ctx context.Context, token, newPassword string) error {
	return nil
}
