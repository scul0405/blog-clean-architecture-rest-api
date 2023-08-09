package repository

const (
	createBlogQuery = `INSERT INTO blogs (author_id, title, content, image_url, category, created_at) 
					VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), now()) 
					RETURNING *`

	getBlogByIDQuery = `SELECT b.blog_id,
						   b.title,
						   b.content,
						   b.updated_at,
						   b.image_url,
						   b.category,
						   CONCAT(u.first_name, ' ', u.last_name) as author,
						   u.user_id as author_id
					FROM blogs b
							 LEFT JOIN users u on u.user_id = b.author_id
					WHERE blog_id = $1`

	updateBlogQuery = `UPDATE blogs 
					SET title = COALESCE(NULLIF($1, ''), title),
						content = COALESCE(NULLIF($2, ''), content), 
					    image_url = COALESCE(NULLIF($3, ''), image_url), 
					    category = COALESCE(NULLIF($4, ''), category), 
					    updated_at = now() 
					WHERE blog_id = $5
					RETURNING *`

	deleteBlogQuery = `DELETE FROM blogs WHERE blog_id = $1`

	listBlogsQuery = `SELECT b.blog_id, b.title, b.content, b.image_url, b.category, b.updated_at, b.created_at,  CONCAT(u.first_name, ' ', u.last_name) as author, u.user_id as author_id 
				FROM blogs b
					LEFT JOIN users u on u.user_id = b.author_id
				ORDER BY b.created_at, b.updated_at OFFSET $1 LIMIT $2`

	getTotalCountQuery = `SELECT COUNT(blog_id) FROM blogs`
)
