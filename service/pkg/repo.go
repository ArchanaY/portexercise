package repo

import (
	"context"
	"errors"
	"log"
	"portexercise/service/domain"
	"sync"
)

type InMemoryStore struct {
	mutex sync.RWMutex
	data  map[string]*domain.Port
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]*domain.Port),
	}
}

func (is *InMemoryStore) Insert(ctx context.Context, pi domain.Port) error {
	if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
		log.Println("context is cancelled")
		return nil
	}

	is.mutex.Lock()
	defer is.mutex.Unlock()

	is.data[pi.Key] = &domain.Port{
		Name:        pi.Name,
		City:        pi.City,
		Country:     pi.Country,
		Province:    pi.Province,
		Code:        pi.Code,
		Alias:       pi.Alias,
		Regions:     pi.Regions,
		Coordinates: pi.Coordinates,
		Timezone:    pi.Timezone,
		Unlocs:      pi.Unlocs,
	}

	// b, _ := json.MarshalIndent(is.data, "", "  ")
	// fmt.Println(string(b))
	return nil
}

func (is *InMemoryStore) Fetch(ctx context.Context, key string) (domain.Port, error) {
	if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
		log.Println("context is cancelled")
		return domain.Port{}, nil
	}

	is.mutex.RLock()
	defer is.mutex.RUnlock()

	port := is.data[key]
	if port == nil {
		return domain.Port{}, errors.New("Port not found")
	}

	return *port, nil
}
