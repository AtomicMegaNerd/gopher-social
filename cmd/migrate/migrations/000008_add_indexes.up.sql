-- Enable the pg_trgm extension for trigram indexing
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Description: This script creates a GIN index on the 'content' column of the 
-- 'comments' table using the trigram operator for efficient text search.
-- This will make writes slower but significantly speed up read operations involving text search.
CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING gin(content gin_trgm_ops);

-- Description: This script creates GIN indexes on the 'title' and 'tags' columns of the 
-- 'posts' table. The 'title' column uses the trigram operator for efficient text search,
-- while the 'tags' column uses a standard GIN index for array containment queries.
-- This will improve read performance for searches on these columns at the cost of slower writes.
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin(title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin(tags);

-- Description: This script creates B-tree indexes on the 'username' and 'user_id'
-- columns of the 'users' table and the 'post_id' column of the 'comments' table.
-- These indexes will enhance the performance of queries filtering or sorting by these columns.
-- The trade-off is that write operations will be slightly slower due to the overhead 
-- of maintaining the indexes. We are not using gin indexes here as these columns
-- are not suitable for GIN indexing.
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

