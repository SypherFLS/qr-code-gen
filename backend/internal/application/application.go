package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"qrgen/internal/qrgen"
	"qrgen/internal/repository"
	"qrgen/internal/validation"
)

type Application struct {
	repo repository.Repository
	ttl  time.Duration
}

func NewApplication(repo repository.Repository, ttl time.Duration) *Application {
	return &Application{
		repo: repo,
		ttl:  ttl,
	}
}

func (s *Application) CreateQRApplication(
	ctx context.Context,
	rawURL string,
) (string, []byte, error) {
	if err := validation.ValidateHTTPURL(rawURL); err != nil {
		return "", nil, err
	}

	id := uuid.NewString()

	png, err := qrgen.Gen(rawURL)
	if err != nil {
		return "", nil, fmt.Errorf("generate qr: %w", err)
	}

	if err := s.repo.Save(ctx, id, png, s.ttl); err != nil {
		return "", nil, fmt.Errorf("save qr: %w", err)
	}

	return id, png, nil
}

func (s *Application) GetQRApplication(
	ctx context.Context,
	id string,
) ([]byte, error) {
	return s.repo.Get(ctx, id)
}
