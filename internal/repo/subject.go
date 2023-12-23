package repo

import (
	"context"
	"digital-book/internal/types"
)

func (r *repo) GetSubjects(ctx context.Context) (res []types.Subject, err error) {
	err = r.db.WithContext(ctx).Raw("SELECT * FROM subject ORDER BY id").Scan(&res).Error
	return
}

func (r *repo) CreateSubject(ctx context.Context, title string) error {
	return r.db.WithContext(ctx).Exec("INSERT INTO subject (ID, title) VALUES(DEFAULT, ?)", title).Error
}

func (r *repo) UpdateSubject(ctx context.Context, title, id string) error {
	return r.db.WithContext(ctx).Exec("UPDATE subject SET title = ? WHERE id = ?", title, id).Error
}

func (r *repo) DeleteSubject(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM subject WHERE id = ?", id).Error
}
