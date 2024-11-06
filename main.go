package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

// Animal represents the schema for the animal data.
type Animal struct {
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Species  string    `json:"species"`
	Family   string    `json:"family"`
	Gender   string    `json:"gender"`
	BornDate string `json:"born_date"`
}

func main() {
	// Define project ID and topic ID
	projectID := "zoo-poc"
	topicID := "animals-topic"

	// Initialize Pub/Sub client
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer client.Close()

	// Reference the Pub/Sub topic
	topic := client.Topic(topicID)
	defer topic.Stop()

	animals := []Animal{
		{"Lion", 5, "Panthera leo", "Felidae", "Male", "2018-06-01"},
		{"Elephant", 12, "Loxodonta", "Elephantidae", "Female", "2011-02-15"},
		{"Tiger", 7, "Panthera tigris", "Felidae", "Female", "2016-08-21"},
		{"Kangaroo", 3, "Macropus", "Macropodidae", "Male", "2020-04-12"},
		{"Panda", 6, "Ailuropoda melanoleuca", "Ursidae", "Female", "2017-09-23"},
		{"Wolf", 4, "Canis lupus", "Canidae", "Male", "2019-05-18"},
		{"Giraffe", 9, "Giraffa camelopardalis", "Giraffidae", "Female", "2014-11-05"},
		{"Penguin", 2, "Aptenodytes forsteri", "Spheniscidae", "Male", "2022-03-08"},
		{"Zebra", 8, "Equus quagga", "Equidae", "Female", "2015-12-29"},
		{"Cheetah", 5, "Acinonyx jubatus", "Felidae", "Male", "2018-07-13"},
	}


	// Publish each animal as a message
	for _, animal := range animals {
		data, err := json.Marshal(animal)
		if err != nil {
			log.Fatalf("Failed to marshal animal data: %v", err)
		}

		// Publish message
		result := topic.Publish(ctx, &pubsub.Message{
			Data: data,
		})

		// Check for publish errors
		_, err = result.Get(ctx)
		if err != nil {
			log.Printf("Failed to publish animal: %v", err)
		} else {
			fmt.Printf("Published animal: %v\n", string(data))
		}
	}
}
