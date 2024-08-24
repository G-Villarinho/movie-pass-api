package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/GSVillas/movie-pass-api/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSessionService_Create_WhenSuccessful_ShouldReturnToken(t *testing.T) {
	var err error
	config.Env.PrivateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRepositoryMock := mock.NewMockSessionRepository(ctrl)
	sessionService := &sessionService{
		sessionRepository: sessionRepositoryMock,
	}

	user := &domain.User{
		ID:        uuid.New(),
		FirstName: "test",
		LastName:  "teste",
		Email:     "test@example.com",
	}

	sessionRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	token, err := sessionService.Create(context.Background(), *user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
