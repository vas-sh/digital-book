package integrationtests_test

import (
	"context"
	"digital-book/internal/types"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateSubject(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.Subject{
			Title: "Math",
		}
		if err := s.CreateSubject(ctx, &created); err != nil {
			t.Errorf("subject is not created: " + err.Error())
			return
		}
		subject, err := s.GetSubject(ctx, fmt.Sprint(created.ID))
		if err != nil {
			t.Errorf("subject is not created: " + err.Error())
			return
		}
		if subject != created {
			t.Errorf("invalid title: want %s, got %s", created.Title, subject.Title)
			return
		}
	})
}

func TestUpdateSubject(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.Subject{
			Title: "math",
		}
		if err := s.CreateSubject(ctx, &created); err != nil {
			t.Errorf("subject is not created: " + err.Error())
			return
		}

		updatedTitle := "biology"

		if err := s.UpdateSubject(ctx, updatedTitle, fmt.Sprint(created.ID)); err != nil {
			t.Errorf("error updating user: " + err.Error())
			return
		}

		subject, err := s.GetSubject(ctx, fmt.Sprint(created.ID))
		if err != nil {
			t.Errorf("error getting user: " + err.Error())
			return
		}

		if subject.Title != updatedTitle {
			t.Errorf("invalid title after update: %s", cmp.Diff(subject, created))
			return
		}
	})
}

func TestDeleteSubject(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.Subject{
			Title: "Math"}

		if err := s.CreateSubject(ctx, &created); err != nil {
			t.Errorf("subject is not created: " + err.Error())
			return
		}
		if err := s.DeleteSubject(ctx, fmt.Sprint(created.ID)); err != nil {
			t.Errorf("delete failed: %v", err)
			return
		}
		subject, err := s.GetSubject(ctx, fmt.Sprint(created.ID))
		if err == nil {
			t.Errorf("subject should be deleted, but it still exists: %+v", subject)
			return
		}
	})
}
