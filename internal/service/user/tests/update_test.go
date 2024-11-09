package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/Paul1k96/microservices_course_auth/internal/errs"
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"github.com/Paul1k96/microservices_course_auth/internal/repository/mocks"
	"github.com/Paul1k96/microservices_course_auth/internal/service"
	"github.com/Paul1k96/microservices_course_auth/internal/service/user"
	tm "github.com/Paul1k96/microservices_course_auth/internal/testmodel"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db/transaction"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestUpdateUserSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserSuite))
}

type UpdateUserSuite struct {
	suite.Suite
	*require.Assertions
	ctrl *gomock.Controller

	userRepo  *mocks.MockUsersRepository
	userCache *mocks.MockUsersCache

	service service.UserService
}

func (t *UpdateUserSuite) SetupTest() {
	t.Assertions = require.New(t.T())
	t.ctrl = gomock.NewController(t.T())

	t.userRepo = mocks.NewMockUsersRepository(t.ctrl)
	t.userCache = mocks.NewMockUsersCache(t.ctrl)

	t.service = user.NewService(t.userRepo, t.userCache, transaction.NewNopTxManager())
}

func (t *UpdateUserSuite) TearDownTest() {
	t.ctrl.Finish()
}

type UpdateUserArgs struct {
	ctx  context.Context
	user *model.User
}

type UpdateUserWant struct {
	err error
}

func (t *UpdateUserSuite) do(args UpdateUserArgs, want UpdateUserWant) {
	err := t.service.Update(args.ctx, args.user)

	if want.err == nil {
		t.Require().NoError(err)
	} else {
		t.Require().ErrorContains(err, want.err.Error())
	}
}

func (t *UpdateUserSuite) TestUpdateUser_OkChangeName() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Name = gofakeit.Name()

	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: nil,
	}

	t.userRepo.EXPECT().Update(args.ctx, args.user).Return(want.err)

	t.userCache.EXPECT().Set(args.ctx, args.user).Return(nil)

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_OkEmptyName() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Name = ""

	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: nil,
	}

	t.userRepo.EXPECT().Update(args.ctx, args.user).Return(want.err)

	t.userCache.EXPECT().Set(args.ctx, args.user).Return(nil)

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_TooShortName() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Name = "a"
	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: errors.New("name must be at least 2 characters"),
	}

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_TooLongName() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Name = strings.Repeat("a", 101)
	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: errors.New("name must be at most 100 characters"),
	}

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_NameContainsRestrictedSymbols() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Name = changeUser.Name + "!"
	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: errors.New("name contains restricted symbols"),
	}

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_OkEmailIsEmpty() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Email = ""
	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: nil,
	}

	t.userRepo.EXPECT().Update(args.ctx, args.user).Return(nil)

	t.userCache.EXPECT().Set(args.ctx, args.user).Return(nil)

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_EmailTooShort() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Email = "a"
	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: errors.New("email must be at least 5 characters"),
	}

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_EmailTooLong() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Email = strings.Repeat("a", 101) + "@gmail.com"
	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: errors.New("email must be at most 100 characters"),
	}

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_EmailInvalid() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Email = gofakeit.URL()
	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: errors.New("email has incorrect format"),
	}

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_RoleInvalid() {
	usr := tm.NewUser()

	changeUser := usr
	changeUser.Role = 100
	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: errors.New("role is not valid"),
	}

	t.do(args, want)
}

func (t *UpdateUserSuite) TestUpdateUser_RepoError() {
	usr := tm.NewUser()

	changeUser := usr
	args := UpdateUserArgs{
		ctx:  context.Background(),
		user: changeUser,
	}

	want := UpdateUserWant{
		err: errs.ErrUserNotFound,
	}

	t.userRepo.EXPECT().Update(args.ctx, args.user).Return(want.err)

	t.do(args, want)
}
