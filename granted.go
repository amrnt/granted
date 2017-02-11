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
	lock    *sync.RWMutex
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

// DefaultConfig ...
var DefaultConfig = NewConfig()

// Default ...
var Default = New(DefaultConfig)

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		lock:    &sync.RWMutex{},
		TagName: "grant",
		Roles:   make(map[string]Role),
	}
}

// New ...
func New(c *Config) *Authorize {
	a := Authorize{c}

	// Set default roles
	a.addRole(Role{
		Name: "ALL",
		CheckFunc: func(i Instance) bool {
			return true
		},
	})

	a.addRole(Role{
		Name: "NONE",
		CheckFunc: func(i Instance) bool {
			return false
		},
	})

	return &a
}

// AddRole ...
func AddRole(r Role) {
	Default.addRole(r)
}

func (a *Authorize) addRole(r Role) {
	a.Config.lock.Lock()
	defer a.Config.lock.Unlock()
	a.Config.Roles[r.Name] = r
}

func (a *Authorize) canAccess(v reflect.Value, t string) bool {
	a.Config.lock.RLock()
	defer a.Config.lock.RUnlock()

	i, ok := a.Config.Roles[t]
	if !ok {
		return false
	}

	return i.CheckFunc(v.Interface().(Instance))
}
