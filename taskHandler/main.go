package main

import (
    "fmt"
    "sync"
    "time"
)

type Task struct {
    ID     string
    Deps   []string
    Run    func() error
    Status string
}

type TaskQueue struct {
    tasks map[string]*Task
    mu    sync.Mutex
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{tasks: make(map[string]*Task)}
}

func (q *TaskQueue) AddTask(t *Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks[t.ID] = t
}

func (q *TaskQueue) GetTask(id string) (*Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	t, ok := q.tasks[id]
	return t, ok
}


type TaskRunner struct {
	taskQueue *TaskQueue
	wg        sync.WaitGroup
	workers   int
}

func NewTaskRunner(tq *TaskQueue) *TaskRunner {
	return &TaskRunner{taskQueue: tq, workers: 2}
}


func (r *TaskRunner) canRun(task *Task) bool {
	for _, depID := range task.Deps {
		depTask, ok := r.taskQueue.GetTask(depID)
		if !ok || depTask.Status != "completed" {
			return false
		}
	}
	return true
}

func (r *TaskRunner) RunTask(t *Task) {
	defer r.wg.Done()
	if !r.canRun(t) {
		fmt.Printf("Task %s: dependencies not met\n", t.ID)
		return
	}
	err := t.Run()
	if err != nil {
		fmt.Printf("Task %s failed: %v\n", t.ID, err)
		t.Status = "failed"
	} else {
		fmt.Printf("Task %s completed\n", t.ID)
		t.Status = "completed"
	}
}


func (r *TaskRunner) Execute() {
	tasksChan := make(chan *Task)

	for i := 0; i < r.workers; i++ {
		go func() {
			for task := range tasksChan {
				r.RunTask(task)
			}
		}()
	}

	for {
		allDone := true
		var runnableTasks []*Task

		r.taskQueue.mu.Lock()
		for _, task := range r.taskQueue.tasks {
			if task.Status == "pending" {
				allDone = false
				runnableTasks = append(runnableTasks, task)
			}
			if task.Status != "completed" && task.Status != "pending" {
				allDone = false
			}
		}
		r.taskQueue.mu.Unlock()

		for _, task := range runnableTasks {
			if r.canRun(task) {
				r.taskQueue.mu.Lock()
				if task.Status == "pending" {
					task.Status = "running"
					r.wg.Add(1)
					tasksChan <- task
				}
				r.taskQueue.mu.Unlock()
			}
		}

		if allDone {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	r.wg.Wait()
	close(tasksChan)
}

func main() {
	queue := NewTaskQueue()

	taskA := &Task{
		ID:   "1",
		Deps: []string{},
		Run: func() error {
			fmt.Println("Running A")
			time.Sleep(1 * time.Second)
			return nil
		},
		Status: "pending",
	}

	taskB := &Task{
		ID:   "2",
		Deps: []string{},
		Run: func() error {
			fmt.Println("Running B")
			time.Sleep(1 * time.Second)
			return nil
		},
		Status: "pending",
	}

	taskC := &Task{
		ID:   "3",
		Deps: []string{"2"},
		Run: func() error {
			fmt.Println("Running C")
			time.Sleep(1 * time.Second)
			return nil
		},
		Status: "pending",
	}

	queue.AddTask(taskA)
    queue.AddTask(taskC)
	queue.AddTask(taskB)
	

	runner := NewTaskRunner(queue)
	runner.Execute()
}