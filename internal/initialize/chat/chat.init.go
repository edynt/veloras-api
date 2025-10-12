package chat

import (
	"github.com/edynt/chogiare/veloras-api/internal/chat/application/service"
	"github.com/edynt/chogiare/veloras-api/internal/chat/controller/http"
	"github.com/edynt/chogiare/veloras-api/internal/chat/infrastructure/persistence/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitChat(db *pgxpool.Pool) *http.ChatHandler {
	// Initialize repository
	chatRepo := repository.NewChatRepository(db)

	// Initialize service
	chatService := service.NewChatService(chatRepo)

	// Initialize handler
	chatHandler := http.NewChatHandler(chatService)

	return chatHandler
}
