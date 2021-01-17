package domain

import "time"

const (
	UserEventUnassigned UserEventType = iota
	UserEventCreated
	UserEventUpdated
	UserEventDeleted
)

type UserID int64

type UserEventType int

type UserEvent interface {
	//instead of type casting or reflection in for loop  i suggest you to use Type method
	Type() UserEventType
}

type User struct {
	ID       UserID `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Country  string `json:"country"`
}

type UserCreatedEvent struct {
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
}

type UserUpdateEvent struct {
	CreatedAt time.Time `json:"created_at"`
	NewUser   User      `json:"new_user"`
	OldUser   User      `json:"old_user"`
}

type UserDeletedEvent struct {
	CreatedAt time.Time `json:"created_at"`
	UserID    UserID    `json:"user_id"`
}

func (u UserDeletedEvent) Type() UserEventType {
	return UserEventDeleted
}

func (u UserUpdateEvent) Type() UserEventType {
	return UserEventUpdated
}

func (u UserCreatedEvent) Type() UserEventType {
	return UserEventCreated
}

func NewUserCreatedEvent(user User) UserCreatedEvent {
	return UserCreatedEvent{User: user, CreatedAt: time.Now()}
}

func NewUserUpdateEvent(newUser User, oldUser User) UserUpdateEvent {
	return UserUpdateEvent{NewUser: newUser, OldUser: oldUser, CreatedAt: time.Now()}
}

func NewUserDeletedEvent(userID UserID) UserDeletedEvent {
	return UserDeletedEvent{UserID: userID, CreatedAt: time.Now()}
}
