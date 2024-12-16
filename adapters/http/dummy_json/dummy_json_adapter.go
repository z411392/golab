package dummy_json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DummyJsonAdapter struct {
	httpClient *http.Client
}

const baseURL = "https://dummyjson.com"

func NewDummyJsonAdapter() *DummyJsonAdapter {
	adapter := &DummyJsonAdapter{
		httpClient: &http.Client{},
	}
	return adapter
}

func (adapter *DummyJsonAdapter) SignIn(username string, password string) (*Credentials, error) {
	const URI = "/auth/login"
	requestJsonString, err := json.Marshal(map[string]interface{}{
		"username":      username,
		"password":      password,
		"expiresInMins": 60,
	})
	if err != nil {
		return nil, err
	}
	requestBody := io.NopCloser(bytes.NewBuffer(requestJsonString))
	request, err := http.NewRequest("POST", fmt.Sprintf("%s%s", baseURL, URI), requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := adapter.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseJsonString, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	credentials := &Credentials{}
	err = json.Unmarshal(responseJsonString, &credentials)
	if err != nil {
		return nil, err
	}
	if credentials.AccessToken == "" {
		return nil, nil
	}
	return credentials, err
}

func (adapter *DummyJsonAdapter) GetAuthUser(token string) (*User, error) {
	const URI = "/auth/me"
	request, err := http.NewRequest("GET", fmt.Sprintf("%s%s", baseURL, URI), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := adapter.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseJsonString, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = json.Unmarshal(responseJsonString, &user)
	if user.Id == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, err
}
