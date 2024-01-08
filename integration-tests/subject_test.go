package integrationtests_test

import (
	"context"
	"digital-book/internal/types"
	"fmt"
	"testing"
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
		if subject.Title != created.Title {
			t.Errorf("invalid title: want %s, got %s", created.Title, subject.Title)
			return
		}
	})
}
