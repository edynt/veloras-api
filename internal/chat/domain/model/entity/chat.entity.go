package entity

import (
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type Conversation struct {
	ID        string
	Type      string
	Title     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConversationParticipant struct {
	ID             string
	ConversationID string
	UserID         int32
	Role           string
	JoinedAt       time.Time
}

type ChatMessage struct {
	ID             string
	ConversationID string
	SenderID       int32
	MessageType    string
	Content        string
	IsRead         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ConversationWithDetails struct {
	Conversation
	Participants []ConversationParticipant
	LastMessage  *ChatMessage
	UnreadCount  int32
}

type ChatMessageWithDetails struct {
	ChatMessage
	SenderName   *string
	SenderEmail  *string
	SenderAvatar *string
}

type ChatStats struct {
	TotalConversations  int64
	TotalMessages       int64
	UnreadMessages      int64
	ActiveConversations int64
}

// Convert from SQLC model to domain entity
func FromSQLCConversation(c gen.Conversation) Conversation {
	return Conversation{
		ID:        c.ID.String(),
		Type:      "direct", // Default type, not in generated model
		Title:     nil,      // Not in generated model
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
}

// Convert from SQLC model to domain entity
func FromSQLCConversationParticipant(cp gen.ConversationParticipant) ConversationParticipant {
	return ConversationParticipant{
		ID:             cp.ID.String(),
		ConversationID: cp.ConversationID.String(),
		UserID:         cp.UserID,
		Role:           "participant", // Default role, not in generated model
		JoinedAt:       cp.JoinedAt.Time,
	}
}

// Convert from SQLC model to domain entity
func FromSQLCChatMessage(cm gen.ChatMessage) ChatMessage {
	return ChatMessage{
		ID:             cm.ID.String(),
		ConversationID: cm.ConversationID.String(),
		SenderID:       cm.SenderID,
		MessageType:    cm.Type, // Using Type field from generated model
		Content:        cm.Content,
		// Attachments field not in domain entity
		IsRead:    cm.IsRead.Bool,
		CreatedAt: cm.CreatedAt.Time,
		UpdatedAt: cm.CreatedAt.Time, // UpdatedAt not in generated model, using CreatedAt
	}
}

// Convert from SQLC model with details to domain entity
func FromSQLCConversationWithDetails(c gen.GetConversationWithDetailsRow) ConversationWithDetails {
	return ConversationWithDetails{
		Conversation: Conversation{
			ID:        c.ID.String(),
			Type:      "direct", // Default type
			Title:     nil,      // Not in generated model
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
		},
		Participants: []ConversationParticipant{}, // Would need to be populated separately
		LastMessage:  nil,                         // Would need to be populated separately
		UnreadCount:  0,                           // Would need to be populated separately
	}
}

// Convert from SQLC model with details to domain entity
func FromSQLCChatMessageWithDetails(cm gen.GetChatMessageWithDetailsRow) ChatMessageWithDetails {
	return ChatMessageWithDetails{
		ChatMessage: ChatMessage{
			ID:             cm.ID.String(),
			ConversationID: cm.ConversationID.String(),
			SenderID:       cm.SenderID,
			MessageType:    cm.Type,
			Content:        cm.Content,
			IsRead:         cm.IsRead.Bool,
			CreatedAt:      cm.CreatedAt.Time,
			UpdatedAt:      cm.CreatedAt.Time,
		},
		SenderName:   convertInterfacePtr(cm.SenderName),
		SenderEmail:  convertInterfacePtr(cm.SenderEmail),
		SenderAvatar: nil, // Not in generated model
	}
}

// Helper function to convert pgtype.Text to *string
func convertTextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

// Helper function to convert interface{} to *string
func convertInterfacePtr(i interface{}) *string {
	if i == nil {
		return nil
	}
	if str, ok := i.(string); ok {
		return &str
	}
	return nil
}

// Convert from SQLC ListUserConversationsRow to domain entity
func FromSQLCConversationWithDetailsFromList(c gen.ListUserConversationsRow) ConversationWithDetails {
	return ConversationWithDetails{
		Conversation: Conversation{
			ID:        c.ID.String(),
			Type:      "direct", // Default type, not in generated model
			Title:     nil,      // Title not in generated model
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
		},
		// ParticipantCount not available in current schema
		LastMessage: nil, // LastMessage not available in current schema
		UnreadCount: 0,   // UnreadCount not available in current schema
	}
}

// Convert from SQLC GetConversationParticipantsRow to domain entity
func FromSQLCConversationParticipantFromGet(p gen.GetConversationParticipantsRow) ConversationParticipant {
	return ConversationParticipant{
		ID:             p.ID.String(),
		ConversationID: p.ConversationID.String(),
		UserID:         p.UserID,
		Role:           "member", // Default role, not in generated model
		JoinedAt:       p.JoinedAt.Time,
	}
}

// Convert from SQLC ListConversationMessagesRow to domain entity
func FromSQLCChatMessageWithDetailsFromList(m gen.ListConversationMessagesRow) ChatMessageWithDetails {
	return ChatMessageWithDetails{
		ChatMessage: ChatMessage{
			ID:             m.ID.String(),
			ConversationID: m.ConversationID.String(),
			SenderID:       m.SenderID,
			MessageType:    m.Type,
			Content:        m.Content,
			IsRead:         m.IsRead.Bool,
			CreatedAt:      m.CreatedAt.Time,
			UpdatedAt:      m.CreatedAt.Time,
		},
		SenderName:  convertInterfacePtr(m.SenderName),
		SenderEmail: convertTextPtr(m.SenderEmail),
	}
}
