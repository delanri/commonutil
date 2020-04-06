package main

import (
	"fmt"
	cache "github.com/delanri/commonutil/cache/redis"
	"log"
	"sort"
	"time"
)

func ch1() {
	c, err := cache.New(&cache.Option{
		Address:  "localhost:6379",
		PoolSize: 4,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	ps, _ := c.Subscribe("hello-world")

	if err := ps.Receive(); err != nil {
		log.Fatal(err)
	}

	ch := ps.Channel()

	go func() {
		for {
			select {
			case m := <-ch:
				if m != nil {
					log.Printf("ch 1 got message %s", m.Payload)
				}
			}
		}
	}()

	for i := 0; i < 10; i++ {
		if err := ps.Publish(fmt.Sprintf("hello world from %d", i)); err != nil {
			log.Fatal(err)
		}
	}

}

func ch2() {
	c, err := cache.New(&cache.Option{
		Address:  "localhost:6379",
		PoolSize: 4,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	ps, _ := c.Subscribe("hello-world")

	if err := ps.Receive(); err != nil {
		log.Fatal(err)
	}

	ch := ps.Channel()

	go func() {
		for {
			select {
			case m := <-ch:
				if m != nil {
					log.Printf("ch 2 got message %s", m.Payload)
				}
			}
		}
	}()

	for i := 0; i < 10; i++ {
		if err := ps.Publish(fmt.Sprintf("hello world from %d", i)); err != nil {
			log.Fatal(err)
		}
	}
}

func ok() {
	go ch1()
	go ch2()

	time.Sleep(3 * time.Second)
}

func main() {
	slices1 := []int64{3,5,1}
	sort.Slice(slices1, func(l, r int) bool {
		return slices1[l] < slices1[r]
	})

	slices2 := []int{3,5,1}
	sort.Ints(slices2)

	log.Println(slices1)
	log.Println(slices2)
}
