package granted

import (
	"reflect"
	"sync"
)

// Authorize ..
type Authorize struct {
	*Config
}

// Config ...
type Config struct {
	lock    sync.RWMutex
	TagName string
	Roles   map[string]Role
}

// Role ...
type Role struct {
	Name      string
	CheckFunc func(i Instance) bool
}

// Instance ...
type Instance interface{}

// New ...
func New(c *Config) *Authorize {
	if len(c.TagName) == 0 {
		c.TagName = "grant"
	}

	if len(c.Roles) == 0 {
		c.Roles = make(map[string]Role)
	}

	a := Authorize{
		Config: c,
	}

	// Set default roles
	a.AddRole(Role{
		Name: "ALL",
		CheckFunc: func(i Instance) bool {
			return true
		},
	})

	a.AddRole(Role{
		Name: "NONE",
		CheckFunc: func(i Instance) bool {
			return false
		},
	})

	return &a
}

// AddRole ...
func (a *Authorize) AddRole(r Role) {
	a.Config.lock.Lock()
	a.Config.Roles[r.Name] = r
	a.Config.lock.Unlock()
}

func (a *Authorize) hasAccess(v reflect.Value, t string) bool {
	a.Config.lock.Lock()
	defer a.Config.lock.Unlock()
	return a.Config.Roles[t].CheckFunc(v.Interface().(Instance))
}
