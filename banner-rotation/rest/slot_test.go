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

func TestSlot(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	conf := &config.Configuration{}
	//OS variable BANNER_ROTATION_REPOSITORY_DSN is required
	if err := conf.LoadFromEvironment(); err != nil {
		t.Fatal(err)
	}

	slotRepo := sql.NewSlotRepository(ctx, conf.RepositoryDSN)
	slotUsecase := usecase.NewSlotUsecase(slotRepo)

	slotHandler := &EntityHandler{
		Accessor: &Slot{ucase: slotUsecase},
		Handler:  endpoint.Handler{Name: "Slot", ServiceName: "Test"},
	}

	var sourceSlot *model.Slot
	form := testify.Form{}

	t.Run("Create",
		func(t *testing.T) {
			caption := "Test Slot #1"

			formData := url.Values{}
			formData.Set("caption", caption)

			responseBody, err := form.EmulatePost("/slot/create", formData, slotHandler.Create)
			if err != nil {
				t.Fatal(err)
			}

			resultSlot := &model.Slot{}
			if err := json.NewDecoder(responseBody).Decode(resultSlot); err != nil {
				t.Fatal(err)
			}

			if resultSlot.ID < 0 {
				t.Errorf("handler returned unexpected banner ID: got %d", resultSlot.ID)
			}

			if resultSlot.Caption != caption {
				t.Errorf("handler returned unexpected banner caption: got %s want %s", resultSlot.Caption, caption)
			}

			sourceSlot = resultSlot
		})

	t.Run("Read",
		func(t *testing.T) {
			if sourceSlot == nil {
				t.Skip()
			}

			formData := url.Values{}
			formData.Set("ID", strconv.FormatInt(sourceSlot.ID, 10))

			responseBody, err := form.EmulatePost("/slot/read", formData, slotHandler.Read)
			if err != nil {
				t.Fatal(err)
			}

			resultSlot := &model.Slot{}
			if err := json.NewDecoder(responseBody).Decode(resultSlot); err != nil {
				t.Fatal(err)
			}

			if *sourceSlot != *resultSlot {
				t.Error(testify.NewObjectNotMatchedError(sourceSlot, resultSlot))
			}
		})

	t.Run("Update",
		func(t *testing.T) {
			if sourceSlot == nil {
				t.Skip()
			}

			sourceSlot.Caption = "Test Slot #2"

			formData := url.Values{}
			formData.Set("ID", strconv.FormatInt(sourceSlot.ID, 10))
			formData.Set("caption", sourceSlot.Caption)

			responseBody, err := form.EmulatePost("/slot/update", formData, slotHandler.Update)
			if err != nil {
				t.Fatal(err)
			}

			resultSlot := &model.Slot{}
			if err := json.NewDecoder(responseBody).Decode(resultSlot); err != nil {
				t.Fatal(err)
			}

			if *sourceSlot != *resultSlot {
				t.Error(testify.NewObjectNotMatchedError(sourceSlot, resultSlot))
			}
		})

	t.Run("Delete",
		func(t *testing.T) {
			if sourceSlot == nil {
				t.Skip()
			}

			formData := url.Values{}
			formData.Set("ID", strconv.FormatInt(sourceSlot.ID, 10))

			_, err := form.EmulatePost("/slot/delete", formData, slotHandler.Delete)
			if err != nil {
				t.Fatal(err)
			}
		})
}
