package repository

const (
	createCommentQuery = `INSERT INTO comments (author_id, blog_id, message) VALUES ($1, $2, $3) RETURNING *`

	updateCommentQuery = `UPDATE comments SET message = $1, updated_at = CURRENT_TIMESTAMP WHERE comment_id = $2 RETURNING *`

	deleteCommentQuery = `DELETE FROM comments WHERE comment_id = $1`

	getCommentByIDQuery = `SELECT concat(u.first_name, ' ', u.last_name) as author, u.avatar as avatar_url, c.message, count(uc.comment_id) as likes, c.updated_at, c.author_id, c.comment_id, c.blog_id	
						FROM comments c
        				LEFT JOIN users u on c.author_id = u.user_id LEFT JOIN user_comments uc on c.comment_id = uc.comment_id 
						WHERE c.comment_id = $1
						GROUP BY u.first_name, u.last_name, u.avatar, c.message, c.updated_at, c.author_id, c.comment_id, c.blog_id`

	getTotalCountByBlogIDQuery = `SELECT COUNT(comment_id) FROM comments WHERE blog_id = $1`

	listCommentsByBlogIDQuery = `SELECT concat(u.first_name, ' ', u.last_name) as author, u.avatar as avatar_url, c.message, count(uc.comment_id) as likes, c.updated_at, c.author_id, c.comment_id, c.blog_id
							FROM comments c
        					LEFT JOIN users u on c.author_id = u.user_id LEFT JOIN user_comments uc on c.comment_id = uc.comment_id 
        					WHERE c.blog_id = $1 
        					GROUP BY u.first_name, u.last_name, u.avatar, c.message, c.updated_at, c.author_id, c.comment_id, c.blog_id
							ORDER BY updated_at OFFSET $2 LIMIT $3`
)
