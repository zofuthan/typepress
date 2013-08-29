package blog

import (
	"net/http"
	"strings"

	"github.com/achun/db"

	. "controllers"
	g "global"
	"meta"

	_ "blog/install"
	_ "blog/root"
	_ "blog/sign"
)

var site = meta.Users{}

func init() {
	g.OnRouteBefore(InitSiteAndUser)
}
func InitSiteAndUser(wr http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == "/install/" {
		return true
	}
	sess, err := g.GetSession(r)
	if err == nil && site.User_id == 0 {
		var c db.Collection
		var item db.Item
		c, err = g.Db.Collection("users")
		if err == nil {
			item, err = c.Find(db.Cond{"User_id": g.BlogId})
		}
		if err == nil {
			g.ToStruct(item, &site)
			site.User_login = ""
			site.User_pass = ""
		}
	}
	if Error(wr, r, 500, err) {
		return false
	}
	i := sess.Values["user"]
	if i != nil {
		user := i.(meta.Users)
		user.User_login = ""
		user.User_pass = ""
		g.SetViewDat(r, "user", user)
		if g.Domain == "" {
			g.SetViewDat(r, "admin", "/admin")
		} else {
			admin := "//" + user.Site + "." + g.Domain
			if g.Port != "80" {
				admin += ":" + g.Port
			}
			admin += "/admin"
			g.SetViewDat(r, "admin", admin)
		}
	}

	if g.Domain == "" {
		g.SetViewDat(r, "site", site)
	} else if pos := strings.Index(r.Host, g.Domain); pos >= 0 {
		sitename := r.Host[0:pos]
		if sitename != "" {
			var c db.Collection
			var item db.Item
			usersite := meta.Users{}
			c, err = g.Db.Collection("users")
			if err == nil {
				item, err = c.Find(db.Cond{"User_site": sitename})
			}

			if err == nil {
				g.ToStruct(item, &usersite)
				usersite.User_login = ""
				usersite.User_pass = ""
				g.SetViewDat(r, "site", usersite)
			}
		}
	}
	return true
}
