package repository

import (
	"Art-Design-Backend/internal/repository/db"
)

type ConversationRepo struct {
	*db.ConversationDB
	*db.MessageDB
}
