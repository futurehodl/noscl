package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func showProfile(opts docopt.Opts) {
	verbose, _ := opts.Bool("--verbose")
	jsonformat, _ := opts.Bool("--json")
	key := nip19.TranslatePublicKey(opts["<pubkey>"].(string))
	if key == "" {
		log.Println("Profile key is empty! Exiting.")
		return
	}

	initNostr()

	_, all := pool.Sub(nostr.Filters{{Authors: []string{key}, Kinds: []int{1}}})
	for event := range nostr.Unique(all) {
		printEvent(event, nil, verbose, jsonformat)
	}
}

func profileToTxtFile(opts docopt.Opts) {
	// verbose, _ := opts.Bool("--verbose")
	// jsonformat, _ := opts.Bool("--json")
	key := nip19.TranslatePublicKey(opts["<pubkey>"].(string))
	if key == "" {
		log.Println("Profile key is empty! Exiting.")
		return
	}

	initNostr()

	_, all := pool.Sub(nostr.Filters{{Authors: []string{key}, Kinds: []int{1}}})

	// process new events and add them to the text file
	for event := range nostr.Unique(all) {
		// printEvent(event, nil, verbose, jsonformat)
		log.Println(event.ID)
		appendEventToFile(event)
	}
}

func appendEventToFile(event nostr.Event) {
	// Read existing events from file
	data, err := ioutil.ReadFile(config.EventFilepath)
	if err != nil {
		log.Fatal(err)
	}

	var existingEvents []nostr.Event

	// Unmarshal existing events from JSON data
	if len(data) > 0 {
		err = json.Unmarshal(data, &existingEvents)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Check if event id already exists
	idExists := false
	for _, e := range existingEvents {
		if e.ID == event.ID {
			idExists = true
			break
		}
	}

	// Append event if id does not exist
	if !idExists {
		fmt.Println("Event ID", event.ID, "does not exist in file. Appending...")
		existingEvents = append(existingEvents, event)
		//printEvent(event, nil, true, true)
	} else {
		fmt.Println("Event ID", event.ID, "already exists in file. Skipping...")
	}

	// Write updated events to file
	updatedData, err := json.Marshal(existingEvents)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(config.EventFilepath, updatedData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Event appended to file successfully!")
}

func follow(opts docopt.Opts) {
	key := nip19.TranslatePublicKey(opts["<pubkey>"].(string))
	if key == "" {
		log.Println("Follow key is empty! Exiting.")
		return
	}

	name, err := opts.String("--name")
	if err != nil {
		name = ""
	}

	config.Following[key] = Follow{
		Key:  key,
		Name: name,
	}
	fmt.Printf("Followed %s.\n", key)
}

func unfollow(opts docopt.Opts) {
	key := nip19.TranslatePublicKey(opts["<pubkey>"].(string))
	if key == "" {
		log.Println("No unfollow key provided! Exiting.")
		return
	}

	delete(config.Following, key)
	fmt.Printf("Unfollowed %s.\n", key)
}

func following(opts docopt.Opts) {
	if len(config.Following) == 0 {
		fmt.Println("You aren't following anyone yet.")
		return
	}
	for _, profile := range config.Following {
		fmt.Println(profile.Key, profile.Name)
	}
}
