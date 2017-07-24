package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	cloudflare "github.com/cloudflare/cloudflare-go"
)

func main() {
	domains := os.Getenv("DNS_NAMES")
	if domains == "" {
		log.Fatal(fmt.Errorf("no env DNS_NAMES given"))
	}
	ipAddresses := os.Getenv("IP_ADDRESS")
	if ipAddresses == "" {
		log.Fatal(fmt.Errorf("no env IP_ADDRESS given"))
	}
	// Construct a new API object
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}

	// Fetch user details on the account
	u, err := api.UserDetails()
	if err != nil {
		log.Fatal(err)
	}
	// Print user details
	log.Infof("CF User info: %+v", u)

	recordIDs := map[string][]string{}

	for _, domain := range strings.Split(domains, ",") {
		// Fetch the zone ID
		zoneNameSplit := strings.Split(domain, ".")
		zoneName := strings.Join(zoneNameSplit[len(zoneNameSplit)-2:], ".")
		log.Infof("Domain Name: %+v", zoneName)
		id, err := api.ZoneIDByName(zoneName)
		if err != nil {
			log.Fatal(err)
		}
		ipAddressesSplit := strings.Split(ipAddresses, ",")
		for _, ipAddress := range ipAddressesSplit {
			records, err := api.DNSRecords(id, cloudflare.DNSRecord{
				Name:    domain,
				Content: ipAddress,
			})
			if err != nil {
				log.Fatal(err)
			}
			if len(records) == 0 {
				recordType := "A"
				if strings.Contains(ipAddress, ":") {
					recordType = "AAAA"
				}
				record, err := api.CreateDNSRecord(id, cloudflare.DNSRecord{
					Name:    domain,
					Content: ipAddress,
					Proxied: false,
					Type:    recordType,
				})
				if err != nil {
					log.Fatal(err)
				}
				recordIDs[id] = append(recordIDs[id], record.Result.ID)
				log.WithFields(log.Fields{
					"zoneID": id,
					"ip":     ipAddress,
				}).Infof("DNS Record created: '%s'", domain)
			} else {
				for _, record := range records {
					recordIDs[id] = append(recordIDs[id], record.ID)
				}
				log.WithFields(log.Fields{
					"zoneID": id,
					"ip":     ipAddress,
				}).Infof("DNS Record already exists: '%s'", domain)
			}
		}
	}
	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func(recordIDs map[string][]string) {
		sig := <-sigs
		log.Info(sig)
		for zoneID, records := range recordIDs {
			for _, recordID := range records {
				if err := api.DeleteDNSRecord(zoneID, recordID); err != nil {
					log.Fatal(err)
				}
				log.WithFields(log.Fields{
					"zoneID":   zoneID,
					"recordID": recordID,
				}).Infof("DNS Record deleted")
			}
			log.WithFields(log.Fields{
				"zoneID": zoneID,
			}).Infof("DNS Records for zone deleted")
		}
		done <- true
	}(recordIDs)

	log.Info("awaiting signal")
	<-done
	log.Info("exiting")
}
