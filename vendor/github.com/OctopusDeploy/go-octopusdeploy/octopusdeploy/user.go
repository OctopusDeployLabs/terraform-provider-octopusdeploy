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

type Users struct {
	Items []User `json:"Items"`
	PagedResults
}

type User struct {
	ID                         string `json:"Id"`
	UserName                   string `json:"UserName"`
	DisplayName                string `json:"DisplayName"`
	SortOrder                  int    `json:"SortOrder"`
	UseGuidedFailure           bool   `json:"UseGuidedFailure"`
	AllowDynamicInfrastructure bool   `json:"AllowDynamicInfrastructure"`
}

func (t *User) Validate() error {
	validate := validator.New()

	err := validate.Struct(t)

	if err != nil {
		return err
	}

	return nil
}

func NewUser(UserName, DisplayName string, useguidedfailure bool) *User {
	return &User{
		UserName:             UserName,
		DisplayName:          DisplayName,
		UseGuidedFailure: useguidedfailure,
	}
}

func (s *UserService) Get(Userid string) (*User, error) {
	path := fmt.Sprintf("Users/%s", Userid)
	resp, err := apiGet(s.sling, new(User), path)

	if err != nil {
		return nil, err
	}

	return resp.(*User), nil
}

func (s *UserService) GetAll() (*[]User, error) {
	var p []User

	path := "Users"

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
		if project.Name == UserName {
			return &project, nil
		}
	}

	return &foundUser, fmt.Errorf("no User found with User name %s", UserName)
}

func (s *UserService) Add(User *User) (*User, error) {
	resp, err := apiAdd(s.sling, User, new(User), "Users")

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

func (s *UserService) Update(User *User) (*User, error) {
	path := fmt.Sprintf("Users/%s", User.ID)
	resp, err := apiUpdate(s.sling, User, new(User), path)

	if err != nil {
		return nil, err
	}

	return resp.(*User), nil
}
