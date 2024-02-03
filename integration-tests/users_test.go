package integrationtests_test

import (
	"context"
	"digital-book/internal/types"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.User{
			Name:  "Vasul",
			Class: "10",
			Login: "VasulTest",
		}
		if err := s.CreateUser(ctx, &created); err != nil {
			t.Errorf("user is not created: " + err.Error())
			return
		}
		user, err := s.GetUser(ctx, fmt.Sprint(created.ID))
		if err != nil {
			t.Errorf("user is not created: " + err.Error())
			return
		}
		if user.Name != created.Name || user.Class != created.Class {
			t.Errorf("invalid name, class: want %s, %s, got %s, %s", created.Name, created.Class, user.Name, user.Class)
			return
		}
	})
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.User{
			Name:  "Vasul",
			Class: "10",
			Login: "VasulTest",
		}
		if err := s.CreateUser(ctx, &created); err != nil {
			t.Errorf("user is not created: " + err.Error())
			return
		}

		created.Name = "John"
		created.Class = "11"

		if err := s.UpdateUser(ctx, created.Name, created.Class, fmt.Sprint(created.ID), created.Login); err != nil {
			t.Errorf("error updating user: " + err.Error())
			return
		}

		user, err := s.GetUser(ctx, fmt.Sprint(created.ID))
		if err != nil {
			t.Errorf("error getting user: " + err.Error())
			return
		}

		if user != created {
			t.Errorf("invalid update: %s", cmp.Diff(user, created))
			return
		}
	})
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	runInTransaction(func(s srv) {
		created := types.User{
			Name:  "Vasul",
			Class: "10",
			Login: "VasulTest",
		}
		if err := s.CreateUser(ctx, &created); err != nil {
			t.Errorf("user is not created: " + err.Error())
			return
		}
		if err := s.DeleteUser(ctx, fmt.Sprint(created.ID)); err != nil {
			t.Errorf("delete failed: %v", err)
			return
		}
		user, err := s.GetUser(ctx, fmt.Sprint(created.ID))
		if err == nil {
			t.Errorf("user should be deleted, but it still exists: %+v", user)
			return
		}
	})
}
