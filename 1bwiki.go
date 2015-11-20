package main

import (
	"fmt"
	"net/http"

	m "1bwiki/model"
	t "1bwiki/tmpl"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func wikiPage(c *echo.Context) error {
	fmt.Println(c.Request().URL)
	p := m.GetPage("sss")
	if p != nil {
		return c.String(http.StatusOK, "Page Exists")
	}
	return c.HTML(http.StatusOK, t.Newpage())
}

func init() {
	db, err := sqlx.Connect("sqlite3", "./1bwiki.db")
	if err != nil {
		panic(err)
	}
	db.Exec(`create table if not exists text (id integer primary KEY, text blob)`)
}

func main() {
	e := echo.New()
	e.Use(mw.Logger())
	e.StripTrailingSlash()
	e.Static("/static", "static")

	e.Get("/*", wikiPage)
	e.Run(":8000")
}
