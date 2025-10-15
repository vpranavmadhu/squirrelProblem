package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
)

type Entry struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

type Counts struct {
	n00 uint
	n01 uint
	n10 uint
	n11 uint
}

type BestWorse struct {
	bestEvent  string
	bestCorr   float64
	worseEvent string
	worseCorr  float64
}

func preprocess(journalEntries []Entry) []Entry {
	var journal []Entry
	for _, entry := range journalEntries {
		hasPeanut := slices.Contains(entry.Events, "peanuts")
		notBrushedTeeth := !slices.Contains(entry.Events, "brushed teeth")

		if hasPeanut && notBrushedTeeth {
			entry.Events = append(entry.Events, "dirty teeth")
			journal = append(journal, entry)
		} else {
			journal = append(journal, entry)
		}
	}
	return journal
}

func phi(counts Counts) float64 {

	n00 := float64(counts.n00)
	n01 := float64(counts.n01)
	n10 := float64(counts.n10)
	n11 := float64(counts.n11)

	num := (n11 * n00) - (n10 * n01)
	nx1 := n11 + n10
	nx0 := n00 + n01
	ny1 := n11 + n01
	ny0 := n10 + n00

	den := math.Sqrt(nx0 * nx1 * ny0 * ny1)
	if den == 0 {
		return 0
	}

	return num / den
}

func getCounts(entries []Entry, event string) Counts {

	var n00, n01, n10, n11 uint

	for _, entry := range entries {

		if slices.Contains(entry.Events, event) { // event true
			if entry.Squirrel { //event true & squirrel true
				n11++
			} else { //event true && sqirrel false
				n10++
			}
		} else { //event false
			if entry.Squirrel { //event false & squirrel true
				n01++
			} else { //both false
				n00++
			}
		}

	}

	counts := Counts{
		n00: n00,
		n01: n01,
		n10: n10,
		n11: n11,
	}

	return counts

}

func getCorrelations(journal []Entry) map[string]float64 {
	correlations := make(map[string]float64)

	for _, entries := range journal {
		for _, event := range entries.Events {
			counts := getCounts(journal, event)
			corr := phi(counts)
			correlations[event] = corr
		}
	}

	return correlations
}

func getBestWorseCorrelation(correlations map[string]float64) BestWorse {
	var bestEvent string
	var bestCorr float64
	var worseEvent string
	var worseCorr float64

	for key, value := range correlations {

		if value > bestCorr {
			bestCorr = value
			bestEvent = key
		}

		if value < worseCorr {
			worseCorr = value
			worseEvent = key
		}

	}

	result := BestWorse{
		bestEvent:  bestEvent,
		bestCorr:   bestCorr,
		worseEvent: worseEvent,
		worseCorr:  worseCorr,
	}
	return result
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

	journal = preprocess(journal)

	correlations := getCorrelations(journal)

	wb := getBestWorseCorrelation(correlations)

	fmt.Println("Most correlated: ", wb.bestCorr, "for event: ", wb.bestEvent)
	fmt.Println("Least correlated: ", wb.worseCorr, "for event: ", wb.worseEvent)

}
