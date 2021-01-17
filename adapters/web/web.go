package web

import (
	"context"
	"encoding/json"
	"errors"
	"faceittt/application"
	"faceittt/domain"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type createUserFunc func(context.Context, application.CreateUserParams) (domain.User, error)
type updateUserFunc func(context.Context, domain.UserID, application.UpdateUserParams) (domain.User, error)
type deleteUserFunc func(context.Context, domain.UserID) error
type findUsersFunc func(ctx context.Context, opts ...application.UserFindOption) ([]domain.User, error)
type findUserFunc func(ctx context.Context, ID domain.UserID) (domain.User, error)

type WebParams struct {
	CreateUser createUserFunc
	UpdateUser updateUserFunc
	DeleteUser deleteUserFunc
	FindUsers  findUsersFunc
	FindUser   findUserFunc
}

func Router(wp WebParams) http.Handler {
	r := chi.NewRouter()
	r.Get("/users", FindUsersHandler(wp.FindUsers))
	r.Post("/users", CreateUserHandler(wp.CreateUser))
	r.Get("/users/{id}", FindUserHandler(wp.FindUser))
	r.Put("/users/{id}", UpdateUserHandler(wp.UpdateUser))
	r.Delete("/users/{id}", DeleteUserHandler(wp.DeleteUser))
	return r
}

func CreateUserHandler(f createUserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logErr := func(err error) {
			log.Println(fmt.Errorf("CreateUserHandler: %w", err))
		}
		params := application.CreateUserParams{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := f(r.Context(), params)
		if err != nil {
			if errors.Is(err, application.ErrAlreadyExists) {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			logErr(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			logErr(err)
			return
		}
	}
}

func UpdateUserHandler(f updateUserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logErr := func(err error) {
			log.Println(fmt.Errorf("CreateUserHandler: %w", err))
		}
		params := application.UpdateUserParams{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		uID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			//we certainly don't have that resource
			http.NotFound(w, r)
			return
		}
		user, err := f(r.Context(), domain.UserID(uID), params)
		if err != nil {
			if errors.Is(err, application.ErrNotFound) {
				http.NotFound(w, r)
				return
			}
			logErr(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			logErr(err)
			return
		}
	}
}

func DeleteUserHandler(f deleteUserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			//we certainly don't have that resource
			http.NotFound(w, r)
			return
		}
		if err := f(r.Context(), domain.UserID(uID)); err != nil {
			log.Println("DeleteUserHandler: %w", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func FindUsersHandler(f findUsersFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logErr := func(err error) {
			log.Println(fmt.Errorf("FindUsersHandler: %w", err))
		}
		query := r.URL.Query()
		filterFunc := func(finder application.UserFinder) application.UserFinder {
			if v := query.Get("email"); v != "" {
				finder = finder.WithEmail(v)
			}
			if v := query.Get("nickname"); v != "" {
				finder = finder.WithNickname(v)
			}
			if v := query.Get("name"); v != "" {
				finder = finder.WithName(v)
			}
			if v := query.Get("country"); v != "" {
				finder = finder.WithCountry(v)
			}
			return finder
		}
		users, err := f(r.Context(), filterFunc)
		if err != nil {
			logErr(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			logErr(err)
		}
	}
}

func FindUserHandler(f findUserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logErr := func(err error) {
			log.Println(fmt.Errorf("FindUserHandler: %w", err))
		}
		uID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			//we certainly don't have that resource
			http.NotFound(w, r)
			return
		}
		user, err := f(r.Context(), domain.UserID(uID))
		if err != nil {
			if errors.Is(err, application.ErrNotFound) {
				http.NotFound(w, r)
				return
			}
			logErr(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			logErr(err)
		}
	}
}
