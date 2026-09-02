-- User's table
CREATE TABLE
	IF NOT EXISTS users (
		id CHAR(36) PRIMARY KEY, 
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL, 
		nick_name TEXT NOT NULL UNIQUE,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+0 hours')) 
	);

CREATE TABLE
	IF NOT EXISTS sessions (
		id CHAR(36) PRIMARY KEY,
		user_id CHAR(36) REFERENCES users (id) ON DELETE CASCADE,
		expire_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+24 hours')),
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+0 hours'))
	);

CREATE TABLE
	IF NOT EXISTS posts (
		id CHAR(36) PRIMARY KEY,
		author_id CHAR(36) REFERENCES users (id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+0 hours'))
	);
CREATE TABLE
	IF NOT EXISTS categories (
		id CHAR(36) PRIMARY KEY,
		title TEXT NOT NULL,
		icon TEXT NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+0 hours'))
	);

-- A post can carry several categories. The composite primary key rejects
-- duplicate pairs, and both foreign keys make an unknown category id a
-- constraint failure rather than an orphan row.
CREATE TABLE
	IF NOT EXISTS post_categories (
		post_id CHAR(36) NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
		category_id CHAR(36) NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
		PRIMARY KEY (post_id, category_id)
	);

CREATE TABLE
	IF NOT EXISTS comments (
		id CHAR(36) PRIMARY KEY,
		author_id CHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
		parent_id CHAR(36) NOT NULL,
		content TEXT NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+0 hours'))
	);


CREATE TABLE
	IF NOT EXISTS reactions (
		id CHAR(36) PRIMARY KEY,
		author_id CHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
		parent_id CHAR(36) NOT NULL,	
		is_like BOOLEAN NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+0 hours'))
	);

CREATE TABLE
	IF NOT EXISTS conversations (
		id CHAR(36) PRIMARY KEY
	);

CREATE TABLE
	IF NOT EXISTS conversation_participants (
		id CHAR(36) PRIMARY KEY,
		user_id CHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
		conversation_id CHAR(36) NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
		UNIQUE(user_id, conversation_id)
	);

CREATE TABLE
	IF NOT EXISTS messages (
		id CHAR(36) PRIMARY KEY,
		sender_id CHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE, -- actuallt this must be reviewed if a user delete whta s the correct practice 
		conversation_id CHAR(36) NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at INTEGER NOT NULL
	);

CREATE INDEX IF NOT EXISTS poster_idx ON posts (author_id);

CREATE INDEX IF NOT EXISTS session_idx ON sessions (user_id);

CREATE INDEX IF NOT EXISTS participants_idx ON conversation_participants (user_id);

CREATE INDEX IF NOT EXISTS sender_idx ON messages (sender_id);

CREATE INDEX IF NOT EXISTS conversation_idx ON messages (conversation_id);

CREATE INDEX IF NOT EXISTS post_idx ON comments (parent_id);

CREATE INDEX IF NOT EXISTS post_categories_category_idx ON post_categories (category_id);

-- Matches the listing's ORDER BY exactly, so cursor paging seeks straight to the
-- page instead of sorting the whole table on every request.
CREATE INDEX IF NOT EXISTS posts_feed_idx ON posts (created_at DESC, id DESC);


CREATE TRIGGER IF NOT EXISTS comment_satisfy
BEFORE INSERT ON comments 
BEGIN
	SELECT CASE	
		WHEN NOT EXISTS(SELECT * FROM comments WHERE id = NEW.parent_id)
		AND NOT EXISTS(SELECT * FROM posts WHERE id = NEW.parent_id) 
		THEN RAISE(ABORT, "parent not exists")
	END;
END;
CREATE TRIGGER IF NOT EXISTS likes_satisfy
BEFORE INSERT ON reactions 
BEGIN
	SELECT CASE	
		WHEN NOT EXISTS(SELECT * FROM comments WHERE id = NEW.parent_id)
		AND NOT EXISTS(SELECT * FROM posts WHERE id = NEW.parent_id) 
		THEN RAISE(ABORT, "parent not exists")
	END; 
END;
CREATE TRIGGER IF NOT EXISTS likes_insert
AFTER INSERT ON reactions 
BEGIN
	UPDATE posts  SET likes_count = likes_count+1 
	WHERE id = NEW.parent_id AND NEW.is_like = TRUE;   

	UPDATE posts  SET dislikes_count = dislikes_count+1 
	WHERE id = NEW.parent_id AND NEW.is_like = FALSE;   

END;
CREATE TRIGGER IF NOT EXISTS likes_update
AFTER UPDATE ON reactions 
BEGIN
	UPDATE posts  SET likes_count = likes_count+1, 
	dislikes_count = dislikes_count-1 
	WHERE id = NEW.parent_id AND NEW.is_like = TRUE ;  

	UPDATE posts  SET dislikes_count = dislikes_count+1, 
	likes_count = likes_count-1
	WHERE id = NEW.parent_id AND NEW.is_like = FALSE;   	
END;
CREATE TRIGGER IF NOT EXISTS likes_delete
AFTER DELETE ON reactions 
BEGIN
	UPDATE posts  SET likes_count = likes_count-1 
	WHERE id = OLD.parent_id AND OLD.is_like = TRUE ;  

	UPDATE posts  SET dislikes_count = dislikes_count-1 
	WHERE id = OLD.parent_id AND OLD.is_like = FALSE;   	
END;

CREATE TRIGGER IF NOT EXISTS cascade_delete_poosts
AFTER DELETE ON posts
BEGIN
	DELETE FROM comments WHERE old.id = parent_id;
	DELETE FROM reactions WHERE old.id = parent_id;
END;

CREATE TRIGGER IF NOT EXISTS cascade_delete_comments 
AFTER DELETE ON comments 
BEGIN
	DELETE FROM comments WHERE old.id = parent_id;	
	DELETE FROM reactions WHERE old.id = parent_id;	

END;
