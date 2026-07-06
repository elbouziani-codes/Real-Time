
-- User's table
CREATE TABLE IF NOT EXISTS  users  (
	id BLOB PRIMARY KEY, -- I will use binary format([]byte) because sqlite has no reserved type for uuid
	nick_name TEXT NOT NULL UNIQUE,
	first_name TEXT NOT NULL , 
	last_name TEXT  NOT NULL , 
	email TEXT NOT NULL  UNIQUE, 
	password_hash TEXT NOT NULL, --length enforced by hashing algorithm in app level (some for other fields)
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT
); 


CREATE TABLE IF NOT EXISTS sessions  (
	id BLOB PRIMARY KEY,
	user_id BLOB REFERENCES users(id) ON DELETE CASCADE, 
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT
);

CREATE TABLE IF NOT EXISTS posts  (
	id BLOB PRIMARY KEY,
	user_id BLOB REFERENCES users(id) ON DELETE CASCADE, 
	content TEXT  NOT NULL , 
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT
);

CREATE TABLE IF NOT EXISTS comments  (
	id BLOB PRIMARY KEY,
	user_id BLOB NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
	post_id BLOB NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
	parent_id BLOB NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
	content TEXT  NOT NULL, 
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT
);


CREATE TABLE IF NOT EXISTS likes  (
	id BLOB PRIMARY KEY,
	user_id BLOB NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
	post_id BLOB REFERENCES posts(id) ON DELETE CASCADE,
	comment_id BLOB  REFERENCES comments(id) ON DELETE CASCADE,
	is_like  BOOLEAN NOT NULL, 
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT,

		CHECK(
			(comment_id IS NOT NULL AND post_id IS NULL) OR
			(post_id IS NOT NULL AND comment_id IS NULL)			
		)
	
);


CREATE TABLE IF NOT EXISTS conversations  (
	id BLOB PRIMARY KEY,
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT
);

CREATE TABLE IF NOT EXISTS conversation_participants  (
	id BLOB PRIMARY KEY,
	user_id BLOB NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	conversation_id BLOB NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
	joined_at TEXT DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE IF NOT EXISTS messages  (
	id BLOB PRIMARY KEY,
	sender_id BLOB NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- actuallt this must be reviewed if a user delete whta s the correct practice 
	conversation_id BLOB NOT NULL REFERENCES conversations(id) ON DELETE CASCADE, 
	content TEXT NOT NULL,
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT
);





CREATE INDEX IF NOT EXISTS poster_idx ON posts(user_id);
CREATE INDEX IF NOT EXISTS session_idx ON sessions(user_id);
CREATE INDEX IF NOT EXISTS  participants_idx ON conversation_participants(user_id);
CREATE INDEX IF NOT EXISTS  sender_idx ON messages(sender_id);
CREATE INDEX IF NOT EXISTS  conversation_idx ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS  post_idx ON comments(post_id);





CREATE TRIGGER IF NOT EXISTS update_date 
BEFORE UPDATE ON users
FOR EACH ROW 
		WHEN NEW.updated_at = OLD.updated_at OR NEW.updated_at IS  NULL 
				BEGIN
					UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE NEW.id = id;
				END;

