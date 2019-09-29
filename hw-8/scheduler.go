package hw8

import (
	"errors"
	"sync"
)

//Scheduler implementation of simple task scheduler
//Tasks will be executed simultaneously in batch. The batch size is defined by SimultaneousTasks property.
//Processing will be terminated if errors count will be equal to or greater than threshold. The threshold is defined by MaxErrors property.
type Scheduler struct {
	SimultaneousTasks int
	MaxErrors         int

	errors             []error
	scheduledTaskCount int
	completedTaskCount int
	startChannel       chan bool

	errorMutex sync.RWMutex
	taskMutex  sync.RWMutex

	wait *sync.WaitGroup

	isWaiting bool
	isRunning bool
}

//Run executes tasks specified in parameter.
//Task will be executed simultaneously in batch. The batch size is defined by SimultaneousTasks property.
func (s *Scheduler) Run(tasks []func() error) error {

	if s.isRunning {
		return errors.New("already running")
	}

	if len(tasks) == 0 {
		return errors.New("task list is empty")
	}

	s.isRunning = true

	s.errors = nil
	s.scheduledTaskCount = 0
	s.completedTaskCount = 0
	s.wait = &sync.WaitGroup{}

	for {
		//Terminate if count of errors be equal to or greater than threshold.
		if s.MaxErrors > 0 && s.MaxErrors <= s.GetErrorCount() {
			break
		}

		//Check if tasks batch is ready then start and wait for completion.
		if s.SimultaneousTasks > 0 && s.SimultaneousTasks <= s.GetScheduledTaskCount() {
			s.start()
			s.waitForCompletion()

			continue
		}

		//Start last scheduled tasks if required and terminate.
		if len(tasks) == 0 {
			s.start()
			s.waitForCompletion()

			break
		}

		//Process the next task
		task := tasks[0]
		tasks = tasks[1:]

		s.registerTask()
		go func(task func() error, start <-chan bool) {
			if start != nil {
				<-start
			}

			err := task()
			s.unregisterTask(err)
		}(task, s.getStartChannel())
	}

	s.isRunning = false
	return nil
}

//getStartChannel returns channel that used for simultaneous start of scheduled tasks
//This method is not synchronized, not intended for concurrent calls and should be called from main goroutine only
func (s *Scheduler) getStartChannel() <-chan bool {
	if s.startChannel == nil {
		s.startChannel = make(chan bool)
	}
	return s.startChannel
}

//start starts all scheduled tasks
//This method is not synchronized, not intended for concurrent calls and should be called from main goroutine only
func (s *Scheduler) start() {
	if s.startChannel != nil {
		close(s.startChannel)
		s.startChannel = nil
	}
}

//waitForCompletion waits until all scheduled tasks completed
//Waiting will be terminated if errors count will be equal to or greater than threshold
//This method should be called from main goroutine only
func (s *Scheduler) waitForCompletion() {
	if s.isWaiting {
		return
	}

	s.isWaiting = true
	s.wait.Wait()
	s.isWaiting = false
}

//registerError register error occurred in completed task
//Method returns true if count of errors less than threshold or false overwise
func (s *Scheduler) registerError(err error) bool {
	s.errorMutex.Lock()
	defer s.errorMutex.Unlock()

	s.errors = append(s.errors, err)
	return (s.MaxErrors == 0 || s.MaxErrors > len(s.errors))
}

//GetErrorCount returns count of errors occurred from last call or Run
func (s *Scheduler) GetErrorCount() int {
	s.errorMutex.RLock()
	defer s.errorMutex.RUnlock()

	return len(s.errors)
}

//GetErrors returns errors occurred from last call or Run
func (s *Scheduler) GetErrors() []error {
	s.errorMutex.RLock()
	defer s.errorMutex.RUnlock()

	errors := make([]error, len(s.errors))
	copy(errors, s.errors)
	return errors
}

//registerTask register scheduled task
//Method increment tasks counter
//Method returns true if count of scheduled tasks less than batch size or false overwise
func (s *Scheduler) registerTask() bool {
	s.taskMutex.Lock()
	defer s.taskMutex.Unlock()

	s.scheduledTaskCount++
	//fmt.Printf("register: %d/%d\n", s.scheduledTaskCount, s.completedTaskCount)
	s.wait.Add(1)

	return (s.SimultaneousTasks == 0 || s.SimultaneousTasks > s.scheduledTaskCount)
}

//unregisterTask unregister completed task
//Method decrement task counter and register error occurred in completed task
func (s *Scheduler) unregisterTask(err error) {
	s.taskMutex.Lock()
	defer s.taskMutex.Unlock()

	if s.scheduledTaskCount <= 0 {
		return
	}

	s.completedTaskCount++
	s.scheduledTaskCount--
	//fmt.Printf("unregister: %d/%d\n", s.scheduledTaskCount, s.completedTaskCount)
	s.wait.Add(-1)

	//Register error and stop processing if error threshold is reached
	if err != nil && !s.registerError(err) && s.scheduledTaskCount > 0 {
		//fmt.Println("unregister termination")
		s.wait.Add(-1 * s.scheduledTaskCount)
		s.scheduledTaskCount = 0
	}
}

//GetScheduledTaskCount returns count of active tasks
func (s *Scheduler) GetScheduledTaskCount() int {
	s.taskMutex.RLock()
	defer s.taskMutex.RUnlock()

	return s.scheduledTaskCount
}

//GetCompletedTaskCount returns count of finished tasks
func (s *Scheduler) GetCompletedTaskCount() int {
	s.taskMutex.RLock()
	defer s.taskMutex.RUnlock()

	return s.completedTaskCount
}

//IsRunning checks if scheduler is running
func (s *Scheduler) IsRunning() bool {
	return s.isRunning
}
