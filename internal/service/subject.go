package service

import (
	"context"
	"digital-book/internal/types"
)

func (s *srv) CreateSubject(ctx context.Context, subject *types.Subject) error {
	return s.repo.CreateSubject(ctx, subject)
}

func (s *srv) GetSubject(ctx context.Context, id string) (res types.Subject, err error) {
	return s.repo.GetSubject(ctx, id)
}
