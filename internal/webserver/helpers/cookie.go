package helpers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func WriteCookie(c echo.Context, name string, value string, expiryInSeconds int) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(time.Second * time.Duration(expiryInSeconds))
	c.SetCookie(cookie)
}

func ReadAllCookies(c echo.Context) map[string]string {
	cookies := make(map[string]string)
	for _, cookie := range c.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	return cookies
}

func ReadCookie(c echo.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func DeleteCookie(c echo.Context, name string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour * 24)
	c.SetCookie(cookie)
}
