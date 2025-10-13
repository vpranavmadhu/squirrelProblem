package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Entry struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

func phi(n11, n00, n10, n01 float64) float64 {
	num := (n11 * n00) - (n10 * n01)
	den := math.Sqrt((n11 + n10) * (n00 + n01) * (n11 + n01) * (n10 + n00))
	if den == 0 {
		return 0
	}
	return num / den
}

func main() {
	data, err := os.ReadFile("journal.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var journal []Entry
	err = json.Unmarshal(data, &journal)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	eventSet := make(map[string]bool)

	for _, e := range journal {
		for _, ev := range e.Events {
			eventSet[ev] = true
		}
	}

	var bestEvent string
	var bestCorr float64
	var worseCorr float64
	var worseEvent string

	for event := range eventSet {
		var n11, n00, n10, n01 float64

		for _, entry := range journal {
			hasEvent := false
			for _, ev := range entry.Events {
				if ev == event {
					hasEvent = true
					break
				}
			}

			if hasEvent && entry.Squirrel {
				n11++
			} else if hasEvent && !entry.Squirrel {
				n10++
			} else if !hasEvent && entry.Squirrel {
				n01++
			} else {
				n00++
			}
		}
		corr := phi(n11, n00, n10, n01)
		fmt.Printf("%s ----> %.4f \n", event, corr)

		if corr > bestCorr {
			bestCorr = corr
			bestEvent = event
		}

		if corr < worseCorr {
			worseCorr = corr
			worseEvent = event
		}
	}

	fmt.Printf("Most correlated event: %s (ϕ = %.4f)\n", bestEvent, bestCorr)
	fmt.Printf("Least correlated event: %s (ϕ = %.4f)\n", worseEvent, worseCorr)

}
