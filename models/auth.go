package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/senko/clog"
	"golang.org/x/crypto/bcrypt"
)

type authRepo struct {
	db *sql.DB
}

func newAuthRepo(sqlDB *sql.DB) newsRepo {
	return authRepo{db: sqlDB}
}

type Authenticate struct {
	Id         int       `db:"id"`
	UserId     int       `db:"user_id"`
	ValidUntil time.Time `db:"valid_until"`
	Token      string    `db:"token"`
}

// errors
var (
	errAuthNotCreated = errors.New("Authentification not created")
	errAuthNotDeleted = errors.New("Authentification not deleted")
)

func (a *authRepo) IsAuth(token string) bool {
	var res *sql.Row
	res = a.db.QueryRow(
		`
		SELECT id, user_id, TO_CHAR(valid_until, 'YYYY-MM-DD HH24:MI:SS') valid_until, token
		FROM public.authentification WHERE token=$1
		`,
		token)

	var timestampStr string
	auth := Authenticate{}

	err := res.Scan(&auth.Id, &auth.UserId, &timestampStr, &auth.Token)
	if err != nil {
		clog.Warningf("Auth DB scan error: %s", err)
		return false
	}

	auth.ValidUntil, err = time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		clog.Warningf("Timestamp parse error: %s", err)
		return false
	}

	if time.Now().Unix() > auth.ValidUntil.Unix() {
		clog.Infof("Auth expired, id:%d, expiration:%s", auth.Id, auth.ValidUntil)
		return false
	}

	return true
}

func (a *authRepo) AuthUser(user User, password string) (Authenticate, error) {
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		clog.Debug("pssst password doesn't match")
		return Authenticate{}, UserNoMatch
	}

	token := uuid.NewV4().String()

	a := Authenticate{
		Token:      token,
		UserId:     user.Id,
		ValidUntil: time.Now().Add(time.Hour * 24),
	}
	err = createAuth(a)
	if err != nil {
		clog.Warningf("Auth not created: %s", err)
		return Authenticate{}, err
	}

	return a, nil
}

func (a *authRepo) createAuth(a Authenticate) error {
	res, err := a.db.Exec(
		`
		INSERT INTO public.authentification
		(user_id, valid_until, token)
		VALUES ($1, $2, $3)"
		`,
		a.UserId,
		pq.FormatTimestamp(a.ValidUntil),
		a.Token)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errAuthNotCreated
	}

	return nil
}

func (a *authRepo) RemoveAuth(token string) error {
	auth := Authenticate{}

	res := a.db.QueryRow(
		`
		SELECT id, user_id, TO_CHAR(valid_until, 'YYYY-MM-DD HH24:MI:SS') valid_until, token
		FROM public.authentification WHERE token=$1
		`,
		token)

	var validUntil string
	err := res.Scan(&auth.Id, &auth.UserId, &validUntil, &auth.Token)
	if err != nil {
		clog.Warningf("%s", err)
		return err
	}

	auth.ValidUntil = time.Now()

	// remove auth
	resRows, err := a.db.Exec(
		`
		UPDATE public.authentification
		SET valid_until=$1
		WHERE id=$2
		`,
		pq.FormatTimestamp(auth.ValidUntil),
		auth.Id)
	if err != nil {
		clog.Errorf("Auth removal error: %s", err)
		return err
	}

	rows, err := resRows.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errAuthNotDeleted
	}

	return nil
}
