package service

import (
	"context"
	"digital-book/internal/types"
)

func (s *srv) CreateStudent(ctx context.Context, student *types.Student) error {
	return s.repo.CreateStudent(ctx, student)
}

func (s *srv) GetStudents(ctx context.Context) (res []types.Student, err error) {
	return s.repo.GetStudents(ctx)
}

func (s *srv) GetStudent(ctx context.Context, id string) (res types.Student, err error) {
	return s.repo.GetStudent(ctx, id)
}

func (s *srv) DeleteStudent(ctx context.Context, id string) error {
	return s.repo.DeleteStudent(ctx, id)
}

func (s *srv) UpdateStudent(ctx context.Context, name, class, id string) error {
	return s.repo.UpdateStudent(ctx, name, class, id)
}
