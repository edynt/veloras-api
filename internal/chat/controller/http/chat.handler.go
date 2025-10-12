package http

import (
	"net/http"
	"strconv"

	"github.com/edynt/chogiare/veloras-api/internal/chat/application/service"
	appDto "github.com/edynt/chogiare/veloras-api/internal/chat/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/chat/controller/dto"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ChatHandler struct {
	service service.ChatService
}

func NewChatHandler(service service.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

// CreateConversation
// @Summary Create a new conversation
// @Description Create a new conversation with participants
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation body dto.CreateConversationRequest true "Conversation data"
// @Success 201 {object} dto.ConversationResponse "Conversation created successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /chat/conversations [post]
func (h *ChatHandler) CreateConversation(ctx *gin.Context) (res interface{}, err error) {
	var req dto.CreateConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Validate the request
	validation, exists := ctx.Get("validation")
	if !exists {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Validation middleware not found", "")
	}
	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		return nil, apiErr
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	// Convert to application DTO
	appReq := &appDto.CreateConversationAppDTO{
		Type:    req.Type,
		Title:   req.Title,
		UserIDs: req.UserIDs,
	}

	conversation, err := h.service.CreateConversation(ctx, appReq, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToConversationResponse(conversation), nil
}

// GetConversation
// @Summary Get conversation by ID
// @Description Get conversation details by ID
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Success 200 {object} dto.ConversationResponse "Conversation details"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Conversation not found"
// @Router /chat/conversations/{id} [get]
func (h *ChatHandler) GetConversation(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Conversation ID is required")
	}

	conversation, err := h.service.GetConversation(ctx, id)
	if err != nil {
		return nil, err
	}

	return h.convertToConversationResponse(conversation), nil
}

// ListUserConversations
// @Summary List user's conversations
// @Description Get a paginated list of conversations for the current user
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.ConversationListResponse "List of conversations"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /chat/conversations [get]
func (h *ChatHandler) ListUserConversations(ctx *gin.Context) (res interface{}, err error) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	conversations, err := h.service.ListUserConversations(ctx, userID, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToConversationListResponse(conversations), nil
}

// UpdateConversation
// @Summary Update conversation
// @Description Update conversation details
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Param conversation body dto.UpdateConversationRequest true "Conversation update data"
// @Success 200 {object} dto.ConversationResponse "Conversation updated successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Conversation not found"
// @Router /chat/conversations/{id} [put]
func (h *ChatHandler) UpdateConversation(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Conversation ID is required")
	}

	var req dto.UpdateConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Convert to application DTO
	appReq := &appDto.UpdateConversationAppDTO{
		Type:  req.Type,
		Title: req.Title,
	}

	conversation, err := h.service.UpdateConversation(ctx, id, appReq)
	if err != nil {
		return nil, err
	}

	return h.convertToConversationResponse(conversation), nil
}

// DeleteConversation
// @Summary Delete conversation
// @Description Delete a conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Success 204 "Conversation deleted successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Conversation not found"
// @Router /chat/conversations/{id} [delete]
func (h *ChatHandler) DeleteConversation(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Conversation ID is required")
	}

	err = h.service.DeleteConversation(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// AddParticipant
// @Summary Add participant to conversation
// @Description Add a user as a participant to a conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Param user_id path string true "User ID"
// @Param role body string true "Participant role"
// @Success 200 {object} dto.ConversationParticipantResponse "Participant added successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Conversation not found"
// @Router /chat/conversations/{id}/participants/{user_id} [post]
func (h *ChatHandler) AddParticipant(ctx *gin.Context) (res interface{}, err error) {
	conversationID := ctx.Param("id")
	if conversationID == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Conversation ID is required")
	}

	userIDStr := ctx.Param("user_id")
	if userIDStr == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "User ID is required")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Invalid user ID")
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	participant, err := h.service.AddParticipant(ctx, conversationID, int32(userID), req.Role)
	if err != nil {
		return nil, err
	}

	return h.convertToConversationParticipantResponse(participant), nil
}

// RemoveParticipant
// @Summary Remove participant from conversation
// @Description Remove a user from a conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Param user_id path string true "User ID"
// @Success 204 "Participant removed successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Conversation not found"
// @Router /chat/conversations/{id}/participants/{user_id} [delete]
func (h *ChatHandler) RemoveParticipant(ctx *gin.Context) (res interface{}, err error) {
	conversationID := ctx.Param("id")
	if conversationID == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Conversation ID is required")
	}

	userIDStr := ctx.Param("user_id")
	if userIDStr == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "User ID is required")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Invalid user ID")
	}

	err = h.service.RemoveParticipant(ctx, conversationID, int32(userID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// CreateChatMessage
// @Summary Create a new chat message
// @Description Create a new message in a conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body dto.CreateChatMessageRequest true "Message data"
// @Success 201 {object} dto.ChatMessageResponse "Message created successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /chat/messages [post]
func (h *ChatHandler) CreateChatMessage(ctx *gin.Context) (res interface{}, err error) {
	var req dto.CreateChatMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Validate the request
	validation, exists := ctx.Get("validation")
	if !exists {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Validation middleware not found", "")
	}
	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		return nil, apiErr
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	// Convert to application DTO
	appReq := &appDto.CreateChatMessageAppDTO{
		ConversationID: req.ConversationID,
		MessageType:    req.MessageType,
		Content:        req.Content,
	}

	message, err := h.service.CreateChatMessage(ctx, appReq, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToChatMessageResponse(message), nil
}

// GetChatMessage
// @Summary Get chat message by ID
// @Description Get chat message details by ID
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Success 200 {object} dto.ChatMessageResponse "Message details"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Message not found"
// @Router /chat/messages/{id} [get]
func (h *ChatHandler) GetChatMessage(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Message ID is required")
	}

	message, err := h.service.GetChatMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return h.convertToChatMessageResponse(message), nil
}

// ListConversationMessages
// @Summary List conversation messages
// @Description Get a paginated list of messages in a conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 50, max: 100)"
// @Success 200 {object} dto.ChatMessageListResponse "List of messages"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /chat/conversations/{id}/messages [get]
func (h *ChatHandler) ListConversationMessages(ctx *gin.Context) (res interface{}, err error) {
	conversationID := ctx.Param("id")
	if conversationID == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Conversation ID is required")
	}

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "50")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize < 1 {
		pageSize = 50
	}

	if pageSize > 100 {
		pageSize = 100
	}

	messages, err := h.service.ListConversationMessages(ctx, conversationID, int32(page), int32(pageSize))
	if err != nil {
		return nil, err
	}

	return h.convertToChatMessageListResponse(messages), nil
}

// MarkMessageAsRead
// @Summary Mark message as read
// @Description Mark a specific message as read
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Success 200 "Message marked as read"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Message not found"
// @Router /chat/messages/{id}/read [post]
func (h *ChatHandler) MarkMessageAsRead(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Message ID is required")
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	err = h.service.MarkMessageAsRead(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// MarkConversationAsRead
// @Summary Mark conversation as read
// @Description Mark all messages in a conversation as read
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Success 200 "Conversation marked as read"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Conversation not found"
// @Router /chat/conversations/{id}/read [post]
func (h *ChatHandler) MarkConversationAsRead(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Conversation ID is required")
	}

	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	err = h.service.MarkConversationAsRead(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteChatMessage
// @Summary Delete chat message
// @Description Delete a chat message
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Success 204 "Message deleted successfully"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "Message not found"
// @Router /chat/messages/{id} [delete]
func (h *ChatHandler) DeleteChatMessage(ctx *gin.Context) (res interface{}, err error) {
	id := ctx.Param("id")
	if id == "" {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request", "Message ID is required")
	}

	err = h.service.DeleteChatMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetChatStats
// @Summary Get chat statistics
// @Description Get chat statistics for the current user
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ChatStatsResponse "Chat statistics"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /chat/stats [get]
func (h *ChatHandler) GetChatStats(ctx *gin.Context) (res interface{}, err error) {
	// Get user ID from context (this would be set by auth middleware)
	userID := int32(1) // Placeholder - should come from JWT token

	stats, err := h.service.GetChatStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return h.convertToChatStatsResponse(stats), nil
}

// Conversion methods
func (h *ChatHandler) convertToConversationResponse(conversation *appDto.ConversationAppDTO) *dto.ConversationResponse {
	result := &dto.ConversationResponse{
		ID:          conversation.ID,
		Type:        conversation.Type,
		Title:       conversation.Title,
		UnreadCount: conversation.UnreadCount,
		CreatedAt:   conversation.CreatedAt,
		UpdatedAt:   conversation.UpdatedAt,
	}

	// Convert participants
	result.Participants = make([]dto.ConversationParticipantResponse, len(conversation.Participants))
	for i, p := range conversation.Participants {
		result.Participants[i] = dto.ConversationParticipantResponse{
			ID:             p.ID,
			ConversationID: p.ConversationID,
			UserID:         p.UserID,
			Role:           p.Role,
			JoinedAt:       p.JoinedAt,
		}
	}

	// Convert last message
	if conversation.LastMessage != nil {
		result.LastMessage = &dto.ChatMessageResponse{
			ID:             conversation.LastMessage.ID,
			ConversationID: conversation.LastMessage.ConversationID,
			SenderID:       conversation.LastMessage.SenderID,
			MessageType:    conversation.LastMessage.MessageType,
			Content:        conversation.LastMessage.Content,
			IsRead:         conversation.LastMessage.IsRead,
			SenderName:     conversation.LastMessage.SenderName,
			SenderEmail:    conversation.LastMessage.SenderEmail,
			SenderAvatar:   conversation.LastMessage.SenderAvatar,
			CreatedAt:      conversation.LastMessage.CreatedAt,
			UpdatedAt:      conversation.LastMessage.UpdatedAt,
		}
	}

	return result
}

func (h *ChatHandler) convertToConversationListResponse(conversations *appDto.ConversationListAppDTO) *dto.ConversationListResponse {
	result := &dto.ConversationListResponse{
		Total:      conversations.Total,
		Page:       conversations.Page,
		PageSize:   conversations.PageSize,
		TotalPages: conversations.TotalPages,
	}

	result.Conversations = make([]dto.ConversationResponse, len(conversations.Conversations))
	for i, conversation := range conversations.Conversations {
		result.Conversations[i] = *h.convertToConversationResponse(&conversation)
	}

	return result
}

func (h *ChatHandler) convertToConversationParticipantResponse(participant *appDto.ConversationParticipantAppDTO) *dto.ConversationParticipantResponse {
	return &dto.ConversationParticipantResponse{
		ID:             participant.ID,
		ConversationID: participant.ConversationID,
		UserID:         participant.UserID,
		Role:           participant.Role,
		JoinedAt:       participant.JoinedAt,
	}
}

func (h *ChatHandler) convertToChatMessageResponse(message *appDto.ChatMessageAppDTO) *dto.ChatMessageResponse {
	return &dto.ChatMessageResponse{
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

func (h *ChatHandler) convertToChatMessageListResponse(messages *appDto.ChatMessageListAppDTO) *dto.ChatMessageListResponse {
	result := &dto.ChatMessageListResponse{
		Total:      messages.Total,
		Page:       messages.Page,
		PageSize:   messages.PageSize,
		TotalPages: messages.TotalPages,
	}

	result.Messages = make([]dto.ChatMessageResponse, len(messages.Messages))
	for i, message := range messages.Messages {
		result.Messages[i] = *h.convertToChatMessageResponse(&message)
	}

	return result
}

func (h *ChatHandler) convertToChatStatsResponse(stats *appDto.ChatStatsAppDTO) *dto.ChatStatsResponse {
	return &dto.ChatStatsResponse{
		TotalConversations:  stats.TotalConversations,
		TotalMessages:       stats.TotalMessages,
		UnreadMessages:      stats.UnreadMessages,
		ActiveConversations: stats.ActiveConversations,
	}
}
