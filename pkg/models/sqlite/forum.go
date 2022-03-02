package mysql

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

	_, err = m.DB.Exec(stmt, user.Email, hashedPassword, user.Nickname)
	if err != nil {
		return err
	}
	return nil
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *ForumModel) LogInUser(user *models.User) (*models.User, error) {
	stmt := `SELECT "password" 
					FROM "users" 
					WHERE "email" = ?`

	u := &models.User{}

	row := m.DB.QueryRow(stmt, user.Email)
	err := row.Scan(&u.Password)
	if err != nil {
		// Специально для этого случая, мы проверим при помощи функции errors.Is()
		// если запрос был выполнен с ошибкой. Если ошибка обнаружена, то
		// возвращаем нашу ошибку из модели models.ErrNoRecord.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		return nil, err
	}
	return u, nil
}
