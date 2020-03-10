package customer

import (
	"errors"
	"sync"
)

// Service is a simple CRUD interface for customers
type Service interface {
	Store(c *Customer) error
	Find(id CustomerID) (*Customer, error)
}

//CustomerID exported for access directly other package
type CustomerID int

//Customer store user information
type Customer struct {
	CustomerID CustomerID
	Name       string
	Email      string
}

// ErrUnknown is used when a cargı could not be found
var ErrUnknown = errors.New("unknown cargo")

type inmemService struct {
	mtx       sync.RWMutex
	customers map[CustomerID]*Customer
}

//NewCustomerService export for using inmem repo the other package
func NewCustomerService() Service {
	return &inmemService{
		customers: make(map[CustomerID]*Customer),
	}
}

func (r *inmemService) Store(c *Customer) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	c.CustomerID = CustomerID(len(r.customers))
	r.customers[c.CustomerID] = c
	return nil
}
func (r *inmemService) Find(id CustomerID) (*Customer, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.customers[id]; ok {
		return val, nil
	}
	return nil, ErrUnknown
}
