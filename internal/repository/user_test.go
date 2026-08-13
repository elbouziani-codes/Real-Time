package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"realTime/crypto"
	"realTime/internal/domain"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db")+"?_foreign_keys=on")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	schema, err := os.ReadFile("../../schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatal(err)
	}
	return db
}

func mustUUID(t *testing.T) crypto.UUID {
	t.Helper()
	id, err := crypto.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func insertUser(t *testing.T, db *sql.DB, id crypto.UUID, nickName string) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO users (id, email, password_hash, nick_name, first_name, last_name, age, gender)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id.Value, nickName+"@test.io", "hash", nickName, nickName, nickName, 20, "man")
	if err != nil {
		t.Fatal(err)
	}
}

// insertSharedMessage opens a conversation between the two users and drops one
// message at createdAt, which is what last_message_at reads.
func insertSharedMessage(t *testing.T, db *sql.DB, sender, other crypto.UUID, createdAt int) {
	t.Helper()
	conversationID := mustUUID(t)
	if _, err := db.Exec(`INSERT INTO conversations (id) VALUES (?)`, conversationID.Value); err != nil {
		t.Fatal(err)
	}
	for _, userID := range []crypto.UUID{sender, other} {
		if _, err := db.Exec(
			`INSERT INTO conversation_participants (id, user_id, conversation_id) VALUES (?, ?, ?)`,
			mustUUID(t).Value, userID.Value, conversationID.Value); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(
		`INSERT INTO messages (id, sender_id, conversation_id, content, created_at) VALUES (?, ?, ?, ?, ?)`,
		mustUUID(t).Value, sender.Value, conversationID.Value, "hello", createdAt); err != nil {
		t.Fatal(err)
	}
}

func TestGetUsersCursorPaging(t *testing.T) {
	db := openTestDB(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	requester := mustUUID(t)
	userA := mustUUID(t)
	userB := mustUUID(t)
	userC := mustUUID(t)
	userD := mustUUID(t)

	insertUser(t, db, requester, "requester")
	insertUser(t, db, userA, "auser")
	insertUser(t, db, userB, "buser")
	insertUser(t, db, userC, "cuser")
	insertUser(t, db, userD, "duser")

	insertSharedMessage(t, db, userA, requester, 300)
	insertSharedMessage(t, db, userB, requester, 200)
	// Same timestamp as B: the nick_name tie-breaker must put buser first.
	insertSharedMessage(t, db, userC, requester, 200)

	// Expected order: auser(300), buser(200), cuser(200), duser(0).
	firstPage, err := repo.GetUsers(ctx, requester, 2, crypto.Nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(firstPage) != 2 {
		t.Fatalf("first page: got %d users, want 2", len(firstPage))
	}
	if firstPage[0].NickName != "auser" || firstPage[1].NickName != "buser" {
		t.Fatalf("first page order: got %s, %s", firstPage[0].NickName, firstPage[1].NickName)
	}

	secondPage, err := repo.GetUsers(ctx, requester, 2, firstPage[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(secondPage) != 2 {
		t.Fatalf("second page: got %d users, want 2", len(secondPage))
	}
	if secondPage[0].NickName != "cuser" || secondPage[1].NickName != "duser" {
		t.Fatalf("second page order: got %s, %s", secondPage[0].NickName, secondPage[1].NickName)
	}

	// The cursor user must not repeat and the pages must cover everyone once.
	thirdPage, err := repo.GetUsers(ctx, requester, 2, secondPage[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(thirdPage) != 0 {
		t.Fatalf("third page: got %d users, want 0", len(thirdPage))
	}
}

func TestGetUsersUnknownCursor(t *testing.T) {
	db := openTestDB(t)
	repo := NewUserRepo(db)

	requester := mustUUID(t)
	insertUser(t, db, requester, "requester")

	_, err := repo.GetUsers(context.Background(), requester, 10, mustUUID(t))
	var domainErr domain.Error
	if !errors.As(err, &domainErr) || domainErr.Code != domain.NotFoundCode {
		t.Fatalf("want NotFoundCode domain error, got %v", err)
	}
}
