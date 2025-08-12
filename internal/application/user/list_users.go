package user

import (
    "context"

    dom "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"
)

type ListUsersUseCase struct {
    svc    *dom.Service
    logger logger.Logger
}
func NewListUsersUseCase(s *dom.Service, l logger.Logger) *ListUsersUseCase {
    return &ListUsersUseCase{svc: s, logger: l}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, query string, page, limit int) (*dom.ListResult, error) {
    if limit <= 0 || limit > 100 { limit = 20 }
    if page <= 0 { page = 1 }
    params := dom.ListParams{
        Query:  query,
        Limit:  limit,
        Offset: (page-1) * limit,
    }
    return uc.svc.List(ctx, params)
}