package main

import "sync"

// Light is different from TrafficPole.

type TrafficLightState string

const (
	Red    TrafficLightState = "red"
	Yellow TrafficLightState = "yellow"
	Green  TrafficLightState = "green"
)

type TrafficLight struct {
	ID     string
	State  TrafficLightState
	PoleID string
}

func (tl *TrafficLight) SetLightState(state TrafficLightState) {
	tl.State = state
}

func (tl *TrafficLight) GetLightState() TrafficLightState {
	return tl.State
}

func NewTrafficLight(id, poleId string) *TrafficLight {
	return &TrafficLight{
		ID:     id,
		State:  Red, // Default state is red
		PoleID: poleId,
	}
}

type RoadDirection string

const (
	North RoadDirection = "north"
	South RoadDirection = "south"
	East  RoadDirection = "east"
	West  RoadDirection = "west"
)

type TrafficPole struct {
	PoleID    string          // Unique identifier for the traffic pole
	Lights    []*TrafficLight // All lights on the pole
	RoadID    string
	Direction RoadDirection
	Mu        sync.Mutex
}

func (t *TrafficPole) SetAllState(state TrafficLightState) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	for _, light := range t.Lights {
		light.SetLightState(state)
	}
}

func NewTrafficPole(poleID string, roadId string, direction RoadDirection, lights []*TrafficLight) *TrafficPole {
	return &TrafficPole{
		RoadID:    roadId,
		Direction: direction,
		Lights:    lights,
	}
}
