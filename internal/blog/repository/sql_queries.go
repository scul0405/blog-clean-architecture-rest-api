package repository

const (
	createBlog = `INSERT INTO blogs (author_id, title, content, image_url, category, created_at) 
					VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($4, ''), now()) 
					RETURNING *`

	getBlogByID = `SELECT n.blog_id,
						   n.title,
						   n.content,
						   n.updated_at,
						   n.image_url,
						   n.category,
						   CONCAT(u.first_name, ' ', u.last_name) as author,
						   u.user_id as author_id
					FROM news n
							 LEFT JOIN users u on u.user_id = n.author_id
					WHERE news_id = $1`

	updateBlog = `UPDATE blogs 
					SET title = COALESCE(NULLIF($1, ''), title),
						content = COALESCE(NULLIF($2, ''), content), 
					    image_url = COALESCE(NULLIF($3, ''), image_url), 
					    category = COALESCE(NULLIF($4, ''), category), 
					    updated_at = now() 
					WHERE blog_id = $5
					RETURNING *`

	deleteBlog = `DELETE FROM blogs WHERE blog_id = $1`

	getBlogs = `SELECT blog_id, author_id, title, content, image_url, category, updated_at, created_at 
				FROM blogs 
				ORDER BY created_at, updated_at OFFSET $1 LIMIT $2`

	getTotalCount = `SELECT COUNT(blog_id) FROM blogs`
)
