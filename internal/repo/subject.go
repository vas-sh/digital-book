package repo

import (
	"context"
	"digital-book/internal/types"
)

func (r *repo) GetSubjects(ctx context.Context) (res []types.Subject, err error) {
	err = r.db.WithContext(ctx).Raw("SELECT * FROM subject ORDER BY id").Scan(&res).Error
	return
}
