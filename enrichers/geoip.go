package enrichers

import (
	"fmt"
	"net"

	geoip2 "github.com/oschwald/geoip2-golang"
	"github.com/sheacloud/goflow-addons/utils"
)

type GeoIPMetadata struct {
	CityName            string        `json:",omitempty"`
	ContinentCode       string        `json:",omitempty"`
	ContinentName       string        `json:",omitempty"`
	CountryIsoCode      string        `json:",omitempty"`
	CountryName         string        `json:",omitempty"`
	CountryInEU         bool          `json:",omitempty"`
	Latitude            float64       `json:",omitempty"`
	Longitude           float64       `json:",omitempty"`
	MetroCode           uint          `json:",omitempty"`
	TimeZone            string        `json:",omitempty"`
	PostalCode          string        `json:",omitempty"`
	Subdivisions        []Subdivision `json:",omitempty"`
	IsAnonymousProxy    bool          `json:",omitempty"`
	IsSatelliteProvider bool          `json:",omitempty"`
}

type Subdivision struct {
	IsoCode string
	Name    string
}

type GeoIPEnricher struct {
	Language string
	db       *geoip2.Reader
}

func (e *GeoIPEnricher) Initialize() {
	db, err := geoip2.Open("./geoip-database/GeoLite2-City.mmdb")
	if err != nil {
		fmt.Println(err)
		return
	}
	e.db = db
}

func (e *GeoIPEnricher) FlattenCity(city *geoip2.City) GeoIPMetadata {
	metadata := GeoIPMetadata{}
	if city.City.Names != nil {
		metadata.CityName = city.City.Names[e.Language]
	}
	if city.Continent.Names != nil {
		metadata.ContinentCode = city.Continent.Code
		metadata.ContinentName = city.Continent.Names[e.Language]
	}
	if city.Country.Names != nil {
		metadata.CountryIsoCode = city.Country.IsoCode
		metadata.CountryName = city.Country.Names[e.Language]
		metadata.CountryInEU = city.Country.IsInEuropeanUnion
	}
	if city.Location.TimeZone != "" {
		metadata.Latitude = city.Location.Latitude
		metadata.Longitude = city.Location.Longitude
		metadata.MetroCode = city.Location.MetroCode
		metadata.TimeZone = city.Location.TimeZone
	}
	if city.Postal.Code != "" {
		metadata.PostalCode = city.Postal.Code
	}
	subdivisions := []Subdivision{}
	for _, sub := range city.Subdivisions {
		subdivision := Subdivision{
			IsoCode: sub.IsoCode,
			Name:    sub.Names[e.Language],
		}
		subdivisions = append(subdivisions, subdivision)
	}
	if subdivisions != nil {
		metadata.Subdivisions = subdivisions
	}

	return metadata
}

func (e *GeoIPEnricher) Enrich(msgs []*utils.ExtendedFlowMessage) {
	for _, msg := range msgs {
		srcIP := net.IP(msg.SrcAddr)
		dstIP := net.IP(msg.DstAddr)

		srcCityData, _ := e.db.City(srcIP)
		dstCityData, _ := e.db.City(dstIP)

		msg.Metadata["SrcGeoIPData"] = e.FlattenCity(srcCityData)
		msg.Metadata["DstGeoIPData"] = e.FlattenCity(dstCityData)
	}
}
