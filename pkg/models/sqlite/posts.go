package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum/pkg/models"
)

const (
	RFC822 = "02 Jan 06 15:04 MST"
)

func (m *ForumModel) AddPost(post *models.Post) (int, error) {
	var err error
	stmt := `INSERT INTO "post" ("user_id", "title", "content", "date") 
	VALUES(?, ?, ?, ?)`

	tag := &models.Tag{}

	result, err := m.DB.Exec(stmt, post.User_id, post.Title, post.Content, time.Now())
	if err != nil {
		return 0, err
	}
	post_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, val := range post.Tags {
		tag.Post_id = int(post_id)
		tag.Name = val
		err = m.AddTag(tag)
		if err != nil {
			return 0, err
		}
	}

	return int(post_id), nil
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

	tags, err := m.GetTagsByPostID(post_id)
	fmt.Println(tags)
	if err != nil {
		return nil, err
	}

	post.Tags = tags

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

func (m *ForumModel) GetPosts(user_id int) ([]*models.Post, error) {
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

		tags, err := m.GetTagsByPostID(p.ID)
		fmt.Println(tags)
		if err != nil {
			return nil, err
		}

		p.Tags = tags

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

func (m *ForumModel) GetUserPosts(session_user_id int, by_user_id int) ([]*models.Post, error) {
	stmt := `SELECT u.username, p.user_id, p.id, p.title, p.content, p.date
					FROM post AS p
					INNER JOIN user AS u
					ON p.user_id = u.id
					WHERE p.user_id = ?`

	rows, err := m.DB.Query(stmt, by_user_id)
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

		tags, err := m.GetTagsByPostID(p.ID)
		fmt.Println(tags)
		if err != nil {
			return nil, err
		}

		p.Tags = tags

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
		if session_user_id != 0 {
			isLike, err = m.isVotedByUser(session_user_id, vote)
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
		if session_user_id != 0 {
			isDislike, err = m.isVotedByUser(session_user_id, vote)
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

func (m *ForumModel) GetLikedUserPosts(session_user_id int, by_user int) ([]*models.Post, error) {
	stmt := `SELECT u.username, p.user_id, p.id, p.title, p.content, p.date
					FROM post AS p
					INNER JOIN user AS u
					ON p.user_id = u.id
					INNER JOIN vote AS v
					ON v.post_id = p.id
					WHERE p.user_id = ?
					AND v.comment_id = 0
					AND v.user_id = ?`

	rows, err := m.DB.Query(stmt, by_user, by_user)
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

		tags, err := m.GetTagsByPostID(p.ID)
		fmt.Println(tags)
		if err != nil {
			return nil, err
		}

		p.Tags = tags

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
		if session_user_id != 0 {
			isLike, err = m.isVotedByUser(session_user_id, vote)
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
		if session_user_id != 0 {
			isDislike, err = m.isVotedByUser(session_user_id, vote)
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

func (m *ForumModel) GetPostsSortedByLikes(user_id int) ([]*models.Post, error) {
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

		tags, err := m.GetTagsByPostID(p.ID)
		fmt.Println(tags)
		if err != nil {
			return nil, err
		}

		p.Tags = tags

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

func (m *ForumModel) GetPostsSortedByDate(user_id int) ([]*models.Post, error) {
	stmt := `SELECT u.username, p.user_id, p.id, p.title, p.content, p.date, COUNT(v.id) AS likes
					FROM post AS p
					INNER JOIN user AS u
					ON p.user_id = u.id
					LEFT JOIN "vote" AS "v"
					ON v.vote_type = 1
					AND v.post_id = p.id
					AND v.comment_id = 0
					GROUP BY p.id
					ORDER BY
					p.date DESC`

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

		tags, err := m.GetTagsByPostID(p.ID)
		fmt.Println(tags)
		if err != nil {
			return nil, err
		}

		p.Tags = tags

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

func (m *ForumModel) GetPostsByTag(user_id int, tag_name string) ([]*models.Post, error) {
	stmt := `SELECT u.username, p.user_id, p.id, p.title, p.content, p.date , COUNT(v.id) AS likes
					FROM post AS p
					INNER JOIN user AS u
					ON p.user_id = u.id
					INNER JOIN post_tag AS pt
					ON pt.post_id = p.id
					INNER JOIN tag AS t
					ON pt.tag_id = t.id
					LEFT JOIN "vote" AS "v"
					ON v.vote_type = 1
					AND v.post_id = p.id
					AND v.comment_id = 0
					WHERE t.name = ?
					GROUP BY p.id`

	rows, err := m.DB.Query(stmt, tag_name)
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

		tags, err := m.GetTagsByPostID(p.ID)
		fmt.Println(tags)
		if err != nil {
			return nil, err
		}

		p.Tags = tags

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
