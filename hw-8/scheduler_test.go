package hw8

import (
	"fmt"
	"testing"
	"time"
)

type TestTask struct {
	Name      string
	SleepTime time.Duration
	IsFail    bool
}

func (t *TestTask) Run(owner *testing.T) error {
	owner.Logf("\"%s\" is started\n", t.Name)
	time.Sleep(t.SleepTime * time.Millisecond)
	if t.IsFail {
		owner.Logf("\"%s\" is completed with error\n", t.Name)
		return fmt.Errorf("task \"%s\" completed with error", t.Name)
	}

	owner.Logf("\"%s\" is successfully completed\n", t.Name)
	return nil
}

type SchedulerTest struct {
	Title             string
	SimultaneousTasks int
	MaxErrors         int
	Tasks             []TestTask
}

func (t *SchedulerTest) Run(owner *testing.T) error {
	scheduler := &Scheduler{SimultaneousTasks: t.SimultaneousTasks, MaxErrors: t.MaxErrors}
	callbacks := make([]func() error, 0, len(t.Tasks))
	for _, task := range t.Tasks {
		callbacks = append(
			callbacks,
			func(task TestTask) func() error {
				return func() error { return task.Run(owner) }
			}(task),
		)
	}

	owner.Logf("--- Start of test %s ---", t.Title)
	err := scheduler.Run(callbacks)
	if err != nil {
		return err
	}

	if scheduler.GetErrorCount() != t.MaxErrors {
		return NewUnexpectedFailureCountError(t.MaxErrors, scheduler.GetErrorCount())
	}

	return nil
}

func NewUnexpectedFailureCountError(expected, received int) error {
	return fmt.Errorf("quantity of errors expected to be %d, but received %d", expected, received)
}

func TestScheduler(t *testing.T) {

	tasks := []TestTask{
		TestTask{Name: "TASK #1", SleepTime: 100, IsFail: false},
		TestTask{Name: "TASK #2", SleepTime: 200, IsFail: true},
		TestTask{Name: "TASK #3", SleepTime: 100, IsFail: false},
		TestTask{Name: "TASK #4", SleepTime: 300, IsFail: true},
		TestTask{Name: "TASK #5", SleepTime: 1000, IsFail: false},
		TestTask{Name: "TASK #6", SleepTime: 600, IsFail: false},
		TestTask{Name: "TASK #7", SleepTime: 1000, IsFail: false},
		TestTask{Name: "TASK #8", SleepTime: 100, IsFail: true},
		TestTask{Name: "TASK #9", SleepTime: 300, IsFail: false},
		TestTask{Name: "TASK #10", SleepTime: 200, IsFail: false},
		TestTask{Name: "TASK #11", SleepTime: 500, IsFail: false},
		TestTask{Name: "TASK #12", SleepTime: 1000, IsFail: false},
		TestTask{Name: "TASK #13", SleepTime: 200, IsFail: false},
		TestTask{Name: "TASK #14", SleepTime: 400, IsFail: false},
		TestTask{Name: "TASK #15", SleepTime: 600, IsFail: false},
		TestTask{Name: "TASK #16", SleepTime: 500, IsFail: true},
	}

	tests := []SchedulerTest{
		SchedulerTest{Title: "[ SimultaneousTasks = 5, MaxErrors = 1 ]", SimultaneousTasks: 5, MaxErrors: 1, Tasks: tasks},
		SchedulerTest{Title: "[ SimultaneousTasks = 5, MaxErrors = 2 ]", SimultaneousTasks: 5, MaxErrors: 2, Tasks: tasks},
		SchedulerTest{Title: "[ SimultaneousTasks = 10, MaxErrors = 3 ]", SimultaneousTasks: 10, MaxErrors: 3, Tasks: tasks},
		SchedulerTest{Title: "[ SimultaneousTasks = 10, MaxErrors = 4 ]", SimultaneousTasks: 10, MaxErrors: 4, Tasks: tasks},
	}

	for _, test := range tests {
		err := test.Run(t)
		if err != nil {
			t.Error(err)
		}
	}
}
