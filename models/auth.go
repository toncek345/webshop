package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/senko/clog"
	"golang.org/x/crypto/bcrypt"
)

type authRepo struct {
	db *sqlx.DB
}

func newAuthRepo(sqlDB *sqlx.DB) authRepo {
	return authRepo{db: sqlDB}
}

type Authentication struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	ValidUntil time.Time `db:"valid_until"`
	Token      string    `db:"token"`
}

func (ar *authRepo) IsAuth(token string) (bool, error) {
	var auth Authentication
	if err := ar.db.Get(
		&auth,
		`
		SELECT id, user_id, valid_until, token
		FROM authentications WHERE token=$1
		`,
		token); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("models/auth: error getting auth: %w", err)
	}

	if time.Now().Unix() > auth.ValidUntil.Unix() {
		clog.Infof("Auth expired, id:%d, expiration:%s", auth.ID, auth.ValidUntil)
		return false, nil
	}

	return true, nil
}

func (ar *authRepo) AuthUser(user User, password string) (Authentication, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		clog.Debug("pssst password doesn't match")
		return Authentication{}, fmt.Errorf("models/auth: user not found")
	}

	token := uuid.NewV4().String()

	a := Authentication{
		Token:      token,
		UserID:     user.ID,
		ValidUntil: time.Now().Add(time.Hour * 24),
	}

	if err := ar.createAuth(a); err != nil {
		clog.Warningf("Auth not created: %s", err)
		return Authentication{}, err
	}

	return a, nil
}

func (ar *authRepo) createAuth(a Authentication) error {
	if _, err := ar.db.Exec(
		`
		INSERT INTO authentications
		(user_id, valid_until, token)
		VALUES ($1, $2, $3)
		`,
		a.UserID,
		pq.FormatTimestamp(a.ValidUntil),
		a.Token); err != nil {
		return fmt.Errorf("models/auth: error inserting auth: %w", err)
	}

	return nil
}

func (ar *authRepo) RemoveAuth(token string) error {
	if _, err := ar.db.Exec(
		"UPDATE authentications SET valid_until = now() WHERE token = $1",
		token); err != nil {
		return fmt.Errorf("models/auth: error invalidating authorization: %w", err)
	}

	return nil
}
