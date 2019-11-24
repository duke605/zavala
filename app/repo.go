package app

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

// Repo is used to interact with the application's datastore
type Repo struct {
	db *sqlx.DB
}

// NewRepo creates a new repo using the driver and DSN provided. If an error occurs
// when connecting to the database, the function will panic
func NewRepo(driver, dsn string) *Repo {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1)

	return &Repo{
		db: db,
	}
}

// Guild represents a discord guild in the database
type Guild struct {
	ID      uint64 `db:"id"`
	GroupID *int64 `db:"group_id"`
}

// User represets a discord user that has connected their Destiny 2 account
type User struct {
	ID             uint64    `db:"id"`
	MembershipType int       `db:"membership_type"`
	MembershipID   int64     `db:"membership_id"`
	AccessToken    string    `db:"access_token"`
	RefreshToken   string    `db:"refresh_token"`
	Expiry         time.Time `db:"expiry"`
}

// Token creates and returns an oauth2.Token for the user
func (u *User) Token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
		Expiry:       u.Expiry,
		TokenType:    "Bearer",
	}
}

// Transaction begins a transaction and calls the provided function. The context passed to fn should be passed
// to all function calls made to repo to ensure that they are made against the database transaction and not the
// database itself. Any errors that occur in fn should be returned so that the transaction is rolled back instead
// of commited
func (r Repo) Transaction(ctx context.Context, fn func(context.Context) error) error {
	ctx = ensureContext(ctx)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// Recovering from panic to rollback the transaction before panicing again
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// Something went wrong in fn, rolling back
			tx.Rollback()
		} else {
			// Nothing when wrong, commiting transaction
			tx.Commit()
		}
	}()

	// Storing transaction in context
	ctx = context.WithValue(ctx, transaction, tx)

	err = fn(ctx)
	return err
}

// SyncGuilds deletes guilds from the database that do not appear in the provided guildIDs array
func (r *Repo) SyncGuilds(ctx context.Context, guildIDs []uint64) error {
	ctx = ensureContext(ctx)
	sql, args, err := sq.Delete("guilds").Where(sq.NotEq{
		"id": guildIDs,
	}).ToSql()
	if err != nil {
		return err
	}

	execer := execerFromContext(ctx, r.db)
	_, err = execer.ExecContext(ctx, sql, args...)
	return err
}

// GetGuildByID gets a guild from the DB by it's ID.
//
// If the guild is not found in the DB sql.ErrNoRows will be returned
func (r *Repo) GetGuildByID(ctx context.Context, gid uint64) (Guild, error) {
	ctx = ensureContext(ctx)
	g := Guild{}
	sql, args, err := sq.Select("*").From("guilds").Where(sq.Eq{
		"id": gid,
	}).ToSql()
	if err != nil {
		return g, nil
	}

	execer := execerFromContext(ctx, r.db)
	if err = execer.GetContext(ctx, &g, sql, args...); err != nil {
		return g, err
	}

	return g, err
}

// GetUserByID gets a user from the DB by their DB.
//
// If the user is not found in the DB sql.ErrNoRows will be returned
func (r *Repo) GetUserByID(ctx context.Context, id uint64) (User, error) {
	ctx = ensureContext(ctx)
	u := User{}
	sql, args, err := sq.Select("*").From("users").Where(sq.Eq{
		"id": id,
	}).ToSql()
	if err != nil {
		return u, nil
	}

	execer := execerFromContext(ctx, r.db)
	if err = execer.GetContext(ctx, &u, sql, args...); err != nil {
		return u, err
	}

	return u, err
}

// GetUserByMembershipID gets a user from the DB by their destiny 2 membership ID.
//
// If the user is not found in the DB sql.ErrNoRows will be returned
func (r *Repo) GetUserByMembershipID(ctx context.Context, id int64) (User, error) {
	ctx = ensureContext(ctx)
	u := User{}
	sql, args, err := sq.Select("*").From("users").Where(sq.Eq{
		"membership_id": id,
	}).ToSql()
	if err != nil {
		return u, nil
	}

	execer := execerFromContext(ctx, r.db)
	if err = execer.GetContext(ctx, &u, sql, args...); err != nil {
		return u, err
	}

	return u, err
}

// key is used to store values in context and retrieve them
type key int

const (

	// transaction is used to store a transaction object in a context
	transaction key = iota
)

type execer interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	GetContext(context.Context, interface{}, string, ...interface{}) error
}

// execerFromContext attempts to pull a transaction from the context provided and return it. If the
// context does not contain a transaction, db will be returned instead
func execerFromContext(ctx context.Context, db *sqlx.DB) execer {
	if e := ctx.Value(transaction); e != nil {
		return e.(execer)
	}

	return db
}

// ensureContext returns a new backround context if the provided on is nil. Otherwise the provided
// context will be returned
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}
