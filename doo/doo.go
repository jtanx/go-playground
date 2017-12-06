package main

import (
	"fmt"
	"sync"
	"time"
)

type entry struct {
	when  time.Time
	stamp uint64
	value int
}

type Doo struct {
	q     [][]*entry
	stamp uint64
	mu    sync.Mutex
}

func minEntry(q [][]*entry) int {
	minIndex := 0
	for i := range q {
		if len(q[minIndex]) == 0 {
			minIndex = i
		} else if i != minIndex && len(q[i]) > 0 && (q[minIndex][0].when.After(q[i][0].when) ||
			(q[minIndex][0].when.Equal(q[i][0].when) && q[minIndex][0].stamp > q[i][0].stamp)) {
			minIndex = i
		}
	}
	return minIndex
}

func (d *Doo) Uncrustify(force bool) {
	d.mu.Lock()
	if !force {
		for _, v := range d.q {
			if len(v) == 0 {
				d.mu.Unlock()
				return
			}
		}
	}

	for true {
		i := minEntry(d.q)
		if len(d.q[i]) == 0 {
			break
		}
		fmt.Printf("%p(%p) --> ", d.q[i], &d.q[i][len(d.q[i])-1])
		ent := d.q[i][0]
		d.q[i] = d.q[i][1:]
		fmt.Printf("%p\n", d.q[i])
		fmt.Printf("DBG: %v(%d) -> %d\n", ent.when, ent.stamp, ent.value)
		if len(d.q[i]) == 0 && !force {
			break
		}
	}
	d.mu.Unlock()
}

func (d *Doo) Crustify(idx, val int) {
	if idx < 0 || idx >= len(d.q) {
		panic("you ok")
	}

	d.mu.Lock()
	d.stamp++
	d.q[idx] = append(d.q[idx], &entry{
		when:  time.Now(),
		stamp: d.stamp,
		value: val,
	})
	d.mu.Unlock()
}

func main() {
	doo := &Doo{
		q: make([][]*entry, 4),
	}

	go func() {
		cc := time.NewTicker(50 * time.Millisecond)
		for range cc.C {
			doo.Uncrustify(false)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	doo.Crustify(1, 1)
	doo.Crustify(1, 1)
	doo.Crustify(1, 1)
	doo.Crustify(3, 3)
	doo.Crustify(3, 3)
	doo.Crustify(3, 3)
	doo.Crustify(0, 0)
	doo.Crustify(0, 0)
	doo.Crustify(0, 0)
	doo.Crustify(2, 2)
	doo.Crustify(2, 2)
	doo.Crustify(2, 2)

	fmt.Printf("%#v\n", *doo)
	doo.Uncrustify(true)
}
