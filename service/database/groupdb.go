package database

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func (db *appdbimpl) CreateGroupConversation(conversationID string, memberIDs []string, name string, photo []byte, creatorID string) error {
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
        INSERT INTO conversations (id, name, type, created_at, conversationPhoto)
        VALUES (?, ?, 'group', ?, ?)
    `, conversationID, name, time.Now().Format(time.RFC3339), photo)
	if err != nil {
		return fmt.Errorf("error creating new conversation: %w", err)
	}
	for _, memberID := range memberIDs {
		role := "member"
		if memberID == creatorID {
			role = "admin"
		}
		_, err = tx.Exec(`
            INSERT INTO conversation_members (conversationId, userId, role)
            VALUES (?, ?, ?)
        `, conversationID, memberID, role)
		if err != nil {
			return fmt.Errorf("error adding member %s to conversation_members: %w", memberID, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetMyGroups(userID string) ([]Conversation, error) {
	query := `
    SELECT 
        c.id,
        c.name,
        c.conversationPhoto as photo
    FROM conversations c
    JOIN conversation_members cm ON c.id = cm.conversationId
    WHERE cm.userId = ? AND c.type = 'group'
    ORDER BY c.created_at DESC;
    `
	rows, err := db.c.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching groups: %w", err)
	}
	defer rows.Close()
	var groups []Conversation
	for rows.Next() {
		var group Conversation
		var photo sql.NullString
		err := rows.Scan(
			&group.Id,
			&group.Name,
			&photo,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning group: %w", err)
		}
		if photo.Valid {
			group.ConversationPhoto.String = base64.StdEncoding.EncodeToString([]byte(photo.String))
			group.ConversationPhoto.Valid = true
		} else {
			group.ConversationPhoto = sql.NullString{String: "", Valid: false}
		}
		groups = append(groups, group)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning groups: %w", err)
	}
	return groups, nil
}

func (db *appdbimpl) GetGroupInfo(groupID string) (Conversation, error) {
	var group Conversation
	var photo []byte
	var membersCSV sql.NullString
	err := db.c.QueryRow(`
        SELECT 
            c.id,
            c.name,
            c.conversationPhoto,
            (SELECT GROUP_CONCAT(userId) FROM conversation_members WHERE conversationId = c.id) AS members
        FROM conversations c
        WHERE c.id = ? AND c.type = 'group'`,
		groupID,
	).Scan(
		&group.Id,
		&group.Name,
		&photo,
		&membersCSV,
	)
	if err == sql.ErrNoRows {
		return Conversation{}, ErrGroupDoesNotExist
	}
	if err != nil {
		return Conversation{}, fmt.Errorf("error fetching group by ID: %w", err)
	}
	if len(photo) > 0 {
		group.ConversationPhoto = sql.NullString{
			String: base64.StdEncoding.EncodeToString(photo),
			Valid:  true,
		}
	} else {
		group.ConversationPhoto = sql.NullString{Valid: false}
	}
	if membersCSV.Valid {
		group.Members = strings.Split(membersCSV.String, ",")
	} else {
		group.Members = []string{}
	}
	return group, nil
}

func (db *appdbimpl) UpdateGroupName(groupId, newName string) error {
	res, err := db.c.Exec(`UPDATE conversations SET name=? WHERE id=?`, newName, groupId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrUserDoesNotExist
	}
	return nil
}

func (db *appdbimpl) UpdateGroupPhoto(groupID string, photo []byte) error {
	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM conversations WHERE id=?)`, groupID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrGroupDoesNotExist
	}
	_, err = db.c.Exec(`UPDATE conversations SET conversationPhoto=? WHERE id=?`, photo, groupID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) LeaveGroup(groupID, userID string) error {
	_, err := db.c.Exec(`
	DELETE FROM conversation_members WHERE conversationId = ? AND userId = ?
	`, groupID, userID)
	if err != nil {
		return fmt.Errorf("error leaving group: %w", err)
	}
	return nil
}

func (db *appdbimpl) AddUserToGroup(conversationID string, userID string) error {
	_, err := db.c.Exec(
		"INSERT INTO conversation_members (conversationId, userId) VALUES (?, ?)",
		conversationID, userID,
	)
	if err != nil {
		return fmt.Errorf("error adding user to group: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetGroupMemberDetails(groupID string) ([]GroupMember, error) {
	query := `
		SELECT u.id, u.name, u.photo, cm.role
		FROM users u
		JOIN conversation_members cm ON u.id = cm.userId
		WHERE cm.conversationId = ?
		ORDER BY CASE cm.role WHEN 'admin' THEN 0 ELSE 1 END, u.name
	`

	rows, err := db.c.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("error fetching group member details: %w", err)
	}
	defer rows.Close()

	var members []GroupMember
	for rows.Next() {
		var m GroupMember
		if err := rows.Scan(&m.Id, &m.Name, &m.Photo, &m.Role); err != nil {
			return nil, fmt.Errorf("error scanning member details: %w", err)
		}
		members = append(members, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating member rows: %w", err)
	}
	return members, nil
}

func (db *appdbimpl) IsGroupAdmin(groupID, userID string) (bool, error) {
	var role string
	err := db.c.QueryRow(
		`SELECT role FROM conversation_members WHERE conversationId = ? AND userId = ?`,
		groupID, userID,
	).Scan(&role)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return role == "admin", nil
}

func (db *appdbimpl) RemoveUserFromGroup(groupID, userID string) error {
	_, err := db.c.Exec(
		`DELETE FROM conversation_members WHERE conversationId = ? AND userId = ?`,
		groupID, userID,
	)
	return err
}

func (db *appdbimpl) AddGroupMember(groupID string, newMemberID string) error {
	// Check if user is already in group
	var exists bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM conversation_members
			WHERE conversationId = ? AND userId = ?
		)
	`, groupID, newMemberID).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking member existence: %w", err)
	}

	if exists {
		return fmt.Errorf("user is already a member of this group")
	}

	// Add new member
	_, err = db.c.Exec(`
		INSERT INTO conversation_members (conversationId, userId)
		VALUES (?, ?)
	`, groupID, newMemberID)

	if err != nil {
		return fmt.Errorf("error in adding member to group: %w", err)
	}

	return nil
}

func (db *appdbimpl) UpdateGroupInfo(groupID string, name string, photo []byte) error {
	var query string
	var result sql.Result
	var err error

	if name != "" && photo != nil {
		query = `UPDATE conversations SET name = ?, conversationPhoto = ? WHERE id = ? AND type = 'group'`
		result, err = db.c.Exec(query, name, photo, groupID)
	} else if name != "" {
		query = `UPDATE conversations SET name = ? WHERE id = ? AND type = 'group'`
		result, err = db.c.Exec(query, name, groupID)
	} else if photo != nil {
		query = `UPDATE conversations SET conversationPhoto = ? WHERE id = ? AND type = 'group'`
		result, err = db.c.Exec(query, photo, groupID)
	} else {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error updating group info: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking update result: %w", err)
	}

	if rows == 0 {
		return ErrGroupDoesNotExist
	}

	return nil
}
