package gum

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
)

type WorkUnit interface {
	Spawn(UnitManager)
	Shutdown()
}

type UnitManager interface {
	ShouldStop() <-chan bool
	Done()
	Panic(err error)
}

type WorkUnitManager struct {
	stop       chan bool
	workerQuit chan bool
	unit       WorkUnit
	panic      chan error
	isPaniced  bool
}

func (w *WorkUnitManager) ShouldStop() <-chan bool {
	return w.stop
}

func (w *WorkUnitManager) Done() {
	w.workerQuit <- true
}

func (w *WorkUnitManager) Panic(err error) {
	w.panic <- err
	w.isPaniced = true
	w.workerQuit <- true
	close(w.stop)
}

type Manager struct {
	signal chan os.Signal

	workers map[string]*WorkUnitManager

	Quit chan bool

	panic chan error // Used for panicing goroutines
}

func (m *Manager) Start() {
	log.Println("Starting manager ...")

	for unitName, w := range m.workers {
		log.Printf("Starting <%s>\n", unitName)
		go w.unit.Spawn(w)
	}

	for {
		select {
		case sig := <-m.signal:
			if sig != os.Interrupt {
				break
			}

			log.Println("shutting event received ... ")

			// send shutdown event to all worker units
			for name, w := range m.workers {
				log.Printf("shutting down <%s>\n", name)
				w.stop <- true
			}

			// Wait for all units to quit
			for name, w := range m.workers {
				<-w.workerQuit
				log.Printf("<%s> down", name)
			}

			// All workers have shutdown
			log.Println("All workers have shutdown, shutting down manager ...")

			m.Quit <- true

			//TODO: wait on Quit signal from any WorkerUnit
		case p := <-m.panic:
			//TODO: refactor shutdown procedure and reuse here
			// send shutdown event to all worker units

			for name, w := range m.workers {
				if w.isPaniced {
					log.Printf("Panicing for <%s>: %s", name, p)
				}
			}

			for name, w := range m.workers {
				log.Printf("shuting down <%s>\n", name)
				if !w.isPaniced {
					w.stop <- true
				}
			}

			// Wait for all units to quit
			for name, w := range m.workers {
				<-w.workerQuit
				log.Printf("<%s> down", name)
			}

			// All workers have shutdown
			log.Println("All workers have shutdown, shutting down manager ...")

			m.Quit <- true

		}
	}
}

func (m *Manager) ShutdownOn(sig os.Signal) {
	signal.Notify(m.signal, sig)
}

func (m *Manager) AddUnit(unit WorkUnit) {

	workUnitManager := &WorkUnitManager{
		workerQuit: make(chan bool, 1),
		stop:       make(chan bool, 1),
		unit:       unit,
		panic:      m.panic,
	}

	unitType := reflect.TypeOf(unit)
	unitName := strings.Split(unitType.String(), ".")[1]

	log.Println("Adding unit ", unitName)

	m.workers[unitName] = workUnitManager
	log.Println(m.workers)
}

func NewManager() *Manager {
	return &Manager{
		signal:  make(chan os.Signal, 1),
		Quit:    make(chan bool, 1),
		workers: make(map[string]*WorkUnitManager),
		panic:   make(chan error, 1),
	}
}
