package main

import (
	"encoding/json"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
)

const eventid = 75463782

var log = stdlog.New(os.Stdout, os.Args[0], 0)

func filter(m map[string]interface{}) (result []string) {
	for _, v := range m["results"].([]interface{}) {
		result = append(result, v.(map[string]interface{})["member"].(map[string]interface{})["name"].(string))
	}
	return 
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
	punters := filter(result)
	fmt.Println(punters)
}
