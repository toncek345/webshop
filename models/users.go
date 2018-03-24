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

// errors
var (
	UserNoMatch           = errors.New("Username or password does not match")
	AuthNotAuthenticates  = errors.New("User not authenticated")
	UserNotFoundByIdError = errors.New("User not found by id")
	UserNotCreatedError   = errors.New("User not created")
)

// sql-s
var (
	// select
	getUserById                  = "SELECT * FROM public.user WHERE id=$1"
	getUserByPasswordAndUsername = "SELECT * FROM public.user u WHERE u.username=$1"

	// update
	updateUser = "UPDATE public.user " +
		"SET username=$1, public.user.password=$2 " +
		"WHERE id=$3"

	// create
	createUser = "INSERT INTO public.user (username, password) " +
		"VALUES ($1, $2)"
)

type User struct {
	Id       int
	Username string
	Password string
}

// Creating user table & adding admin
func initUser() (err error) {
	sql := `CREATE TABLE public.user (
	id serial NOT NULL PRIMARY KEY,
    username varchar(250),
    password varchar(60)
)`

	_, err = sqlDB.Query(sql)
	if err != nil {
		clog.Errorf("%s", err)
		return
	}

	// creating admin
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		clog.Errorf("%s", err)
		return
	}

	admin := User{
		Username: "admin",
		Password: string(hash),
	}

	err = CreateUser(admin)
	if err != nil {
		clog.Errorf("%s", err)
		return
	}

	return
}

func findUser(username string) (u User, err error) {
	var res *sql.Row
	res = sqlDB.QueryRow(getUserByPasswordAndUsername, username)
	err = res.Scan(&u.Id, &u.Username, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			err = UserNoMatch
			return
		}
		return
	}

	return
}

func UpdateUser(id int, u User) error {
	res, err := sqlDB.Exec(updateUser, u.Username, u.Password, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return UserNotFoundByIdError
	}

	return nil
}

func CreateUser(u User) error {
	res, err := sqlDB.Exec(createUser, u.Username, u.Password)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return UserNotCreatedError
	}

	return nil
}

/* ******* auth ******* */

type Authenticate struct {
	Id         int
	UserId     int
	ValidUntil time.Time
	Token      string
}

func initAuth() (err error) {
	sql := `CREATE TABLE public.authentification (
	id serial NOT NULL PRIMARY KEY,
    user_id int,
    valid_until timestamp,
    token uuid,
    FOREIGN KEY (user_id) REFERENCES public.user(id)
)`
	_, err = sqlDB.Query(sql)
	return
}

// sql-s
var (
	// select
	selectAuthSession = "SELECT id, user_id, TO_CHAR(valid_until, 'YYYY-MM-DD HH24:MI:SS') valid_until, token " +
		"FROM public.authentification WHERE token=$1"

	// create
	createAuthSession = "INSERT INTO public.authentification " +
		"(user_id, valid_until, token) " +
		"VALUES ($1, $2, $3)"

	// update
	removeAuthSession = "UPDATE public.authentification " +
		"SET valid_until=$1 " +
		"WHERE id=$2"
)

// errors
var (
	errAuthNotCreated = errors.New("Authentification not created")
	errAuthNotDeleted = errors.New("Authentification not deleted")
)

func IsAuth(token string) bool {
	var res *sql.Row
	res = sqlDB.QueryRow(selectAuthSession, token)

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

	if time.Now().Second() > auth.ValidUntil.Second() {
		clog.Infof("Auth expired, id:%d, expiration:%s", auth.Id, auth.ValidUntil)
		return false
	}

	return true
}

func AuthUser(username, password string) (Authenticate, error) {
	u, err := findUser(username)
	if err != nil {
		clog.Infof("find user: %s", err)
		return Authenticate{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		clog.Debug("pssst password doesn't match")
		return Authenticate{}, UserNoMatch
	}

	token := uuid.Must(uuid.NewV4()).String()

	a := Authenticate{
		Token:      token,
		UserId:     u.Id,
		ValidUntil: time.Now().Add(time.Hour * 24),
	}
	err = createAuth(a)
	if err != nil {
		clog.Warningf("Auth not created: %s", err)
		return Authenticate{}, err
	}

	return a, nil
}

func createAuth(a Authenticate) error {
	res, err := sqlDB.Exec(createAuthSession,
		a.UserId, pq.FormatTimestamp(a.ValidUntil), a.Token)
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

func RemoveAuth(token string) error {
	auth := Authenticate{}

	res := sqlDB.QueryRow(selectAuthSession, token)

	var validUntil string
	err := res.Scan(&auth.Id, &auth.UserId, &validUntil, &auth.Token)
	if err != nil {
		clog.Warningf("%s", err)
		return err
	}

	auth.ValidUntil = time.Now()

	// remove auth
	resRows, err := sqlDB.Exec(removeAuthSession, pq.FormatTimestamp(auth.ValidUntil), auth.Id)
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
