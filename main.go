package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	// http://www.meetup.com/golang-syd/events/75463782/
	eventid = 75463782

	// how many exchanges
	shuffles = 9001
)

// filter names from the json blob meetup gives us
func filter(m map[string]interface{}) (result []string) {
	for _, v := range m["results"].([]interface{}) {
		result = append(result, v.(map[string]interface{})["member"].(map[string]interface{})["name"].(string))
	}
	return
}

// shuffle the deck
func shuffle(punters []string) {
	for i := 0; i < shuffles; i++ {
		a, b := rand.Intn(len(punters)), rand.Intn(len(punters))
		punters[a], punters[b] = punters[b], punters[a]
	}
}

// because things are funnier when they are piped through figlet
func figlet(name string) error {
	cmd := exec.Command("/usr/bin/figlet", name)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func main() {
	log.Println("Contacting meetup.com")
	resp, err := http.Get(fmt.Sprintf("http://api.meetup.com/2/rsvps?key=%s&event_id=%d&order=name&rsvp=yes", os.Getenv("MEETUP_API_KEY"), eventid))
	if err != nil {
		log.Fatalf("Could not contact meetup.com: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad http response %v", resp.Status)
	}
	result := make(map[string]interface{})
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		log.Fatalf("Aww snap, that doesn't look like json: %v", err)
	}
	log.Println("Collating results")
	time.Sleep(1500 * time.Millisecond)
	punters := filter(result)

	log.Println("Seeding random numbers")
	time.Sleep(1 * time.Second)
	rand.Seed(time.Now().Unix())

	log.Println("Shuffling")
	time.Sleep(2 * time.Second)
	shuffle(punters)

	log.Println("Pausing for dramatic effect")
	time.Sleep(3 * time.Second)

	fmt.Println("And the winner is ...")
	time.Sleep(1500 * time.Millisecond)
	figlet(punters[0])
}
