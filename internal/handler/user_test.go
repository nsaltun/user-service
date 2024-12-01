package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/nsaltun/userapi/internal/mocks/service"
	"github.com/nsaltun/userapi/internal/model"
	"github.com/nsaltun/userapi/pkg/lib/errwrap"
	"github.com/nsaltun/userapi/pkg/lib/middleware"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	noSetup := func(*mocks.UserService, interface{}) {}
	tests := []struct {
		name        string
		req         interface{}
		setup       func(*mocks.UserService, interface{})
		assertError require.ValueAssertionFunc
		assertResp  require.ValueAssertionFunc
	}{
		{
			name: "service returns success",
			req:  &model.User{FirstName: "t_firstName", NickName: "t_nickname", Email: "t@email.com", Country: "CA"},
			setup: func(s *mocks.UserService, u interface{}) {
				s.On("CreateUser", mock.Anything, u).Return(u, nil).Once()
			},
			assertError: require.NotNil,
			assertResp: func(tt require.TestingT, respBytes interface{}, statusCode ...interface{}) {
				require.Contains(tt, statusCode, 201)
				res, ok := respBytes.([]byte)
				require.True(tt, ok)

				var actual model.User
				require.NoError(tt, json.Unmarshal(res, &actual))
				expected := model.User{FirstName: "t_firstName", NickName: "t_nickname", Email: "t@email.com", Country: "CA"}
				require.Equal(tt, expected, actual)
			},
		},
		{
			name:  "invalid json",
			req:   "test",
			setup: noSetup,
			assertError: func(tt require.TestingT, respBytes interface{}, statusCode ...interface{}) {
				res, ok := respBytes.([]byte)
				require.True(tt, ok)

				var actual errwrap.ErrorResponse
				require.NoError(tt, json.Unmarshal(res, &actual))
				expected := errwrap.ErrorResponse{
					Message: "invalid payload format",
					Code:    "400"}
				require.Equal(tt, expected, actual)
			},
			assertResp: require.NotNil, //just to ignore
		},
		{
			name:  "invalid request: empty user",
			req:   &model.Meta{},
			setup: noSetup,
			assertError: func(tt require.TestingT, respBytes interface{}, statusCode ...interface{}) {
				res, ok := respBytes.([]byte)
				require.True(tt, ok)

				var actual errwrap.ErrorResponse
				require.NoError(tt, json.Unmarshal(res, &actual))
				expected := errwrap.ErrorResponse{
					Message: "firstName can't be empty;;email can't be empty;;nickName can't be empty;;country can't be empty",
					Code:    "400"}
				require.Equal(tt, expected, actual)
			},
			assertResp: require.NotNil, //just to ignore
		},
		{
			name:  "invalid request: empty email",
			req:   &model.User{FirstName: "t_firstName", NickName: "t_nickname", Country: "CA"},
			setup: noSetup,
			assertError: func(tt require.TestingT, respBytes interface{}, _ ...interface{}) {
				res, ok := respBytes.([]byte)
				require.True(tt, ok)

				var actual errwrap.ErrorResponse
				require.NoError(tt, json.Unmarshal(res, &actual))
				expected := errwrap.ErrorResponse{
					Message: "email can't be empty",
					Code:    "400"}
				require.Equal(tt, expected, actual)
			},
			assertResp: require.NotNil, //just to ignore
		},
		{
			name: "service returns error - conflict error",
			req:  &model.User{FirstName: "t_firstName", NickName: "t_nickname", Email: "t@email.com", Country: "CA"},
			setup: func(s *mocks.UserService, u interface{}) {
				s.On("CreateUser", mock.Anything, u).Return(nil, errwrap.ErrConflict.SetMessage("mock service err")).Once()
			},
			assertError: func(tt require.TestingT, respBytes interface{}, _ ...interface{}) {
				res, ok := respBytes.([]byte)
				require.True(tt, ok)

				var actual errwrap.ErrorResponse
				require.NoError(tt, json.Unmarshal(res, &actual))
				expected := errwrap.ErrorResponse{
					Message: "mock service err",
					Code:    "409"}
				require.Equal(tt, expected, actual)
			},
			assertResp: require.NotNil, //just to ignore
		},
		{
			name: "service returns error - unknown error type",
			req:  &model.User{FirstName: "t_firstName", NickName: "t_nickname", Email: "t@email.com", Country: "CA"},
			setup: func(s *mocks.UserService, u interface{}) {
				s.On("CreateUser", mock.Anything, u).Return(nil, errors.New("unexpected error")).Once()
			},
			assertError: func(tt require.TestingT, respBytes interface{}, _ ...interface{}) {
				res, ok := respBytes.([]byte)
				require.True(tt, ok)

				var actual errwrap.ErrorResponse
				require.NoError(tt, json.Unmarshal(res, &actual))
				expected := errwrap.ErrorResponse{
					Message: "internal server error",
					Code:    "500"}
				require.Equal(tt, expected, actual)
			},
			assertResp: require.NotNil, //just to ignore
		},
	}
	for _, tCase := range tests {
		t.Run(tCase.name, func(tt *testing.T) {
			//setup
			userSvcMock := mocks.NewUserService(t)
			h := NewUserHandler(userSvcMock)
			tCase.setup(userSvcMock, tCase.req)

			// prepare httptest
			jsonData, err := json.Marshal(tCase.req)
			require.NoError(tt, err)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
			w := httptest.NewRecorder()

			//execute
			httpCtx := &middleware.HttpContext{Response: w, Request: req}
			h.CreateUser(httpCtx)
			b, err := io.ReadAll(w.Body)

			//assert
			require.NoError(tt, err)
			tCase.assertResp(tt, b, w.Result().StatusCode)
			tCase.assertError(tt, b, w.Result().StatusCode)
		})
	}
}

func TestUpdateUserById(t *testing.T) {
	//TODO
}

func TestDeleteUserById(t *testing.T) {
	//TODO
}

func TestListUsers(t *testing.T) {
	//TODO
}
