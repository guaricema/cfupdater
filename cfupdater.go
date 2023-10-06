package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type Config struct {
	ApiKey string `json:"api_key"`
	Zones []Zone `json:"zones"`
}

type Zone struct {
	Name string `json:"name"`
	Records []Record `json:"records"`
}

type Record struct {
	Record string `json:"record"`
	Type string `json:"type"`
	IP string `json:"ip"`
}

type IP struct {
    Query string
}

// main cycle
func main() {
	configFile, err := os.Open("config.json")
	defer configFile.Close()

	if err != nil { log.Fatal("Couldn't open config file at config.json!") }

	data, _ := ioutil.ReadAll(configFile)
	var config Config

	json.Unmarshal(data, &config)
	api_key := config.ApiKey

	if api_key == "" {
		log.Fatal("Nothing to do. Check your config file.")
	}

	// detect public ip
	current_ip := getIP()

	// go for all
	for _, zone := range config.Zones {
		for _, record := range zone.Records {
			ip := current_ip
			if record.IP != "%" { ip = record.IP }
			
			// update if possible
			updateCFRecord(api_key, zone.Name, record.Record, record.Type, ip)
		}
	}
}

// get your public IP
func getIP() string {
    req, err := http.Get("http://ip-api.com/json/")
    if err != nil { return err.Error() }
    defer req.Body.Close()

    body, err := ioutil.ReadAll(req.Body)
    if err != nil { return err.Error() }

    var ipddress IP
    json.Unmarshal(body, &ipddress)

    return ipddress.Query
}

// update DNS record using Cloudflare API v4
func updateCFRecord(key string, zone string, subdomain string, dtype string, ip string) {
	dname := subdomain + "." + zone;

	// init api
	api, err := cloudflare.NewWithAPIToken(key)
	if err != nil { log.Fatal(err) }

	ctx := context.Background()

	// get zone id
	zoneID, err := api.ZoneIDByName(zone)
	if err != nil { log.Fatal(err) }

	// get records with certain type & name to detect their ID
	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{ Name: dname, Type: dtype })
	if err != nil { log.Fatal(err) }

	// nothing was found
	if len(records) == 0 {
		fmt.Println("DNS record not found.")
		return 
	}

	// update every record we've got
	for _, r := range records {
		// check IP address
		if r.Content == ip {
			fmt.Println("Nothing to change (there is same IP address).")
			return 
		}

		// params used in UpdateDNSRecord
		params := cloudflare.UpdateDNSRecordParams{ Type: dtype, Name: dname, Content: ip, ID: r.ID }

		// update records
		_, err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), params)
		if err != nil { log.Fatal(err) }

		fmt.Printf("%s was updated by IP %s\n", dname, ip)
	}
}