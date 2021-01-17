package application

import (
	"context"
	"errors"
	"faceittt/domain"
	"fmt"
	"log"
)

var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
)

type UserEventFunc func(e domain.UserEvent) error

func handleUserEvents(e domain.UserEvent, consumers ...UserEventFunc) {
	for _, c := range consumers {
		if err := c(e); err != nil {
			//TODO: deal with error? just log it?
			log.Println(fmt.Errorf("handleUserEvents: %w", err))
		}
	}
}

type CreateUserParams struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Country  string `json:"country"`
	Password string `json:"password"`
}

type UpdateUserParams struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Country  string `json:"country"`
}

type UserRepository interface {
	Create(context.Context, CreateUserParams) (domain.User, error)
	Update(context.Context, domain.UserID, UpdateUserParams) (domain.User, error)
	Delete(context.Context, domain.UserID) error
}

type UserFinder interface {
	WithID(domain.UserID) UserFinder
	WithEmail(string) UserFinder
	WithName(string) UserFinder
	WithNickname(string) UserFinder
	WithCountry(string) UserFinder
	Find(context.Context) ([]domain.User, error)
	FindOne(context.Context) (domain.User, error)
}

func CreateUser(repo UserRepository, hooks ...UserEventFunc) func(context.Context, CreateUserParams) (domain.User, error) {
	return func(ctx context.Context, params CreateUserParams) (domain.User, error) {
		user, err := repo.Create(ctx, params)
		if err != nil {
			return domain.User{}, fmt.Errorf("CreateUser: %w", err)
		}
		handleUserEvents(domain.NewUserCreatedEvent(user), hooks...)
		return user, nil
	}
}

func UpdateUser(repo UserRepository, finder UserFinder, hooks ...UserEventFunc) func(context.Context, domain.UserID, UpdateUserParams) (domain.User, error) {
	return func(ctx context.Context, uID domain.UserID, params UpdateUserParams) (domain.User, error) {
		handleErr := func(err error) (domain.User, error) {
			return domain.User{}, fmt.Errorf("UpdateUser: %w", err)
		}
		oldUser, err := finder.WithID(uID).FindOne(ctx)
		if err != nil {
			return handleErr(err)
		}
		user, err := repo.Update(ctx, uID, params)
		if err != nil {
			return handleErr(err)
		}
		handleUserEvents(domain.NewUserUpdateEvent(user, oldUser), hooks...)
		return user, nil
	}
}

func DeleteUser(repo UserRepository, hooks ...UserEventFunc) func(context.Context, domain.UserID) error {
	return func(ctx context.Context, ID domain.UserID) error {
		if err := repo.Delete(ctx, ID); err != nil {
			return fmt.Errorf("DeleteUser: %w", err)
		}
		handleUserEvents(domain.NewUserDeletedEvent(ID), hooks...)
		return nil
	}
}

// Can be predefined as preset for a number of columns, can be adhoc closure in caller code
type UserFindOption func(UserFinder) UserFinder

func FindUsers(finder UserFinder) func(ctx context.Context, opts ...UserFindOption) ([]domain.User, error) {
	return func(ctx context.Context, opts ...UserFindOption) ([]domain.User, error) {
		for _, opt := range opts {
			finder = opt(finder)
		}
		u, err := finder.Find(ctx)
		if err != nil {
			return nil, fmt.Errorf("FindUser: %w", err)
		}
		return u, nil
	}
}

func FindUser(finder UserFinder) func(ctx context.Context, ID domain.UserID) (domain.User, error) {
	return func(ctx context.Context, ID domain.UserID) (domain.User, error) {
		u, err := finder.WithID(ID).FindOne(ctx)
		if err != nil {
			return domain.User{}, fmt.Errorf("FindUser: %w", err)
		}
		return u, nil
	}
}
