package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"database/sql"
	"errors"
	"example/internal/dataloaders"
	"example/internal/db"
	"example/internal/graph/model"
	"example/internal/helper"
	"example/internal/middleware"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

// Name is the resolver for the name field.
func (r *chatResolver) Name(ctx context.Context, obj *db.Chat) (*string, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	if obj.Type == db.ChatTypePrivate {
		var user db.User
		err = r.DB.NewSelect().
			Column("id", "first_name", "last_name").
			Model(&user).
			Join("LEFT JOIN chat_users ON chat_users.user_id = \"user\".id").
			Where("chat_users.user_id <> ?", currentUser.ID).
			Where("chat_users.chat_id = ?", obj.ID).
			Where("\"user\".organisation_id = ?", currentUser.OrganisationID).
			Scan(ctx)
		if err != nil {
			return nil, err
		}

		fullName := user.FirstName + " " + user.LastName

		return &fullName, nil
	}

	return &obj.Name.String, nil
}

// Users is the resolver for the users field.
func (r *chatResolver) Users(ctx context.Context, obj *db.Chat) ([]*db.User, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var users []*db.User
	err = r.DB.NewSelect().
		Model(&users).
		Join("JOIN chat_users ON chat_users.user_id = \"user\".id").
		Where("chat_users.chat_id = ?", obj.ID).
		Where("\"user\".organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Type is the resolver for the type field.
func (r *chatResolver) Type(ctx context.Context, obj *db.Chat) (db.ChatType, error) {
	//return obj.Type, nil
	// return as uppercase
	return db.ChatType(strings.ToUpper(string(obj.Type))), nil
}

// Messages is the resolver for the messages field.
func (r *chatResolver) Messages(ctx context.Context, obj *db.Chat) ([]*db.ChatMessage, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var messages []*db.ChatMessage
	err = r.DB.NewSelect().
		Model(&messages).
		Where("chat_id = ?", obj.ID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// LastMessage is the resolver for the lastMessage field.
func (r *chatResolver) LastMessage(ctx context.Context, obj *db.Chat) (*db.ChatMessage, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var chatMessage db.ChatMessage
	err = r.DB.NewSelect().
		Model(&chatMessage).
		Where("chat_id = ?", obj.ID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &chatMessage, nil
}

// DeletedAt is the resolver for the deletedAt field.
func (r *chatResolver) DeletedAt(ctx context.Context, obj *db.Chat) (*time.Time, error) {
	panic(fmt.Errorf("not implemented: DeletedAt - deletedAt"))
}

// UnreadMessageCount is the resolver for the unreadMessageCount field.
func (r *chatResolver) UnreadMessageCount(ctx context.Context, obj *db.Chat) (int, error) {
	if (rand.Int() % 100) < 70 {
		return 0, nil
	} else {
		return rand.Int() % 200, nil
	}
}

// UserCount is the resolver for the userCount field.
func (r *chatResolver) UserCount(ctx context.Context, obj *db.Chat) (int, error) {
	count, err := r.DB.NewSelect().
		Model(&db.ChatUser{}).
		Where("chat_id = ?", obj.ID).
		Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Chat is the resolver for the chat field.
func (r *chatMessageResolver) Chat(ctx context.Context, obj *db.ChatMessage) (*db.Chat, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var chat db.Chat
	err = r.DB.NewSelect().
		Model(&chat).
		Where("chat.id = ?", obj.ChatID).
		Join("LEFT JOIN chat_users ON chat_users.chat_id = chat.id").
		Where("chat.organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	return &chat, nil
}

// User is the resolver for the user field.
func (r *chatMessageResolver) User(ctx context.Context, obj *db.ChatMessage) (*db.User, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	return dataloaders.GetUser(ctx, obj.UserID, currentUser)
}

// IsEdited is the resolver for the isEdited field.
func (r *chatMessageResolver) IsEdited(ctx context.Context, obj *db.ChatMessage) (bool, error) {
	if obj.UpdatedAt.IsZero() {
		return false, nil
	}

	return true, nil
}

// Chat is the resolver for the chat field.
func (r *chatUserResolver) Chat(ctx context.Context, obj *db.ChatUser) (*db.Chat, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var chat db.Chat
	err = r.DB.NewSelect().
		Model(&chat).
		Where("id = ?", obj.ChatID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// User is the resolver for the user field.
func (r *chatUserResolver) User(ctx context.Context, obj *db.ChatUser) (*db.User, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	return dataloaders.GetUser(ctx, obj.UserID, currentUser)
}

// CreatedAt is the resolver for the createdAt field.
func (r *chatUserResolver) CreatedAt(ctx context.Context, obj *db.ChatUser) (*time.Time, error) {
	panic(fmt.Errorf("not implemented: CreatedAt - createdAt"))
}

// DeletedAt is the resolver for the deletedAt field.
func (r *chatUserResolver) DeletedAt(ctx context.Context, obj *db.ChatUser) (*time.Time, error) {
	panic(fmt.Errorf("not implemented: DeletedAt - deletedAt"))
}

// CreateChat is the resolver for the createChat field.
func (r *mutationResolver) CreateChat(ctx context.Context, input model.CreateChatInput) (*db.Chat, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var chat db.Chat
	chat.Name = sql.NullString{
		String: *input.Name,
		Valid:  true,
	}

	chat.Type = db.ChatTypeGroup
	chat.OrganisationID = currentUser.OrganisationID

	err = r.DB.NewInsert().
		Model(&chat).
		Returning("*").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var chatUser db.ChatUser
	chatUser.ChatID = chat.ID
	chatUser.UserID = currentUser.ID
	chatUser.OrganisationID = currentUser.OrganisationID

	err = r.DB.NewInsert().
		Model(&chatUser).
		Returning("*").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// DeleteChat is the resolver for the deleteChat field.
func (r *mutationResolver) DeleteChat(ctx context.Context, input model.DeleteChatInput) (*db.Chat, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var chat db.Chat
	err = r.DB.NewSelect().
		Model(&chat).
		Where("chat.id = ?", input.ID).
		Join("INNER JOIN chat_users ON chat_users.chat_id = chat.id AND chat_users.user_id = ?", currentUser.ID).
		Where("chat.organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	err = r.DB.NewDelete().
		Model(&chat).
		Where("id = ?", input.ID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Returning("*").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// CreatePrivatChat is the resolver for the createPrivatChat field.
func (r *mutationResolver) CreatePrivatChat(ctx context.Context, userID string) (*db.Chat, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	// check if a privat chat already exists
	var existingChat db.Chat
	err = r.DB.NewSelect().
		Model(&existingChat).
		Where("chat.type = ?", db.ChatTypePrivate).
		Join("INNER JOIN chat_users ON chat_users.chat_id = chat.id AND chat_users.user_id = ?", currentUser.ID).
		Join("INNER JOIN chat_users AS chat_users2 ON chat_users2.chat_id = chat.id AND chat_users2.user_id = ?", userID).
		Where("chat.organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err == nil {
		return &existingChat, nil
	}

	t, err := r.DB.BeginTx(ctx, nil)

	var user db.User
	err = t.NewSelect().
		Model(&user).
		Where("id = ?", userID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		_ = t.Rollback()
		return nil, errors.New("user not found")
	}

	var chat db.Chat
	chat.Type = db.ChatTypePrivate
	chat.OrganisationID = currentUser.OrganisationID

	err = t.NewInsert().
		Model(&chat).
		Returning("*").
		Scan(ctx)
	if err != nil {
		_ = t.Rollback()
		return nil, errors.New("unable to create chat")
	}

	var chatUser1 db.ChatUser
	chatUser1.ChatID = chat.ID
	chatUser1.UserID = currentUser.ID
	chatUser1.OrganisationID = currentUser.OrganisationID

	err = t.NewInsert().
		Model(&chatUser1).
		Returning("*").
		Scan(ctx)
	if err != nil {
		_ = t.Rollback()
		return nil, errors.New("unable to create chat")
	}

	var user2 db.User
	err = t.NewSelect().
		Model(&user2).
		Where("id = ?", userID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		_ = t.Rollback()
		return nil, errors.New("other user not found")
	}

	var chatUser2 db.ChatUser
	chatUser2.ChatID = chat.ID
	chatUser2.UserID = user2.ID
	chatUser2.OrganisationID = currentUser.OrganisationID

	err = t.NewInsert().
		Model(&chatUser2).
		Returning("*").
		Scan(ctx)
	if err != nil {
		_ = t.Rollback()
		return nil, errors.New("unable to create chat")
	}

	err = t.Commit()
	if err != nil {
		return nil, errors.New("unable to create chat - commit failed")
	}

	return &chat, nil
}

// AddUserToChat is the resolver for the addUserToChat field.
func (r *mutationResolver) AddUserToChat(ctx context.Context, input model.AddUserToChatInput) (*db.ChatUser, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	if input.UserID == currentUser.ID {
		return nil, errors.New("you cannot add yourself to a chat")
	}

	var chat db.Chat
	err = r.DB.NewSelect().
		Model(&chat).
		Where("chat.id = ?", input.ChatID).
		Join("INNER JOIN chat_users ON chat_users.chat_id = chat.id AND chat_users.user_id = ?", currentUser.ID).
		Where("chat.organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	var user db.User
	err = r.DB.NewSelect().
		Model(&user).
		Where("id = ?", input.UserID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	var chatUser3 db.ChatUser
	chatUser3.ChatID = input.ChatID
	chatUser3.UserID = input.UserID
	chatUser3.OrganisationID = currentUser.OrganisationID

	err = r.DB.NewInsert().
		Model(&chatUser3).
		Returning("*").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &chatUser3, nil
}

// RemoveUserFromChat is the resolver for the removeUserFromChat field.
func (r *mutationResolver) RemoveUserFromChat(ctx context.Context, input model.RemoveUserFromChatInput) (*db.ChatUser, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var chat db.Chat
	err = r.DB.NewSelect().
		Model(&chat).
		Where("chat.id = ?", input.ChatID).
		Join("INNER JOIN chat_users ON chat_users.chat_id = chat.id AND chat_users.user_id = ?", currentUser.ID).
		Where("chat.organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	var chatUser db.ChatUser
	err = r.DB.NewDelete().
		Model(&chatUser).
		Where("chat_id = ?", input.ChatID).
		Where("user_id = ?", input.UserID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Returning("*").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &chatUser, nil
}

// SendMessage is the resolver for the sendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, input model.SendMessageInput) (*db.ChatMessage, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	// check that the message is less than 4096 characters
	if len(input.Message) > 4096 {
		return nil, errors.New("message too long, max 4096 characters")
	}

	var chat db.Chat
	err = r.DB.NewSelect().
		Model(&chat).
		Where("chat.id = ?", input.ChatID).
		Join("INNER JOIN chat_users ON chat_users.chat_id = chat.id AND chat_users.user_id = ?", currentUser.ID).
		Where("chat.organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	var chatMessage db.ChatMessage
	chatMessage.ChatID = input.ChatID
	chatMessage.UserID = currentUser.ID
	chatMessage.OrganisationID = currentUser.OrganisationID
	chatMessage.Message = input.Message
	chatMessage.CreatedAt = time.Now()

	err = r.DB.NewInsert().
		Model(&chatMessage).
		Returning("*").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// background handling of sending message to other users
	go r.SubscriptionHandler.PublishMessage(&chatMessage)
	go r.ChatMessageProcessor.NewMessage(chatMessage)

	return &chatMessage, nil
}

// EditChatMessage is the resolver for the editChatMessage field.
func (r *mutationResolver) EditChatMessage(ctx context.Context, input model.EditChatMessageInput) (*db.ChatMessage, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var chatMessage db.ChatMessage
	err = r.DB.NewSelect().
		Model(&chatMessage).
		Where("id = ?", input.ID).
		Where("user_id = ?", currentUser.ID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("message not found")
	}

	// check that the message is less than 4096 characters
	if len(input.Message) > 4096 {
		return nil, errors.New("message too long, max 4096 characters")
	}

	chatMessage.Message = input.Message
	chatMessage.UpdatedAt = bun.NullTime{Time: time.Now()}

	err = r.DB.NewUpdate().
		Model(&chatMessage).
		Where("id = ?", input.ID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Returning("*").
		Scan(ctx)
	if err != nil {
		return nil, errors.New("unable to update message")
	}

	go r.SubscriptionHandler.PublishMessage(&chatMessage)

	return &chatMessage, nil
}

// UpdateChat is the resolver for the updateChat field.
func (r *mutationResolver) UpdateChat(ctx context.Context, input model.UpdateChatInput) (*db.Chat, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var chat db.Chat
	err = r.DB.NewSelect().
		Model(&chat).
		Where("chat.id = ?", input.ID).
		Join("INNER JOIN chat_users ON chat_users.chat_id = chat.id AND chat_users.user_id = ?", currentUser.ID).
		Where("chat.organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	chat.Name = sql.NullString{
		String: *input.Name,
		Valid:  true,
	}

	err = r.DB.NewUpdate().
		Model(&chat).
		Where("id = ?", input.ID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Returning("*").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// Chat is the resolver for the chat field.
func (r *queryResolver) Chat(ctx context.Context, id string) (*db.Chat, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	var chat db.Chat
	err = r.DB.NewSelect().
		Model(&chat).
		Where("chat.id = ?", id).
		Join("INNER JOIN chat_users ON chat_users.chat_id = chat.id AND chat_users.user_id = ?", currentUser.ID).
		Where("chat.organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	return &chat, nil
}

// Chats is the resolver for the chats field.
func (r *queryResolver) Chats(ctx context.Context, limit *int, offset *int) (*model.ChatConnection, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	pageLimit, pageOffset := helper.SetPageLimits(limit, offset)

	var count int
	var chats []*ChatWithLastMessage
	count, err = r.DB.
		NewSelect().
		Model(&chats).
		ColumnExpr("chat.*").
		ColumnExpr("MAX(cm.created_at) as last_message_at").
		Join("LEFT JOIN chat_users ON chat_users.chat_id = chat.id").
		Join("LEFT JOIN chat_messages cm ON cm.chat_id = chat.id").
		TableExpr("chats AS chat").
		Where("chat_users.user_id = ?", currentUser.ID).
		Where("chat.organisation_id = ?", currentUser.OrganisationID).
		Where("chat.deleted_at IS NULL").
		Limit(pageLimit).
		Offset(pageOffset).
		Group("chat.id").
		Order("last_message_at DESC").
		ScanAndCount(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return &model.ChatConnection{}, nil
	}
	if err != nil {
		return nil, err
	}

	pageInfo, err := helper.CreatePageInfo(pageOffset, pageLimit, count)
	if err != nil {
		return nil, err
	}

	var returnChats []*db.Chat
	for _, chat := range chats {
		returnChats = append(returnChats, chat.Chat)
	}

	return &model.ChatConnection{
		Edges:      returnChats,
		TotalCount: count,
		PageInfo:   pageInfo,
	}, nil
}

// MessageAdded is the resolver for the messageAdded field.
func (r *subscriptionResolver) MessageAdded(ctx context.Context, chatID string) (<-chan *db.ChatMessage, error) {
	currentUser, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, nil
	}

	// check if user is participant of chat
	var chatUser db.ChatUser
	err = r.DB.NewSelect().
		Model(&chatUser).
		Where("chat_id = ?", chatID).
		Where("user_id = ?", currentUser.ID).
		Where("organisation_id = ?", currentUser.OrganisationID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	_ = r.SubscriptionHandler.AddChatChannel(chatID)

	channel := r.SubscriptionHandler.ChatRooms[chatID][len(r.SubscriptionHandler.ChatRooms[chatID])-1]

	go func() {
		<-ctx.Done()
		_ = r.SubscriptionHandler.RemoveChatChannel(chatID, channel)
	}()

	return channel, nil
}

// Chat returns ChatResolver implementation.
func (r *Resolver) Chat() ChatResolver { return &chatResolver{r} }

// ChatMessage returns ChatMessageResolver implementation.
func (r *Resolver) ChatMessage() ChatMessageResolver { return &chatMessageResolver{r} }

// ChatUser returns ChatUserResolver implementation.
func (r *Resolver) ChatUser() ChatUserResolver { return &chatUserResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type chatResolver struct{ *Resolver }
type chatMessageResolver struct{ *Resolver }
type chatUserResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

type ChatWithLastMessage struct {
	bun.BaseModel `bun:"table:chats"`

	*db.Chat
	LastMessage time.Time `bun:"last_message_at"`
}
