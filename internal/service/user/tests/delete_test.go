package tests

import (
	"context"
	"testing"

	"github.com/Paul1k96/microservices_course_auth/internal/repository/mocks"
	"github.com/Paul1k96/microservices_course_auth/internal/service"
	"github.com/Paul1k96/microservices_course_auth/internal/service/user"
	mocks2 "github.com/Paul1k96/microservices_course_platform_common/pkg/client/db/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestDeleteUserSuite(t *testing.T) {
	suite.Run(t, new(DeleteUserSuite))
}

type DeleteUserSuite struct {
	suite.Suite
	*require.Assertions
	ctrl *gomock.Controller

	userRepo  *mocks.MockUsersRepository
	userCache *mocks.MockUsersCache

	service service.UserService
}

func (t *DeleteUserSuite) SetupTest() {
	t.Assertions = require.New(t.T())
	t.ctrl = gomock.NewController(t.T())

	t.userRepo = mocks.NewMockUsersRepository(t.ctrl)
	t.userCache = mocks.NewMockUsersCache(t.ctrl)

	t.service = user.NewService(t.userRepo, t.userCache, mocks2.NewMockTxManager(t.ctrl))
}

func (t *DeleteUserSuite) TearDownTest() {
	t.ctrl.Finish()
}

type DeleteUserArgs struct {
	ctx context.Context
	id  int64
}

type DeleteUserWant struct {
	err error
}

func (t *DeleteUserSuite) do(args DeleteUserArgs, want DeleteUserWant) {
	err := t.service.Delete(args.ctx, args.id)

	if want.err == nil {
		t.Require().NoError(err)
	} else {
		t.Require().ErrorContains(err, want.err.Error())
	}
}

func (t *DeleteUserSuite) TestDeleteUser_Ok() {
	args := DeleteUserArgs{
		ctx: context.Background(),
		id:  gofakeit.Int64(),
	}

	want := DeleteUserWant{}

	t.userRepo.EXPECT().Delete(args.ctx, args.id).Return(nil)
	t.userCache.EXPECT().Delete(args.ctx, args.id).Return(nil)

	t.do(args, want)
}

func (t *DeleteUserSuite) TestDeleteUser_RepoError() {
	args := DeleteUserArgs{
		ctx: context.Background(),
		id:  gofakeit.Int64(),
	}

	want := DeleteUserWant{
		err: gofakeit.Error(),
	}

	t.userRepo.EXPECT().Delete(args.ctx, args.id).Return(want.err)

	t.do(args, want)
}
