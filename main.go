package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type OpTime struct {
	Timestamp primitive.Timestamp `bson:"ts"`
	Term      int64               `bson:"t"`
}

type HelloLastWrite struct {
	OpTime            OpTime    `bson:"opTime"`
	LastWriteDate     time.Time `bson:"lastWriteDate"`
	MajorityOpTime    OpTime    `bson:"majorityOpTime"`
	MajorityWriteDate time.Time `bson:"majorityWriteDate"`
}

type HelloResponseCore struct {
	IsWritablePrimary            bool      `bson:"isWritablePrimary"`
	MaxBsonObjectSize            int       `bson:"maxBsonObjectSize"`
	MaxMessageSizeBytes          int       `bson:"maxMessageSizeBytes"`
	MaxWriteBatchSize            int       `bson:"maxWriteBatchSize"`
	LocalTime                    time.Time `bson:"localTime"`
	LogicalSessionTimeoutMinutes int       `bson:"logicalSessionTimeoutMinutes"`
	ConnectionId                 int       `bson:"connectionId"`
	MinWireVersion               int       `bson:"minWireVersion"`
	MaxWireVersion               int       `bson:"maxWireVersion"`
	ReadOnly                     bool      `bson:"readOnly"`
	OK                           int       `bson:"ok"`
}

type HelloResponseReplicaSets struct {
	Hosts      []string           `bson:"hosts"`
	SetName    string             `bson:"setName"`
	SetVersion int                `bson:"setVersion"`
	Secondary  bool               `bson:"secondary"`
	Passives   []string           `bson:"passives"`
	Arbiters   []string           `bson:"arbiters"`
	Passive    bool               `bson:"passive"`
	Hidden     bool               `bson:"hidden"`
	Me         string             `bson:"me"`
	ElectionId primitive.ObjectID `bson:"electionId"`
	LastWrite  HelloLastWrite     `bson:"lastWrite"`
}

type HelloResponse struct {
	HelloResponseCore        `bson:",inline"`
	HelloResponseReplicaSets `bson:",inline"`
}

func createClient() (*mongo.Client, context.Context) {
	var err error
	port := 27017
	portString := os.Getenv("MONGODB_PORT_NUMBER")
	if portString != "" {
		if port, err = strconv.Atoi(portString); err != nil {
			log.Fatalf("invalid port specified: %v", err)
		}
	}

	clientOpts := options.Client().SetHosts([]string{fmt.Sprintf("localhost:%d", port)})
	client, err := mongo.NewClient(clientOpts)

	if err != nil {
		log.Fatalf("failed to create mongo client: %v", err)
	}

	ctx := context.Background()
	log.Println("attempting to connect")
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	return client, ctx
}

func livenessProbe() {
	c, ctx := createClient()

	log.Println("running ping")
	if err := c.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("ping failed: %v", err)
	}

	log.Println("ping successful")
}

func readinessProbe() {
	c, ctx := createClient()

	log.Println("running hello in admin database")
	db := c.Database("admin")
	result := db.RunCommand(ctx, bson.M{"hello": 1})

	if result.Err() != nil {
		log.Fatalf("hello failed: %v", result.Err())
	}

	hello := &HelloResponse{}
	err := result.Decode(hello)
	if err != nil {
		log.Fatalf("hello failed: %v", result.Err())
	}

	if !hello.IsWritablePrimary && !hello.Secondary {
		log.Fatalf("not writeable or secondary")
	}

	log.Println("hello successful")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("specify probe type to run")
		os.Exit(1)
	}

	probe := os.Args[1]
	switch probe {
	case "liveness":
		livenessProbe()
	case "readiness":
		readinessProbe()
	case "startup":
		readinessProbe()
	default:
		fmt.Printf("unknown probe type '%s', supported: liveness, readiness, startup\n", probe)
		os.Exit(1)
	}
}
