package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/pkg/models"
	"strconv"
)

func (m *ForumModel) GetVoteType(vote *models.Vote) (int, error) {
	stmt_select := `SELECT "vote_type" 
									FROM "vote"
									WHERE user_id = ?
									AND post_id = ?
									AND comment_id = ?`

	row := m.DB.QueryRow(stmt_select, vote.User_id, vote.Post_id, vote.Comment_id)

	var vote_type int
	err := row.Scan(&vote_type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // there is no votes yet
		} else {
			return 0, err
		}
	}
	return vote_type, nil
}

func (m *ForumModel) AddVote(vote *models.Vote) error {
	stmt_insert := `INSERT INTO "vote" ("post_id", "comment_id", "user_id", "vote_type") 
									VALUES(?, ?, ?, ?)`

	_, err := m.DB.Exec(stmt_insert, vote.Post_id, vote.Comment_id, vote.User_id, vote.Vote_type)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("added %v vote type by %v user, %v post to comment %v\n", vote.Vote_type, vote.User_id, vote.Post_id, vote.Comment_id)
	return nil
}

func (m *ForumModel) DeleteVote(vote *models.Vote) error {
	stmt_delete := `DELETE FROM "vote"
									WHERE user_id = ?
									AND post_id = ?
									AND comment_id = ?`
	_, err := m.DB.Exec(stmt_delete, vote.User_id, vote.Post_id, vote.Comment_id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("deleted vote of %v user, %v post, comment %v\n", vote.User_id, vote.Post_id, vote.Comment_id)
	return nil
}

func (m *ForumModel) GetVotes(vote *models.Vote) (int, error) {
	stmt_select := `SELECT COUNT(id)
									FROM vote
									WHERE vote_type = ?
									AND post_id = ?
									AND comment_id = ?`

	var output string
	row := m.DB.QueryRow(stmt_select, vote.Vote_type, vote.Post_id, vote.Comment_id)
	fmt.Printf("searching %v, %v, %v, %v\n", vote.Vote_type, vote.Post_id, vote.Comment_id)

	err := row.Scan(&output)
	// Catch errors
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No likes found.\n")
		return 0, nil
	case err != nil:
		fmt.Printf("%s", err)
		return 0, err
	default:
		fmt.Printf("Counted %s likes \n", output)
		result, err := strconv.Atoi(output)
		if err != nil {
			return 0, err
		}
		return result, nil
	}
}

func (m *ForumModel) isVotedByUser(user_id int, vote *models.Vote) (bool, error) {
	stmt_select := `SELECT vote_type
									FROM vote
									WHERE vote_type = ?
									AND post_id = ?
									AND comment_id = ?
									AND user_id = ?`

	var vote_type int
	row := m.DB.QueryRow(stmt_select, vote.Vote_type, vote.Post_id, vote.Comment_id, user_id)
	fmt.Printf("searching user by %v, %v, %v, %v\n", vote.Vote_type, vote.Post_id, vote.Comment_id, vote.User_id)

	err := row.Scan(&vote_type)
	// Catch errors
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No votes of that user found.\n")
		return false, nil
	case err != nil:
		fmt.Printf("%s", err)
		return false, err
	default:
		return true, nil
	}
}
