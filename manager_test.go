package gum

import (
	"log"
	"os"
	"syscall"
	"testing"
	"time"
)

var WorkerID int

type Worker struct{}

// Example loop, it will be spwaned in a goroutine
func (w *Worker) Spawn(um UnitManager) {
	ticker := time.NewTicker(time.Second)

	// Worker's loop
	for {
		select {
		case <-ticker.C:
			log.Println("tick")

		// Read from channel if this worker unit should stop
		case <-um.ShouldStop():

			// Shutdown work for current unit
			w.Shutdown()

			// Notify manager that this unit is done.
			um.Done()
		}
	}
}

func (w *Worker) Shutdown() {
	// Do shutdown procedure for worker
	return
}

func NewWorker() *Worker {
	return &Worker{}
}

func DoRunMain(pid chan int, quit chan<- bool) {

	pid <- os.Getpid()

	// Create a unit manager
	manager := NewManager()

	// Shutdown all units on SIGINT
	manager.ShutdownOn(os.Interrupt)

	// NewWorker returns a type implementing WorkUnit interface unit :=
	worker1 := NewWorker()
	worker2 := NewWorker()

	// Register the unit with the manager
	manager.AddUnit(worker1)
	manager.AddUnit(worker2)

	// Start the manager
	go manager.Start()

	// Wait for all units to shutdown gracefully through their `Shutdown` method
	<-manager.Quit
	quit <- true
}

func TestRunMain(t *testing.T) {
	mainPid := make(chan int, 1)
	quit := make(chan bool)
	go DoRunMain(mainPid, quit)

	time.Sleep(3 * time.Second)

	syscall.Kill(<-mainPid, syscall.SIGINT)
	<-quit

}
