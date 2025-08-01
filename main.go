package main

import "time"

func main() {

	// Create all the traffic poles.
	northPoleId := "north-road-pole"
	NorthStraightTrafficLight := NewTrafficLight("north-road-1-straight", northPoleId)
	NorthLeftTrafficLight := NewTrafficLight("north-road-1-left", northPoleId)
	NorthRightTrafficLight := NewTrafficLight("north-road-1-right", northPoleId)

	NorthTrafficPole := NewTrafficPole(northPoleId, "north-road", North, []*TrafficLight{NorthStraightTrafficLight, NorthLeftTrafficLight, NorthRightTrafficLight})

	southPoleId := "south-road-pole"
	SouthStraightTrafficLight := NewTrafficLight("south-road-1-straight", southPoleId)
	SouthLeftTrafficLight := NewTrafficLight("south-road-1-left", southPoleId)
	SouthRightTrafficLight := NewTrafficLight("south-road-1-right", southPoleId)
	SouthTrafficPole := NewTrafficPole(southPoleId, "south-road", South, []*TrafficLight{SouthStraightTrafficLight, SouthLeftTrafficLight, SouthRightTrafficLight})

	eastPoleId := "east-road-pole"
	EastStraightTrafficLight := NewTrafficLight("east-road-1-straight", eastPoleId)
	EastLeftTrafficLight := NewTrafficLight("east-road-1-left", eastPoleId)
	EastRightTrafficLight := NewTrafficLight("east-road-1-right", eastPoleId)
	EastTrafficPole := NewTrafficPole(eastPoleId, "east-road", East, []*TrafficLight{EastStraightTrafficLight, EastLeftTrafficLight, EastRightTrafficLight})

	westPoleId := "west-road-pole"
	WestStraightTrafficLight := NewTrafficLight("west-road-1-straight", westPoleId)
	WestLeftTrafficLight := NewTrafficLight("west-road-1-left", westPoleId)
	WestRightTrafficLight := NewTrafficLight("west-road-1-right", westPoleId)
	WestTrafficPole := NewTrafficPole(westPoleId, "west-road", West, []*TrafficLight{WestStraightTrafficLight, WestLeftTrafficLight, WestRightTrafficLight})

	// Phase 1
	phase1 := NewPhase("northGreen", "North Green", 3*time.Second, 5*time.Second, 2*time.Second)
	phase1.SetTrafficLights([]*TrafficLight{NorthStraightTrafficLight})
	phase1.SetConcurrentPhases([]string{"southGreen", "northLeft", "southLeft"})

	// Phase 2
	phase2 := NewPhase("southGreen", "South Green", 3*time.Second, 5*time.Second, 2*time.Second)
	phase2.SetTrafficLights([]*TrafficLight{SouthStraightTrafficLight})
	phase2.ConcurrentPhases = []string{"northGreen", "southLeft", "northLeft"}

	// Phase 3
	phase3 := NewPhase("northLeft", "North Left", 20*time.Second, 5*time.Second, 2*time.Second)
	phase3.SetTrafficLights([]*TrafficLight{NorthLeftTrafficLight})
	phase3.ConcurrentPhases = []string{"northGreen", "southGreen", "southLeft"}

	// Phase 4
	phase4 := NewPhase("southLeft", "South Left", 20*time.Second, 5*time.Second, 2*time.Second)
	phase4.SetTrafficLights([]*TrafficLight{SouthLeftTrafficLight})
	phase4.ConcurrentPhases = []string{"northGreen", "southGreen", "northLeft"}

	// Phase 5
	phase5 := NewPhase("northRight", "North Right", 3*time.Second, 5*time.Second, 2*time.Second)
	phase5.SetTrafficLights([]*TrafficLight{NorthRightTrafficLight})
	phase5.SetConcurrentPhases([]string{"southRight", "westLeft", "eastLeft"})

	// Phase 6
	phase6 := NewPhase("southRight", "South Right", 3*time.Second, 5*time.Second, 2*time.Second)
	phase6.SetTrafficLights([]*TrafficLight{SouthRightTrafficLight})
	phase6.ConcurrentPhases = []string{"northRight", "westLeft", "eastLeft"}

	// Phase 7
	phase8 := NewPhase("eastGreen", "East Green", 3*time.Second, 5*time.Second, 2*time.Second)
	phase8.SetTrafficLights([]*TrafficLight{EastStraightTrafficLight})
	phase8.ConcurrentPhases = []string{"westGreen", "eastLeft", "westLeft"}

	// Phase 8
	phase7 := NewPhase("westGreen", "West Green", 3*time.Second, 5*time.Second, 2*time.Second)
	phase7.SetTrafficLights([]*TrafficLight{WestStraightTrafficLight})
	phase7.ConcurrentPhases = []string{"eastGreen", "westLeft", "eastLeft"}

	// Phase 9
	phase9 := NewPhase("eastLeft", "East Left", 3*time.Second, 5*time.Second, 2*time.Second)
	phase9.SetTrafficLights([]*TrafficLight{EastLeftTrafficLight})
	phase9.ConcurrentPhases = []string{"eastGreen", "westGreen", "westLeft"}

	// Phase 10
	phase10 := NewPhase("westLeft", "West Left", 3*time.Second, 5*time.Second, 2*time.Second)
	phase10.SetTrafficLights([]*TrafficLight{WestLeftTrafficLight})
	phase10.ConcurrentPhases = []string{"eastGreen", "eastLeft", "westGreen"}

	// Phase 11
	phase11 := NewPhase("eastRight", "East right", 3*time.Second, 5*time.Second, 2*time.Second)
	phase11.SetTrafficLights([]*TrafficLight{EastRightTrafficLight})
	phase11.ConcurrentPhases = []string{"westRight", "northLeft", "southLeft"}

	// Phase 12
	phase12 := NewPhase("westRight", "West right", 3*time.Second, 5*time.Second, 2*time.Second)
	phase12.SetTrafficLights([]*TrafficLight{WestRightTrafficLight})
	phase12.ConcurrentPhases = []string{"eastRight", "northLeft", "southLeft"}

	fourLanePhases := []*Phase{
		phase1, phase2, phase3, phase4, phase5, phase6, phase7, phase8, phase9, phase10, phase11, phase12,
	}

	// Create the intersection strategy with the four lane phases.
	intersectionConfig := &IntersectionConfig{
		EmergencyDuration: 10 * time.Minute, // Duration for emergency green light
	}

	overlappingStrategy := NewOverlappingStrategy(fourLanePhases, []*TrafficPole{NorthTrafficPole, SouthTrafficPole, EastTrafficPole, WestTrafficPole}, intersectionConfig)

	controller := NewIntersectionController(overlappingStrategy)
	err := controller.Start()
	if err != nil {
		panic(err)
	}
}
