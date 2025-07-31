package utils

import (
	// standard Packages

	"log"
	"sort"
	"strconv"
	"strings"
)

// read list of times provided in the config and convert to a list of ints
// each values represented a number of days left

func ParseNotifTimes(str string) []int {
	var intSlice []int

	if str == "" {
		return []int{}
	}

	// use fields() and join() to get rid of all whitespace
	// split by delimeter ,
	// try to convert each value to an int, error out if failed
	// sort final slice of ints

	parsedString := strings.Split(strings.Join(strings.Fields(str), ""), ",")
	for _, val := range parsedString {
		converted, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("[ERROR] Failed to parse notification time: %s", err)
		}
		intSlice = append(intSlice, converted)
	}

	sort.Ints(intSlice)

	return intSlice
}

// read grace period value provided in the config and convert it to an int

func ParseGracePeriod(value string) int {
	// Atoi means ASCII to Integer

	days, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("[ERROR] Failed to parse grace period value: %s", err)
	} else if days < 1 {
		log.Fatal("[ERROR] For safety reasons, grace period cannot be lower than one day.")
	}
	return days
}

func ParseStrList(str string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(strings.Join(strings.Fields(str), ""), ",")
}
