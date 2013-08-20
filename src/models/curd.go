package models

import (
	"errors"
	//"net/url"
	//"reflect"
	//
	"github.com/achun/db"
	//"github.com/gosexy/to"
	//
	g "global"
	"meta"
)

var Tables = struct {
	Terms             meta.Terms
	Termtaxonomy      meta.Termtaxonomy
	Termrelationships meta.Termrelationships
	Commentmeta       meta.Commentmeta
	Comments          meta.Comments
	Links             meta.Links
	Postmeta          meta.Postmeta
	Posts             meta.Posts
	Users             meta.Users
	Usermeta          meta.Usermeta
	Members           meta.Members
	Sitemeta          meta.Sitemeta
}{}

type Cond struct {
	Fields []string // fetch fields
	Exists []string // must be exists
}
type CurdCond struct {
	name string // Collection name
	id   string // id named field
	A    *Cond  // Append Cond
	U    *Cond  // Update Cond
	Q    *Cond  // Query Cond
	D    *Cond  // Delete Cond
}

func NewCond() *Cond {
	return &Cond{}
}

func NewCurdCond(name string, id string) *CurdCond {
	return &CurdCond{
		name: name,
		id:   id,
		A:    NewCond(),
		U:    NewCond(),
		Q:    NewCond(),
		D:    NewCond(),
	}
}

var (
	Terms             = NewCurdCond("Terms", "Term_id")
	Termtaxonomy      = NewCurdCond("Termtaxonomy", "Termtaxonomy_id")
	Termrelationships = NewCurdCond("Termrelationships", "")
	Commentmeta       = NewCurdCond("Commentmeta", "Commentmeta_id")
	Comments          = NewCurdCond("Comments", "Comment_id")
	Links             = NewCurdCond("Links", "Link_id")
	Postmeta          = NewCurdCond("Postmeta", "Postmeta_id")
	Posts             = NewCurdCond("Posts", "Post_id")
	Users             = NewCurdCond("Users", "User_id")
	Usermeta          = NewCurdCond("Usermeta", "Usermeta_id")
	Members           = NewCurdCond("Members", "")
	Sitemeta          = NewCurdCond("Sitemeta", "Sitemeta_id")
)

func init() {
	InitCond()
}

func InitCond() {
	Users.A.Exists = []string{"User_login", "User_pass", "User_nicename", "Site"}
	Users.U.Fields = []string{"User_login", "User_pass", "User_nicename", "Site", "User_status"}
	Users.U.Exists = []string{"User_id"}
}

var MustBeExists = "field must be exists: "
var NotEnoughArgumentsFor = "not enough arguments for: "
var InValidField = "invalid field : "

func (p *CurdCond) Name() string {
	return p.name
}

func (p *CurdCond) IdField() string {
	return p.id
}

func (p *CurdCond) Append(mp map[string]interface{}) ([]db.Id, error) {
	c, err := g.Db.Collection(p.name)
	if err != nil {
		return nil, err
	}

	fetch := map[string]interface{}{}
	f := p.A
	for _, field := range f.Fields {
		v, ok := mp[field]
		if ok {
			fetch[field] = v
		}
	}
	for _, field := range f.Exists {
		v, ok := mp[field]
		if !ok {
			return nil, errors.New(MustBeExists + field + " on " + p.name + ".Append()")
		}
		fetch[field] = v
	}
	if len(fetch) == 0 {
		return nil, errors.New(NotEnoughArgumentsFor + " on " + p.name + ".Append()")
	}
	errstr := meta.Valid.IsValidMap(fetch)
	if errstr != "" {
		return nil, errors.New(InValidField + errstr + " on " + p.name + ".Append()")
	}
	return c.Append(fetch)
}
