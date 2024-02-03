package service

import (
	"context"
	"digital-book/internal/types"
)

func (s *srv) CreateUser(ctx context.Context, user *types.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *srv) GetUsers(ctx context.Context) (res []types.User, err error) {
	return s.repo.GetUsers(ctx)
}

func (s *srv) GetUser(ctx context.Context, id string) (res types.User, err error) {
	return s.repo.GetUser(ctx, id)
}

func (s *srv) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *srv) UpdateUser(ctx context.Context, name, class, id, login string) error {
	return s.repo.UpdateUser(ctx, name, class, id, login)
}
