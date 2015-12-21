package main

import (
	"net/http"
	"strings"

	m "1bwiki/model"
	"1bwiki/tmpl/page"
	"1bwiki/tmpl/special"

	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
)

func edit(c *echo.Context) error {
	n, t := parseTitle(c.Query("title"))

	ut := strings.ToLower(t)
	if strings.HasPrefix(ut, noEditArea) {
		return echo.NewHTTPError(http.StatusForbidden, "Editing of special pages disallowed")
	}

	ct := cleanTitle(t)
	if ct != t {
		if n != "" {
			n += ":"
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/special/edit?title="+n+ct)
	}
	pv := m.GetPageView(n, t)
	if pv.NiceTitle == "" {
		if n != "" {
			n += ":"
		}
		pv.NameSpace = n
		pv.Title = t
		pv.NiceTitle = strings.Replace(n+t, "_", " ", -1)
	}
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, page.Editpage(val.(*m.User), pv))
}

func history(c *echo.Context) error {
	n, t := parseTitle(c.Query("title"))
	ct := cleanTitle(t)
	if ct != t {
		if n != "" {
			n += ":"
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/special/history?title="+n+ct)
	}
	revs, err := m.GetPageRevisions(c.Query("title"))
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	niceTitle := m.NiceTitle(c.Query("title"))
	return c.HTML(http.StatusOK, page.History(val.(*m.User), niceTitle, revs))
}

func recentChanges(c *echo.Context) error {
	revs, err := m.GetRevisions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Recentchanges(val.(*m.User), revs))
}

func pages(c *echo.Context) error {
	p, err := m.GetPages()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Pages(val.(*m.User), p))
}
