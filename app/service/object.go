package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/madshov/bitburst/app"
)

type service struct {
	logger *log.Logger
	repo   app.ObjectStorage
}

func NewService(logger *log.Logger, r app.ObjectStorage) ObjectService {
	return &service{
		logger: logger,
		repo:   r,
	}
}

type ObjectService interface {
	StoreObjects(ctx context.Context, objs []int)
	DeleteObjects(ctx context.Context)
}

func (cs *service) StoreObjects(ctx context.Context, objectIDs []int) {
	objs := make([]app.Object, len(objectIDs))
	resCh := make(chan app.Object)
	errCh := make(chan error)

	defer close(resCh)
	defer close(errCh)

	for _, id := range objectIDs {
		go func(id int, resCh chan app.Object, errCh chan error) {
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(fmt.Sprintf("http://bitburst-tester-service:9010/objects/%d", id))
			if err != nil {
				errCh <- err
				return
			}

			var ret app.Object
			err = json.NewDecoder(resp.Body).Decode(&ret)
			if err != nil {
				errCh <- err
				return
			}

			_ = resp.Body.Close()
			resCh <- ret

		}(id, resCh, errCh)
	}

	for i := 0; i < len(objectIDs); i++ {
		select {
		case result := <-resCh:
			objs[i] = result
		case err := <-errCh:
			cs.logger.Println("msg", "error fetching objects", "reason", err)
		}
	}

	cs.logger.Println("msg", "creating objects", "count", len(objs))

	err := cs.repo.CreateObjects(objs)
	if err != nil {
		cs.logger.Println("msg", "error creating objects", "reason", err)
	}
}

func (cs *service) DeleteObjects(ctx context.Context) {
	cs.logger.Println("msg", "deleting expired objects")

	err := cs.repo.DeleteObjects()
	if err != nil {
		cs.logger.Println("msg", "error deleting objects", "reason", err)
	}
}
