package main

import (
	"fmt"
	"sync"
	"time"
)

// Strategy that allows phases to overlap, meaning phases which are not conflicting can be active at the same time.
// Mainly used so that I remember strategy pattern.

type OverlappingStrategy struct {
	StopChan           chan bool
	rwMutex            sync.RWMutex        // Mutex to protect shared state
	Phases             map[string]*Phase   // Map of phases by ID
	TrafficPoles       []*TrafficPole      // Traffic poles controlled by this strategy
	IntersectionConfig *IntersectionConfig // Configuration for the intersection
}

func (os *OverlappingStrategy) Start(wg *sync.WaitGroup) error {
	defer wg.Done()
	// Start all phases in the overlapping strategy
	for _, phase := range os.Phases {
		select {
		case <-os.StopChan:
			fmt.Println("STOPPING OVERLAPPING STRATEGY DUE TO EMERGENCY SIGNAL")
			return nil // Stop if the StopChan has got a signal
		default:
			// Start the phase
			os.rwMutex.Lock()
			fmt.Printf("Starting phase %s with state %s\n", phase.ID, phase.GetPhaseState())
			err := os.handlePhase(phase)
			if err != nil {
				os.rwMutex.Unlock()
				return err // Return error if handling the phase fails
			}
			os.rwMutex.Unlock()
		}
	}
	return nil
}

func (os *OverlappingStrategy) handlePhase(newPhase *Phase) error {
	// Overlapping phases, where we allow previous phases to continue running unless they are conflicting.

	var conflictingPhases []*Phase

	for phaseID, phase := range os.Phases {
		isPhaseConcurrent := newPhase.HasPhaseIdInConcurrentPhases(phaseID)
		if !isPhaseConcurrent {
			conflictingPhases = append(conflictingPhases, phase)
		}
	}

	if len(conflictingPhases) > 0 {
		err := os.turnPhasesToInActive(conflictingPhases)
		if err != nil {
			return err
		}
	}
	var timeRemainingInActivePhase time.Duration = 0

	if newPhase.GetPhaseState() == PhaseActive {
		// If the new phase is already active, we check how long it has been active.
		timeInActivePhase := time.Since(newPhase.StartTime)
		if timeInActivePhase < newPhase.GreenTime {
			// If the new phase has been active for less than its green time, we need to wait for the remaining time.
			timeRemainingInActivePhase = newPhase.GreenTime - timeInActivePhase
		}
	} else {
		// If the new phase is not active, we need to set it to active.
		timeRemainingInActivePhase = newPhase.GreenTime
	}

	// Set all concurrent phases to active.
	newPhase.SetPhaseState(PhaseActive)
	for _, concurrentPhaseID := range newPhase.ConcurrentPhases {
		if concurrentPhase, exists := os.Phases[concurrentPhaseID]; exists {
			concurrentPhase.SetPhaseState(PhaseActive)
		}
	}
	// Sleep for the remaining time in the active phase.
	time.Sleep(timeRemainingInActivePhase)

	return nil

}

func (os *OverlappingStrategy) turnPhasesToInActive(phasesList []*Phase) error {
	var maxYellowTime time.Duration = 0
	// Set all conflicting phases to yellow.
	for _, phase := range phasesList {
		if phase.GetPhaseState() == PhaseActive {
			if phase.YellowTime > maxYellowTime {
				maxYellowTime = phase.YellowTime
			}
			phase.SetPhaseState(PhaseYellowing)
		}
	}
	// Sleep for the maximum yellow time of conflicting phases to ensure they have time to finish.
	time.Sleep(maxYellowTime)

	// Set all conflicting phases to red.
	for _, phase := range phasesList {
		phase.SetPhaseState(PhaseInactive)
	}
	return nil
}

func (os *OverlappingStrategy) EmergencyGreen(roadId string, direction RoadDirection, wg *sync.WaitGroup) error {
	defer wg.Done()

	os.StopChan <- true // Stop the current strategy if it is running

	// Handle emergency green light logic here
	// This could involve setting the traffic lights to green for the specified road and direction
	os.rwMutex.Lock()
	defer os.rwMutex.Unlock()

	var allPhases []*Phase
	for _, phase := range os.Phases {
		allPhases = append(allPhases, phase)
	}
	err := os.turnPhasesToInActive(allPhases)
	if err != nil {
		return err // Return error if turning phases to inactive fails
	}

	return nil
}

func NewOverlappingStrategy(phases []*Phase, trafficPoles []*TrafficPole, config *IntersectionConfig) IntersectionStrategy {
	res := &OverlappingStrategy{
		StopChan:           make(chan bool),
		TrafficPoles:       trafficPoles,
		rwMutex:            sync.RWMutex{},
		IntersectionConfig: config,
	}
	res.Phases = make(map[string]*Phase)
	for _, phase := range phases {
		res.Phases[phase.ID] = phase
	}
	return res
}
