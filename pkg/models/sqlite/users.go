package sqlite

import (
	"database/sql"
	"errors"

	"github.com/erdauletbatalov/forum/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

type ForumModel struct {
	DB *sql.DB
}

func (m *ForumModel) AddUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO "user"(
		"username",
		"email",
		"password"
	) VALUES (?, ?, ?)`

	_, err = m.DB.Exec(stmt, user.Username, user.Email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (m *ForumModel) PasswordCompare(login, password string) error {
	s := `SELECT "password" FROM "user" 
	WHERE "username"=? OR "email"=?`
	row := m.DB.QueryRow(s, login, login)
	u := &models.User{}
	err := row.Scan(&u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecord
		}
		return err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (m *ForumModel) GetUserInfo(login string) (*models.User, error) {
	statement := `SELECT "id","username","email" FROM user 
					WHERE "username" = ? OR "email" = ?`
	row := m.DB.QueryRow(statement, login, login)
	u := &models.User{}
	return m.userQueryScan(row, u)
}

func (m *ForumModel) GetUserByID(id int) (*models.User, error) {
	stmt := `SELECT "id", "username", "email"
					FROM "user"
					WHERE "id" = ?`
	u := &models.User{}
	row := m.DB.QueryRow(stmt, id)
	return m.userQueryScan(row, u)
}

func (m *ForumModel) userQueryScan(row *sql.Row, user *models.User) (*models.User, error) {
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}
	return user, nil
}
