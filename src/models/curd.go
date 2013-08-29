package models

import (
	"errors"
	"net/http"

	"github.com/achun/db"

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
	Users.Q.Fields = []string{"User_id", "User_login", "User_pass", "Site", "User_nicename", "User_status"}
	Users.A.Fields = []string{"User_login", "User_pass", "User_nicename", "Site"}
	Users.A.Exists = Users.A.Fields
	Users.U.Fields = []string{"User_pass", "User_nicename", "Site", "User_status"}
	Users.U.Exists = []string{"User_id"}
}

func MustBeExists(r *http.Request, s1, s2, s3 string) string {
	return g.I18n(r, "%s must be exists on %s %s", s1, s2, s3)
}

func NotEnoughArguments(r *http.Request, s1, s3 string) string {
	return g.I18n(r, "not enough arguments on %s %s", s1, s3)
}

func InValidField(r *http.Request, s1, s2, s3 string) string {
	return g.I18n(r, "%s is invalid on %s %s", s1, s2, s3)
}

func (p *CurdCond) Name() string {
	return p.name
}

func (p *CurdCond) IdField() string {
	return p.id
}

func (p *CurdCond) Append(r *http.Request) ([]db.Id, error) {
	c, err := g.Db.Collection(p.name)
	if err != nil {
		return nil, err
	}

	f := p.A
	vs := r.Form
	for _, key := range f.Exists {
		_, ok := vs[key]
		if !ok {
			return nil, errors.New(MustBeExists(r, key, p.name, "Append()"))
		}
	}
	item := db.Item{}
	for _, key := range f.Fields {
		_, ok := vs[key]
		if !ok {
			continue
		}
		v := vs.Get(key)
		if !meta.Valid.IsValid(key, v) {
			return nil, errors.New(InValidField(r, key, p.name, "Append()"))
		}

		item[key] = v

	}

	if len(item) == 0 {
		return nil, errors.New(NotEnoughArguments(r, p.name, "Append()"))
	}

	res, err := c.Append(item)
	return res, g.DbError(409, r, err)
}

func (p *CurdCond) Find(r *http.Request, offset uint) (db.Item, error) {
	c, err := g.Db.Collection(p.name)
	if err != nil {
		return nil, err
	}
	f := p.Q
	vs := r.Form
	for _, key := range f.Exists {
		_, ok := vs[key]
		if !ok {
			return nil, errors.New(MustBeExists(r, key, p.name, "Find()"))
		}
	}
	and := db.And{}
	for _, key := range f.Fields {
		_, ok := vs[key]
		if !ok {
			continue
		}
		v := vs.Get(key)
		if !meta.Valid.IsValid(key, v) {
			return nil, errors.New(InValidField(r, key, p.name, "Find()"))
		}

		and = append(and, db.Cond{key: v})

	}

	if len(and) == 0 {
		return nil, errors.New(NotEnoughArguments(r, p.name, "Find()"))
	}

	res, err := c.Find(and, db.Offset(offset))
	if err == db.ErrNoMoreRows {
		return res, nil
	}
	return res, g.DbError(409, r, err)
}
func (p *CurdCond) FindAll(r *http.Request, offset, limit uint) ([]db.Item, error) {
	c, err := g.Db.Collection(p.name)
	if err != nil {
		return nil, err
	}
	f := p.Q
	vs := r.Form
	for _, key := range f.Exists {
		_, ok := vs[key]
		if !ok {
			return nil, errors.New(MustBeExists(r, key, p.name, "FindAll()"))
		}
	}
	and := db.And{}
	for _, key := range f.Fields {
		_, ok := vs[key]
		if !ok {
			continue
		}
		v := vs.Get(key)
		if !meta.Valid.IsValid(key, v) {
			return nil, errors.New(InValidField(r, key, p.name, "FindAll()"))
		}

		and = append(and, db.Cond{key: v})

	}

	if len(and) == 0 {
		return nil, errors.New(NotEnoughArguments(r, p.name, "FindAll()"))
	}

	res, err := c.FindAll(and, db.Offset(offset), db.Limit(limit))
	if err == db.ErrNoMoreRows {
		return res, nil
	}
	return res, g.DbError(409, r, err)
}
func (p *CurdCond) Query(r *http.Request, offset, limit uint) (db.Result, error) {
	c, err := g.Db.Collection(p.name)
	if err != nil {
		return nil, err
	}
	f := p.Q
	vs := r.Form
	for _, key := range f.Exists {
		_, ok := vs[key]
		if !ok {
			return nil, errors.New(MustBeExists(r, key, p.name, "Query()"))
		}
	}
	and := db.And{}
	for _, key := range f.Fields {
		_, ok := vs[key]
		if !ok {
			continue
		}
		v := vs.Get(key)
		if !meta.Valid.IsValid(key, v) {
			return nil, errors.New(InValidField(r, key, p.name, "Query()"))
		}

		and = append(and, db.Cond{key: v})

	}

	if len(and) == 0 {
		return nil, errors.New(NotEnoughArguments(r, p.name, "Query()"))
	}

	res, err := c.Query(and, db.Offset(offset), db.Limit(limit))
	if err == db.ErrNoMoreRows {
		return res, nil
	}
	return res, g.DbError(409, r, err)
}

func (p *CurdCond) Update(r *http.Request) error {
	c, err := g.Db.Collection(p.name)
	if err != nil {
		return err
	}
	f := p.U
	vs := r.Form

	and := db.And{}
	for _, key := range f.Exists {
		_, ok := vs[key]
		if !ok {
			return errors.New(MustBeExists(r, key, p.name, "Update()"))
		}
		v := vs.Get(key)
		if !meta.Valid.IsValid(key, v) {
			return errors.New(InValidField(r, key, p.name, "Update()"))
		}

		and = append(and, db.Cond{key: v})
	}
	if len(and) == 0 {
		return errors.New(NotEnoughArguments(r, p.name, "Update()"))
	}

	set := db.Set{}
	for _, key := range f.Fields {
		_, ok := vs[key]
		if !ok {
			continue
		}
		v := vs.Get(key)
		if !meta.Valid.IsValid(key, v) {
			return errors.New(InValidField(r, key, p.name, "Update()"))
		}
		set[key] = v
	}

	if len(set) == 0 {
		return errors.New(NotEnoughArguments(r, p.name, "Update()"))
	}

	return g.DbError(409, r, c.Update(and, set))
}
