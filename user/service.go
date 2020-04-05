package user

import (
	"errors"
	"sync"
)

// Service is a simple CRUD interface for users
type Service interface {
	Store(c *User) error
	Find(id UserID) (*User, error)
	List() ([]User, error)
	FindWithUserName(string) (*User, error)
}

//UserID exported for access directly other package
type UserID int

//User store user information
type User struct {
	UserID   UserID
	Name     string
	Password string
}

// ErrUnknown is used when a cargı could not be found
var ErrUnknown = errors.New("unknown user")

type inmemService struct {
	mtx   sync.RWMutex
	users map[UserID]*User
}

//NewUserService export for using inmem repo the other package
func NewUserService() Service {
	return &inmemService{
		users: make(map[UserID]*User),
	}
}

func (r *inmemService) Store(c *User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	c.UserID = UserID(len(r.users))
	r.users[c.UserID] = c
	return nil
}
func (r *inmemService) Find(id UserID) (*User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.users[id]; ok {
		return val, nil
	}
	return nil, ErrUnknown
}

func (r *inmemService) List() ([]User, error) {
	list := []User{}
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	for _, i := range r.users {
		list = append(list, *i)
	}
	return list, nil
}

// FindWithUserName created for token service
func (r *inmemService) FindWithUserName(name string) (*User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	for _, u := range r.users {
		if u.Name == name {
			return u, nil
		}
	}
	return nil, ErrUnknown
}
