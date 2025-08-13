package user

import (
    "context"

    "github.com/google/uuid"
)


type Repository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetByPhone(ctx context.Context, phone string) (*User, error)
    List(ctx context.Context, params ListParams) (*ListResult, error)
    Update(ctx context.Context, user *User) error
    Exists(ctx context.Context, phone string) (bool, error)
}


type ListParams struct {
    Query  string
    Limit  int
    Offset int
}


type ListResult struct {
    Users      []*User `json:"users"`
    Total      int     `json:"total"`
    Page       int     `json:"page"`
    PerPage    int     `json:"per_page"`
    TotalPages int     `json:"total_pages"`
}