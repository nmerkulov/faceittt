package postgres

import (
	"context"
	"errors"
	"faceittt/application"
	"faceittt/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	pqDuplicateKeyCode = "23505"
)

type PGUserRepo struct {
	db *sqlx.DB
}

func (pr PGUserRepo) Create(ctx context.Context, p application.CreateUserParams) (domain.User, error) {
	handleErr := func(err error) (domain.User, error) {
		return domain.User{}, fmt.Errorf("PGUserRepo: %w", err)

	}
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return handleErr(err)
	}
	p.Password = string(hash)

	var u domain.User
	if err := pr.db.GetContext(
		ctx,
		&u,
		`
		insert into users (name, last_name, nickname, email, password, country)
		values($1, $2, $3, $4, $5, $6) returning id, name, last_name, nickname, email, country`,
		p.Name, p.LastName, p.Nickname, p.Email, p.Password, p.Country); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if string(pqErr.Code) == pqDuplicateKeyCode {
				return handleErr(application.ErrAlreadyExists)
			}
		}
		return handleErr(err)
	}
	return u, nil
}

func (pr PGUserRepo) Update(ctx context.Context, ID domain.UserID, p application.UpdateUserParams) (domain.User, error) {
	var u domain.User
	if err := pr.db.GetContext(
		ctx,
		&u, `
		update users set
		name=$2, last_name=$3, country=$4
		where id=$1
		returning id, name, last_name, nickname, email, country`,
		ID, p.Name, p.LastName, p.Country); err != nil {
		return domain.User{}, fmt.Errorf("PGUserRepo.Update: %w", err)
	}
	return u, nil
}

func (pr PGUserRepo) Delete(ctx context.Context, ID domain.UserID) error {
	if _, err := pr.db.QueryContext(ctx, `delete from users where id=$1`, ID); err != nil {
		return fmt.Errorf("PGUserRepo.Delete: %w", err)
	}
	return nil
}

func (pr PGUserRepo) MustMigrate(ctx context.Context) {
	if _, err := pr.db.QueryContext(ctx, initialMigration); err != nil {
		panic(fmt.Errorf("PGUserRepo.Migrate: %w", err))
	}
}

func NewPGUserRepo(cs string) (PGUserRepo, error) {
	db, err := sqlx.Connect("postgres", cs)
	if err != nil {
		return PGUserRepo{}, fmt.Errorf("NewPGUserRepo#Connect: %w", err)
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return PGUserRepo{db: db}, nil
}
