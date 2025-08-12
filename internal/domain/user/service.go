package user

import (
    "context"
    "github.com/google/uuid"
    appErr "github.com/shbrx98/Decamond_Tsak/internal/pkg/errors"
    "time"
)

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) GetOrCreateByPhone(ctx context.Context, phone string) (*User, error) {
    u, err := s.repo.GetByPhone(ctx, phone)
    if err == nil && u != nil {
        return u, nil
    }
    nu := NewUser(phone)
    if err := s.repo.Create(ctx, nu); err != nil {
        return nil, err
    }
    return nu, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, params ListParams) (*ListResult, error) {
    // guardrails for pagination
    if params.Limit <= 0 || params.Limit > 100 {
        params.Limit = 20
    }
    if params.Offset < 0 {
        params.Offset = 0
    }
    return s.repo.List(ctx, params)
}

func (u *User) Touch() { u.UpdatedAt = time.Now() }

var ErrInvalidPhone = appErr.NewValidationError("invalid phone")