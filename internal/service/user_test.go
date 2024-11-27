package service

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	repomocks "github.com/nsaltun/userapi/internal/mocks/repository"
	"github.com/nsaltun/userapi/internal/model"
	"github.com/nsaltun/userapi/pkg/lib/errwrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	type setupFunc func(*repomocks.UserRepository, *model.User)
	noSetup := func(*repomocks.UserRepository, *model.User) {}

	tests := []struct {
		name        string
		userRequest *model.User
		setup       setupFunc
		assertResp  require.ValueAssertionFunc
		assertErr   require.ErrorAssertionFunc
	}{
		{
			name:        "repository returns success",
			userRequest: &model.User{Password: "test_password_123"},
			setup: func(r *repomocks.UserRepository, u *model.User) {
				r.On("Create", mock.Anything, u).Return(nil).Once()
			},
			assertResp: require.NotNil,
			assertErr:  require.NoError,
		},
		{
			name:        "password is too long",
			userRequest: &model.User{Password: strings.Repeat("a", 73)},
			setup:       noSetup,
			assertResp:  require.Nil,
			assertErr: func(t require.TestingT, err error, _ ...interface{}) {
				require.Equal(t, err, errwrap.ErrBadRequest.SetMessage("password is too long"))
			},
		},
		{
			name:        "repository returns error",
			userRequest: &model.User{Password: "test_password_123"},
			setup: func(r *repomocks.UserRepository, u *model.User) {
				r.On("Create", mock.Anything, u).Return(errwrap.ErrConflict.SetMessage("test create error")).Once()
			},
			assertResp: require.Nil,
			assertErr: func(t require.TestingT, err error, _ ...interface{}) {
				require.Equal(t, err, errwrap.ErrConflict.SetMessage("test create error"))
			},
		},
	}
	for _, tCase := range tests {
		t.Run(tCase.name, func(tt *testing.T) {
			//test setup
			ctx := context.TODO()
			mockRepo := new(repomocks.UserRepository)
			svc := NewUserService(mockRepo)
			tCase.setup(mockRepo, tCase.userRequest)

			//execution
			res, err := svc.CreateUser(ctx, tCase.userRequest)

			//assertion
			tCase.assertErr(tt, err)
			tCase.assertResp(tt, res)

			//assert mocking calls
			assert.True(tt, mockRepo.AssertExpectations(tt))
		})
	}
}

func TestUpdateUserById(t *testing.T) {
	type request struct {
		id   string
		user *model.User
	}
	tests := []struct {
		name       string
		req        *request
		setup      func(*repomocks.UserRepository, *request)
		assertResp require.ValueAssertionFunc
		assertErr  require.ErrorAssertionFunc
	}{
		{
			name: "repository returns success",
			req:  &request{id: uuid.NewString(), user: &model.User{FirstName: "test_firstName"}},
			setup: func(r *repomocks.UserRepository, u *request) {
				r.On("Update", mock.Anything, u.id, u.user).Return(&model.User{FirstName: "test_firstName"}, nil).Once()
			},
			assertResp: func(t require.TestingT, actual interface{}, _ ...interface{}) {
				expected := &model.User{FirstName: "test_firstName"}
				require.Equal(t, expected, actual)
			},
			assertErr: require.NoError,
		},
		{
			name: "repository returns error",
			req:  &request{id: uuid.NewString(), user: &model.User{FirstName: "test_firstName"}},
			setup: func(r *repomocks.UserRepository, u *request) {
				r.On("Update", mock.Anything, u.id, u.user).Return(nil, errwrap.ErrConflict.SetMessage("test update error")).Once()
			},
			assertResp: require.Nil,
			assertErr: func(t require.TestingT, err error, _ ...interface{}) {
				require.Equal(t, errwrap.ErrConflict.SetMessage("test update error"), err)
			},
		},
	}
	for _, tCase := range tests {
		t.Run(tCase.name, func(tt *testing.T) {
			//test setup
			ctx := context.TODO()
			mockRepo := new(repomocks.UserRepository)
			svc := NewUserService(mockRepo)
			tCase.setup(mockRepo, tCase.req)

			//execution
			res, err := svc.UpdateUserById(ctx, tCase.req.id, *tCase.req.user)

			//assertion
			tCase.assertErr(tt, err)
			tCase.assertResp(tt, res)

			//assert mocking calls
			assert.True(tt, mockRepo.AssertExpectations(tt))
		})
	}
}

func TestDeleteUserById(t *testing.T) {
	tests := []struct {
		name      string
		req       string
		setup     func(*repomocks.UserRepository, string)
		assertErr require.ErrorAssertionFunc
	}{
		{
			name: "repository returns success",
			req:  uuid.NewString(),
			setup: func(r *repomocks.UserRepository, id string) {
				r.On("Delete", mock.Anything, id).Return(nil).Once()
			},
			assertErr: require.NoError,
		},
		{
			name: "repository returns error",
			req:  uuid.NewString(),
			setup: func(r *repomocks.UserRepository, id string) {
				r.On("Delete", mock.Anything, id).Return(errwrap.ErrInternal.SetMessage("test delete error")).Once()
			},
			assertErr: func(t require.TestingT, err error, _ ...interface{}) {
				require.Equal(t, errwrap.ErrInternal.SetMessage("test delete error"), err)
			},
		},
	}
	for _, tCase := range tests {
		t.Run(tCase.name, func(tt *testing.T) {
			//test setup
			ctx := context.TODO()
			mockRepo := new(repomocks.UserRepository)
			svc := NewUserService(mockRepo)
			tCase.setup(mockRepo, tCase.req)

			//execution
			err := svc.DeleteUserById(ctx, tCase.req)

			//assertion
			tCase.assertErr(tt, err)

			//assert mocking calls
			assert.True(tt, mockRepo.AssertExpectations(tt))
		})
	}
}

func TestListUsers(t *testing.T) {
	type request struct {
		filter model.UserFilter
		limit  int
		offset int
	}

	users := []model.User{
		{
			Id:        uuid.NewString(),
			FirstName: "testFirstName_1",
		},
		{
			Id:        uuid.NewString(),
			FirstName: "testFirstName_2",
		},
	}

	tests := []struct {
		name       string
		req        *request
		setup      func(*repomocks.UserRepository, *request)
		assertResp require.ValueAssertionFunc
		assertErr  require.ErrorAssertionFunc
	}{
		{
			name: "repository returns success",
			req:  &request{limit: 10, offset: 0, filter: model.UserFilter{FirstName: "test_firstName"}},
			setup: func(r *repomocks.UserRepository, req *request) {
				r.On("ListByFilter", mock.Anything, req.filter.ToBson(), req.limit, req.offset).Return(users, int64(2), nil).Once()
			},
			assertResp: func(t require.TestingT, actual interface{}, _ ...interface{}) {
				expected := &model.Pagination{
					TotalRecords: int64(len(users)),
					Limit:        10,
					Offset:       0,
					HasNext:      false,
					HasPrevious:  false,
					Items:        users,
				}
				require.Equal(t, expected, actual)
			},
			assertErr: require.NoError,
		},
		{
			name: "hasNext and hasPrevious true",
			req:  &request{limit: 10, offset: 2, filter: model.UserFilter{FirstName: "test_firstName"}},
			setup: func(r *repomocks.UserRepository, req *request) {
				r.On("ListByFilter", mock.Anything, req.filter.ToBson(), req.limit, req.offset).Return(users, int64(100), nil).Once()
			},
			assertResp: func(t require.TestingT, actual interface{}, _ ...interface{}) {
				expected := &model.Pagination{
					TotalRecords: int64(100),
					Limit:        10,
					Offset:       2,
					HasNext:      true,
					HasPrevious:  true,
					Items:        users,
				}
				require.Equal(t, expected, actual)
			},
			assertErr: require.NoError,
		},
		{
			name: "repository returns error",
			req:  &request{limit: 10, offset: 0, filter: model.UserFilter{FirstName: "test_firstName"}},
			setup: func(r *repomocks.UserRepository, req *request) {
				r.On("ListByFilter", mock.Anything, req.filter.ToBson(), req.limit, req.offset).Return(nil, int64(0), errwrap.ErrInternal.SetMessage("test list error")).Once()
			},
			assertResp: require.Nil,
			assertErr: func(t require.TestingT, err error, _ ...interface{}) {
				require.Equal(t, errwrap.ErrInternal.SetMessage("test list error"), err)
			},
		},
	}
	for _, tCase := range tests {
		t.Run(tCase.name, func(tt *testing.T) {
			//test setup
			ctx := context.TODO()
			mockRepo := new(repomocks.UserRepository)
			svc := NewUserService(mockRepo)
			tCase.setup(mockRepo, tCase.req)

			//execution
			res, err := svc.ListUsers(ctx, tCase.req.filter, tCase.req.limit, tCase.req.offset)

			//assertion
			tCase.assertErr(tt, err)
			tCase.assertResp(tt, res)

			//assert mocking calls
			assert.True(tt, mockRepo.AssertExpectations(tt))
		})
	}
}
