package repository

const (
	createUserCommentQuery = `INSERT INTO user_comments (user_id, comment_id) VALUES ($1, $2)`

	getUserCommentQuery = `SELECT *	
						FROM user_comments
						WHERE user_id = $1 AND comment_id = $2`

	deleteUserCommentQuery = `DELETE FROM user_comments WHERE user_id = $1 AND comment_id = $2`
)
