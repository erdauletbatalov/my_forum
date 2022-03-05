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
func (m *ForumModel) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	stmt := `INSERT INTO "users" ("email", "password", "nickname") 
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
					FROM "users" 
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

func (m *ForumModel) GetUser(email string) (*models.User, error) {
	stmt := `SELECT "id", "email", "username", "password"
					FROM "users"
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
