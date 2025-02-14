package message

import (
	"context"
	"encoding/json"
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

type MessageSendFunc func(userID string, message []byte) error

type MessageService interface {
	List(ctx context.Context, userID string, friendID string) ([]tmessage.Message, error)
	ListRoom(ctx context.Context, roomID string) ([]tgroupmessage.GroupMessage, error)
	SendToUser(ctx context.Context, userID string, message []byte, send MessageSendFunc) error
	SendToRoom(ctx context.Context, roomID string, message []byte, send MessageSendFunc) error
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

func (m *MessageServiceImpl) SendToUser(ctx context.Context, userID string, message []byte, send MessageSendFunc) error {
	sender := opctx.GetUserID(ctx)

	tMsg := tmessage.Message{
		Sender:    sender,
		Receiver:  userID,
		Content:   string(message),
		CreatedAt: time.Now(),
	}

	msg := common.Msg{
		MsgType:  common.MsgTypeUser,
		Sender:   sender,
		Receiver: userID,
		RoomID:   "",
		Content:  string(message),
	}

	if err := warpSend(ctx, send, userID, msg); err != nil {
		log.Error(ctx, "send message to user", zap.Error(err), zap.String("user_id", userID), zap.String("msg", string(message)))

		tMsg.Status = tmessage.StatusUnread
	} else {
		tMsg.Status = tmessage.StatusRead
	}

	// only save message to db when sender is self
	if sender == userID {
		if err := m.db.MessageDB.Insert(ctx, tMsg); err != nil {
			log.Error(ctx, "insert message", zap.Error(err))
		}
	}

	return nil
}

func (m *MessageServiceImpl) SendToRoom(ctx context.Context, roomID string, message []byte, send MessageSendFunc) error {
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
		Content:  string(message),
	}

	for _, user := range room.Users {
		if user == sender {
			continue
		}

		if err := warpSend(ctx, send, user, msg); err != nil {
			log.Error(ctx, "send message to room", zap.Error(err), zap.String("room_id", roomID), zap.String("msg", string(message)))
		}
	}

	if err := m.db.GroupMessageDB.Insert(ctx, groupMsg); err != nil {
		log.Error(ctx, "insert group message", zap.Error(err))
	}

	return nil
}

func warpSend(_ context.Context, send MessageSendFunc, userID string, v interface{}) error {
	if s, ok := v.(string); ok {
		return send(userID, []byte(s))
	}

	if b, ok := v.([]byte); ok {
		return send(userID, b)
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return send(userID, jsonBytes)
}
