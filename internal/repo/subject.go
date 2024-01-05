package repo

import (
	"context"
	"digital-book/internal/types"
)

func (r *repo) GetSubjects(ctx context.Context) (res []types.Subject, err error) { //done
	err = r.db.WithContext(ctx).Find(&res).Error
	return
}

func (r *repo) CreateSubject(ctx context.Context, subject *types.Subject) error { //done
	return r.db.WithContext(ctx).Create(subject).Error
}

func (r *repo) UpdateSubject(ctx context.Context, title, id string) error { //done
	return r.db.WithContext(ctx).Model(&types.Subject{}).Where("id = ?", id).Update("title", title).Error
}

func (r *repo) DeleteSubject(ctx context.Context, id string) error { //done
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&types.Subject{}).Error
}

func (r *repo) GetSubject(ctx context.Context, id string) (res types.Subject, err error) { // done
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}
