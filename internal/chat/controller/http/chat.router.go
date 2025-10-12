package http

import (
	"github.com/edynt/chogiare/veloras-api/internal/middleware"
	"github.com/edynt/chogiare/veloras-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterChatRoutes(r *gin.RouterGroup, handler *ChatHandler) {
	chat := r.Group("/chat")
	{
		// Protected routes
		chat.Use(middleware.AuthenMiddleware())
		{
			// Conversation operations
			chat.POST("/conversations", response.Wrap(handler.CreateConversation))
			chat.GET("/conversations", response.Wrap(handler.ListUserConversations))
			chat.GET("/conversations/:id", response.Wrap(handler.GetConversation))
			chat.PUT("/conversations/:id", response.Wrap(handler.UpdateConversation))
			chat.DELETE("/conversations/:id", response.Wrap(handler.DeleteConversation))

			// Conversation participant operations
			chat.POST("/conversations/:id/participants/:user_id", response.Wrap(handler.AddParticipant))
			chat.DELETE("/conversations/:id/participants/:user_id", response.Wrap(handler.RemoveParticipant))

			// Chat message operations
			chat.POST("/messages", response.Wrap(handler.CreateChatMessage))
			chat.GET("/messages/:id", response.Wrap(handler.GetChatMessage))
			chat.GET("/conversations/:id/messages", response.Wrap(handler.ListConversationMessages))
			chat.POST("/messages/:id/read", response.Wrap(handler.MarkMessageAsRead))
			chat.POST("/conversations/:id/read", response.Wrap(handler.MarkConversationAsRead))
			chat.DELETE("/messages/:id", response.Wrap(handler.DeleteChatMessage))

			// Statistics
			chat.GET("/stats", response.Wrap(handler.GetChatStats))
		}
	}
}
