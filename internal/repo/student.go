package repo

import (
	"context"
	"digital-book/internal/types"
)

func (r *repo) GetStudents(ctx context.Context) (res []types.Student, err error) { //done
	err = r.db.WithContext(ctx).Find(&res).Error
	return
}

func (r *repo) CreateStudent(ctx context.Context, student *types.Student) error { //done
	return r.db.WithContext(ctx).Create(student).Error
}

func (r *repo) UpdateStudent(ctx context.Context, name, class, id string) error { //done
	return r.db.WithContext(ctx).Model(&types.Student{}).Where("id = ?", id).Updates(types.Student{Name: name, Class: class}).Error
}

func (r *repo) GetStudent(ctx context.Context, id string) (res types.Student, err error) { // done
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}

func (r *repo) DeleteStudent(ctx context.Context, id string) error { //done
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&types.Student{}).Error
}
