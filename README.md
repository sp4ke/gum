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
4. Start the manager and wait on it's `Quit` channel

```golang
func main() {
    // Create a unit manager
    manager := gum.NewManager()

    // Subscribe to SIGINT
    manager.SubscribeTo(os.Interrupt)

    // NewWorker returns a type implementing WorkUnit interface unit :=
    NewWorker()

    // Register the unit with the manager
    manager.AddUnit(scheduler)

    // Start the manager
    go manager.Start()


    // Wait for all units to shutdown gracefully through their `Shutdown` method
    <-manager.Quit
}

```
