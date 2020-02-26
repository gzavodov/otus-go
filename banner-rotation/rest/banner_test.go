package rest

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gzavodov/otus-go/banner-rotation/config"
	"github.com/gzavodov/otus-go/banner-rotation/internal/sql"
	"github.com/gzavodov/otus-go/banner-rotation/internal/testify"
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
)

func TestBanner(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	conf := &config.Configuration{}
	//OS variable BANNER_ROTATION_REPOSITORY_DSN is required
	if err := conf.LoadFromEvironment(); err != nil {
		t.Fatal(err)
	}

	bannerRepo := sql.NewBannerRepository(ctx, conf.RepositoryDSN)
	//slotRepo := sql.NewSlotRepository(ctx, conf.RepositoryDSN)
	bindingRepo := sql.NewBindingRepository(ctx, conf.RepositoryDSN)
	//groupRepo := sql.NewGroupRepository(ctx, conf.RepositoryDSN)
	statisticsRepo := sql.NewStatisticsRepository(ctx, conf.RepositoryDSN)

	bannerUsecase := usecase.NewBannerUsecase(bannerRepo, bindingRepo, statisticsRepo, conf.AlgorithmTypeID)
	bannerHandler := NewBannerHandler(bannerUsecase, "Test", nil, nil)

	var sourceBanner *model.Banner
	form := testify.Form{}

	t.Run("Create",
		func(t *testing.T) {
			caption := "Test Banner #1"

			formData := url.Values{}
			formData.Set("caption", caption)

			responseBody, err := form.EmulatePost("/banner/create", formData, bannerHandler.Create)
			if err != nil {
				t.Fatal(err)
			}

			resultBanner := &model.Banner{}
			if err := json.NewDecoder(responseBody).Decode(resultBanner); err != nil {
				t.Fatal(err)
			}

			if resultBanner.ID < 0 {
				t.Errorf("handler returned unexpected banner ID: got %d", resultBanner.ID)
			}

			if resultBanner.Caption != caption {
				t.Errorf("handler returned unexpected banner caption: got %s want %s", resultBanner.Caption, caption)
			}

			sourceBanner = resultBanner
		})

	t.Run("Read",
		func(t *testing.T) {
			if sourceBanner == nil {
				t.Skip()
			}

			formData := url.Values{}
			formData.Set("ID", strconv.FormatInt(sourceBanner.ID, 10))

			responseBody, err := form.EmulatePost("/banner/read", formData, bannerHandler.Read)
			if err != nil {
				t.Fatal(err)
			}

			resultBanner := &model.Banner{}
			if err := json.NewDecoder(responseBody).Decode(resultBanner); err != nil {
				t.Fatal(err)
			}

			if *sourceBanner != *resultBanner {
				t.Error(testify.NewObjectNotMatchedError(sourceBanner, resultBanner))
			}
		})

	t.Run("Update",
		func(t *testing.T) {
			if sourceBanner == nil {
				t.Skip()
			}

			sourceBanner.Caption = "Test Banner #2"

			formData := url.Values{}
			formData.Set("ID", strconv.FormatInt(sourceBanner.ID, 10))
			formData.Set("caption", sourceBanner.Caption)

			responseBody, err := form.EmulatePost("/banner/update", formData, bannerHandler.Update)
			if err != nil {
				t.Fatal(err)
			}

			resultBanner := &model.Banner{}
			if err := json.NewDecoder(responseBody).Decode(resultBanner); err != nil {
				t.Fatal(err)
			}

			if *sourceBanner != *resultBanner {
				t.Error(testify.NewObjectNotMatchedError(sourceBanner, resultBanner))
			}
		})

	t.Run("Delete",
		func(t *testing.T) {
			if sourceBanner == nil {
				t.Skip()
			}

			formData := url.Values{}
			formData.Set("ID", strconv.FormatInt(sourceBanner.ID, 10))

			_, err := form.EmulatePost("/banner/delete", formData, bannerHandler.Delete)
			if err != nil {
				t.Fatal(err)
			}
		})
}
