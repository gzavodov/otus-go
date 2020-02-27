package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/cucumber/godog/gherkin"
	"github.com/gzavodov/otus-go/banner-rotation/config"
	"github.com/gzavodov/otus-go/banner-rotation/internal/rabbitmq"
	"github.com/gzavodov/otus-go/banner-rotation/internal/testify"
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/queue"
	"golang.org/x/sync/errgroup"
)

//DefaultNotificationTimeout Notification Timeout
const DefaultNotificationTimeout = 5

type BannerChoise struct {
	BannerID int64
	SlotID   int64
	GroupID  int64
}

//NewFeatureTest creates new banner rotation test according to configuration
func NewFeatureTest(ctx context.Context, conf *config.Configuration) (*FeatureTest, error) {
	queueChannel := rabbitmq.NewChannel(ctx, conf.AMPQName, conf.AMPQAddress)
	notifications := testify.NewNotificationReceiver()

	notificationListener := queue.NewNotificationClient(
		ctx,
		queueChannel,
		notifications,
		nil,
	)

	return &FeatureTest{
			HTTPAddress:          conf.HTTPAddress,
			NotificationListener: notificationListener,
			Notifications:        notifications,
			NotificationTimeout:  DefaultNotificationTimeout,
		},
		nil
}

type FeatureTest struct {
	HTTPAddress          string
	Notifications        *testify.NotificationReceiver
	NotificationListener *queue.NotificationClient
	NotificationTimeout  int
	Choises              []*BannerChoise
}

//Start run services and clients are required for test process
func (t *FeatureTest) Start(outline *gherkin.Feature) {
	go func(client *queue.NotificationClient) {
		if client == nil {
			return
		}

		if err := client.Start(); err != nil {
			log.Fatalf("failed to start notification client: %v", err)
		}
	}(t.NotificationListener)
}

//Stop halt services and clients are required for test process
func (t *FeatureTest) Stop(outline *gherkin.Feature) {
	go func(client *queue.NotificationClient) {
		if client != nil {
			client.Stop()
		}
	}(t.NotificationListener)
}

func (t *FeatureTest) WaitForBannerChoiseNotification() error {
	for _, choise := range t.Choises {
		if err := t.WaitForNotification(choise.BannerID, queue.EventChoice); err != nil {
			return err
		}
	}
	return nil
}

func (t *FeatureTest) WaitForBannerClickNotification() error {
	for _, choise := range t.Choises {
		if err := t.WaitForNotification(choise.BannerID, queue.EventClick); err != nil {
			return err
		}
	}
	return nil
}

//WaitForNotification waits for banner event notification
func (t *FeatureTest) WaitForNotification(bannerID int64, eventType string) error {
	return t.Notifications.Wait(bannerID, eventType, t.NotificationTimeout)
}

//RunEntityActionWithResult
func (t *FeatureTest) RunEntityActionWithResult(entityType string, actionName string, data url.Values, result interface{}) error {
	form := testify.Form{}
	response, err := form.Post("http://"+t.HTTPAddress+"/"+entityType+"/"+actionName, data)
	if err != nil {
		return err
	}

	return json.Unmarshal(response, result)

}

//RunEntityAction
func (t *FeatureTest) RunEntityAction(entityType string, actionName string, data url.Values) error {
	form := testify.Form{}
	if _, err := form.Post("http://"+t.HTTPAddress+"/"+entityType+"/"+actionName, data); err != nil {
		return err
	}

	return nil
}

//CreateEntityFromTable creates entity collection by the data table
func (t *FeatureTest) CreateEntityFromTable(entityType string, captionColumnName string, table *gherkin.DataTable) error {
	for i := 1; i < len(table.Rows); i++ {
		data := url.Values{}
		data.Set(captionColumnName, table.Rows[i].Cells[0].Value)

		if err := t.RunEntityAction(entityType, "create", data); err != nil {
			return err
		}
	}
	return nil
}

//CreateSocialGroupsFromTable creates banner social groups collection by the data table
func (t *FeatureTest) CreateSocialGroupsFromTable(table *gherkin.DataTable) error {
	return t.CreateEntityFromTable("group", "caption", table)
}

//CreateSlotsFromTable creates banner slots collection by the data table
func (t *FeatureTest) CreateSlotsFromTable(table *gherkin.DataTable) error {
	for i := 1; i < len(table.Rows); i++ {
		data := url.Values{}
		data.Set("caption", table.Rows[i].Cells[0].Value)

		err := t.RunEntityAction("slot", "create", data)
		if err != nil {
			return err
		}
	}

	return nil
}

//CreateBannersFromTable
func (t *FeatureTest) CreateBannersFromTable(table *gherkin.DataTable) error {
	for i := 1; i < len(table.Rows); i++ {
		data := url.Values{}
		data.Set("caption", table.Rows[i].Cells[0].Value)

		banner := &model.Banner{}
		if err := t.RunEntityActionWithResult("banner", "create", data, banner); err != nil {
			return err
		}

		data.Set("caption", table.Rows[i].Cells[1].Value)
		slot := &model.Slot{}
		if err := t.RunEntityActionWithResult("slot", "get-by-caption", data, slot); err != nil {
			return err
		}

		data = url.Values{}
		data.Set("bannerId", strconv.FormatInt(banner.ID, 10))
		data.Set("slotId", strconv.FormatInt(slot.ID, 10))
		err := t.RunEntityAction("banner", "add-to-slot", data)
		if err != nil {
			return err
		}
	}

	return nil
}

//VerifyEntityFromTable ensures that entities specified by the data table exits
func (t *FeatureTest) VerifyEntityFromTable(entityType string, table *gherkin.DataTable) error {
	for i := 1; i < len(table.Rows); i++ {
		data := url.Values{}
		data.Set("caption", table.Rows[i].Cells[0].Value)

		if err := t.RunEntityAction(entityType, "get-by-caption", data); err != nil {
			return err
		}
	}
	return nil
}

//VerifySocialGroupsFromTable ensures that banner social groups specified by the data table exit
func (t *FeatureTest) VerifySocialGroupsFromTable(table *gherkin.DataTable) error {
	return t.VerifyEntityFromTable("group", table)
}

//VerifySlotsFromTable ensures that banner slots specified by the data table exit
func (t *FeatureTest) VerifySlotsFromTable(table *gherkin.DataTable) error {
	return t.VerifyEntityFromTable("slot", table)
}

func (t *FeatureTest) VerifyBannersFromTable(table *gherkin.DataTable) error {
	for i := 1; i < len(table.Rows); i++ {
		data := url.Values{}

		data.Set("caption", table.Rows[i].Cells[0].Value)
		banner := &model.Banner{}
		if err := t.RunEntityActionWithResult("banner", "get-by-caption", data, banner); err != nil {
			return err
		}

		data.Set("caption", table.Rows[i].Cells[1].Value)
		slot := &model.Slot{}
		if err := t.RunEntityActionWithResult("slot", "get-by-caption", data, slot); err != nil {
			return err
		}

		data.Del("caption")
		data.Set("bannerId", strconv.FormatInt(banner.ID, 10))
		data.Set("slotId", strconv.FormatInt(slot.ID, 10))
		isInSlot := false
		if err := t.RunEntityActionWithResult("banner", "is-in-slot", data, &isInSlot); err != nil {
			return err
		}

		if !isInSlot {
			return errors.New("could not find banner slot")
		}

	}
	return nil
}

//ChooseBanner selects banners for show by the data table
func (t *FeatureTest) ChooseBanner(table *gherkin.DataTable) error {
	g, _ := errgroup.WithContext(context.Background())
	for i := 1; i < len(table.Rows); i++ {
		data := url.Values{}

		data.Set("caption", table.Rows[i].Cells[0].Value)
		slot := &model.Slot{}
		if err := t.RunEntityActionWithResult("slot", "get-by-caption", data, slot); err != nil {
			return err
		}

		data.Set("caption", table.Rows[i].Cells[1].Value)
		group := &model.Group{}
		if err := t.RunEntityActionWithResult("group", "get-by-caption", data, group); err != nil {
			return err
		}

		data.Set("slotId", strconv.FormatInt(slot.ID, 10))
		data.Set("groupId", strconv.FormatInt(group.ID, 10))

		g.Go(
			func() error {
				var bannerID int64
				if err := t.RunEntityActionWithResult("banner", "choose", data, &bannerID); err != nil {
					return err
				}

				if bannerID <= 0 {
					return errors.New("failed to get banner for show")
				}

				t.Choises = append(
					t.Choises,
					&BannerChoise{BannerID: bannerID, SlotID: slot.ID, GroupID: group.ID},
				)

				return nil
			},
		)
		time.Sleep(15 + time.Duration(rand.Int31n(75))*time.Millisecond)
	}
	return g.Wait()
}

//RegisterBannerClick selects banners for show by the data table
func (t *FeatureTest) RegisterBannerClick() error {
	for _, choise := range t.Choises {
		data := url.Values{}

		data.Set("bannerId", strconv.FormatInt(choise.BannerID, 10))
		data.Set("groupId", strconv.FormatInt(choise.GroupID, 10))

		if err := t.RunEntityAction("banner", "register-click", data); err != nil {
			return err
		}
	}
	return nil
}
