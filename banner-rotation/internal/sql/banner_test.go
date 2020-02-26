package sql

import (
	"context"
	"testing"
	"time"

	"github.com/gzavodov/otus-go/banner-rotation/config"
	testify "github.com/gzavodov/otus-go/banner-rotation/internal/testify"
	"github.com/gzavodov/otus-go/banner-rotation/model"
)

func TestBannerRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	configuration := &config.Configuration{}
	//OS variable BANNER_ROTATION_REPOSITORY_DSN is required
	if err := configuration.LoadFromEvironment(); err != nil {
		t.Fatal(err)
	}

	repo := NewBannerRepository(ctx, configuration.RepositoryDSN)

	t.Run("BannerRepository::Create",
		func(t *testing.T) {
			source := &model.Banner{BaseReference: model.BaseReference{Caption: "Creation Test"}}

			if err := repo.Create(source); err != nil {
				t.Fatal(err)
			}

			ok, err := repo.IsExists(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				t.Error(testify.NewObjectNotFoundError())
			}

			if err = repo.Delete(source.ID); err != nil {
				t.Fatal(err)
			}
		})

	t.Run("BannerRepository::Read",
		func(t *testing.T) {
			source := &model.Banner{BaseReference: model.BaseReference{Caption: "Reading Test"}}

			if err := repo.Create(source); err != nil {
				t.Fatal(err)
			}

			result, err := repo.Read(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *source != *result {
				t.Error(testify.NewObjectNotMatchedError(source, result))
			}

			if err = repo.Delete(source.ID); err != nil {
				t.Fatal(err)
			}
		})

	t.Run("BannerRepository::Update",
		func(t *testing.T) {
			source := &model.Banner{BaseReference: model.BaseReference{Caption: "Test"}}

			err := repo.Create(source)
			if err != nil {
				t.Fatal(err)
			}

			source.Caption = "Modification Test"
			err = repo.Update(source)
			if err != nil {
				t.Fatal(err)
			}

			result, err := repo.Read(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *source != *result {
				t.Error(testify.NewObjectNotMatchedError(source, result))
			}

			if err = repo.Delete(source.ID); err != nil {
				t.Fatal(err)
			}
		})

	t.Run("BannerRepository::Delete",
		func(t *testing.T) {
			source := &model.Banner{BaseReference: model.BaseReference{Caption: "Deletion Test"}}

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
				t.Fatal(testify.NewObjectNotDeletedError())
			}
		})
}
