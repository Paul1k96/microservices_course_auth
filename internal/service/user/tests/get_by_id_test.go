package tests

import (
	"context"
	"testing"

	"github.com/Paul1k96/microservices_course_auth/internal/errs"
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"github.com/Paul1k96/microservices_course_auth/internal/repository/mocks"
	"github.com/Paul1k96/microservices_course_auth/internal/service"
	"github.com/Paul1k96/microservices_course_auth/internal/service/user"
	tm "github.com/Paul1k96/microservices_course_auth/internal/testmodel"
	mocks2 "github.com/Paul1k96/microservices_course_platform_common/pkg/client/db/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestGetUserSuite(t *testing.T) {
	suite.Run(t, new(GetUserSuite))
}

type GetUserSuite struct {
	suite.Suite
	*require.Assertions
	ctrl *gomock.Controller

	userRepo  *mocks.MockUsersRepository
	userCache *mocks.MockUsersCache

	service service.UserService
}

func (t *GetUserSuite) SetupTest() {
	t.Assertions = require.New(t.T())
	t.ctrl = gomock.NewController(t.T())

	t.userRepo = mocks.NewMockUsersRepository(t.ctrl)
	t.userCache = mocks.NewMockUsersCache(t.ctrl)

	t.service = user.NewService(t.userRepo, t.userCache, mocks2.NewMockTxManager(t.ctrl))
}

func (t *GetUserSuite) TearDownTest() {
	t.ctrl.Finish()
}

type GetUserArgs struct {
	ctx context.Context
	id  int64
}

type GetUserWant struct {
	user *model.User
	err  error
}

func (t *GetUserSuite) do(args GetUserArgs, want GetUserWant) {
	usr, err := t.service.GetByID(args.ctx, args.id)

	t.Require().Equal(want.user, usr)

	if want.err == nil {
		t.Require().NoError(err)
	} else {
		t.Require().ErrorContains(err, want.err.Error())
	}
}

func (t *GetUserSuite) TestGetUser_OkFromCache() {
	usr := tm.NewUser()

	args := GetUserArgs{
		ctx: context.Background(),
		id:  usr.ID,
	}
	want := GetUserWant{
		user: usr,
		err:  nil,
	}

	t.userCache.EXPECT().Get(args.ctx, args.id).Return(usr, nil)

	t.do(args, want)
}

func (t *GetUserSuite) TestGetUser_OkFromRepoAndSetCache() {
	usr := tm.NewUser()

	args := GetUserArgs{
		ctx: context.Background(),
		id:  usr.ID,
	}
	want := GetUserWant{
		user: usr,
		err:  nil,
	}

	t.userCache.EXPECT().Get(args.ctx, args.id).Return(nil, errs.ErrUserNotFound)
	t.userRepo.EXPECT().GetByID(args.ctx, args.id).Return(usr, nil)
	t.userCache.EXPECT().Set(args.ctx, usr).Return(nil)

	t.do(args, want)
}

func (t *GetUserSuite) TestGetUser_UserNotFound() {
	args := GetUserArgs{
		ctx: context.Background(),
		id:  0,
	}
	want := GetUserWant{
		user: nil,
		err:  errs.ErrUserNotFound,
	}

	t.userCache.EXPECT().Get(args.ctx, args.id).Return(nil, errs.ErrUserNotFound)
	t.userRepo.EXPECT().GetByID(args.ctx, args.id).Return(nil, errs.ErrUserNotFound)

	t.do(args, want)
}

func (t *GetUserSuite) TestGetUser_FailedToGetFromCache() {
	args := GetUserArgs{
		ctx: context.Background(),
		id:  0,
	}
	want := GetUserWant{
		user: nil,
		err:  gofakeit.Error(),
	}

	t.userCache.EXPECT().Get(args.ctx, args.id).Return(nil, want.err)

	t.do(args, want)
}

func (t *GetUserSuite) TestGetUser_FailedToGetFromRepo() {
	args := GetUserArgs{
		ctx: context.Background(),
		id:  0,
	}
	want := GetUserWant{
		user: nil,
		err:  gofakeit.Error(),
	}

	t.userCache.EXPECT().Get(args.ctx, args.id).Return(nil, errs.ErrUserNotFound)
	t.userRepo.EXPECT().GetByID(args.ctx, args.id).Return(nil, want.err)

	t.do(args, want)
}

func (t *GetUserSuite) TestGetUser_RepoNilErrorButEmptyUser() {
	args := GetUserArgs{
		ctx: context.Background(),
		id:  0,
	}
	want := GetUserWant{
		user: nil,
		err:  errs.ErrUserNotFound,
	}

	t.userCache.EXPECT().Get(args.ctx, args.id).Return(nil, errs.ErrUserNotFound)
	t.userRepo.EXPECT().GetByID(args.ctx, args.id).Return(nil, nil)

	t.do(args, want)
}
