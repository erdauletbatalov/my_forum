package sqlite

import (
	"database/sql"
	"fmt"
	"forum/pkg/models"
)

type ForumModel struct {
	DB *sql.DB
}

func (m *ForumModel) AddComment(comment *models.Comment) error {
	stmt := `INSERT INTO "comment" ("user_id", "post_id", "content") 
	VALUES(?, ?, ?)`

	_, err := m.DB.Exec(stmt, comment.User_id, comment.Post_id, comment.Content)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (m *ForumModel) GetCommentsByPostID(post_id int) ([]*models.Comment, error) {
	stmt := `SELECT c.id, c.user_id, u.username, c.post_id, c.content
					FROM "comment" AS "c"
					INNER JOIN "user" AS "u"
					on u.id = c.user_id; `

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.Comment

	for rows.Next() {
		c := &models.Comment{}
		err = rows.Scan(&c.ID, &c.User_id, &c.Author, &c.Post_id, &c.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
