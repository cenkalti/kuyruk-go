package kuyruk

import (
	"log"
	"sync"
	"testing"
)

func TestConcurrentSend(t *testing.T) {
	k, err := New(DefaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	task := k.Task("m", "f", "q")
	go k.Run()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go f(task, &wg)
	}
	wg.Wait()
}

func f(task *Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		err := task.SendToQueue(nil, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
