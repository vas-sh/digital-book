package repo

import (
	"context"
	"digital-book/internal/types"
)

func (r *repo) GetStudents(ctx context.Context) (res []types.Student, err error) {
	err = r.db.WithContext(ctx).Raw("SELECT * FROM student ORDER BY id ASC").Scan(&res).Error
	return
}

func (r *repo) CreateStudent(ctx context.Context, name, class string) error {
	return r.db.WithContext(ctx).Exec("INSERT INTO student (name, class) VALUES (?, ?)", name, class).Error
}

func (r *repo) UpdateStudent(ctx context.Context, name, class, id string) error {
	return r.db.WithContext(ctx).Exec("UPDATE student SET name = ?, class = ? WHERE id = ?", name, class, id).Error
}

func (r *repo) GetStudent(ctx context.Context, id string) (res types.Student, err error) {
	err = r.db.WithContext(ctx).Raw("SELECT * FROM student WHERE id = ?", id).Scan(&res).Error
	return
}

func (r *repo) DeleteStudent(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM student WHERE id = ?", id).Error
}
