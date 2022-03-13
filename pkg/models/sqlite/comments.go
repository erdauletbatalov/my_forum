package sqlite

import (
	"fmt"
	"forum/pkg/models"
)

const Post = 0
const Comment = 1
const Like = 1
const Dislike = -1

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

func (m *ForumModel) GetCommentsByPostID(post_id int, user_id int) ([]*models.Comment, error) {
	stmt := `SELECT c.id, c.user_id, u.username, c.post_id, c.content
					FROM "comment" AS "c"
					INNER JOIN "user" AS "u"
					on u.id = c.user_id
					WHERE c.post_id = ?`

	rows, err := m.DB.Query(stmt, post_id)
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
		vote := &models.Vote{
			Vote_obj:   Comment,
			Post_id:    post_id,
			Comment_id: c.ID,
			Vote_type:  Like,
		}
		likes, err := m.GetVotes(vote)
		if err != nil {
			return nil, err
		}
		// checking if user liked
		var isLike bool
		if user_id != 0 {
			isLike, err = m.isVote(user_id, vote)
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
			isDislike, err = m.isVote(user_id, vote)
			if err != nil {
				return nil, err
			}
		}

		c.Likes = likes
		c.Dislikes = dislikes
		c.IsLike = isLike
		c.IsDislike = isDislike

		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
