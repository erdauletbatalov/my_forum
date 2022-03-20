package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/pkg/models"
)

func (m *ForumModel) AddTag(tag *models.Tag) error {
	stmtInsertTag := `INSERT INTO "tag" ("name")
										VALUES(?)`

	stmtInsertPostTag := `INSERT INTO "post_tag" ("post_id", "tag_id") 
												VALUES(?, ?)`

	tag_id, err := m.getTagID(tag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Inserting new Tag name into Tag table

			_, err := m.DB.Exec(stmtInsertTag, tag.Name)
			if err != nil {
				return err
			}
			tag_id, _ = m.getTagID(tag)
		} else {
			return err
		}
	}

	_, err = m.DB.Exec(stmtInsertPostTag, tag.Post_id, tag_id)
	if err != nil {
		return err
	}
	return nil
}

func (m *ForumModel) GetTagsByPostID(post_id int) ([]string, error) {
	stmt := `SELECT tag.name, post_tag.post_id
					FROM tag
					INNER JOIN post_tag
					ON post_tag.tag_id = tag.id
					WHERE post_tag.post_id = ?`

	rows, err := m.DB.Query(stmt, post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []string
	count := 0
	for rows.Next() {
		var tag string
		var temp int

		fmt.Println(count)
		count++
		err = rows.Scan(&tag, &temp)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
		fmt.Println(tag, temp)
	}
	return tags, nil
}

func (m *ForumModel) GetTags() ([]string, error) {
	stmt := `SELECT tag.name
					FROM tag`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []string
	count := 0
	for rows.Next() {
		var tag string

		fmt.Println(count)
		count++
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (m *ForumModel) getTagID(tag *models.Tag) (int, error) {
	var tag_id int
	stmtSelect := `SELECT ("id")
								FROM "tag"
								WHERE "name" = ?`
	row := m.DB.QueryRow(stmtSelect, tag.Name)
	err := row.Scan(&tag_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, sql.ErrNoRows
		} else {
			return 0, err
		}
	}
	return tag_id, nil
}
