package repo

import (
	"context"
	"digital-book/internal/types"
)

func (r *repo) CreateMark(ctx context.Context, studentID, subjectID, value string) error {
	return r.db.WithContext(ctx).Exec("INSERT INTO mark (ID, student_id, subject_id, value) VALUES(DEFAULT, ?, ?, ?)",
		studentID, subjectID, value).Error
}

func (r *repo) UpdateMark(ctx context.Context, studentID, subjectID, value, id string) error {
	return r.db.WithContext(ctx).Exec("UPDATE mark SET student_id = ?, subject_id = ?, value = ? WHERE id = ?",
		studentID, subjectID, value, id).Error
}

func (r *repo) GetMarks(ctx context.Context) (res []types.MarkResponse, err error) {
	err = r.db.WithContext(ctx).Raw(`
       SELECT mark.id, student.name as student_name, subject.title as subject_title, mark.value 
        FROM mark
            INNER JOIN student 
            ON mark.student_id = student.id 
            INNER JOIN subject 
            ON mark.subject_id = subject.id
            ORDER BY id ASC`).Scan(&res).Error
	return
}

func (r *repo) DeleteMark(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM mark WHERE id = ?", id).Error
}

func (r *repo) GetMark(ctx context.Context, id string) (res types.Mark, err error) {
	err = r.db.WithContext(ctx).Raw("SELECT * FROM mark WHERE id = ?", id).Scan(&res).Error
	return
}
