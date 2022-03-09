package sqlite

import (
	"database/sql"
	"errors"

	"github.com/erdauletbatalov/forum.git/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// SnippetModel - Определяем тип который обертывает пул подключения sql.DB
type ForumModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *ForumModel) AddUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	stmt := `INSERT INTO "user" ("email", "password", "nickname") 
	VALUES(?, ?, ?)`

	_, err = m.DB.Exec(stmt, user.Email, hashedPassword, user.Username)
	if err != nil {
		return err
	}
	return nil
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *ForumModel) LogInUser(user *models.User) error {
	stmt := `SELECT "password" 
					FROM "user" 
					WHERE "email" = ?`

	u := &models.User{}

	row := m.DB.QueryRow(stmt, user.Email)
	err := row.Scan(&u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecord
		} else {
			return err
		}
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		return err
	}
	return nil
}

func (m *ForumModel) GetUserByID(id int) (*models.User, error) {
	stmt := `SELECT "id", "email", "username", "password"
					FROM "user"
					WHERE "id" = ?`
	user := &models.User{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return user, nil
}

func (m *ForumModel) GetUserByEmail(email string) (*models.User, error) {
	stmt := `SELECT "id", "email", "username", "password"
					FROM "user"
					WHERE "email" = ?`
	user := &models.User{}

	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return user, nil
}

func (m *ForumModel) AddPost(post *models.Post) (int, error) {
	stmt := `INSERT INTO "post" ("user_id", "title", "content") 
	VALUES(?, ?, ?)`

	result, err := m.DB.Exec(stmt, post.ID, post.Title, post.Content)
	if err != nil {
		return 0, err
	}
	// Используем метод LastInsertId(), чтобы получить последний ID
	// созданной записи из таблицу snippets.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *ForumModel) GetPostByID(id int) (*models.Post, error) {
	stmt := `SELECT "id", "user_id", "title", "content"
					FROM "post"
					WHERE "id" = ?`
	post := &models.Post{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&post.ID, &post.User_id, &post.Title, &post.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return post, nil
}
