-- User's table
CREATE TABLE
	IF NOT EXISTS users (
		id TEXT PRIMARY KEY, 
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
		id TEXT PRIMARY KEY,
		user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
		expiry_date TEXT DEFAULT (datetime (CURRENT_TIMESTAMP, '+24 hours')),
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE TABLE
	IF NOT EXISTS posts (
		id TEXT PRIMARY KEY,
		user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE TABLE
	IF NOT EXISTS comments (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
		post_id TEXT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
		parent_id TEXT NOT NULL REFERENCES comments (id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE TABLE
	IF NOT EXISTS likes (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
		post_id TEXT REFERENCES posts (id) ON DELETE CASCADE,
		comment_id TEXT REFERENCES comments (id) ON DELETE CASCADE,
		is_like BOOLEAN NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		CHECK (
			(
				comment_id IS NOT NULL
				AND post_id IS NULL
			)
			OR (
				post_id IS NOT NULL
				AND comment_id IS NULL
			)
		)
	);

CREATE TABLE
	IF NOT EXISTS conversations (
		id TEXT PRIMARY KEY,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE TABLE
	IF NOT EXISTS conversation_participants (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
		conversation_id TEXT NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
		joined_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE TABLE
	IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		sender_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE, -- actuallt this must be reviewed if a user delete whta s the correct practice 
		conversation_id TEXT NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);

CREATE INDEX IF NOT EXISTS poster_idx ON posts (user_id);

CREATE INDEX IF NOT EXISTS session_idx ON sessions (user_id);

CREATE INDEX IF NOT EXISTS participants_idx ON conversation_participants (user_id);

CREATE INDEX IF NOT EXISTS sender_idx ON messages (sender_id);

CREATE INDEX IF NOT EXISTS conversation_idx ON messages (conversation_id);

CREATE INDEX IF NOT EXISTS post_idx ON comments (post_id);

CREATE TRIGGER IF NOT EXISTS update_date BEFORE
UPDATE ON users FOR EACH ROW WHEN NEW.updated_at = OLD.updated_at
OR NEW.updated_at IS NULL BEGIN
UPDATE users
SET
	updated_at = CURRENT_TIMESTAMP
WHERE
	NEW.id = id;

END;
