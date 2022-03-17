package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"forum/pkg/models"
)

const (
	RFC822 = "02 Jan 06 15:04 MST"
)

func (m *ForumModel) AddPost(post *models.Post) (int, error) {
	stmt := `INSERT INTO "post" ("user_id", "title", "content", "date") 
	VALUES(?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, post.User_id, post.Title, post.Content, time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *ForumModel) GetPostByID(post_id int, user_id int) (*models.Post, error) {
	stmt := `SELECT p.id, p.user_id, u.username, p.title, p.content, p.date
					FROM post AS p
					INNER JOIN user AS u
					ON p.user_id = u.id
					WHERE p.id = ?;`
	post := &models.Post{}

	row := m.DB.QueryRow(stmt, post_id)

	err := row.Scan(&post.ID, &post.User_id, &post.Author, &post.Title, &post.Content, &post.Date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	vote := &models.Vote{
		Post_id:    post.ID,
		Comment_id: 0,
		Vote_type:  Like,
	}
	likes, err := m.GetVotes(vote)
	if err != nil {
		return nil, err
	}
	// checking if current user liked
	var isLike bool
	if user_id != 0 {
		isLike, err = m.isVotedByUser(user_id, vote)
		if err != nil {
			return nil, err
		}
	}

	vote.Vote_type = Dislike
	dislikes, err := m.GetVotes(vote)
	if err != nil {
		return nil, err
	}

	var isDislike bool
	if user_id != 0 {
		isDislike, err = m.isVotedByUser(user_id, vote)
		if err != nil {
			return nil, err
		}
	}

	post.Likes = likes
	post.Dislikes = dislikes
	post.IsLike = isLike
	post.IsDislike = isDislike

	return post, nil
}

func (m *ForumModel) GetAllPosts(user_id int) ([]*models.Post, error) {
	stmt := `SELECT u.username, p.user_id, p.id, p.title, p.content, p.date
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

		err = rows.Scan(&p.Author, &p.User_id, &p.ID, &p.Title, &p.Content, &p.Date)
		if err != nil {
			return nil, err
		}

		vote := &models.Vote{
			Post_id:    p.ID,
			Comment_id: 0,
			Vote_type:  Like,
		}
		likes, err := m.GetVotes(vote)
		if err != nil {
			return nil, err
		}
		// checking if current user liked
		var isLike bool
		if user_id != 0 {
			isLike, err = m.isVotedByUser(user_id, vote)
			if err != nil {
				return nil, err
			}
		}

		vote.Vote_type = Dislike
		dislikes, err := m.GetVotes(vote)
		if err != nil {
			return nil, err
		}

		var isDislike bool
		if user_id != 0 {
			isDislike, err = m.isVotedByUser(user_id, vote)
			if err != nil {
				return nil, err
			}
		}

		p.Likes = likes
		p.Dislikes = dislikes
		p.IsLike = isLike
		p.IsDislike = isDislike

		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *ForumModel) GetAllPostsSortedByLikes(user_id int) ([]*models.Post, error) {
	stmt := `SELECT u.username, p.user_id, p.id, p.title, p.content, p.date , COUNT(v.id) AS likes
					FROM post AS p
					INNER JOIN user AS u
					ON p.user_id = u.id
					LEFT JOIN "vote" AS "v"
					ON v.vote_type = 1
					AND v.post_id = p.id
					AND v.comment_id = 0
					GROUP BY p.id
					ORDER BY
					likes DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*models.Post

	for rows.Next() {
		p := &models.Post{}
		err = rows.Scan(&p.Author, &p.User_id, &p.ID, &p.Title, &p.Content, &p.Date, &p.Likes)
		if err != nil {
			return nil, err
		}

		vote := &models.Vote{
			Post_id:    p.ID,
			Comment_id: 0,
			Vote_type:  Like,
		}
		likes, err := m.GetVotes(vote)
		if err != nil {
			return nil, err
		}
		// checking if current user liked
		var isLike bool
		if user_id != 0 {
			isLike, err = m.isVotedByUser(user_id, vote)
			if err != nil {
				return nil, err
			}
		}

		vote.Vote_type = Dislike
		dislikes, err := m.GetVotes(vote)
		if err != nil {
			return nil, err
		}

		var isDislike bool
		if user_id != 0 {
			isDislike, err = m.isVotedByUser(user_id, vote)
			if err != nil {
				return nil, err
			}
		}

		p.Likes = likes
		p.Dislikes = dislikes
		p.IsLike = isLike
		p.IsDislike = isDislike

		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
