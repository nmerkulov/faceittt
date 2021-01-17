package postgres

const initialMigration = `create table IF NOT EXISTS users
(
	name varchar(100),
	nickname varchar(100),
	last_name varchar(100),
	password text,
	email varchar(100),
	country varchar(100),
	id serial
		constraint users_pk
			primary key
);

create unique index IF NOT EXISTS user_email_uindex 
	on users (email);

create unique index IF NOT EXISTS user_nickname_uindex
	on users (nickname);

`
