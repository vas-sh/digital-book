package integrationtests_test

import (
	"context"
	"digital-book/internal/types"
	"fmt"
	"testing"
)

func TestCreateMark(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.Mark{
			StudentID: 1,
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
		if mark.StudentID != created.StudentID || mark.SubjectID != created.SubjectID || mark.Value != created.Value {
			t.Errorf("invalid name, subject, value: want %d, %d,%d, got %d, %d, %d ",
				created.StudentID, created.SubjectID, created.Value, mark.StudentID, mark.SubjectID, mark.Value)
			return
		}
	})
}

func TestDeleteMark(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.MarkResponse{
			StudentName:  "Vasul",
			SubjectTitle: "math",
			Value:        10,
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
