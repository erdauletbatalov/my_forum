package mysql

import (
	"database/sql"

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
	stmt := `INSERT INTO "user" ("email", "password", "nickname") 
	VALUES(?, ?, ?)`

	_, err = m.DB.Exec(stmt, user.Email, hashedPassword, user.Nickname)
	if err != nil {
		return err
	}
	return nil
}
