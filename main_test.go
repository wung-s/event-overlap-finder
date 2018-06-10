package main

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

func TestOverlappingEvents(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		events := []Event{}

		got := overlappingEvents(events)
		if len(got) != 0 {
			t.Errorf("got %v, wanted: %v", got, map[uuid.UUID][]Event{})
		}
	})

	t.Run("no overlaps", func(t *testing.T) {
		id1 := uuid.Must(uuid.NewV4())
		id2 := uuid.Must(uuid.NewV4())
		id3 := uuid.Must(uuid.NewV4())

		e1 := Event{id1, "2018-06-09 23:16 UTC", "2018-06-09 23:20 UTC"}
		e2 := Event{id2, "2018-06-09 11:00 UTC", "2018-06-09 12:00 UTC"}
		e3 := Event{id3, "2018-06-09 01:00 UTC", "2018-06-09 02:00 UTC"}
		events := []Event{e1, e2, e3}

		want := map[uuid.UUID][]Event{}
		want[id1] = []Event{}
		want[id2] = []Event{}
		want[id3] = []Event{}

		got := overlappingEvents(events)
		for _, v := range got {
			if len(v) != 0 {
				t.Errorf("got %v, wanted: %v", got, want)
			}
		}
	})

	t.Run("single overlaps", func(t *testing.T) {
		id1 := uuid.Must(uuid.NewV4())
		id2 := uuid.Must(uuid.NewV4())
		id3 := uuid.Must(uuid.NewV4())

		e1 := Event{id1, "2018-06-09 01:00 UTC", "2018-06-09 02:00 UTC"}
		e2 := Event{id2, "2018-06-09 23:16 UTC", "2018-06-09 23:20 UTC"}
		e3 := Event{id3, "2018-06-09 23:19 UTC", "2018-06-09 23:30 UTC"}
		events := []Event{e1, e2, e3}

		got := overlappingEvents(events)
		want := map[uuid.UUID][]Event{}
		want[id1] = []Event{}
		want[id2] = []Event{e3}
		want[id3] = []Event{e2}

		if !groupMatches(want, got) {
			t.Errorf("got %v, wanted: %v", got, want)
		}
	})

	t.Run("multiple overlaps", func(t *testing.T) {
		id1 := uuid.Must(uuid.NewV4())
		id2 := uuid.Must(uuid.NewV4())
		id3 := uuid.Must(uuid.NewV4())
		id4 := uuid.Must(uuid.NewV4())

		e1 := Event{id1, "2018-06-09 01:00 UTC", "2018-06-09 02:00 UTC"}
		e2 := Event{id2, "2018-06-09 23:16 UTC", "2018-06-09 23:40 UTC"}
		e3 := Event{id3, "2018-06-09 23:19 UTC", "2018-06-09 23:55 UTC"}
		e4 := Event{id4, "2018-06-09 23:20 UTC", "2018-06-09 23:30 UTC"}
		events := []Event{e1, e2, e3, e4}

		got := overlappingEvents(events)
		want := map[uuid.UUID][]Event{}
		want[id1] = []Event{}
		want[id2] = []Event{e3, e4}
		want[id3] = []Event{e2, e4}
		want[id4] = []Event{e2, e3}

		if !groupMatches(want, got) {
			t.Errorf("got %v, wanted: %v", got, want)
		}
	})
}

func groupMatches(a, b map[uuid.UUID][]Event) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		if !equal(a[k], b[k]) {
			return false
		}
	}
	return true
}

func equal(a []Event, b []Event) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].ID != b[i].ID {
			return false
		}
		if a[i].Start != b[i].Start || a[i].End != b[i].End {
			return false
		}
	}

	return true
}

func TestIsPresent(t *testing.T) {
	layout := "2006-01-02 15:04 MST"
	id1 := uuid.Must(uuid.NewV4())
	id2 := uuid.Must(uuid.NewV4())

	t1 := time.Now()
	t2 := time.Now().Add(2 * time.Minute)

	e1 := Event{id1, t1.Format(layout), t2.Add(2 * time.Minute).Format(layout)}
	e2 := Event{id2, time.Now().Format(layout), time.Now().Add(3 * time.Minute).Format(layout)}

	t.Run("true", func(t *testing.T) {
		events := []Event{e1, e2}
		got := isPresent(events, EventTime{id1, t1, t2})
		if !got {
			t.Errorf("want %v, got %v", true, got)
		}
	})

	t.Run("false", func(t *testing.T) {
		events := []Event{e1}
		got := isPresent(events, EventTime{id2, t1, t2})
		if got {
			t.Errorf("want %v, got %v", true, got)
		}
	})
}
