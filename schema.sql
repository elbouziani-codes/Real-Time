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
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')), 
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')) 
	);

CREATE TABLE
	IF NOT EXISTS sessions (
		id CHAR(36) PRIMARY KEY,
		user_id CHAR(36) REFERENCES users (id) ON DELETE CASCADE,
		expire_at TEXT DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+24 hours')),
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE TABLE
	IF NOT EXISTS posts (
		id CHAR(36) PRIMARY KEY,
		author_id CHAR(36) REFERENCES users (id) ON DELETE CASCADE,
		category CHAR(36) REFERENCES categories (id) ON DELETE CASCADE,
		title TEXT NOT NULL, 
		content TEXT NOT NULL,
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);
CREATE TABLE
	IF NOT EXISTS categories (
		id CHAR(36) PRIMARY KEY,
		title TEXT NOT NULL, 
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE TABLE
	IF NOT EXISTS comments (
		id CHAR(36) PRIMARY KEY,
		author_id CHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
		parent_id CHAR(36) NOT NULL,
		content TEXT NOT NULL,
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);


CREATE TABLE
	IF NOT EXISTS reactions (
		id CHAR(36) PRIMARY KEY,
		author_id CHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
		parent_id CHAR(36) NOT NULL,	
		is_like BOOLEAN NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))	
	);

CREATE TABLE
	IF NOT EXISTS conversations (
		id CHAR(36) PRIMARY KEY,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE TABLE
	IF NOT EXISTS conversation_participants (
		id CHAR(36) PRIMARY KEY,
		user_id CHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
		conversation_id CHAR(36) NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
		joined_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE TABLE
	IF NOT EXISTS messages (
		id CHAR(36) PRIMARY KEY,
		sender_id CHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE, -- actuallt this must be reviewed if a user delete whta s the correct practice 
		conversation_id CHAR(36) NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE INDEX IF NOT EXISTS poster_idx ON posts (author_id);

CREATE INDEX IF NOT EXISTS session_idx ON sessions (user_id);

CREATE INDEX IF NOT EXISTS participants_idx ON conversation_participants (user_id);

CREATE INDEX IF NOT EXISTS sender_idx ON messages (sender_id);

CREATE INDEX IF NOT EXISTS conversation_idx ON messages (conversation_id);

CREATE INDEX IF NOT EXISTS post_idx ON comments (parent_id);


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
	UPDATE comments  SET likes_count = likes_count+1 
	WHERE id = NEW.parent_id AND NEW.is_like = TRUE; 

	UPDATE posts  SET likes_count = likes_count+1 
	WHERE id = NEW.parent_id AND NEW.is_like = TRUE;   

	UPDATE comments  SET dislikes_count = dislikes_count+1 
	WHERE id = NEW.parent_id AND NEW.is_like = FALSE;   

	UPDATE posts  SET dislikes_count = dislikes_count+1 
	WHERE id = NEW.parent_id AND NEW.is_like = FALSE;   

END;
CREATE TRIGGER IF NOT EXISTS likes_update
AFTER UPDATE ON reactions 
BEGIN
	UPDATE comments  SET likes_count = likes_count+1, 
	dislikes_count = dislikes_count-1 
	WHERE id = NEW.parent_id AND NEW.is_like = TRUE ;  

	UPDATE posts  SET likes_count = likes_count+1, 
	dislikes_count = dislikes_count-1 
	WHERE id = NEW.parent_id AND NEW.is_like = TRUE ;  

	UPDATE comments  SET dislikes_count = dislikes_count+1, 
	likes_count = likes_count-1 
	WHERE id = NEW.parent_id AND NEW.is_like = FALSE;   

	UPDATE posts  SET dislikes_count = dislikes_count+1, 
	likes_count = likes_count-1
	WHERE id = NEW.parent_id AND NEW.is_like = FALSE;   	
END;
CREATE TRIGGER IF NOT EXISTS likes_delete
AFTER DELETE ON reactions 
BEGIN
	UPDATE comments  SET likes_count = likes_count-1 
	WHERE id = OLD.parent_id AND OLD.is_like = TRUE ;  

	UPDATE posts  SET likes_count = likes_count-1 
	WHERE id = OLD.parent_id AND OLD.is_like = TRUE ;  

	UPDATE comments  SET dislikes_count = dislikes_count-1 
	WHERE id = OLD.parent_id AND OLD.is_like = FALSE;   

	UPDATE posts  SET dislikes_count = dislikes_count-1 
	WHERE id = OLD.parent_id AND OLD.is_like = FALSE;   	
END;

CREATE TRIGGER IF NOT EXISTS cascade_delete_poosts 
AFTER DELETE ON posts 
BEGIN
	DELETE FROM comments WHERE old.id = parent_id;	
	DELETE FROM reaction WHERE old.id = parent_id;	
END;

CREATE TRIGGER IF NOT EXISTS cascade_delete_comments 
AFTER DELETE ON comments 
BEGIN
	DELETE FROM comments WHERE old.id = parent_id;	
	DELETE FROM reactions WHERE old.id = parent_id;	

END;

CREATE TRIGGER IF NOT EXISTS update_date BEFORE
UPDATE ON users FOR EACH ROW WHEN NEW.updated_at = OLD.updated_at
OR NEW.updated_at IS NULL BEGIN
UPDATE users
SET
	updated_at = CURRENT_TIMESTAMP
WHERE
	NEW.id = id;

END;
