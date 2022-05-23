package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
func BoolPointer(b bool) *bool {
	return &b
}
func main() {
	api_token := os.Getenv("CF_API_TOKEN")
	if api_token == "" {
		log.Fatal("CF_API_TOKEN is not set")
	}

	zone_name := os.Getenv("CF_ZONE_NAME")
	if zone_name == "" {
		log.Fatal("CF_ZONE_NAME is not set")
	}

	record_prefix := os.Getenv("CF_RECORD_PREFIX")
	if record_prefix == "" {
		record_prefix = "controller"
	}
	proxy_mode := os.Getenv("CF_PROXY_MODE")
	if proxy_mode == "" {
		log.Fatal("CF_PROXY_MODE is not set")
	}
	CloudFlareProxyMode := false
	if proxy_mode == "true" {
		CloudFlareProxyMode = true
	} else {
		CloudFlareProxyMode = false
	}

	record_name := record_prefix + "." + zone_name

	api, err := cloudflare.NewWithAPIToken(api_token)
	if err != nil {
		log.Fatal("error creating api object", err)
	}

	// Most API calls require a Context
	ctx := context.Background()
	// Fetch the zone ID
	zone_id, err := api.ZoneIDByName(zone_name)
	if err != nil {
		log.Fatal("error retrieving zone id", err)
	}
	update_data := cloudflare.DNSRecord{
		Type:    "A",
		Name:    record_name,
		Content: GetOutboundIP().String(),
		TTL:     60,
		Proxied: BoolPointer(CloudFlareProxyMode),
	}

	// Fetch records for the requested record.
	needle := cloudflare.DNSRecord{Name: record_name}
	recs, err := api.DNSRecords(ctx, zone_id, needle)
	if err != nil {
		log.Fatal("error retrieving records: ", err)
	}
	if recs == nil {
		log.Println("no records found, will create.")
		create_record, create_err := api.CreateDNSRecord(ctx, zone_id, update_data)
		if create_err != nil {
			log.Fatal("error creating record: ", create_err)
		}
		log.Println("created record: ", create_record.Messages)

	} else {
		if recs[0].Content != update_data.Content {
			log.Println("record found, but IP is not correct, will update.")
			update_err := api.UpdateDNSRecord(ctx, zone_id, recs[0].ID, update_data)
			if update_err != nil {
				log.Fatal("error updating record: ", update_err)
			}
			log.Println("updated record: ", update_data.Content)
		} else {
			log.Println("record found, and IP is correct, ignoring.")
		}
	}
}
