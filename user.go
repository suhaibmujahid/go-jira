package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// UserService handles users for the JIRA instance / API.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user
type UserService struct {
	client *Client
}

// User represents a JIRA user.
type User struct {
	Self        string `json:"self,omitempty" structs:"self,omitempty"`
	AccountID   string `json:"accountId,omitempty" structs:"accountId,omitempty"`
	AccountType string `json:"accountType,omitempty" structs:"accountType,omitempty"`
	// TODO: name & key are deprecated, see:
	// https://developer.atlassian.com/cloud/jira/platform/api-changes-for-user-privacy-announcement/
	Name            string     `json:"name,omitempty" structs:"name,omitempty"`
	Key             string     `json:"key,omitempty" structs:"key,omitempty"`
	Password        string     `json:"-"`
	EmailAddress    string     `json:"emailAddress,omitempty" structs:"emailAddress,omitempty"`
	AvatarUrls      AvatarUrls `json:"avatarUrls,omitempty" structs:"avatarUrls,omitempty"`
	DisplayName     string     `json:"displayName,omitempty" structs:"displayName,omitempty"`
	Active          bool       `json:"active,omitempty" structs:"active,omitempty"`
	TimeZone        string     `json:"timeZone,omitempty" structs:"timeZone,omitempty"`
	Locale          string     `json:"locale,omitempty" structs:"locale,omitempty"`
	ApplicationKeys []string   `json:"applicationKeys,omitempty" structs:"applicationKeys,omitempty"`
}

// UserGroup represents the group list
type UserGroup struct {
	Self string `json:"self,omitempty" structs:"self,omitempty"`
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}

type userSearchParam struct {
	name  string
	value string
}

type userSearch []userSearchParam

type userSearchF func(userSearch) userSearch

// GetWithContext gets user info from JIRA
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-getUser
//
// /!\ Deprecation notice: https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-user-privacy-api-migration-guide/
// By 29 April 2019, we will remove personal data from the API that is used to identify users,
// such as username and userKey, and instead use the Atlassian account ID (accountId).
func (s *UserService) GetWithContext(ctx context.Context, username string) (*User, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user?username=%s", username)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return user, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *UserService) Get(username string) (*User, *Response, error) {
	return s.GetWithContext(context.Background(), username)
}

// GetByAccountIDWithContext gets user info from JIRA
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-getUser
func (s *UserService) GetByAccountIDWithContext(ctx context.Context, accountID string) (*User, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user?accountId=%s", accountID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return user, resp, nil
}

// GetByAccountID wraps GetByAccountIDWithContext using the background context.
func (s *UserService) GetByAccountID(accountID string) (*User, *Response, error) {
	return s.GetByAccountIDWithContext(context.Background(), accountID)
}

// CreateWithContext creates an user in JIRA.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-createUser
func (s *UserService) CreateWithContext(ctx context.Context, user *User) (*User, *Response, error) {
	apiEndpoint := "/rest/api/2/user"
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, user)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	responseUser := new(User)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := fmt.Errorf("could not read the returned data")
		return nil, resp, NewJiraError(resp, e)
	}
	err = json.Unmarshal(data, responseUser)
	if err != nil {
		e := fmt.Errorf("could not unmarshall the data into struct")
		return nil, resp, NewJiraError(resp, e)
	}
	return responseUser, resp, nil
}

// Create wraps CreateWithContext using the background context.
func (s *UserService) Create(user *User) (*User, *Response, error) {
	return s.CreateWithContext(context.Background(), user)
}

// DeleteWithContext deletes an user from JIRA.
// Returns http.StatusNoContent on success.
//
// JIRA API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-api-2-user-delete
func (s *UserService) DeleteWithContext(ctx context.Context, username string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user?username=%s", username)
	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, NewJiraError(resp, err)
	}
	return resp, nil
}

// Delete wraps DeleteWithContext using the background context.
func (s *UserService) Delete(username string) (*Response, error) {
	return s.DeleteWithContext(context.Background(), username)
}

// GetGroupsWithContext returns the groups which the user belongs to
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-getUserGroups
func (s *UserService) GetGroupsWithContext(ctx context.Context, username string) (*[]UserGroup, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user/groups?username=%s", username)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	userGroups := new([]UserGroup)
	resp, err := s.client.Do(req, userGroups)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return userGroups, resp, nil
}

// GetGroups wraps GetGroupsWithContext using the background context.
func (s *UserService) GetGroups(username string) (*[]UserGroup, *Response, error) {
	return s.GetGroupsWithContext(context.Background(), username)
}

// GetSelfWithContext information about the current logged-in user
//
// JIRA API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-api-2-myself-get
func (s *UserService) GetSelfWithContext(ctx context.Context) (*User, *Response, error) {
	const apiEndpoint = "rest/api/2/myself"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	var user User
	resp, err := s.client.Do(req, &user)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return &user, resp, nil
}

// GetSelf wraps GetSelfWithContext using the background context.
func (s *UserService) GetSelf() (*User, *Response, error) {
	return s.GetSelfWithContext(context.Background())
}

// WithMaxResults sets the max results to return
func WithMaxResults(maxResults int) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "maxResults", value: fmt.Sprintf("%d", maxResults)})
		return s
	}
}

// WithStartAt set the start pager
func WithStartAt(startAt int) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "startAt", value: fmt.Sprintf("%d", startAt)})
		return s
	}
}

// WithActive sets the active users lookup
func WithActive(active bool) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "includeActive", value: fmt.Sprintf("%t", active)})
		return s
	}
}

// WithInactive sets the inactive users lookup
func WithInactive(inactive bool) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "includeInactive", value: fmt.Sprintf("%t", inactive)})
		return s
	}
}

// FindWithContext searches for user info from JIRA:
// It can find users by email, username or name
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-findUsers
func (s *UserService) FindWithContext(ctx context.Context, property string, tweaks ...userSearchF) ([]User, *Response, error) {
	search := []userSearchParam{
		{
			name:  "username",
			value: property,
		},
	}
	for _, f := range tweaks {
		search = f(search)
	}

	var queryString = ""
	for _, param := range search {
		queryString += param.name + "=" + param.value + "&"
	}

	apiEndpoint := fmt.Sprintf("/rest/api/2/user/search?%s", queryString[:len(queryString)-1])
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	users := []User{}
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return users, resp, nil
}

// Find wraps FindWithContext using the background context.
func (s *UserService) Find(property string, tweaks ...userSearchF) ([]User, *Response, error) {
	return s.FindWithContext(context.Background(), property, tweaks...)
}
