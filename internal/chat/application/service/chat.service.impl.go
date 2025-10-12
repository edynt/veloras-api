package service

import (
	"context"
	"math"
	"net/http"

	"github.com/edynt/chogiare/veloras-api/internal/chat/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/chat/domain/model/entity"
	"github.com/edynt/chogiare/veloras-api/internal/chat/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
)

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (s *chatService) CreateConversation(ctx context.Context, req *dto.CreateConversationAppDTO, userID int32) (*dto.ConversationAppDTO, error) {
	conversation, err := s.chatRepo.CreateConversation(ctx, repository.CreateConversationParams{
		Type:  req.Type,
		Title: req.Title,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to create conversation", err)
	}

	// Add the creator as a participant
	_, err = s.chatRepo.AddParticipant(ctx, repository.AddParticipantParams{
		ConversationID: conversation.ID,
		UserID:         userID,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to add participant", err)
	}

	// Add other participants
	for _, participantUserID := range req.UserIDs {
		if participantUserID != userID {
			_, err = s.chatRepo.AddParticipant(ctx, repository.AddParticipantParams{
				ConversationID: conversation.ID,
				UserID:         participantUserID,
			})
			if err != nil {
				return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to add participant", err)
			}
		}
	}

	return s.convertToConversationAppDTO(conversation), nil
}

func (s *chatService) GetConversation(ctx context.Context, id string) (*dto.ConversationAppDTO, error) {
	conversation, err := s.chatRepo.GetConversation(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Conversation not found", err)
	}

	return s.convertToConversationAppDTOWithDetails(conversation), nil
}

func (s *chatService) ListUserConversations(ctx context.Context, userID int32, page, pageSize int32) (*dto.ConversationListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	conversations, err := s.chatRepo.ListUserConversations(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list conversations", err)
	}

	total := int64(len(conversations))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ConversationListAppDTO{
		Conversations: s.convertToConversationAppDTOsWithDetails(conversations),
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}, nil
}

func (s *chatService) UpdateConversation(ctx context.Context, id string, req *dto.UpdateConversationAppDTO) (*dto.ConversationAppDTO, error) {
	// Get existing conversation to preserve unchanged fields
	existingConversation, err := s.chatRepo.GetConversation(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Conversation not found", err)
	}

	updateParams := repository.UpdateConversationParams{
		ID:    id,
		Type:  existingConversation.Type,
		Title: existingConversation.Title,
	}

	// Update only provided fields
	if req.Type != nil {
		updateParams.Type = *req.Type
	}
	if req.Title != nil {
		updateParams.Title = req.Title
	}

	conversation, err := s.chatRepo.UpdateConversation(ctx, updateParams)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to update conversation", err)
	}

	return s.convertToConversationAppDTO(conversation), nil
}

func (s *chatService) DeleteConversation(ctx context.Context, id string) error {
	err := s.chatRepo.DeleteConversation(ctx, id)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to delete conversation", err)
	}
	return nil
}

func (s *chatService) AddParticipant(ctx context.Context, conversationID string, userID int32, role string) (*dto.ConversationParticipantAppDTO, error) {
	participant, err := s.chatRepo.AddParticipant(ctx, repository.AddParticipantParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to add participant", err)
	}

	return s.convertToConversationParticipantAppDTO(participant), nil
}

func (s *chatService) RemoveParticipant(ctx context.Context, conversationID string, userID int32) error {
	err := s.chatRepo.RemoveParticipant(ctx, conversationID, userID)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to remove participant", err)
	}
	return nil
}

func (s *chatService) CreateChatMessage(ctx context.Context, req *dto.CreateChatMessageAppDTO, userID int32) (*dto.ChatMessageAppDTO, error) {
	message, err := s.chatRepo.CreateChatMessage(ctx, repository.CreateChatMessageParams{
		ConversationID: req.ConversationID,
		SenderID:       userID,
		MessageType:    req.MessageType,
		Content:        req.Content,
	})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to create message", err)
	}

	return s.convertToChatMessageAppDTO(message), nil
}

func (s *chatService) GetChatMessage(ctx context.Context, id string) (*dto.ChatMessageAppDTO, error) {
	message, err := s.chatRepo.GetChatMessage(ctx, id)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, "Message not found", err)
	}

	return s.convertToChatMessageAppDTOWithDetails(message), nil
}

func (s *chatService) ListConversationMessages(ctx context.Context, conversationID string, page, pageSize int32) (*dto.ChatMessageListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	messages, err := s.chatRepo.ListConversationMessages(ctx, conversationID, pageSize, offset)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to list messages", err)
	}

	total := int64(len(messages))
	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ChatMessageListAppDTO{
		Messages:   s.convertToChatMessageAppDTOsWithDetails(messages),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *chatService) MarkMessageAsRead(ctx context.Context, messageID string, userID int32) error {
	err := s.chatRepo.MarkMessageAsRead(ctx, messageID, userID)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to mark message as read", err)
	}
	return nil
}

func (s *chatService) MarkConversationAsRead(ctx context.Context, conversationID string, userID int32) error {
	err := s.chatRepo.MarkConversationAsRead(ctx, conversationID, userID)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to mark conversation as read", err)
	}
	return nil
}

func (s *chatService) DeleteChatMessage(ctx context.Context, id string) error {
	err := s.chatRepo.DeleteChatMessage(ctx, id)
	if err != nil {
		return response.NewAPIError(http.StatusInternalServerError, "Failed to delete message", err)
	}
	return nil
}

func (s *chatService) GetChatStats(ctx context.Context, userID int32) (*dto.ChatStatsAppDTO, error) {
	stats, err := s.chatRepo.GetChatStats(ctx, userID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to get chat stats", err)
	}

	return &dto.ChatStatsAppDTO{
		TotalConversations:  stats.TotalConversations,
		TotalMessages:       stats.TotalMessages,
		UnreadMessages:      stats.UnreadMessages,
		ActiveConversations: stats.ActiveConversations,
	}, nil
}

// Conversion methods
func (s *chatService) convertToConversationAppDTO(conversation *entity.Conversation) *dto.ConversationAppDTO {
	return &dto.ConversationAppDTO{
		ID:        conversation.ID,
		Type:      conversation.Type,
		Title:     conversation.Title,
		CreatedAt: conversation.CreatedAt,
		UpdatedAt: conversation.UpdatedAt,
	}
}

func (s *chatService) convertToConversationAppDTOWithDetails(conversation *entity.ConversationWithDetails) *dto.ConversationAppDTO {
	result := &dto.ConversationAppDTO{
		ID:          conversation.ID,
		Type:        conversation.Type,
		Title:       conversation.Title,
		UnreadCount: conversation.UnreadCount,
		CreatedAt:   conversation.CreatedAt,
		UpdatedAt:   conversation.UpdatedAt,
	}

	// Convert participants
	result.Participants = make([]dto.ConversationParticipantAppDTO, len(conversation.Participants))
	for i, p := range conversation.Participants {
		result.Participants[i] = dto.ConversationParticipantAppDTO{
			ID:             p.ID,
			ConversationID: p.ConversationID,
			UserID:         p.UserID,
			Role:           p.Role,
			JoinedAt:       p.JoinedAt,
		}
	}

	// Convert last message
	if conversation.LastMessage != nil {
		result.LastMessage = &dto.ChatMessageAppDTO{
			ID:             conversation.LastMessage.ID,
			ConversationID: conversation.LastMessage.ConversationID,
			SenderID:       conversation.LastMessage.SenderID,
			MessageType:    conversation.LastMessage.MessageType,
			Content:        conversation.LastMessage.Content,
			IsRead:         conversation.LastMessage.IsRead,
			CreatedAt:      conversation.LastMessage.CreatedAt,
			UpdatedAt:      conversation.LastMessage.UpdatedAt,
		}
	}

	return result
}

func (s *chatService) convertToConversationAppDTOsWithDetails(conversations []entity.ConversationWithDetails) []dto.ConversationAppDTO {
	result := make([]dto.ConversationAppDTO, len(conversations))
	for i, conversation := range conversations {
		result[i] = *s.convertToConversationAppDTOWithDetails(&conversation)
	}
	return result
}

func (s *chatService) convertToConversationParticipantAppDTO(participant *entity.ConversationParticipant) *dto.ConversationParticipantAppDTO {
	return &dto.ConversationParticipantAppDTO{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		UserID:         participant.UserID,
		Role:           participant.Role,
		JoinedAt:       participant.JoinedAt,
	}
}

func (s *chatService) convertToChatMessageAppDTO(message *entity.ChatMessage) *dto.ChatMessageAppDTO {
	return &dto.ChatMessageAppDTO{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		MessageType:    message.MessageType,
		Content:        message.Content,
		IsRead:         message.IsRead,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,
	}
}

func (s *chatService) convertToChatMessageAppDTOWithDetails(message *entity.ChatMessageWithDetails) *dto.ChatMessageAppDTO {
	return &dto.ChatMessageAppDTO{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		MessageType:    message.MessageType,
		Content:        message.Content,
		IsRead:         message.IsRead,
		SenderName:     message.SenderName,
		SenderEmail:    message.SenderEmail,
		SenderAvatar:   message.SenderAvatar,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,
	}
}

func (s *chatService) convertToChatMessageAppDTOsWithDetails(messages []entity.ChatMessageWithDetails) []dto.ChatMessageAppDTO {
	result := make([]dto.ChatMessageAppDTO, len(messages))
	for i, message := range messages {
		result[i] = *s.convertToChatMessageAppDTOWithDetails(&message)
	}
	return result
}
