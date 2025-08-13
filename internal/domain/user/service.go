package user

import (
    "context"
    "database/sql"
    stdErrors "errors"

    "github.com/google/uuid"
    "github.com/jackc/pgconn"

    appErr "github.com/shbrx98/Decamond_Tsak/internal/pkg/errors"
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
    if err != nil && !isNotFoundErr(err) {
        return nil, err
    }

    nu := NewUser(phone)
    if err := s.repo.Create(ctx, nu); err != nil {
        if isUniqueViolation(err) {
            return s.repo.GetByPhone(ctx, phone)
        }
        return nil, err
    }
    return nu, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, params ListParams) (*ListResult, error) {
    if params.Limit <= 0 || params.Limit > 100 {
        params.Limit = 20
    }
    if params.Offset < 0 {
        params.Offset = 0
    }
    return s.repo.List(ctx, params)
}

// Helpers

func isNotFoundErr(err error) bool {
    if err == nil {
        return false
    }
    if stdErrors.Is(err, sql.ErrNoRows) {
        return true
    }
    var nf *appErr.NotFoundError
    if stdErrors.As(err, &nf) {
        return true
    }
    return false
}

func isUniqueViolation(err error) bool {
    var pgErr *pgconn.PgError
    if stdErrors.As(err, &pgErr) {
        return pgErr.Code == "23505" // unique_violation
    }
    return false
}