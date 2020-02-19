package rest

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gzavodov/otus-go/banner-rotation/config"
	"github.com/gzavodov/otus-go/banner-rotation/endpoint"
	"github.com/gzavodov/otus-go/banner-rotation/internal/sql"
	"github.com/gzavodov/otus-go/banner-rotation/internal/testify"
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
)

func TestGroup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	conf := &config.Configuration{}
	//OS variable BANNER_ROTATION_REPOSITORY_DSN is required
	if err := conf.LoadFromEvironment(); err != nil {
		t.Fatal(err)
	}

	groupRepo := sql.NewGroupRepository(ctx, conf.RepositoryDSN)
	groupUsecase := usecase.NewGroupUsecase(groupRepo)

	groupHandler := &Group{
		Handler: endpoint.Handler{Name: "Group", ServiceName: "Test"},
		ucase:   groupUsecase,
	}

	var sourceGroup *model.Group
	form := testify.Form{}

	t.Run("Create",
		func(t *testing.T) {
			caption := "Test Group #1"

			formData := url.Values{}
			formData.Set("caption", caption)

			responseBody, err := form.EmulatePost("/group/create", formData, groupHandler.Create)
			if err != nil {
				t.Fatal(err)
			}

			resultGroup := &model.Group{}
			if err := json.NewDecoder(responseBody).Decode(resultGroup); err != nil {
				t.Fatal(err)
			}

			if resultGroup.ID < 0 {
				t.Errorf("handler returned unexpected banner ID: got %d", resultGroup.ID)
			}

			if resultGroup.Caption != caption {
				t.Errorf("handler returned unexpected banner caption: got %s want %s", resultGroup.Caption, caption)
			}

			sourceGroup = resultGroup
		})

	t.Run("Read",
		func(t *testing.T) {
			if sourceGroup == nil {
				t.Skip()
			}

			formData := url.Values{}
			formData.Set("ID", strconv.FormatInt(sourceGroup.ID, 10))

			responseBody, err := form.EmulatePost("/group/read", formData, groupHandler.Read)
			if err != nil {
				t.Fatal(err)
			}

			resultGroup := &model.Group{}
			if err := json.NewDecoder(responseBody).Decode(resultGroup); err != nil {
				t.Fatal(err)
			}

			if *sourceGroup != *resultGroup {
				t.Error(testify.NewObjectNotMatchedError(sourceGroup, resultGroup))
			}
		})

	t.Run("Update",
		func(t *testing.T) {
			if sourceGroup == nil {
				t.Skip()
			}

			sourceGroup.Caption = "Test Group #2"

			formData := url.Values{}
			formData.Set("ID", strconv.FormatInt(sourceGroup.ID, 10))
			formData.Set("caption", sourceGroup.Caption)

			responseBody, err := form.EmulatePost("/group/update", formData, groupHandler.Update)
			if err != nil {
				t.Fatal(err)
			}

			resultGroup := &model.Group{}
			if err := json.NewDecoder(responseBody).Decode(resultGroup); err != nil {
				t.Fatal(err)
			}

			if *sourceGroup != *resultGroup {
				t.Error(testify.NewObjectNotMatchedError(sourceGroup, resultGroup))
			}
		})

	t.Run("Delete",
		func(t *testing.T) {
			if sourceGroup == nil {
				t.Skip()
			}

			formData := url.Values{}
			formData.Set("ID", strconv.FormatInt(sourceGroup.ID, 10))

			_, err := form.EmulatePost("/group/delete", formData, groupHandler.Delete)
			if err != nil {
				t.Fatal(err)
			}
		})
}
