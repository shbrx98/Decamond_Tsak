package user

import (
    "context"

    "github.com/google/uuid"

    dom "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"
)

type GetUserUseCase struct {
    svc    *dom.Service
    logger logger.Logger
}

func NewGetUserUseCase(s *dom.Service, l logger.Logger) *GetUserUseCase {
    return &GetUserUseCase{svc: s, logger: l}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id uuid.UUID) (*dom.User, error) {
    return uc.svc.GetByID(ctx, id)
}