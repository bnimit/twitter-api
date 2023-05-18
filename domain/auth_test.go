package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/nimit-bhandari/twitter"
	"github.com/nimit-bhandari/twitter/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	validInput := twitter.RegisterInput{
		Username:        "bob",
		Email:           "bob@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	t.Run("Can register", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		// Add interface method mocks
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("Create", mock.Anything, mock.Anything).
			Return(twitter.User{
				ID:       "123",
				Username: validInput.Username,
				Email:    validInput.Email,
			}, nil)

		service := NewAuthService(userRepo)

		res, err := service.Register(ctx, validInput)
		require.NoError(t, err)

		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.Username)
		require.NotEmpty(t, res.User.Email)
		require.NotEmpty(t, res.User.Username)

		userRepo.AssertExpectations(t)
	})

	t.Run("Username taken", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
		Return(twitter.User{}, nil)

		service := NewAuthService(userRepo)

		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, twitter.ErrUserNameTaken)

		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)
	})

	t.Run("Email taken", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
		Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
		Return(twitter.User{}, nil)

		service := NewAuthService(userRepo)

		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, twitter.ErrEmailTaken)

		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)
	})

	t.Run("create error", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		// Add interface method mocks
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("Create", mock.Anything, mock.Anything).
			Return(twitter.User{}, errors.New("something"))

		service := NewAuthService(userRepo)

		_, err := service.Register(ctx, validInput)
		require.Error(t, err)

		userRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		service := NewAuthService(userRepo)
 
		_, err := service.Register(ctx, twitter.RegisterInput{})
		require.ErrorIs(t, err, twitter.ErrValidation)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)
	})
}
