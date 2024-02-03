package integrationtests_test

import (
	"context"
	"digital-book/internal/types"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateMark(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.Mark{
			UserID:    1,
			SubjectID: 2,
			Value:     10,
		}
		if err := s.CreateMark(ctx, &created); err != nil {
			t.Errorf("mark is not created: " + err.Error())
			return
		}
		mark, err := s.GetMark(ctx, fmt.Sprint(created.ID))
		if err != nil {
			t.Errorf("mark is not created: " + err.Error())
			return
		}
		if mark.UserID != created.UserID || mark.SubjectID != created.SubjectID || mark.Value != created.Value {
			t.Errorf("invalid name, subject, value: want %d, %d,%d, got %d, %d, %d ",
				created.UserID, created.SubjectID, created.Value, mark.UserID, mark.SubjectID, mark.Value)
			return
		}
	})
}

func TestDeleteMark(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.Mark{
			UserID: 1, SubjectID: 2,
			Value: 10}
		if err := s.CreateMark(ctx, &created); err != nil {
			t.Errorf("mark is not created: " + err.Error())
			return
		}
		if err := s.DeleteMark(ctx, fmt.Sprint(created.ID)); err != nil {
			t.Errorf("delete failed: %v", err)
			return
		}
		mark, err := s.GetMark(ctx, fmt.Sprint(created.ID))
		if err == nil {
			t.Errorf("mark should be deleted, but it still exists: %+v", mark)
			return
		}
	})
}
func TestUpdateMark(t *testing.T) {
	ctx := context.Background()
	runInTransaction(func(s srv) {
		created := types.Mark{UserID: 1,
			SubjectID: 2, Value: 10,
		}
		if err := s.CreateMark(ctx, &created); err != nil {
			t.Errorf("mark is not created: " + err.Error())
			return
		}
		updatedUserID := 2
		updatedSubjectID := 3
		updatedValue := 11

		if err := s.UpdateMark(ctx, fmt.Sprint(updatedUserID), fmt.Sprint(updatedSubjectID), fmt.Sprint(updatedValue), fmt.Sprint(created.ID)); err != nil {
			t.Errorf("error updating user: " + err.Error())
			return
		}
		mark, err := s.GetMark(ctx, fmt.Sprint(created.ID))
		if err != nil {
			t.Errorf("error getting mark: " + err.Error())
			return
		}
		if mark.UserID != updatedUserID || mark.SubjectID != updatedSubjectID || mark.Value != updatedValue {
			t.Errorf("invalid values after update: %s", cmp.Diff(mark, created))
			return
		}
	})
}
