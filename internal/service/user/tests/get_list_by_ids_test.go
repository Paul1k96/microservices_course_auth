package tests

import (
	"context"
	"testing"

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
	"golang.org/x/sync/errgroup"
)

func TestGetListByIDsSuite(t *testing.T) {
	suite.Run(t, new(GetListByIDsSuite))
}

type GetListByIDsSuite struct {
	suite.Suite
	*require.Assertions
	ctrl *gomock.Controller

	userRepo  *mocks.MockUsersRepository
	userCache *mocks.MockUsersCache

	service service.UserService
}

func (t *GetListByIDsSuite) SetupTest() {
	t.Assertions = require.New(t.T())
	t.ctrl = gomock.NewController(t.T())

	t.userRepo = mocks.NewMockUsersRepository(t.ctrl)
	t.userCache = mocks.NewMockUsersCache(t.ctrl)

	t.service = user.NewService(t.userRepo, t.userCache, mocks2.NewMockTxManager(t.ctrl))
}

func (t *GetListByIDsSuite) TearDownTest() {
	t.ctrl.Finish()
}

type GetListByIDsArgs struct {
	ctx context.Context
	ids []int64
}

type GetListByIDsWant struct {
	user []*model.User
	err  error
}

func (t *GetListByIDsSuite) do(args GetListByIDsArgs, want GetListByIDsWant) {
	usr, err := t.service.GetListByIDs(args.ctx, args.ids)

	t.Require().ElementsMatch(want.user, usr)

	if want.err == nil {
		t.Require().NoError(err)
	} else {
		t.Require().ErrorContains(err, want.err.Error())
	}
}

func (t *GetListByIDsSuite) TestGetUser_OkFromCache() {
	usrs := tm.NewUsers(3)

	args := GetListByIDsArgs{
		ctx: context.Background(),
		ids: []int64{usrs[0].ID, usrs[1].ID, usrs[2].ID},
	}
	want := GetListByIDsWant{
		user: usrs,
		err:  nil,
	}

	_, errCtx := errgroup.WithContext(args.ctx)
	t.userCache.EXPECT().Get(errCtx, usrs[0].ID).Return(usrs[0], nil)
	t.userCache.EXPECT().Get(errCtx, usrs[1].ID).Return(usrs[1], nil)
	t.userCache.EXPECT().Get(errCtx, usrs[2].ID).Return(usrs[2], nil)

	t.do(args, want)
}

func (t *GetListByIDsSuite) TestGetUser_OkFromRepoAndSetCacheAll() {
	usrs := tm.NewUsers(3)

	args := GetListByIDsArgs{
		ctx: context.Background(),
		ids: []int64{usrs[0].ID, usrs[1].ID, usrs[2].ID},
	}
	want := GetListByIDsWant{
		user: usrs,
		err:  nil,
	}

	_, errCtx := errgroup.WithContext(args.ctx)

	t.userCache.EXPECT().Get(errCtx, usrs[0].ID).Return(nil, gofakeit.Error())
	t.userCache.EXPECT().Get(errCtx, usrs[1].ID).Return(nil, gofakeit.Error())
	t.userCache.EXPECT().Get(errCtx, usrs[2].ID).Return(nil, gofakeit.Error())

	t.userRepo.EXPECT().GetByIDs(
		args.ctx,
		gomock.InAnyOrder([]int64{usrs[0].ID, usrs[1].ID, usrs[2].ID})).Return(usrs, nil)

	t.userCache.EXPECT().Set(gomock.Any(), usrs[0]).Return(nil)
	t.userCache.EXPECT().Set(gomock.Any(), usrs[1]).Return(nil)
	t.userCache.EXPECT().Set(gomock.Any(), usrs[2]).Return(nil)

	t.do(args, want)
}

func (t *GetListByIDsSuite) TestGetUser_OkFromRepoAndSetCachePartially() {
	usrs := tm.NewUsers(3)

	args := GetListByIDsArgs{
		ctx: context.Background(),
		ids: []int64{usrs[0].ID, usrs[1].ID, usrs[2].ID},
	}
	want := GetListByIDsWant{
		user: usrs,
		err:  nil,
	}

	_, errCtx := errgroup.WithContext(args.ctx)

	t.userCache.EXPECT().Get(errCtx, usrs[0].ID).Return(usrs[0], nil)
	t.userCache.EXPECT().Get(errCtx, usrs[1].ID).Return(nil, gofakeit.Error())
	t.userCache.EXPECT().Get(errCtx, usrs[2].ID).Return(nil, gofakeit.Error())

	repoUsrs := []*model.User{usrs[1], usrs[2]}
	t.userRepo.EXPECT().GetByIDs(
		args.ctx,
		gomock.InAnyOrder([]int64{usrs[1].ID, usrs[2].ID})).Return(repoUsrs, nil)

	t.userCache.EXPECT().Set(gomock.Any(), usrs[1]).Return(nil)
	t.userCache.EXPECT().Set(gomock.Any(), usrs[2]).Return(nil)

	t.do(args, want)
}

func (t *GetListByIDsSuite) TestGetUser_OkEmptyUsers() {
	usrs := tm.NewUsers(3)

	args := GetListByIDsArgs{
		ctx: context.Background(),
		ids: []int64{usrs[0].ID, usrs[1].ID, usrs[2].ID},
	}
	want := GetListByIDsWant{
		user: []*model.User{},
		err:  nil,
	}

	_, errCtx := errgroup.WithContext(args.ctx)

	t.userCache.EXPECT().Get(errCtx, usrs[0].ID).Return(nil, gofakeit.Error())
	t.userCache.EXPECT().Get(errCtx, usrs[1].ID).Return(nil, gofakeit.Error())
	t.userCache.EXPECT().Get(errCtx, usrs[2].ID).Return(nil, gofakeit.Error())

	t.userRepo.EXPECT().GetByIDs(
		args.ctx,
		gomock.InAnyOrder([]int64{usrs[0].ID, usrs[1].ID, usrs[2].ID})).Return([]*model.User{}, nil)

	t.do(args, want)
}
