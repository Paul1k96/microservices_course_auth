package user

import (
	"context"
	"fmt"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"golang.org/x/sync/errgroup"
)

const (
	groupLimit = 10
)

// GetListByIDs returns list of users by ids.
func (s *service) GetListByIDs(ctx context.Context, ids []int64) ([]*model.User, error) {
	notInCacheCh := make(chan int64, len(ids))
	notInCache := make([]int64, 0, len(ids))

	inCacheCh := make(chan *model.User, len(ids))

	errGroup, errCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(groupLimit)

	for _, id := range ids {
		id := id
		errGroup.Go(func() error {
			user, err := s.cache.Get(errCtx, id)
			if err != nil {
				notInCacheCh <- id
				return nil
			}

			inCacheCh <- user

			return nil
		})
	}

	_ = errGroup.Wait()
	close(notInCacheCh)
	close(inCacheCh)

	result := make([]*model.User, 0, len(ids))
	for user := range inCacheCh {
		result = append(result, user)
	}

	for id := range notInCacheCh {
		notInCache = append(notInCache, id)
	}

	if len(notInCache) > 0 {
		users, err := s.repo.GetByIDs(ctx, notInCache)
		if err != nil {
			return nil, fmt.Errorf("get list by ids: %w", err)
		}

		for _, user := range users {
			user := user
			errGroup.Go(func() error {
				_ = s.cache.Set(errCtx, user)
				return nil
			})
		}

		_ = errGroup.Wait()

		result = append(result, users...)
	}

	return result, nil
}
