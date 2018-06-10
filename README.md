### Objective

When maintaining a calendar of events, it is important to know if an event overlaps with another event.

Given a sequence of events, each having a start and end time, this program will return the overlapping events for each event

### Setup

    $ mkdir $GOPATH/src/github.com/wung-s
    $ cd $GOPATH/src/github.com/wung-s
    $ git clone git@github.com:wung-s/event-overlap-finder.git && cd event-overlap-finder

### Package Requirement

    $ go get github.com/satori/go.uuid

### Run Test

    $ go test

### Assumptions

- All dates are properly formatted in the same valid format
- Event ID is unique across events
