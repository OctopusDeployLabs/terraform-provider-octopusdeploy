package octopusdeploy

import (
	"fmt"

	"github.com/dghubble/sling"
	"gopkg.in/go-playground/validator.v9"
)

type UserService struct {
	sling *sling.Sling
}

func NewUserService(sling *sling.Sling) *UserService {
	return &UserService{
		sling: sling,
	}
}

type Users struct {
	Items []User `json:"Items"`
	PagedResults
}

type User struct {
	ID                  string `json:"Id"`
	Username            string `json:"Username"`
	DisplayName         string `json:"DisplayName"`
	IsActive            bool   `json:"IsActive"`
	IsService           bool   `json:"IsService"`
	EmailAddress        string `json:"EmailAddress"`
	CanPasswordBeEdited bool   `json:"CanPasswordBeEdited"`
	IsRequestor         bool   `json:"IsRequestor"`
	Links               struct {
		Self        string `json:"Self"`
		Permissions string `json:"Permissions"`
		APIKeys     string `json:"ApiKeys"`
		Avatar      string `json:"Avatar"`
	} `json:"Links"`
}

func (t *User) Validate() error {
	validate := validator.New()

	err := validate.Struct(t)

	if err != nil {
		return err
	}

	return nil
}

func NewUser(Username, DisplayName string) *User {
	return &User{
		Username:    Username,
		DisplayName: DisplayName,
	}
}

func (s *UserService) Get(UserID string) (*User, error) {
	path := fmt.Sprintf("Users/%s", UserID)
	resp, err := apiGet(s.sling, new(User), path)

	if err != nil {
		return nil, err
	}

	return resp.(*User), nil
}

func (s *UserService) GetAll() (*[]User, error) {
	var p []User

	path := "users"

	loadNextPage := true

	for loadNextPage {
		resp, err := apiGet(s.sling, new(Users), path)

		if err != nil {
			return nil, err
		}

		r := resp.(*Users)

		for _, item := range r.Items {
			p = append(p, item)
		}

		path, loadNextPage = LoadNextPage(r.PagedResults)
	}

	return &p, nil
}

func (s *UserService) GetByName(UserName string) (*User, error) {
	var foundUser User
	Users, err := s.GetAll()

	if err != nil {
		return nil, err
	}

	for _, project := range *Users {
		if project.Username == UserName {
			return &project, nil
		}
	}

	return &foundUser, fmt.Errorf("no User found with User name %s", UserName)
}

func (s *UserService) Add(user *User) (*User, error) {
	resp, err := apiAdd(s.sling, user, new(User), "users")

	if err != nil {
		return nil, err
	}

	return resp.(*User), nil
}

func (s *UserService) Delete(Userid string) error {
	path := fmt.Sprintf("Users/%s", Userid)
	err := apiDelete(s.sling, path)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Update(user *User) (*User, error) {
	path := fmt.Sprintf("Users/%s", user.ID)
	resp, err := apiUpdate(s.sling, user, new(user), path)

	if err != nil {
		return nil, err
	}

	return resp.(*User), nil
}
