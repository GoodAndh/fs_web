package user

import (
	"backend/server/utils"
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {

	return &repository{db: db}
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}

	err := r.db.QueryRowContext(ctx, "select id,username,email,password,created_at,last_updated from users where username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return &User{}, utils.ErrNotFound
		}
		return &User{}, err
	}

	return user, nil

}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(ctx, "select id,username,email,password,created_at,last_updated from users where email = ?", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return &User{}, utils.ErrNotFound
		}
		return &User{}, err
	}

	return user, nil
}

func (r *repository) CreateUsers(ctx context.Context, user *User) (int, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	result, err := tx.ExecContext(ctx, "insert into users(username,password,email,created_at,last_updated) values(?,?,?,?,?)", user.Username, user.Password, user.Email, user.CreatedAt, user.LastUpdated)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	ids, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(ids), nil
}

func (r *repository) GetUserProfile(ctx context.Context, userID int) (*UserProfile, error) {
	us := &UserProfile{}
	err := r.db.QueryRowContext(ctx, "select id,userID,url,captions from users_profile where userid=", userID).Scan(&us.ID, &us.UserID, &us.Url, &us.Captions)
	if err != nil {
		if err == sql.ErrNoRows {
			return &UserProfile{}, utils.ErrNotFound
		}
		return &UserProfile{}, err
	}

	return us, nil
}

func (r *repository) CreateUserProfile(ctx context.Context, user *UserProfile) (int,error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0,err
	}
	result, err := tx.ExecContext(ctx, "insert into users_profile(userid,url,captions) values(?,?,?)", user.UserID, user.Url, user.Captions)
	if err != nil {
		tx.Rollback()
		return 0,err
	}

	id,err:=result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0,err
	}

	return int(id),nil
}
