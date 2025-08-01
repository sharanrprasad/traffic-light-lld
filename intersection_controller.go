package main

import (
	"sync"
)

type IntersectionController struct {
	Strategy IntersectionStrategy
	wg       *sync.WaitGroup // WaitGroup to manage goroutines
}

// Start To finish.
func (ic *IntersectionController) Start() error {
	ic.wg.Add(1)

	go func() {
		err := ic.Strategy.Start(ic.wg)
		if err != nil {
			panic(err)
		}
	}()

	//go func() {
	//	time.Sleep(time.Second * 10)
	//	ic.wg.Add(2)
	//	err := ic.Strategy.EmergencyGreen("north-road", North, ic.wg)
	//	if err != nil {
	//		panic(err)
	//	}
	//	// Start again after emergency green is completed.
	//	err = ic.Strategy.Start(ic.wg)
	//}()

	ic.wg.Wait() // Wait for all goroutines to finish
	return nil
}

func NewIntersectionController(strategy IntersectionStrategy) *IntersectionController {
	return &IntersectionController{
		Strategy: strategy,
		wg:       &sync.WaitGroup{},
	}
}
