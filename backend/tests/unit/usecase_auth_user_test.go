package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"cafe-erp-system/backend/internal/delivery/http/middleware"
	"cafe-erp-system/backend/internal/domain"
	"cafe-erp-system/backend/internal/usecase/auth"
	"cafe-erp-system/backend/pkg/crypto"
)

// MockUserRepository implements domain.UserRepository in-memory for testing
type MockUserRepository struct {
	users       map[uuid.UUID]*domain.User
	roles       map[uuid.UUID]string
	permissions map[uuid.UUID][]domain.Permission
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:       make(map[uuid.UUID]*domain.User),
		roles:       make(map[uuid.UUID]string),
		permissions: make(map[uuid.UUID][]domain.Permission),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found by email")
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found by username")
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	var list []domain.User
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func (m *MockUserRepository) GetUserRoleAndPermissions(ctx context.Context, userID uuid.UUID) (string, *uuid.UUID, []domain.Permission, error) {
	roleName, ok := m.roles[userID]
	if !ok {
		roleName = "Branch Manager"
	}
	roleID := uuid.New()
	perms := m.permissions[userID]
	return roleName, &roleID, perms, nil
}

func TestAuthUsecase_Login(t *testing.T) {
	middleware.SetJWTSecretForTesting("test-secret-key-for-unit-testing-32chars!!")
	repo := NewMockUserRepository()
	usecase := auth.NewAuthUsecase(repo)

	rawPassword := "SecurePassword123!"
	hashedPwd, _ := crypto.HashPassword(rawPassword)

	activeUser := &domain.User{
		BaseEntity:   domain.BaseEntity{ID: uuid.New()},
		Username:     "john_barista",
		Email:        "john@cafe.com",
		PasswordHash: hashedPwd,
		FullName:     "John Doe",
		IsActive:     true,
	}
	_ = repo.Create(context.Background(), activeUser)

	t.Run("Happy Path - Login with Email", func(t *testing.T) {
		res, err := usecase.Login(context.Background(), domain.LoginRequest{
			Email:    "john@cafe.com",
			Password: rawPassword,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Token)
		assert.NotEmpty(t, res.RefreshToken)
		assert.Equal(t, activeUser.ID, res.User.ID)
		assert.Equal(t, "Branch Manager", res.Role)
	})

	t.Run("Happy Path - Login with Username", func(t *testing.T) {
		res, err := usecase.Login(context.Background(), domain.LoginRequest{
			Email:    "john_barista",
			Password: rawPassword,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Token)
		assert.Equal(t, activeUser.ID, res.User.ID)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		_, err := usecase.Login(context.Background(), domain.LoginRequest{
			Email:    "john@cafe.com",
			Password: "WrongPassword!",
		})

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("User Not Found", func(t *testing.T) {
		_, err := usecase.Login(context.Background(), domain.LoginRequest{
			Email:    "nonexistent@cafe.com",
			Password: rawPassword,
		})

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("Inactive Account Rejected", func(t *testing.T) {
		inactiveUser := &domain.User{
			BaseEntity:   domain.BaseEntity{ID: uuid.New()},
			Username:     "inactive_user",
			Email:        "inactive@cafe.com",
			PasswordHash: hashedPwd,
			IsActive:     false,
		}
		_ = repo.Create(context.Background(), inactiveUser)

		_, err := usecase.Login(context.Background(), domain.LoginRequest{
			Email:    "inactive@cafe.com",
			Password: rawPassword,
		})

		assert.Error(t, err)
		assert.Equal(t, "account is disabled", err.Error())
	})
}

func TestAuthUsecase_Register(t *testing.T) {
	repo := NewMockUserRepository()
	usecase := auth.NewAuthUsecase(repo)

	t.Run("Successful Registration", func(t *testing.T) {
		req := domain.RegisterRequest{
			Username: "newbie",
			Email:    "newbie@cafe.com",
			Password: "Password123!",
			FullName: "Newbie Staff",
			Phone:    "08123456789",
		}

		user, err := usecase.Register(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "newbie", user.Username)
		assert.Equal(t, "newbie@cafe.com", user.Email)
		assert.True(t, user.IsActive)
		assert.True(t, crypto.CheckPassword("Password123!", user.PasswordHash))
	})

	t.Run("Duplicate Email Registration Fails", func(t *testing.T) {
		req := domain.RegisterRequest{
			Username: "newbie2",
			Email:    "newbie@cafe.com",
			Password: "Password123!",
		}

		_, err := usecase.Register(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "email already in use", err.Error())
	})
}

func TestAuthUsecase_RefreshToken(t *testing.T) {
	middleware.SetJWTSecretForTesting("test-secret-key-for-unit-testing-32chars!!")
	repo := NewMockUserRepository()
	usecase := auth.NewAuthUsecase(repo)

	user := &domain.User{
		BaseEntity: domain.BaseEntity{ID: uuid.New()},
		Username: "refresh_user",
		Email:    "refresh@cafe.com",
		IsActive: true,
	}
	_ = repo.Create(context.Background(), user)

	_, validRefresh, _ := middleware.GenerateTokenPair(user.ID, user.Email, "Cashier", nil, nil)

	t.Run("Valid Refresh Token", func(t *testing.T) {
		res, err := usecase.RefreshToken(context.Background(), validRefresh)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Token)
		assert.NotEmpty(t, res.RefreshToken)
		assert.Equal(t, user.ID, res.User.ID)
	})

	t.Run("Invalid Refresh Token String", func(t *testing.T) {
		_, err := usecase.RefreshToken(context.Background(), "invalid.jwt.refresh")
		assert.Error(t, err)
	})

	t.Run("User Not Found for Valid Token", func(t *testing.T) {
		orphanID := uuid.New()
		_, orphanRefresh, _ := middleware.GenerateTokenPair(orphanID, "orphan@cafe.com", "Cashier", nil, nil)

		_, err := usecase.RefreshToken(context.Background(), orphanRefresh)
		assert.Error(t, err)
	})
}

func TestAuthUsecase_ForgotPassword_And_Reset(t *testing.T) {
	middleware.SetJWTSecretForTesting("test-secret-key-for-unit-testing-32chars!!")
	repo := NewMockUserRepository()
	usecase := auth.NewAuthUsecase(repo)

	user := &domain.User{
		BaseEntity:   domain.BaseEntity{ID: uuid.New()},
		Email:        "forgot@cafe.com",
		PasswordHash: "oldhash",
		IsActive:     true,
	}
	_ = repo.Create(context.Background(), user)

	t.Run("ForgotPassword Existing User", func(t *testing.T) {
		err := usecase.ForgotPassword(context.Background(), "forgot@cafe.com")
		assert.NoError(t, err)
	})

	t.Run("ForgotPassword Nonexistent User", func(t *testing.T) {
		err := usecase.ForgotPassword(context.Background(), "ghost@cafe.com")
		assert.Error(t, err)
		assert.Equal(t, "user with this email does not exist", err.Error())
	})

	t.Run("ResetPassword with Valid Token", func(t *testing.T) {
		resetToken, _, _ := middleware.GenerateTokenPair(user.ID, user.Email, "Cashier", nil, nil)
		err := usecase.ResetPassword(context.Background(), resetToken, "BrandNewPassword123!")
		assert.NoError(t, err)

		updated, _ := repo.GetByID(context.Background(), user.ID)
		assert.True(t, crypto.CheckPassword("BrandNewPassword123!", updated.PasswordHash))
	})

	t.Run("ResetPassword with Invalid Token", func(t *testing.T) {
		err := usecase.ResetPassword(context.Background(), "bogus-token", "NewPass")
		assert.Error(t, err)
	})
}

func TestUserUsecase_CRUD(t *testing.T) {
	repo := NewMockUserRepository()
	usecase := auth.NewUserUsecase(repo)

	t.Run("Create User", func(t *testing.T) {
		user, err := usecase.Create(context.Background(), domain.RegisterRequest{
			Username: "created_user",
			Email:    "create@cafe.com",
			Password: "Pass123Password!",
			FullName: "Created Fullname",
		})
		assert.NoError(t, err)
		assert.Equal(t, "created_user", user.Username)

		// GetByID
		fetched, err := usecase.GetByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, fetched.ID)

		// Update
		updated, err := usecase.Update(context.Background(), user.ID, domain.RegisterRequest{
			Username: "updated_username",
			Email:    "updated@cafe.com",
			FullName: "Updated Fullname",
			Phone:    "089999999",
		})
		assert.NoError(t, err)
		assert.Equal(t, "updated_username", updated.Username)

		// List
		list, total, err := usecase.List(context.Background(), 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, list, 1)

		// AssignRole & GetUsersByBranch stubs
		assert.Nil(t, usecase.AssignRole(context.Background(), user.ID, uuid.New()))
		branchUsers, err := usecase.GetUsersByBranch(context.Background(), uuid.New())
		assert.NoError(t, err)
		assert.Nil(t, branchUsers)

		// Delete
		err = usecase.Delete(context.Background(), user.ID)
		assert.NoError(t, err)

		// Verify deleted
		_, err = usecase.GetByID(context.Background(), user.ID)
		assert.Error(t, err)
	})
}

func TestAuthUsecase_GetProfile(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authUsecase := auth.NewAuthUsecase(mockRepo)

	t.Run("Existing User returns profile with role and permissions", func(t *testing.T) {
		userID := uuid.New()
		mockRepo.users[userID] = &domain.User{
			BaseEntity: domain.BaseEntity{ID: userID},
			Username:   "owner_test",
			Email:      "owner@cafe.com",
			FullName:   "Owner Test",
			IsActive:   true,
		}
		mockRepo.roles[userID] = "Owner"
		mockRepo.permissions[userID] = []domain.Permission{
			{Module: "dashboard", Action: "view"},
		}

		res, err := authUsecase.GetProfile(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, "owner_test", res.User.Username)
		assert.Equal(t, "Owner", res.Role)
		assert.Len(t, res.Permissions, 1)
	})

	t.Run("Non-existent User returns error", func(t *testing.T) {
		_, err := authUsecase.GetProfile(context.Background(), uuid.New())
		assert.Error(t, err)
	})
}


