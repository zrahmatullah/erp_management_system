package auth

import (
	"context"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/domain"
	"cafe-erp-system/backend/pkg/crypto"
)

type UserUsecaseImpl struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &UserUsecaseImpl{userRepo: userRepo}
}

func (u *UserUsecaseImpl) Create(ctx context.Context, req domain.RegisterRequest) (domain.User, error) {
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

func (u *UserUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}

func (u *UserUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req domain.RegisterRequest) (domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	user.Username = req.Username
	user.Email = req.Email
	user.FullName = req.FullName
	user.Phone = req.Phone

	err = u.userRepo.Update(ctx, user)
	return *user, err
}

func (u *UserUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *UserUsecaseImpl) List(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	users, err := u.userRepo.List(ctx, limit, offset)
	return users, len(users), err
}

func (u *UserUsecaseImpl) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return nil
}

func (u *UserUsecaseImpl) GetUsersByBranch(ctx context.Context, branchID uuid.UUID) ([]domain.User, error) {
	return nil, nil
}
