package main

import (
	"sync"
	"time"
)

type IntersectionConfig struct {
	EmergencyDuration time.Duration
}

type IntersectionStrategy interface {
	Start(wg *sync.WaitGroup) error
	EmergencyGreen(roadId string, direction RoadDirection, wg *sync.WaitGroup) error
}
