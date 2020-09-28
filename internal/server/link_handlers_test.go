package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/sergeychur/avito_auto/internal/config"
	"github.com/sergeychur/avito_auto/internal/mocks"
	"github.com/sergeychur/avito_auto/internal/models"
	"github.com/sergeychur/avito_auto/internal/repository"
	"github.com/sergeychur/avito_auto/internal/server"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLinkHandlerSuite(t *testing.T) {
	suite.Run(t, new(LinkHandlerTestSuite))
}

type LinkHandlerTestSuite struct {
	suite.Suite
	repo   *mocks.MockRepository
	validator *mocks.MockIValidator
	underTest  *server.Server
}

func (suite *LinkHandlerTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()
	suite.repo = mocks.NewMockRepository(mockCtrl)
	suite.validator = mocks.NewMockIValidator(mockCtrl)
	serv := new(server.Server)
	conf := new(config.Config)
	conf.Port = ":8090"
	serv.Config = conf
	serv.Repo = suite.repo
	serv.Validator = suite.validator
	suite.underTest = serv
}

type TestGetLinkCase struct {
	link        string
	shortcut string
	response string
	status	int
	method   string
}

func (suite *LinkHandlerTestSuite) TestGetLink() {
	cases := []TestGetLinkCase{
		{
			link: "https://google.com",
			response: "303 See Other",
			shortcut: "google",
			status: repository.OK,
			method: http.MethodGet,
		},
		{
			link: "https://google.com",
			response: "404 Not Found",
			shortcut: "goo",
			status: repository.EMPTY_RESULT,
			method: http.MethodGet,
		},
	}

	for _, item := range cases {
		suite.repo.EXPECT().GetLink(item.shortcut).Return(item.status, item.link)
		r, _ := http.NewRequest(item.method, fmt.Sprintf("/link/%s", item.shortcut), nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("shortcut", item.shortcut)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		suite.underTest.GetLink(w, r)

		response := w.Result()
		suite.Equal(response.Status, item.response)
		if item.response == http.StatusText(http.StatusSeeOther) {
			suite.Equal(item.link, response.Header.Get("Location"))
		}
	}
}

type TestCreateLinkCase struct {
	reqLink  models.Link
	respLink models.Link
	response string
	status int
	method string
	validatorError error
	readFromBodyError error
}

func (suite *LinkHandlerTestSuite) TestCreateLink() {
	cases := []TestCreateLinkCase{
		{
			reqLink: models.Link{
				RealURL:  "https://google.com",
			},
			respLink: models.Link{
				RealURL:  "https://google.com",
				Shortcut: "hagagag",
			},
			response: "201 Created",
			status: repository.CREATED,
			method: http.MethodPost,
			validatorError: nil,
		},

		{
			reqLink: models.Link{
				RealURL:  "https://google.com",
				Shortcut: "google",
			},
			respLink: models.Link{
				RealURL:  "https://google.com",
				Shortcut: "google",
			},
			response: "201 Created",
			status: repository.CREATED,
			method: http.MethodPost,
			validatorError: nil,
			readFromBodyError: nil,
		},

		{
			reqLink: models.Link{
				RealURL:  "https://google.com",
				Shortcut: "google",
			},
			respLink: models.Link{
				RealURL:  "https://google.com",
				Shortcut: "google",
			},
			response: "403 Forbidden",
			status: repository.FORBIDDEN,
			method: http.MethodPost,
			validatorError: nil,
			readFromBodyError: nil,
		},
		{
			reqLink: models.Link{},
			respLink: models.Link{},
			response: "400 Bad Request",
			status: repository.CREATED,
			method: http.MethodPost,
			validatorError: nil,
			readFromBodyError: errors.New("cannot read"),
		},

		{
			reqLink: models.Link{
				RealURL:  "ht/google.com",
				Shortcut: "google",
			},
			respLink: models.Link{
				RealURL:  "https://google.com",
				Shortcut: "google",
			},
			response: "400 Bad Request",
			status: repository.FORBIDDEN,
			method: http.MethodPost,
			validatorError: errors.New("url is invalid"),
			readFromBodyError: nil,
		},
	}

	for _, item := range cases {
		var body []byte
		if item.readFromBodyError == nil {
			suite.repo.EXPECT().InsertLink(item.reqLink).Return(item.status, item.respLink)
			suite.validator.EXPECT().ValidateLink(item.reqLink).Return(item.validatorError)
			body, _ = json.Marshal(item.reqLink)
		}

		r, _ := http.NewRequest(item.method, "/link", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		suite.underTest.CreateLink(w, r)

		response := w.Result()
		suite.Equal(response.Status, item.response)
		if item.response == http.StatusText(http.StatusCreated) {
			result := models.Link{}

			_ = json.NewDecoder(response.Body).Decode(&result)
			_ = response.Body.Close()
			suite.Equal(item.respLink.Shortcut, result)
		}
	}
}
