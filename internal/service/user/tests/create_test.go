package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"github.com/Paul1k96/microservices_course_auth/internal/repository/mocks"
	"github.com/Paul1k96/microservices_course_auth/internal/service"
	"github.com/Paul1k96/microservices_course_auth/internal/service/user"
	tm "github.com/Paul1k96/microservices_course_auth/internal/testmodel"
	infraMocks "github.com/Paul1k96/microservices_course_platform_common/pkg/client/db/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestCreateUserSuite(t *testing.T) {
	suite.Run(t, new(CreateUserSuite))
}

type CreateUserSuite struct {
	suite.Suite
	*require.Assertions
	ctrl *gomock.Controller

	userRepo  *mocks.MockUsersRepository
	userCache *mocks.MockUsersCache

	service service.UserService
}

func (t *CreateUserSuite) SetupTest() {
	t.Assertions = require.New(t.T())
	t.ctrl = gomock.NewController(t.T())

	t.userRepo = mocks.NewMockUsersRepository(t.ctrl)
	t.userCache = mocks.NewMockUsersCache(t.ctrl)

	t.service = user.NewService(t.userRepo, t.userCache, infraMocks.NewMockTxManager(t.ctrl))
}

func (t *CreateUserSuite) TearDownTest() {
	t.ctrl.Finish()
}

type CreateUserArgs struct {
	ctx  context.Context
	user *model.User
}

type CreateUserWant struct {
	id  int64
	err error
}

func (t *CreateUserSuite) do(args CreateUserArgs, want CreateUserWant) {
	id, err := t.service.Create(args.ctx, args.user)

	t.Require().Equal(want.id, id)

	if want.err == nil {
		t.Require().NoError(err)
	} else {
		t.Require().ErrorContains(err, want.err.Error())
	}
}

func (t *CreateUserSuite) TestCreateUser_Ok() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}

	want := CreateUserWant{
		id:  1,
		err: nil,
	}

	t.userRepo.EXPECT().Create(args.ctx, args.user).Return(want.id, want.err)

	t.userCache.EXPECT().Set(args.ctx, args.user).Return(nil)

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_EmptyName() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}
	args.user.Name = ""

	want := CreateUserWant{
		id:  0,
		err: errors.New("name is required"),
	}

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_TooShortName() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}
	args.user.Name = "a"

	want := CreateUserWant{
		id:  0,
		err: errors.New("name must be at least 2 characters"),
	}

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_TooLongName() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}
	args.user.Name = strings.Repeat("a", 101)

	want := CreateUserWant{
		id:  0,
		err: errors.New("name must be at most 100 characters"),
	}

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_NameContainsRestrictedSymbols() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}
	args.user.Name = args.user.Name + "!"

	want := CreateUserWant{
		id:  0,
		err: errors.New("name contains restricted symbols"),
	}

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_EmailIsEmpty() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}
	args.user.Email = ""

	want := CreateUserWant{
		id:  0,
		err: errors.New("email is required"),
	}

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_EmailTooShort() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}
	args.user.Email = "a"

	want := CreateUserWant{
		id:  0,
		err: errors.New("email must be at least 5 characters"),
	}

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_EmailTooLong() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}
	args.user.Email = strings.Repeat("a", 101) + "@gmail.com"

	want := CreateUserWant{
		id:  0,
		err: errors.New("email must be at most 100 characters"),
	}

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_EmailInvalid() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}
	args.user.Email = gofakeit.URL()

	want := CreateUserWant{
		id:  0,
		err: errors.New("email has incorrect format"),
	}

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_RoleInvalid() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}
	args.user.Role = 100

	want := CreateUserWant{
		id:  0,
		err: errors.New("role is not valid"),
	}

	t.do(args, want)
}

func (t *CreateUserSuite) TestCreateUser_RepoError() {
	args := CreateUserArgs{
		ctx:  context.Background(),
		user: tm.NewUser(),
	}

	want := CreateUserWant{
		id:  0,
		err: errors.New("repo error"),
	}

	t.userRepo.EXPECT().Create(args.ctx, args.user).Return(int64(0), want.err)

	t.do(args, want)
}
