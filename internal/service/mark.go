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

func (s *srv) CreateMark(ctx context.Context, mark *types.Mark) error {
	return s.repo.CreateMark(ctx, mark)
}

func (s *srv) GetMark(ctx context.Context, id string) (res types.Mark, err error) {
	return s.repo.GetMark(ctx, id)
}

func (s *srv) UpdateMark(ctx context.Context, userID, subjectID, value, id string) error {
	return s.repo.UpdateMark(ctx, userID, subjectID, value, id)
}

func (s *srv) AvgMarks(ctx context.Context) (res []types.MarkAverege, err error) {
	return s.repo.AvgMarks(ctx)
}
