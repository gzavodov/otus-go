package rpc

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gzavodov/otus-go/calendar/app/logger"
	"github.com/gzavodov/otus-go/calendar/app/sqldb"
	"google.golang.org/grpc"
)

func TestGRPCService(t *testing.T) {
	appLogger, err := logger.Create("../../log.json", "debug")
	if err != nil {
		t.Fatalf("Could not initialize logger: %v", err)
	}
	defer appLogger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dataSourceName, ok := os.LookupEnv("CALENDAR_REPOSITORY_DSN")
	if !ok {
		t.Fatal("The environment variable CALENDAR_REPOSITORY_DSN is reqiured")
	}
	appRepo := sqldb.NewEventRepository(ctx, dataSourceName)

	serverAddress := "127.0.0.1:9090"
	server := NewServer(serverAddress, appRepo, appLogger)

	go func() {
		err := server.Start()
		if err != nil {
			t.Fatalf("Could not start server: %v", err)
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			server.Stop()
		}
	}()

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()
	client := NewEventServiceClient(conn)

	startTime, err := ptypes.TimestampProto(time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Could not parse start time: %v", err)
	}

	endTime, err := ptypes.TimestampProto(time.Date(2019, 10, 1, 13, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Could not parse end time: %v", err)
	}

	event := &Event{
		Title:        "RPC Test Event (2019-10-01T12:00:00)",
		UserID:       1,
		StartTime:    startTime,
		EndTime:      endTime,
		NotifyBefore: int64(30 * time.Minute),
	}

	event, err = client.Create(ctx, event)
	if err != nil {
		t.Fatalf("Could not create event: %v", err)
	}
	t.Log("Event was created successfully")

	event, err = client.Read(ctx, &EventIdentifier{Value: event.ID})
	if err != nil {
		t.Fatalf("Could not read event: %v", err)
	}
	t.Log("Event was received successfully")

	event.Description = "Test event detail description"
	event, err = client.Update(ctx, event)
	if err != nil {
		t.Fatalf("Could not update event: %v", err)
	}
	t.Log("Event was modified successfully")

	_, err = client.Delete(ctx, &EventIdentifier{Value: event.ID})
	if err != nil {
		t.Fatalf("Could not delete event: %v", err)
	}
	t.Log("Event was removed successfully")
}
