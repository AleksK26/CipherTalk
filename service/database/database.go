package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrUserDoesNotExist = errors.New("User does not exist")
var ErrConversationDoesNotExist = errors.New("Conversation does not exist")
var ErrMessageDoesNotExist = errors.New("Message does not exist")
var ErrCommentDoesNotExist = errors.New("Comment does not exist")
var ErrUnauthorizedToDeleteMessage = errors.New("Unauthorized To Delete Message")
var ErrGroupDoesNotExist = errors.New("Group does not exist")
var ErrWrongPassword = errors.New("Wrong password")

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Photo    []byte `json:"photo,omitempty"`
	Password string `json:"-"`
}

type GroupMember struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Photo []byte `json:"photo,omitempty"`
	Role  string `json:"role"`
}

type Group struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Photo []byte `json:"photo,omitempty"`
}

type Conversation struct {
	Id                string         `json:"id"`
	Name              string         `json:"name"`
	Type              string         `json:"type"`
	CreatedAt         string         `json:"createdAt"`
	Members           []string       `json:"members"`
	LastMessage       *Message       `json:"lastMessage,omitempty"`
	Messages          []Message      `json:"messages,omitempty"`
	ConversationPhoto sql.NullString `json:"conversationPhoto,omitempty"`
}

type Message struct {
	Id                string   `json:"id"`
	ConversationId    string   `json:"conversationId"`
	SenderId          string   `json:"senderId"`
	SenderName        string   `json:"senderName"`
	Content           string   `json:"content"`
	Timestamp         string   `json:"timestamp"`
	Attachment        []byte   `json:"attachment"`
	SenderPhoto       string   `json:"senderPhoto,omitempty"`
	ReactionCount     int      `json:"reactionCount"`
	ReactingUserNames []string `json:"reactingUserNames"`
	Status            string   `json:"status"`
	ReplyTo           string   `json:"replyTo,omitempty"`
	ReplyContent      string   `json:"replyContent,omitempty"`
	ReplySenderName   string   `json:"replySenderName,omitempty"`
	ReplyAttachment   []byte   `json:"replyAttachment,omitempty"`
	IsDelivered       bool     `json:"isDelivered"`
	IsRead            bool     `json:"isRead"`
}

type Comment struct {
	Id       string `json:"id"`
	AuthorId string `json:"authorId"`
}

type ReadReceipt struct {
	MessageId   string  `json:"messageId"`
	UserId      string  `json:"userId"`
	DeliveredAt string  `json:"deliveredAt"`
	ReadAt      *string `json:"readAt,omitempty"`
}

type AppDatabase interface {
	GetGroupMemberDetails(groupId string) ([]GroupMember, error)
	IsGroupAdmin(groupID, userID string) (bool, error)
	Ping() error
	GetUserByName(name string) (User, error)
	CreateUser(u User) (User, error)
	UpdateUserName(userId string, newName string) (User, error)
	UpdateUserPhoto(userID string, photo []byte) error
	SearchUsersByName(username string) ([]User, error)
	GetDirectConversation(senderID, recipientID string) (string, error)
	CreateDirectConversation(conversationID, senderID, recipientID string) error
	SaveMessage(conversationID, senderID, messageID, content string, attachment []byte, replyTo string) (Message, error)
	InsertDeliveryReceipt(messageID, userID, deliveredAt string) error
	IsUserInConversation(conversationID, userID string) (bool, error)
	GetConversationDetails(conversationID, currentUserID string) (Conversation, error)
	GetMessagesForConversation(conversationID string) ([]Message, error)
	GetMyConversations(userID string) ([]Conversation, error)
	GetConversationMembers(conversationID string) ([]string, error)
	GetUsersPhoto(userID string) (User, error)
	DeleteMessage(conversationID, messageID, userID string) error
	GetMessage(messageID, userID string) (Message, error)
	CreateGroupConversation(conversationID string, memberIDs []string, name string, photo []byte, creatorID string) error
	GetMyGroups(userID string) ([]Conversation, error)
	GetGroupInfo(groupID string) (Conversation, error)
	UpdateGroupName(groupId, newName string) error
	UpdateGroupPhoto(groupID string, photo []byte) error
	LeaveGroup(groupID, userID string) error
	AddUserToGroup(conversationID string, userID string) error
	RemoveUserFromGroup(groupID, userID string) error
	CommentMessage(commentID, messageID, authorID string) error
	UncommentMessage(messageID, authorID string) error
	MarkMessagesAsRead(conversationID, userID string) error
	GetConversationById(conversationID string) (Conversation, error)
}

type appdbimpl struct {
	c *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}
	// Ensure the database directory exists
	dbPath := dbPathFromDSN(db)
	if dbPath != "" {
		dir := filepath.Dir(dbPath)
		_ = os.MkdirAll(dir, 0755)
	}
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}
	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		usersTable := `CREATE TABLE users (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			photo BLOB
		);`
		conversationsTable := `CREATE TABLE conversations (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			created_at TEXT NOT NULL,
			conversationPhoto BLOB
		);`
		conversationMembersTable := `CREATE TABLE conversation_members (
			conversationId TEXT NOT NULL,
			userId TEXT NOT NULL,
			FOREIGN KEY (conversationId) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE,
			PRIMARY KEY(conversationId, userId)
		);`
		messagesTable := `CREATE TABLE messages (
			id TEXT NOT NULL PRIMARY KEY,
			conversationId TEXT NOT NULL,
			senderId TEXT NOT NULL,
			content TEXT NOT NULL,
			timestamp TEXT NOT NULL,
			attachment BLOB,
			replyTo TEXT,  
			FOREIGN KEY (conversationId) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (senderId) REFERENCES users(id) ON DELETE CASCADE
		);`
		commentsTable := `CREATE TABLE comments (
			id TEXT NOT NULL PRIMARY KEY,
			messageId TEXT NOT NULL,
			authorId TEXT NOT NULL,
			UNIQUE(messageId, authorId),
			FOREIGN KEY (messageId) REFERENCES messages(id) ON DELETE CASCADE,
			FOREIGN KEY (authorId) REFERENCES users(id) ON DELETE CASCADE
		);`
		readReceiptsTable := `CREATE TABLE read_receipts (
			messageId TEXT NOT NULL,
			userId TEXT NOT NULL,
			deliveredAt TEXT NOT NULL,
			readAt TEXT, 
			PRIMARY KEY (messageId, userId),
			FOREIGN KEY (messageId) REFERENCES messages(id) ON DELETE CASCADE,
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);`

		creationQueries := []string{
			usersTable,
			conversationsTable,
			conversationMembersTable,
			messagesTable,
			commentsTable,
			readReceiptsTable,
		}
		for _, q := range creationQueries {
			_, execErr := db.Exec(q)
			if execErr != nil {
				return nil, fmt.Errorf("error creating database structure: %w", execErr)
			}
		}
		// Adding index queries for better performance
		indexQueries := []string{
			`CREATE INDEX IF NOT EXISTS idx_conversation_members_user ON conversation_members(userId)`,
			`CREATE INDEX IF NOT EXISTS idx_conversation_members_conv ON conversation_members(conversationId)`,
			`CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversationId)`,
			`CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(conversationId, timestamp)`,
			`CREATE INDEX IF NOT EXISTS idx_messages_sender ON messages(senderId)`,
			`CREATE INDEX IF NOT EXISTS idx_read_receipts_message ON read_receipts(messageId)`,
			`CREATE INDEX IF NOT EXISTS idx_comments_message ON comments(messageId)`,
			`CREATE INDEX IF NOT EXISTS idx_users_name ON users(name)`,
		}

		for _, q := range indexQueries {
			_, err := db.Exec(q)
			if err != nil {
				return nil, fmt.Errorf("error creating indexes: %w", err)
			}
		}

	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Schema migrations — safe to run on both new and existing databases
	migrations := []string{
		`ALTER TABLE users ADD COLUMN password TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE conversation_members ADD COLUMN role TEXT NOT NULL DEFAULT 'member'`,
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			// SQLite returns "duplicate column name" when column already exists — ignore that
			if !strings.Contains(err.Error(), "duplicate column name") {
				return nil, fmt.Errorf("migration error: %w", err)
			}
		}
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// dbPathFromDSN extracts the file path from the DSN if using SQLite3
func dbPathFromDSN(db *sql.DB) string {
	// This is a workaround: in your main.go, you have the path as cfg.DB.Filename
	// so you should ensure the directory exists there, before calling New().
	// This function is a placeholder in case you want to extract it from DSN.
	return ""
}
