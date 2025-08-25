package services

import (
	"container/heap"
	"log"
	"sync"
	"time"
)

type ScheduleService struct {
	once         sync.Once
	scheduleHeap ScheduleHeap
	mutex        sync.Mutex
}

type ScheduleItem struct {
	StartTime time.Time
	Duration  time.Duration
	DeviceID  uint
}

type ScheduleHeap []ScheduleItem

func NewScheduleService() *ScheduleService {
	return &ScheduleService{
		scheduleHeap: ScheduleHeap{},
	}
}

func (s *ScheduleService) StartScheduler(deviceService *DeviceService) {
	s.once.Do(func() {
		go s.schedulerLoop(deviceService)
	})
}

func (s *ScheduleService) AddScheduleHandler(item ScheduleItem) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	heap.Push(&s.scheduleHeap, item)
}

func (s *ScheduleService) schedulerLoop(deviceService *DeviceService) {
	for {
		s.mutex.Lock()
		if s.scheduleHeap.Len() == 0 {
			s.mutex.Unlock()
			time.Sleep(1 * time.Minute)
			continue
		}

		nextItem := s.scheduleHeap[0]
		now := time.Now()
		if nextItem.StartTime.After(now) {
			s.mutex.Unlock()
			time.Sleep(nextItem.StartTime.Sub(now))
			continue
		}

		// Remove the item from the heap
		heap.Pop(&s.scheduleHeap)
		s.mutex.Unlock()

		// Process the scheduled item
		log.Printf("[SCHEDULE] Processing schedule: %+v", nextItem)
		// TODO: Add logic to activate the device or perform the scheduled action
		err := deviceService.EnqueueActivation(&DeviceRequest{
			DeviceID: nextItem.DeviceID,
			Duration: nextItem.Duration,
		})
		if err != nil {
			log.Printf("[SCHEDULE] Error enqueuing device activation: %v", err)
		}
	}
}

func (h ScheduleHeap) Len() int           { return len(h) }                                // Value receiver
func (h ScheduleHeap) Less(i, j int) bool { return h[i].StartTime.Before(h[j].StartTime) } // Value receiver
func (h ScheduleHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }                      // Value receiver

func (h *ScheduleHeap) Push(x interface{}) {
	*h = append(*h, x.(ScheduleItem))
}

func (h *ScheduleHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}
