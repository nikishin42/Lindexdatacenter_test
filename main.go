package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/bcicen/jstream"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("wrong count of args")
	}
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if strings.HasSuffix(strings.ToLower(fileName), ".json") {
		parseJSON(file)
	} else if strings.HasSuffix(strings.ToLower(fileName), ".csv") {
		parseCSV(file)
	} else {
		log.Fatalln("unknown file format")
	}
}

func parseJSON(file *os.File) {
	var (
		price, rating, maxPrice, maxRating float64
		name, maxPriceName, maxRatingName  string
		fileNotEmpty                       bool
	)
	decoder := jstream.NewDecoder(file, 2).EmitKV()
	for mv := range decoder.Stream() {
		product, ok := mv.Value.(jstream.KV)
		if !ok {
			log.Fatalln("cast error")
		}
		switch product.Key {
		case "name":
			name, ok = product.Value.(string)
			if !ok {
				log.Fatalln("cast error")
			}
		case "price":
			price, ok = product.Value.(float64)
			if !ok {
				log.Fatalln("cast error")
			}
			if price > maxPrice {
				maxPrice = price
				maxPriceName = name
			}
		case "rating":
			rating, ok = product.Value.(float64)
			if !ok {
				log.Fatalln("cast error")
			}
			if rating > maxRating {
				maxRating = rating
				maxRatingName = name
			}
		default:
			log.Fatalln("unknown json field")
		}
		fileNotEmpty = true
	}
	printResult(maxPrice, maxRating, maxPriceName, maxRatingName, fileNotEmpty)
}

func parseCSV(file *os.File) {
	var (
		maxPrice, maxRating         float64
		maxPriceName, maxRatingName string
		fileNotEmpty                bool
	)

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = 3
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalln(err)
		}
		price, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatalln(err)
		}
		rating, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatalln(err)
		}
		if price > maxPrice {
			maxPrice = price
			maxPriceName = record[0]
		}
		if rating > maxRating {
			maxRating = rating
			maxRatingName = record[0]
		}
		fileNotEmpty = true
	}

	printResult(maxPrice, maxRating, maxPriceName, maxRatingName, fileNotEmpty)
}

func printResult(maxPrice, maxRating float64, maxPriceName, maxRatingName string, fileNotEmpty bool) {
	if fileNotEmpty {
		fmt.Printf("самый дорогой продукт: %s [%.2f]\n", maxPriceName, maxPrice)
		fmt.Printf("продукт с самым высоким рейтингом: %s [%.1f]\n", maxRatingName, maxRating)
	} else {
		fmt.Println("Нет данных")
	}
}
