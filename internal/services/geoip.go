package services

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var GeoDB *geoip2.Reader

func InitGeoIP() {
	var err error
	// Load the database from the current working directory
	GeoDB, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Printf("Warning: Failed to load GeoLite2-City.mmdb: %v. Geolocation will be disabled.", err)
	} else {
		log.Println("Loaded MaxMind GeoLite2-City database")
	}
}

func CloseGeoIP() {
	if GeoDB != nil {
		GeoDB.Close()
	}
}

// GetLocation returns the country and city for a given IP string.
func GetLocation(ipStr string) (country string, city string) {
	if GeoDB == nil || ipStr == "" {
		return "Unknown", "Unknown"
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "Unknown", "Unknown"
	}

	record, err := GeoDB.City(ip)
	if err != nil {
		return "Unknown", "Unknown"
	}

	country = record.Country.Names["en"]
	if country == "" {
		country = "Unknown"
	}

	city = record.City.Names["en"]
	if city == "" {
		city = "Unknown"
	}

	return country, city
}
