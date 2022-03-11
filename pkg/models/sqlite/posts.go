package sqlite

import (
	"database/sql"
	"errors"

	"forum/pkg/models"
)

func (m *ForumModel) AddPost(post *models.Post) (int, error) {
	stmt := `INSERT INTO "post" ("user_id", "title", "content") 
	VALUES(?, ?, ?)`

	result, err := m.DB.Exec(stmt, post.User_id, post.Title, post.Content)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *ForumModel) GetPostByID(id int) (*models.Post, error) {
	stmt := `SELECT p.id, p.user_id, u.username, p.title, p.content
					FROM post AS p
					INNER JOIN user AS u
					ON p.user_id = u.id
					WHERE p.id = ?;`
	post := &models.Post{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&post.ID, &post.User_id, &post.Author, &post.Title, &post.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return post, nil
}

func (m *ForumModel) GetAllPosts() ([]*models.Post, error) {
	stmt := `SELECT u.username, p.user_id, p.id, p.title, p.content 
					FROM post AS p
					INNER JOIN user AS u
					ON p.user_id = u.id`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*models.Post

	for rows.Next() {
		p := &models.Post{}
		err = rows.Scan(&p.Author, &p.User_id, &p.ID, &p.Title, &p.Content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
