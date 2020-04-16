package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"simplon.biz/corona/pkg/mocks"
)

func TestRequestCodes(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logEntry := logger.WithField("test", true)
	mockCtrl := gomock.NewController(t)
	testingAuthTokenManager := mocks.NewMockTokenManager(mockCtrl)
	testingAuthTokenManager.EXPECT().StoreToken(gomock.Any()).Times(10)

	app := Application{
		testingAuthTokenManager: testingAuthTokenManager,
	}

	makeRequest := func(count int) *http.Request {
		data := fmt.Sprintf(`{"count": %d}`, count)
		r := httptest.NewRequest("POST", "/api/request-codes", strings.NewReader(data))
		r = r.WithContext(context.WithValue(r.Context(), ctxLog, logEntry))
		return r
	}

	recorder := httptest.NewRecorder()
	app.requestCodes(recorder, makeRequest(0))
	response := recorder.Result()
	if response.StatusCode != http.StatusBadRequest {
		t.Error("Request for 0 codes did not fail")
	}

	recorder = httptest.NewRecorder()
	app.requestCodes(recorder, makeRequest(101))
	response = recorder.Result()
	if response.StatusCode != http.StatusBadRequest {
		t.Error("Request for 0 codes did not fail")
	}

	recorder = httptest.NewRecorder()
	app.requestCodes(recorder, makeRequest(10))
	response = recorder.Result()
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Correct request returned wrong status code: %d", response.StatusCode)
	}
}
