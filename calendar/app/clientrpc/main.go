package main

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gzavodov/otus-go/calendar/app/rpc"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*60*time.Second)

	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer conn.Close()
	client := rpc.NewEventServiceClient(conn)

	startTime, err := ptypes.TimestampProto(time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		log.Fatalf("Cannot parse start time: %v", err)
	}

	endTime, err := ptypes.TimestampProto(time.Date(2019, 10, 1, 13, 0, 0, 0, time.UTC))
	if err != nil {
		log.Fatalf("Cannot parse end time: %v", err)
	}

	event := &rpc.Event{
		Title:     "Test",
		UserID:    1,
		StartTime: startTime,
		EndTime:   endTime,
	}
	event, err = client.Create(ctx, event)
	if err != nil {
		log.Fatalf("Cannot create event: %v", err)
	}

	log.Printf("Created successfully: %d\n", event.ID)
	cancel()
}
