package repository

import (
	"context"

	"github.com/edynt/chogiare/veloras-api/internal/chat/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/chat/domain/repository"
	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type chatRepository struct {
	db      *pgxpool.Pool
	queries *gen.Queries
}

func NewChatRepository(db *pgxpool.Pool) repository.ChatRepository {
	return &chatRepository{
		db:      db,
		queries: gen.New(db),
	}
}

// Helper functions for type conversion
func stringToUUID(s string) pgtype.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}

func stringPtrToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func boolToPgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func (r *chatRepository) CreateConversation(ctx context.Context, req repository.CreateConversationParams) (*entity.Conversation, error) {
	conversation, err := r.queries.CreateConversation(ctx)
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCConversation(conversation)
	return &result, nil
}

func (r *chatRepository) GetConversation(ctx context.Context, id string) (*entity.ConversationWithDetails, error) {
	conversation, err := r.queries.GetConversationWithDetails(ctx, stringToUUID(id))
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCConversationWithDetails(conversation)
	return &result, nil
}

func (r *chatRepository) ListUserConversations(ctx context.Context, userID int32, limit, offset int32) ([]entity.ConversationWithDetails, error) {
	conversations, err := r.queries.ListUserConversations(ctx, gen.ListUserConversationsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ConversationWithDetails, len(conversations))
	for i, c := range conversations {
		result[i] = entity.FromSQLCConversationWithDetailsFromList(c)
	}

	return result, nil
}

func (r *chatRepository) UpdateConversation(ctx context.Context, req repository.UpdateConversationParams) (*entity.Conversation, error) {
	conversation, err := r.queries.UpdateConversation(ctx, stringToUUID(req.ID))
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCConversation(conversation)
	return &result, nil
}

func (r *chatRepository) DeleteConversation(ctx context.Context, id string) error {
	return r.queries.DeleteConversation(ctx, stringToUUID(id))
}

func (r *chatRepository) AddParticipant(ctx context.Context, req repository.AddParticipantParams) (*entity.ConversationParticipant, error) {
	participant, err := r.queries.AddParticipant(ctx, gen.AddParticipantParams{
		ConversationID: stringToUUID(req.ConversationID),
		UserID:         req.UserID,
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCConversationParticipant(participant)
	return &result, nil
}

func (r *chatRepository) RemoveParticipant(ctx context.Context, conversationID string, userID int32) error {
	return r.queries.RemoveParticipant(ctx, gen.RemoveParticipantParams{
		ConversationID: stringToUUID(conversationID),
		UserID:         userID,
	})
}

func (r *chatRepository) GetConversationParticipants(ctx context.Context, conversationID string) ([]entity.ConversationParticipant, error) {
	participants, err := r.queries.GetConversationParticipants(ctx, stringToUUID(conversationID))
	if err != nil {
		return nil, err
	}

	result := make([]entity.ConversationParticipant, len(participants))
	for i, p := range participants {
		result[i] = entity.FromSQLCConversationParticipantFromGet(p)
	}

	return result, nil
}

func (r *chatRepository) CreateChatMessage(ctx context.Context, req repository.CreateChatMessageParams) (*entity.ChatMessage, error) {
	message, err := r.queries.CreateChatMessage(ctx, gen.CreateChatMessageParams{
		ConversationID: stringToUUID(req.ConversationID),
		SenderID:       req.SenderID,
		Content:        req.Content,
		Type:           req.MessageType,
		Attachments:    []string{}, // Attachments not available in current schema
	})
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCChatMessage(message)
	return &result, nil
}

func (r *chatRepository) GetChatMessage(ctx context.Context, id string) (*entity.ChatMessageWithDetails, error) {
	message, err := r.queries.GetChatMessageWithDetails(ctx, stringToUUID(id))
	if err != nil {
		return nil, err
	}

	result := entity.FromSQLCChatMessageWithDetails(message)
	return &result, nil
}

func (r *chatRepository) ListConversationMessages(ctx context.Context, conversationID string, limit, offset int32) ([]entity.ChatMessageWithDetails, error) {
	messages, err := r.queries.ListConversationMessages(ctx, gen.ListConversationMessagesParams{
		ConversationID: stringToUUID(conversationID),
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ChatMessageWithDetails, len(messages))
	for i, m := range messages {
		result[i] = entity.FromSQLCChatMessageWithDetailsFromList(m)
	}

	return result, nil
}

func (r *chatRepository) MarkMessageAsRead(ctx context.Context, messageID string, userID int32) error {
	return r.queries.MarkMessageAsRead(ctx, stringToUUID(messageID))
}

func (r *chatRepository) MarkConversationAsRead(ctx context.Context, conversationID string, userID int32) error {
	return r.queries.MarkConversationAsRead(ctx, gen.MarkConversationAsReadParams{
		ConversationID: stringToUUID(conversationID),
		SenderID:       userID,
	})
}

func (r *chatRepository) DeleteChatMessage(ctx context.Context, id string) error {
	return r.queries.DeleteChatMessage(ctx, stringToUUID(id))
}

func (r *chatRepository) GetChatStats(ctx context.Context, userID int32) (*entity.ChatStats, error) {
	stats, err := r.queries.GetChatStats(ctx)
	if err != nil {
		return nil, err
	}

	result := &entity.ChatStats{
		TotalConversations:  stats.TotalConversations,
		TotalMessages:       stats.TotalMessages,
		UnreadMessages:      stats.UnreadMessages,
		ActiveConversations: 0, // ActiveConversations not available in current schema
	}

	return result, nil
}
