package message

import (
	"context"
	"errors"
	"slices"
	"time"

	"go.uber.org/zap"

	"chat_game/config"
	"chat_game/log"
	"chat_game/models/postgresql"
	tgroupmessage "chat_game/models/postgresql/t_group_message"
	tmessage "chat_game/models/postgresql/t_message"
	"chat_game/services/room"
	"chat_game/utils/common"
	"chat_game/utils/opctx"
)

type MessageSendFunc func(receiver string, msg common.Msg) error

type MessageService interface {
	List(ctx context.Context, userID string, friendID string) ([]tmessage.Message, error)
	ListRoom(ctx context.Context, roomID string) ([]tgroupmessage.GroupMessage, error)
	SendToUser(ctx context.Context, receiver string, message string, send MessageSendFunc) error
	SendToRoom(ctx context.Context, roomID string, message string, send MessageSendFunc) error
}

type MessageServiceImpl struct {
	db          *postgresql.DB
	roomService room.RoomService
}

func NewMessageService() MessageService {
	appConfig := config.GetAppConfig()
	db := postgresql.NewDB(appConfig.Postgres.Dsn)
	roomService := room.NewRoomService()
	return &MessageServiceImpl{db: db, roomService: roomService}
}

func (m *MessageServiceImpl) List(ctx context.Context, userID string, friendID string) ([]tmessage.Message, error) {
	return m.db.MessageDB.List(ctx, userID, friendID)
}

func (m *MessageServiceImpl) ListRoom(ctx context.Context, roomID string) ([]tgroupmessage.GroupMessage, error) {
	groupMessages, err := m.db.GroupMessageDB.List(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return groupMessages, nil
}

func (m *MessageServiceImpl) SendToUser(ctx context.Context, receiver string, message string, send MessageSendFunc) error {
	sender := opctx.GetUserID(ctx)

	tMsg := tmessage.Message{
		Sender:    sender,
		Receiver:  receiver,
		Content:   message,
		CreatedAt: time.Now(),
	}

	msg := common.Msg{
		MsgType:  common.MsgTypeUser,
		Sender:   sender,
		Receiver: receiver,
		RoomID:   "",
		Content:  message,
	}

	if err := send(receiver, msg); err != nil {
		log.Error(ctx, "send message to user", zap.Error(err), zap.String("receiver", receiver), zap.String("msg", message))

		tMsg.Status = tmessage.StatusUnread
	} else {
		tMsg.Status = tmessage.StatusRead
	}

	if err := m.db.MessageDB.Insert(ctx, tMsg); err != nil {
		log.Error(ctx, "insert message", zap.Error(err))
	}

	return nil
}

func (m *MessageServiceImpl) SendToRoom(ctx context.Context, roomID string, message string, send MessageSendFunc) error {
	sender := opctx.GetUserID(ctx)

	groupMsg := tgroupmessage.GroupMessage{
		Sender:    sender,
		RoomID:    roomID,
		Content:   string(message),
		CreatedAt: time.Now(),
	}

	room, err := m.roomService.Detail(ctx, roomID)
	if err != nil {
		log.Error(ctx, "get room detail", zap.Error(err), zap.String("room_id", roomID))
		return err
	}

	if !slices.Contains(room.Users, sender) {
		return errors.New("user not in room")
	}

	msg := common.Msg{
		MsgType:  common.MsgTypeRoom,
		Sender:   sender,
		Receiver: "",
		RoomID:   roomID,
		Content:  message,
	}

	for _, receiver := range room.Users {
		if receiver == sender {
			continue
		}

		if err := send(receiver, msg); err != nil {
			log.Error(ctx, "send message to room", zap.Error(err), zap.String("room_id", roomID), zap.String("msg", message), zap.String("receiver", receiver))
		}
	}

	if err := m.db.GroupMessageDB.Insert(ctx, groupMsg); err != nil {
		log.Error(ctx, "insert group message", zap.Error(err))
	}

	return nil
}
