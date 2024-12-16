package e2e_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	dummyJson "github.com/z411392/golab/adapters/http/dummy_json"
)

func Test_ShouldRetrieveProfile(t *testing.T) {
	t.SkipNow()
	request, _ := http.NewRequest("GET", "/auth/me", nil)
	token := os.Getenv("TOKEN")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	responseJsonString, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	type ResponseData = struct {
		Payload struct {
			User dummyJson.User `json:"user"`
		} `json:"payload"`
	}
	responseData := &ResponseData{}
	err = json.Unmarshal(responseJsonString, &responseData)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if responseData.Payload.User.Id == 0 {
		t.FailNow()
	}
}
