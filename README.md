# gum

Go Unit Manager is a simple Goroutine unit manager for GoLang.


Features:

- Scheduling of multiple goroutines.
- Subscribe to `os.Signal` events.
- Gracefull shutdown of units


## Overview

A unit is a type that implements `WorkUnit` interface. The `Spawn()` method
of registered units are run in goroutines. 

The `Manager` handles communication and synchronized shutdown procedure.


## Usage

1. Create a unit manager
2. Implement the `WorkUnit` on your goroutines
3. Add units to the manager
4. Start the manager and wait on its `Quit` channel

```golang
import (
    "os"
    "log"
    "time"
    "git.sp4ke.com/sp4ke/gum"
)

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

func main() {
    // Create a unit manager
    manager := gum.NewManager()

    // Subscribe to SIGINT
    manager.SubscribeTo(os.Interrupt)

    // NewWorker returns a type implementing WorkUnit interface unit :=
    worker := NewWorker()

    // Register the unit with the manager
    manager.AddUnit(worker)

    // Start the manager
    go manager.Start()


    // Wait for all units to shutdown gracefully through their `Shutdown` method
    <-manager.Quit
}
```

## Issues and Comments
The github repo is just a mirror.

For any question or issues use the repo hosted at
https://git.sp4ke.com/sp4ke/gum. 



