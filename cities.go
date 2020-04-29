package main

import "os"
import "github.com/gocarina/gocsv"

type City struct {
	City       string `csv:"city"`
	CityAscii  string `csv:"city_ascii"`
	Lat        string `csv:"lat"`
	Lng        string `csv:"lng"`
	Country    string `csv:"country"`
	Iso2       string `csv:"iso2"`
	Iso3       string `csv:"iso3"`
	AdminName  string `csv:"admin_name"`
	Capital    string `csv:"capital"`
	Population string `csv:"population"`
	Id         string `csv:"id"`
}

func loadCities() []City {

	clientsFile, err := os.OpenFile("worldcities.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	cities := []City{}

	if err := gocsv.UnmarshalFile(clientsFile, &cities); err != nil { // Load cities from file
		panic(err)
	}

	return cities
}

func createCities(n int) []string {
	c := make([]string, n)

	all := loadCities()

	for i := 0; i < len(c); i++ {
		c[i] = all[i].City
	}

	return c
}
