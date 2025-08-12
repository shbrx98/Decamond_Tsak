package postgres

import (
    "context"
    "strings"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    dom "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
)

type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) dom.Repository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *dom.User) error {
    q := `INSERT INTO users (id, phone, created_at, updated_at) VALUES ($1, $2, $3, $4)`
    _, err := r.db.ExecContext(ctx, q, u.ID, u.Phone, u.CreatedAt, u.UpdatedAt)
    return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*dom.User, error) {
    var u dom.User
    q := `SELECT id, phone, created_at, updated_at FROM users WHERE id=$1`
    if err := r.db.GetContext(ctx, &u, q, id); err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*dom.User, error) {
    var u dom.User
    q := `SELECT id, phone, created_at, updated_at FROM users WHERE phone=$1`
    if err := r.db.GetContext(ctx, &u, q, phone); err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *userRepository) Update(ctx context.Context, u *dom.User) error {
    q := `UPDATE users SET phone=$1, updated_at=$2 WHERE id=$3`
    _, err := r.db.ExecContext(ctx, q, u.Phone, u.UpdatedAt, u.ID)
    return err
}

func (r *userRepository) Exists(ctx context.Context, phone string) (bool, error) {
    var exists bool
    q := `SELECT EXISTS(SELECT 1 FROM users WHERE phone=$1)`
    if err := r.db.GetContext(ctx, &exists, q, phone); err != nil {
        return false, err
    }
    return exists, nil
}

func (r *userRepository) List(ctx context.Context, params dom.ListParams) (*dom.ListResult, error) {
    like := "%" + strings.TrimSpace(params.Query) + "%"
    var total int
    if params.Query == "" {
        _ = r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM users`)
    } else {
        _ = r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM users WHERE phone ILIKE $1`, like)
    }

    var users []*dom.User
    if params.Query == "" {
        q := `SELECT id, phone, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
        if err := r.db.SelectContext(ctx, &users, q, params.Limit, params.Offset); err != nil {
            return nil, err
        }
    } else {
        q := `SELECT id, phone, created_at, updated_at FROM users WHERE phone ILIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
        if err := r.db.SelectContext(ctx, &users, q, like, params.Limit, params.Offset); err != nil {
            return nil, err
        }
    }

    page := params.Offset/params.Limit + 1
    totalPages := (total + params.Limit - 1) / params.Limit
    return &dom.ListResult{
        Users:      users,
        Total:      total,
        Page:       page,
        PerPage:    params.Limit,
        TotalPages: totalPages,
    }, nil
}