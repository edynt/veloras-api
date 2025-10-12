package service

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/chat/application/service/dto"
)

type ChatService interface {
	// Conversation operations
	CreateConversation(ctx context.Context, req *dto.CreateConversationAppDTO, userID int32) (*dto.ConversationAppDTO, error)
	GetConversation(ctx context.Context, id string) (*dto.ConversationAppDTO, error)
	ListUserConversations(ctx context.Context, userID int32, page, pageSize int32) (*dto.ConversationListAppDTO, error)
	UpdateConversation(ctx context.Context, id string, req *dto.UpdateConversationAppDTO) (*dto.ConversationAppDTO, error)
	DeleteConversation(ctx context.Context, id string) error

	// Conversation participant operations
	AddParticipant(ctx context.Context, conversationID string, userID int32, role string) (*dto.ConversationParticipantAppDTO, error)
	RemoveParticipant(ctx context.Context, conversationID string, userID int32) error

	// Chat message operations
	CreateChatMessage(ctx context.Context, req *dto.CreateChatMessageAppDTO, userID int32) (*dto.ChatMessageAppDTO, error)
	GetChatMessage(ctx context.Context, id string) (*dto.ChatMessageAppDTO, error)
	ListConversationMessages(ctx context.Context, conversationID string, page, pageSize int32) (*dto.ChatMessageListAppDTO, error)
	MarkMessageAsRead(ctx context.Context, messageID string, userID int32) error
	MarkConversationAsRead(ctx context.Context, conversationID string, userID int32) error
	DeleteChatMessage(ctx context.Context, id string) error

	// Statistics
	GetChatStats(ctx context.Context, userID int32) (*dto.ChatStatsAppDTO, error)
}
