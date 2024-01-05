package repo

import (
	"context"
	"digital-book/internal/types"
)

func (r *repo) CreateMark(ctx context.Context, mark *types.Mark) error { // при можливості виправити
	return r.db.WithContext(ctx).Create(mark).Error
}

func (r *repo) UpdateMark(ctx context.Context, studentID, subjectID, value, id string) error { //done
	return r.db.WithContext(ctx).Model(&types.Mark{}).Where("id = ?", id).Updates(map[string]interface{}{"StudentID": studentID, "SubjectID": subjectID, "Value": value}).Error
}

func (r *repo) GetMarks(ctx context.Context) (res []types.MarkResponse, err error) { //don`t touch
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

func (r *repo) DeleteMark(ctx context.Context, id string) error { //done
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&types.Mark{}).Error
}

func (r *repo) GetMark(ctx context.Context, id string) (res types.Mark, err error) { //done
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}

func (r *repo) AvgMarks(ctx context.Context) (res []types.MarkAverege, err error) { //don't touch
	err = r.db.WithContext(ctx).Raw(`
    SELECT student.id, student.name, subject.title, AVG(value) AS Value
    FROM mark 
        INNER JOIN student ON mark.student_id = student.id
        INNER JOIN subject ON mark.subject_id = subject.id
        GROUP BY student.id, student.name, subject.title 
        ORDER BY student.id ASC`).Scan(&res).Error
	return
}
