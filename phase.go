package main

import (
	"fmt"
	"sync"
	"time"
)

// PhaseState represents the current state of a phase
type PhaseState string

const (
	PhaseInactive  PhaseState = "INACTIVE"  // Phase not running
	PhaseActive    PhaseState = "ACTIVE"    // Phase running (green)
	PhaseYellowing PhaseState = "YELLOWING" // Phase ending (yellow)
)

// Phase represents a traffic light phase in the intersection controller. Some traffic light companies support up to 16 phases. This is a simplified version.

type Phase struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	State      PhaseState    `json:"state"`
	GreenTime  time.Duration `json:"green_time"`
	YellowTime time.Duration `json:"yellow_time"`
	AllRedTime time.Duration `json:"all_red_time"`
	//ConflictingPhases []string        `json:"conflicting_phases"`
	ConcurrentPhases []string `json:"concurrent_phases"`
	//CompatiblePhases []string        `json:"compatible_phases"`
	StartTime     time.Time       `json:"start_time"`
	LastExtension time.Time       `json:"last_extension"`
	TrafficLights []*TrafficLight `json:"traffic_lights"` // Traffic lights controlled by this phase
	mutex         sync.RWMutex
}

func (p *Phase) SetPhaseState(state PhaseState) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	fmt.Printf("CHANING PHASE %s STATE FROM %s TO %s\n", p.ID, p.State, state)

	if p.State == state {
		p.LastExtension = time.Now()
	}
	if p.State != state {
		p.StartTime = time.Now()
	}
	p.State = state
	for _, light := range p.TrafficLights {
		switch state {
		case PhaseInactive:
			light.SetLightState(Red)
		case PhaseActive:
			light.SetLightState(Green)
		case PhaseYellowing:
			light.SetLightState(Yellow)
		}
	}
}

func (p *Phase) GetPhaseState() PhaseState {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.State
}

func (p *Phase) SetTrafficLights(lights []*TrafficLight) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.TrafficLights = lights
}

func (p *Phase) SetConcurrentPhases(phases []string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ConcurrentPhases = phases
}

func (p *Phase) HasPhaseIdInConcurrentPhases(id string) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for _, phaseId := range p.ConcurrentPhases {
		if phaseId == id {
			return true
		}
	}
	return false
}

func NewPhase(id, name string, greenTime, yellowTime, allRedTime time.Duration) *Phase {
	return &Phase{
		ID:         id,
		Name:       name,
		State:      PhaseInactive,
		GreenTime:  greenTime,
		YellowTime: yellowTime,
		AllRedTime: allRedTime,
		//ConflictingPhases: []string{},
		ConcurrentPhases: []string{},
		//CompatiblePhases:  []string{},
		TrafficLights: []*TrafficLight{},
		mutex:         sync.RWMutex{},
	}
}
