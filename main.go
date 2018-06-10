package main

import (
	"log"
	"sort"
	"time"

	"github.com/satori/go.uuid"
)

// Event hold the schedule
type Event struct {
	ID    uuid.UUID
	Start string
	End   string
}

// EventTime hold the schedule
type EventTime struct {
	ID    uuid.UUID
	Start time.Time
	End   time.Time
}

func main() {
	schedule := []Event{}
	overlappingEvents(schedule)
}

func overlappingEvents(arr []Event) map[uuid.UUID][]Event {
	eventTimes := []EventTime{}
	layout := "2006-01-02 15:04 MST"

	for _, v := range arr {
		s, err := time.Parse(layout, v.Start)
		if err != nil {
			log.Fatal("error parsing start date", err)
		}

		e, err := time.Parse(layout, v.End)
		if err != nil {
			log.Fatal("error parsing end date", err)
		}

		eventTimes = append(eventTimes, EventTime{v.ID, s, e})
	}

	sort.SliceStable(eventTimes, func(i, j int) bool { return eventTimes[i].Start.Before(eventTimes[j].Start) })

	overlapEvents := map[uuid.UUID][]Event{}

	for i := 0; i < len(eventTimes)-1; i++ {
		baseID := eventTimes[i].ID

		if len(overlapEvents[baseID]) == 0 {
			overlapEvents[baseID] = []Event{}
		}

		for j := i + 1; j < len(eventTimes); j++ {
			if isOverlapping(eventTimes[i].End, eventTimes[j].Start) {
				currEvent := eventTimes[j]
				// append only if not present already to avoid duplication
				if !isPresent(overlapEvents[baseID], eventTimes[j]) {
					overlapEvents[baseID] = append(overlapEvents[baseID], Event{currEvent.ID, currEvent.Start.Format(layout), currEvent.End.Format(layout)})
				}

				// append only if not present already to avoid duplication
				if !isPresent(overlapEvents[currEvent.ID], eventTimes[i]) {
					overlapEvents[currEvent.ID] = append(overlapEvents[currEvent.ID], Event{baseID, eventTimes[i].Start.Format(layout), eventTimes[i].End.Format(layout)})
				}
			}
		}
	}

	return overlapEvents
}

func isPresent(arr []Event, e EventTime) bool {
	for _, v := range arr {
		if v.ID == e.ID {
			return true
		}
	}
	return false
}

func isOverlapping(a, b time.Time) bool {
	if a.Before(b) {
		return false
	}
	return true
}
