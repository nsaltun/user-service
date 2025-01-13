package user

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/nsaltun/userapi/internal/mocks/service"
	"github.com/nsaltun/userapi/internal/model"
	"github.com/nsaltun/userapi/pkg/lib/errwrap"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	// noSetup := func(*mocks.UserService, interface{}) {}
	tests := []struct {
		name        string
		req         interface{}
		setup       func(*mocks.UserService, interface{})
		assertError require.ValueAssertionFunc
		assertResp  require.ValueAssertionFunc
	}{
		{
			name: "service returns success",
			req:  &CreateUserRequest{&model.User{FirstName: "t_firstName", NickName: "t_nickname", Email: "t@email.com", Country: "CA"}},
			setup: func(s *mocks.UserService, u interface{}) {
				req, _ := u.(*CreateUserRequest)
				s.On("CreateUser", mock.Anything, req.User).Return(req.User, nil).Once()
			},
			assertError: require.Nil,
			assertResp: func(tt require.TestingT, resp interface{}, statusCode ...interface{}) {
				require.Contains(tt, statusCode, 201)
				actual, ok := resp.(*CreateUserResponse)
				require.True(tt, ok)

				expected := &CreateUserResponse{&model.User{FirstName: "t_firstName", NickName: "t_nickname", Email: "t@email.com", Country: "CA"}}
				require.Equal(tt, expected, actual)
			},
		},
		{
			name: "service returns error - conflict error",
			req:  &CreateUserRequest{&model.User{FirstName: "t_firstName", NickName: "t_nickname", Email: "t@email.com", Country: "CA"}},
			setup: func(s *mocks.UserService, u interface{}) {
				req, _ := u.(*CreateUserRequest)
				s.On("CreateUser", mock.Anything, req.User).Return(nil, errwrap.ErrConflict.SetMessage("mock service err")).Once()
			},
			assertError: func(tt require.TestingT, err interface{}, statusCode ...interface{}) {
				actualErr, ok := err.(error)
				require.True(tt, ok)
				expected := errwrap.ErrConflict.SetMessage("mock service err")
				require.Equal(tt, expected, actualErr)
			},
			assertResp: require.Nil, //just to ignore
		},
		{
			name: "service returns error - unknown error type",
			req:  &CreateUserRequest{&model.User{FirstName: "t_firstName", NickName: "t_nickname", Email: "t@email.com", Country: "CA"}},
			setup: func(s *mocks.UserService, u interface{}) {
				req, _ := u.(*CreateUserRequest)
				s.On("CreateUser", mock.Anything, req.User).Return(nil, errors.New("unexpected error")).Once()
			},
			assertError: func(tt require.TestingT, err interface{}, _ ...interface{}) {
				actualErr, ok := err.(error)
				require.True(tt, ok)

				expected := errors.New("unexpected error")
				require.Equal(tt, expected, actualErr)
			},
			assertResp: require.Nil, //just to ignore
		},
	}
	for _, tCase := range tests {
		t.Run(tCase.name, func(tt *testing.T) {
			//setup
			userSvcMock := mocks.NewUserService(t)
			h := NewUserHandler(userSvcMock)
			tCase.setup(userSvcMock, tCase.req)

			//execute
			resp, statusCode, err := h.CreateUser(context.Background(), tCase.req.(*CreateUserRequest))

			//assert
			tCase.assertResp(tt, resp, statusCode)
			tCase.assertError(tt, err, statusCode)
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
