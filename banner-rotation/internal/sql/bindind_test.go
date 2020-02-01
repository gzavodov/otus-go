package sql

import (
	"context"
	"testing"
	"time"

	"github.com/gzavodov/otus-go/banner-rotation/config"
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/test"
)

func TestBindingRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	configuration := &config.Configuration{}
	//OS variable BANNER_ROTATION_REPOSITORY_DSN is required
	if err := configuration.LoadFromEvironment(); err != nil {
		t.Fatal(err)
	}

	repo := NewBindingRepository(ctx, configuration.RepositoryDSN)

	t.Run("BindingRepository::Create",
		func(t *testing.T) {
			source := &model.Binding{BannerID: 1, SlotID: 1}

			if err := repo.Create(source); err != nil {
				t.Fatal(err)
			}

			ok, err := repo.IsExists(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				t.Error(test.NewObjectNotFoundError())
			}

			if err = repo.Delete(source.ID); err != nil {
				t.Fatal(err)
			}
		})

	t.Run("BindingRepository::Read",
		func(t *testing.T) {
			source := &model.Binding{BannerID: 1, SlotID: 1}

			if err := repo.Create(source); err != nil {
				t.Fatal(err)
			}

			result, err := repo.Read(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *source != *result {
				t.Error(test.NewObjectNotMatchedError(source, result))
			}

			if err = repo.Delete(source.ID); err != nil {
				t.Fatal(err)
			}
		})

	t.Run("BindingRepository::Update",
		func(t *testing.T) {
			source := &model.Binding{BannerID: 1, SlotID: 1}

			err := repo.Create(source)
			if err != nil {
				t.Fatal(err)
			}

			source.BannerID = 2
			source.SlotID = 2

			err = repo.Update(source)
			if err != nil {
				t.Fatal(err)
			}

			result, err := repo.Read(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *source != *result {
				t.Error(test.NewObjectNotMatchedError(source, result))
			}

			if err = repo.Delete(source.ID); err != nil {
				t.Fatal(err)
			}
		})

	t.Run("BindingRepository::Delete",
		func(t *testing.T) {
			source := &model.Binding{BannerID: 1, SlotID: 1}

			err := repo.Create(source)
			if err != nil {
				t.Fatal(err)
			}

			err = repo.Delete(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			ok, err := repo.IsExists(source.ID)
			if err != nil {
				t.Error(err)
			}

			if ok {
				t.Fatal(test.NewObjectNotDeletedError())
			}
		})
}
