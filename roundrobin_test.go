package hqgoroundrobin_test

import (
	"errors"
	"sync"
	"testing"

	roundrobin "github.com/hueristiq/hqgoroundrobin"
)

func TestNew(t *testing.T) {
	t.Parallel()

	_, err := roundrobin.New("item1", "item2", "item3")
	if err != nil {
		t.Errorf("Failed to create a new RoundRobin instance: %s", err)
	}
}

func TestNewWithOptions(t *testing.T) {
	t.Parallel()

	options := roundrobin.Options{
		RotateAmount: 2,
	}

	_, err := roundrobin.NewWithOptions(options, "item1", "item2")
	if err != nil {
		t.Errorf("Failed to create a new RoundRobin instance with options: %s", err)
	}
}

func TestAddAndNext(t *testing.T) {
	t.Parallel()

	rr, _ := roundrobin.New("item1", "item2")

	rr.Add("item3")

	counts := make(map[string]int)

	for i := 0; i < 6; i++ {
		item := rr.Next()

		counts[item.Value()]++
	}

	for _, count := range counts {
		if count != 2 {
			t.Errorf("Item was not retrieved the expected number of times: got %d, want %d", count, 2)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	rr, _ := roundrobin.New("item1", "item2", "item3", "item4")

	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(rbx *roundrobin.RoundRobin, wg *sync.WaitGroup) {
			defer wg.Done()

			for j := 0; j < 3; j++ {
				rbx.Next()
			}
		}(rr, wg)
	}

	wg.Wait()

	/*
		In Roundrobin algo all items have same priority and
		are assinged in circular order Hence test results for 100
		access with 3 iterations from different goroutines should be
		item1=75,item2=75,item3=75,item4=75
	*/

	for _, v := range rr.Items() {
		if v.Statistics.ServesCount != int32(75) {
			t.Errorf("Total item retrieval count was incorrect: got %d, want %d", v.Statistics.ServesCount, 75)
		}
	}
}

func TestStatistics(t *testing.T) {
	t.Parallel()

	rr, _ := roundrobin.New("item1", "item2")
	item := rr.Next()

	if item.Statistics.ServesCount != 1 {
		t.Errorf("Item statistics were not correctly updated: got %d, want %d", item.Statistics.ServesCount, 1)
	}
}

func TestNoItemsError(t *testing.T) {
	t.Parallel()

	_, err := roundrobin.New()
	if !errors.Is(err, roundrobin.ErrNoItems) {
		t.Errorf("Expected ErrNoItems error, got %v", err)
	}
}
