package utils

import (
	"encoding/base64"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// msgType = "error" | "success"
func SetFlashMessage(c echo.Context, msgType, msg string) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	rawMessage := msgType + ";" + msg
	encoded := base64.StdEncoding.EncodeToString([]byte(rawMessage))

	sess.Values["flash_message"] = encoded
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}

func RetrieveFlashMessage(c echo.Context) (string, bool, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return "", false, err
	}

	message, exists := sess.Values["flash_message"].(string)
	delete(sess.Values, "flash_message")
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return "", false, err
	}

	return message, exists, nil
}
