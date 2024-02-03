package repo

import (
	"context"
	"digital-book/internal/types"
)

func (r *repo) GetUsers(ctx context.Context) (res []types.User, err error) {
	err = r.db.WithContext(ctx).Find(&res).Error
	return
}

func (r *repo) CreateUser(ctx context.Context, user *types.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repo) UpdateUser(ctx context.Context, name, class, id, login string) error {
	return r.db.WithContext(ctx).Model(&types.User{}).Where("id = ?", id).Updates(types.User{Name: name, Class: class, Login: login}).Error
}

func (r *repo) GetUser(ctx context.Context, id string) (res types.User, err error) {
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}

func (r *repo) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&types.User{}).Error
}
