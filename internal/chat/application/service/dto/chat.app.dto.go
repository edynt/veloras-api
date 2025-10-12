package dto

import "time"

type ConversationAppDTO struct {
	ID           string                          `json:"id"`
	Type         string                          `json:"type"`
	Title        *string                         `json:"title,omitempty"`
	Participants []ConversationParticipantAppDTO `json:"participants,omitempty"`
	LastMessage  *ChatMessageAppDTO              `json:"lastMessage,omitempty"`
	UnreadCount  int32                           `json:"unreadCount"`
	CreatedAt    time.Time                       `json:"createdAt"`
	UpdatedAt    time.Time                       `json:"updatedAt"`
}

type ConversationParticipantAppDTO struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversationId"`
	UserID         int32     `json:"userId"`
	Role           string    `json:"role"`
	JoinedAt       time.Time `json:"joinedAt"`
}

type ChatMessageAppDTO struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversationId"`
	SenderID       int32     `json:"senderId"`
	MessageType    string    `json:"messageType"`
	Content        string    `json:"content"`
	IsRead         bool      `json:"isRead"`
	SenderName     *string   `json:"senderName,omitempty"`
	SenderEmail    *string   `json:"senderEmail,omitempty"`
	SenderAvatar   *string   `json:"senderAvatar,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ConversationListAppDTO struct {
	Conversations []ConversationAppDTO `json:"conversations"`
	Total         int64                `json:"total"`
	Page          int32                `json:"page"`
	PageSize      int32                `json:"pageSize"`
	TotalPages    int32                `json:"totalPages"`
}

type ChatMessageListAppDTO struct {
	Messages   []ChatMessageAppDTO `json:"messages"`
	Total      int64               `json:"total"`
	Page       int32               `json:"page"`
	PageSize   int32               `json:"pageSize"`
	TotalPages int32               `json:"totalPages"`
}

type CreateConversationAppDTO struct {
	Type    string  `json:"type" binding:"required"`
	Title   *string `json:"title,omitempty"`
	UserIDs []int32 `json:"userIds" binding:"required,min=1"`
}

type UpdateConversationAppDTO struct {
	Type  *string `json:"type,omitempty"`
	Title *string `json:"title,omitempty"`
}

type CreateChatMessageAppDTO struct {
	ConversationID string `json:"conversationId" binding:"required"`
	MessageType    string `json:"messageType" binding:"required"`
	Content        string `json:"content" binding:"required"`
}

type ChatStatsAppDTO struct {
	TotalConversations  int64 `json:"totalConversations"`
	TotalMessages       int64 `json:"totalMessages"`
	UnreadMessages      int64 `json:"unreadMessages"`
	ActiveConversations int64 `json:"activeConversations"`
}
