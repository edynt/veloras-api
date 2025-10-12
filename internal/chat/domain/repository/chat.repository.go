package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/chat/domain/model/entity"
)

type ChatRepository interface {
	// Conversation operations
	CreateConversation(ctx context.Context, req CreateConversationParams) (*entity.Conversation, error)
	GetConversation(ctx context.Context, id string) (*entity.ConversationWithDetails, error)
	ListUserConversations(ctx context.Context, userID int32, limit, offset int32) ([]entity.ConversationWithDetails, error)
	UpdateConversation(ctx context.Context, req UpdateConversationParams) (*entity.Conversation, error)
	DeleteConversation(ctx context.Context, id string) error

	// Conversation participant operations
	AddParticipant(ctx context.Context, req AddParticipantParams) (*entity.ConversationParticipant, error)
	RemoveParticipant(ctx context.Context, conversationID string, userID int32) error
	GetConversationParticipants(ctx context.Context, conversationID string) ([]entity.ConversationParticipant, error)

	// Chat message operations
	CreateChatMessage(ctx context.Context, req CreateChatMessageParams) (*entity.ChatMessage, error)
	GetChatMessage(ctx context.Context, id string) (*entity.ChatMessageWithDetails, error)
	ListConversationMessages(ctx context.Context, conversationID string, limit, offset int32) ([]entity.ChatMessageWithDetails, error)
	MarkMessageAsRead(ctx context.Context, messageID string, userID int32) error
	MarkConversationAsRead(ctx context.Context, conversationID string, userID int32) error
	DeleteChatMessage(ctx context.Context, id string) error

	// Statistics
	GetChatStats(ctx context.Context, userID int32) (*entity.ChatStats, error)
}

type CreateConversationParams struct {
	Type  string
	Title *string
}

type UpdateConversationParams struct {
	ID    string
	Type  string
	Title *string
}

type AddParticipantParams struct {
	ConversationID string
	UserID         int32
	Role           string
}

type CreateChatMessageParams struct {
	ConversationID string
	SenderID       int32
	MessageType    string
	Content        string
}
