package integrationtests_test

import (
	"context"
	"digital-book/internal/types"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateStudent(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.Student{
			Name:  "Vasul",
			Class: "10",
		}
		if err := s.CreateStudent(ctx, &created); err != nil {
			t.Errorf("student is not created: " + err.Error())
			return
		}
		student, err := s.GetStudent(ctx, fmt.Sprint(created.ID))
		if err != nil {
			t.Errorf("student is not created: " + err.Error())
			return
		}
		if student.Name != created.Name || student.Class != created.Class {
			t.Errorf("invalid name, class: want %s, %s, got %s, %s", created.Name, created.Class, student.Name, student.Class)
			return
		}
	})
}

func TestUpdateStudent(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.Student{
			Name:  "Vasul",
			Class: "10",
		}
		if err := s.CreateStudent(ctx, &created); err != nil {
			t.Errorf("student is not created: " + err.Error())
			return
		}

		created.Name = "John"
		created.Class = "11"

		if err := s.UpdateStudent(ctx, created.Name, created.Class, fmt.Sprint(created.ID)); err != nil {
			t.Errorf("error updating student: " + err.Error())
			return
		}

		student, err := s.GetStudent(ctx, fmt.Sprint(created.ID))
		if err != nil {
			t.Errorf("error getting student: " + err.Error())
			return
		}

		if student != created {
			t.Errorf("invalid update: %s", cmp.Diff(student, created))
			return
		}
	})
}

func TestDeleteStudent(t *testing.T) {
	ctx := context.Background()

	runInTransaction(func(s srv) {
		created := types.Student{
			Name:  "Vasul",
			Class: "10",
		}

		if err := s.DeleteStudent(ctx, fmt.Sprint(created.ID)); err != nil {
			t.Errorf("delete failed: %v", err)
			return
		}

		student, err := s.GetStudent(ctx, fmt.Sprint(created.ID))
		if err == nil {
			t.Errorf("student should be deleted, but it still exists: %+v", student)
			return
		}
	})
}
