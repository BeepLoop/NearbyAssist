package session

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Session interface{}

func NewCookieStore() *sessions.CookieStore {
	// TODO: Implement better keyPair
	return sessions.NewCookieStore([]byte("secret"))
}

type CookieSession struct {
	name    string
	options *sessions.Options
}

func NewSession() CookieSession {
	return CookieSession{
		name: "session",
		options: &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
		},
	}
}

func (s CookieSession) Create(ctx echo.Context) error {
	sess, err := session.Get(s.name, ctx)
	if err != nil {
		fmt.Println("error getting session: ", err.Error())
		return err
	}

	sess.Options = s.options
	sess.Values["foo"] = "bar"

	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		fmt.Println("error saving session: ", err.Error())
		return err
	}

	return nil
}
