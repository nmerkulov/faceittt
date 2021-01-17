package postgres

import (
	"context"
	"database/sql"
	"errors"
	"faceittt/application"
	"faceittt/domain"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"log"
	"strings"
)

type PGUserFinder struct {
	db *sqlx.DB
	b  sq.SelectBuilder
}

func (p PGUserFinder) WithID(ID domain.UserID) application.UserFinder {
	p.b = p.b.Where("id=?", ID)
	return p
}

func (p PGUserFinder) WithEmail(e string) application.UserFinder {
	p.b = p.b.Where("email=?", e)
	return p
}

func (p PGUserFinder) WithName(n string) application.UserFinder {
	p.b = p.b.Where("name=?", n)
	return p
}

func (p PGUserFinder) WithNickname(nn string) application.UserFinder {
	p.b = p.b.Where("nickname=?", nn)
	return p
}

func (p PGUserFinder) WithCountry(c string) application.UserFinder {
	p.b = p.b.Where("country=?", c)
	return p
}

func (p PGUserFinder) Find(ctx context.Context) ([]domain.User, error) {
	handleErr := func(err error) ([]domain.User, error) {
		return nil, fmt.Errorf("PGUserFinder.Find: %w", err)
	}
	var users []domain.User
	q, args, err := p.b.ToSql()
	if err != nil {
		return handleErr(err)
	}
	log.Println(q)
	rows, err := p.db.QueryxContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]domain.User, 0, 0), nil
		}
		return handleErr(err)
	}
	var u domain.User
	for rows.Next() {
		if err := rows.StructScan(&u); err != nil {
			return handleErr(err)
		}
		users = append(users, u)

	}
	return users, nil
}

func (p PGUserFinder) FindOne(ctx context.Context) (domain.User, error) {
	handleErr := func(err error) (domain.User, error) {
		return domain.User{}, fmt.Errorf("PGUserFinder.Find: %w", err)
	}
	var user domain.User
	q, args, err := p.b.ToSql()
	if err != nil {
		return handleErr(err)
	}
	if err := p.db.GetContext(ctx, &user, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, application.ErrNotFound
		}
		return handleErr(err)
	}
	return user, nil
}

func NewPGUserFinder(cs string) (PGUserFinder, error) {
	db, err := sqlx.Connect("postgres", cs)
	if err != nil {
		return PGUserFinder{}, fmt.Errorf("NewPGUserFinder#Connect: %w", err)
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "name", "last_name", "nickname", "email", "country").From("users")
	return PGUserFinder{db: db, b: b}, nil
}
