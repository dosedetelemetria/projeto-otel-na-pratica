// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package memory

import (
	"context"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/model"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/store"
	"go.opentelemetry.io/otel"
)

type inMemoryUser struct {
	store map[string]*model.User
}

func NewUserStore(ctx context.Context) store.User {
	_, span := otel.Tracer("user").Start(ctx, "NewUserStore")
	defer span.End()

	return &inMemoryUser{
		store: make(map[string]*model.User),
	}
}

func (u *inMemoryUser) Get(_ context.Context, id string) (*model.User, error) {
	return u.store[id], nil
}

func (u *inMemoryUser) Create(_ context.Context, user *model.User) (*model.User, error) {
	u.store[user.ID] = user
	return user, nil
}

func (u *inMemoryUser) Update(_ context.Context, user *model.User) (*model.User, error) {
	u.store[user.ID] = user
	return user, nil
}

func (u *inMemoryUser) Delete(_ context.Context, id string) error {
	delete(u.store, id)
	return nil
}

func (u *inMemoryUser) List(_ context.Context) ([]*model.User, error) {
	users := make([]*model.User, 0, len(u.store))
	for _, user := range u.store {
		users = append(users, user)
	}
	return users, nil
}
