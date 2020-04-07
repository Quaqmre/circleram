package auth

import (
	"errors"
	"testing"

	"github.com/Quaqmre/circleramkit/user"
)

//Mock User Service
type mockUserService struct {
	storeMetot       bool
	listMetot        bool
	findWithuserName bool
}

func (s *mockUserService) Find(id user.UserID) (*user.User, error) {
	return nil, errors.New("test error")
}
func (s *mockUserService) Store(c *user.User) error {
	s.storeMetot = true
	return nil
}
func (s *mockUserService) List() ([]user.User, error) {
	s.listMetot = true
	return []user.User{user.User{Name: "akif"}}, nil
}
func (s *mockUserService) FindWithUserName(string) (*user.User, error) {
	s.findWithuserName = true
	return &user.User{Name: "akif", Password: "123"}, nil
}

var mockedUserService user.Service = &mockUserService{}
var secretKey string = "test"

var invalidtoken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im9zcGlyIiwiZXhwIjoxNTg2MjY2MzU3fQ.ugllsXCulj5Us0CNhIbEn3fsWVeo35ZdgKiKA0ZnA7w"

var validtoken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFraWYiLCJleHAiOjE1ODYyNzA1MDZ9.Uu_IiiN-IdpClIE8lMuF2BuTCUKfBLHn5bU4O023kvU"

var srv Service = NewAuthService(secretKey, mockedUserService)

func Test_GetUser_with_Invalid_Token_Return_Error(t *testing.T) {
	u, err := srv.GetUser(invalidtoken)
	if err == nil {
		t.Error("expecting err but returning nil error")
	}
	if u != "" {
		t.Error("returned string should be empty string")
	}
}
func Test_GetUser_with_Valid_Token(t *testing.T) {
	u, err := srv.GetUser(validtoken)

	if u != "akif" {
		t.Error("expected akif name but returned wrong")
	}

	if err != nil {
		t.Error("Valid token returned error")
	}

}
