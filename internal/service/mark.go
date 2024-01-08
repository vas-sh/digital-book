package service

import (
	"context"
	"digital-book/internal/types"
)

func (s *srv) GetMarks(ctx context.Context) (res []types.MarkResponse, err error) {
	return s.repo.GetMarks(ctx)
}

func (s *srv) DeleteMark(ctx context.Context, id string) error {
	return s.repo.DeleteMark(ctx, id)
}
