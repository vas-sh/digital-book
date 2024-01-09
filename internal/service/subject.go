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

func (s *srv) GetSubjects(ctx context.Context) (res []types.Subject, err error) {
	return s.repo.GetSubjects(ctx)
}

func (s *srv) DeleteSubject(ctx context.Context, id string) error {
	return s.repo.DeleteSubject(ctx, id)
}

func (s *srv) UpdateSubject(ctx context.Context, title, id string) error {
	return s.repo.UpdateSubject(ctx, title, id)
}
