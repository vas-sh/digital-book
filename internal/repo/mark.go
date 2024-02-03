package repo

import (
	"context"
	"digital-book/internal/types"
)

func (r *repo) CreateMark(ctx context.Context, mark *types.Mark) error {
	return r.db.WithContext(ctx).Create(mark).Error
}

func (r *repo) UpdateMark(ctx context.Context, userID, subjectID, value, id string) error {
	return r.db.WithContext(ctx).Model(&types.Mark{}).Where("id = ?", id).Updates(map[string]any{"user_id": userID, "subject_id": subjectID, "value": value}).Error
}

func (r *repo) GetMarks(ctx context.Context) (res []types.MarkResponse, err error) {
	err = r.db.WithContext(ctx).Raw(`
       SELECT mark.id, "user".name as user_name, subject.title as subject_title, mark.value 
        FROM mark
            INNER JOIN "user" 
            ON mark.user_id = "user".id 
            INNER JOIN subject 
            ON mark.subject_id = subject.id
            ORDER BY id ASC`).Scan(&res).Error
	return
}

func (r *repo) DeleteMark(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&types.Mark{}).Error
}

func (r *repo) GetMark(ctx context.Context, id string) (res types.Mark, err error) {
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}

func (r *repo) AvgMarks(ctx context.Context) (res []types.MarkAverege, err error) {
	err = r.db.WithContext(ctx).Raw(`
    SELECT u.id, u.name, subject.title, AVG(value) AS Value
    FROM mark 
        INNER JOIN public.user as u ON mark.user_id = u.id
        INNER JOIN subject ON mark.subject_id = subject.id
        GROUP BY 1,2,3 
        ORDER BY u.id ASC`).Scan(&res).Error
	return
}
