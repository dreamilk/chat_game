package room

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-uuid"

	"chat_game/models/redis"
)

type Room struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	RoomName  string    `json:"room_name"`
	CreatedAt time.Time `json:"created_at"`
	Users     []string  `json:"users"`
}

const (
	roomKey = "room"
	userKey = "room:%s:users"
)

type RoomService interface {
	Create(ctx context.Context, ownerID string, roomName string) (*Room, error)
	List(ctx context.Context) ([]*Room, error)
	Detail(ctx context.Context, roomID string) (*Room, error)
	Join(ctx context.Context, roomID string, userID string) error
	Leave(ctx context.Context, roomID string, userID string) error
	Delete(ctx context.Context, roomID string) error
}

type roomServiceImpl struct {
	redisClient redis.Client
}

var _ RoomService = &roomServiceImpl{}

func NewRoomService(redisClient redis.Client) RoomService {
	return &roomServiceImpl{
		redisClient: redisClient,
	}
}

// Delete implements RoomService.
func (r *roomServiceImpl) Delete(ctx context.Context, roomID string) error {
	panic("unimplemented")
}

// Detail implements RoomService.
func (r *roomServiceImpl) Detail(ctx context.Context, roomID string) (*Room, error) {
	panic("unimplemented")
}

// Join implements RoomService.
func (r *roomServiceImpl) Join(ctx context.Context, roomID string, userID string) error {
	panic("unimplemented")
}

// Leave implements RoomService.
func (r *roomServiceImpl) Leave(ctx context.Context, roomID string, userID string) error {
	panic("unimplemented")
}

func (r *roomServiceImpl) Create(ctx context.Context, ownerID string, roomName string) (*Room, error) {
	roomID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	room := Room{
		ID:        roomID,
		OwnerID:   ownerID,
		RoomName:  roomName,
		CreatedAt: time.Now(),
		Users:     []string{},
	}

	roomJSON, err := json.Marshal(room)
	if err != nil {
		return nil, err
	}

	if err := r.redisClient.HSet(ctx, roomKey, roomID, string(roomJSON)); err != nil {
		return nil, err
	}

	userKey := fmt.Sprintf(userKey, roomID)
	if err := r.redisClient.SAdd(ctx, userKey, ownerID); err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *roomServiceImpl) List(ctx context.Context) ([]*Room, error) {
	rooms, err := r.redisClient.HGetAll(ctx, roomKey)
	if err != nil {
		return nil, err
	}

	roomList := make([]*Room, 0, len(rooms))
	for roomID, roomJSON := range rooms {
		room := Room{}
		if err := json.Unmarshal([]byte(roomJSON), &room); err != nil {
			return nil, err
		}

		room.ID = roomID
		roomList = append(roomList, &room)
	}

	return roomList, nil
}
