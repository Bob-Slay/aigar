//  Copyright 2018 The go-ethereum Authors
//  Copyright 2019 The go-aigar Authors
//  This file is part of the go-aigar library.
//
//  The go-aigar library is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  The go-aigar library is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with the go-aigar library. If not, see <http://www.gnu.org/licenses/>.

package metrics

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"testing"
	"time"
)

const FANOUT = 128

// Stop the compiler from complaining during debugging.
var (
	_ = ioutil.Discard
	_ = log.LstdFlags
)

func BenchmarkMetrics(b *testing.B) {
	r := NewRegistry()
	c := NewRegisteredCounter("counter", r)
	g := NewRegisteredGauge("gauge", r)
	gf := NewRegisteredGaugeFloat64("gaugefloat64", r)
	h := NewRegisteredHistogram("histogram", r, NewUniformSample(100))
	m := NewRegisteredMeter("meter", r)
	t := NewRegisteredTimer("timer", r)
	RegisterDebugGCStats(r)
	RegisterRuntimeMemStats(r)
	b.ResetTimer()
	ch := make(chan bool)

	wgD := &sync.WaitGroup{}
	/*
		wgD.Add(1)
		go func() {
			defer wgD.Done()
			//log.Println("go CaptureDebugGCStats")
			for {
				select {
				case <-ch:
					//log.Println("done CaptureDebugGCStats")
					return
				default:
					CaptureDebugGCStatsOnce(r)
				}
			}
		}()
	//*/

	wgR := &sync.WaitGroup{}
	//*
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		//log.Println("go CaptureRuntimeMemStats")
		for {
			select {
			case <-ch:
				//log.Println("done CaptureRuntimeMemStats")
				return
			default:
				CaptureRuntimeMemStatsOnce(r)
			}
		}
	}()
	//*/

	wgW := &sync.WaitGroup{}
	/*
		wgW.Add(1)
		go func() {
			defer wgW.Done()
			//log.Println("go Write")
			for {
				select {
				case <-ch:
					//log.Println("done Write")
					return
				default:
					WriteOnce(r, ioutil.Discard)
				}
			}
		}()
	//*/

	wg := &sync.WaitGroup{}
	wg.Add(FANOUT)
	for i := 0; i < FANOUT; i++ {
		go func(i int) {
			defer wg.Done()
			//log.Println("go", i)
			for i := 0; i < b.N; i++ {
				c.Inc(1)
				g.Update(int64(i))
				gf.Update(float64(i))
				h.Update(int64(i))
				m.Mark(1)
				t.Update(1)
			}
			//log.Println("done", i)
		}(i)
	}
	wg.Wait()
	close(ch)
	wgD.Wait()
	wgR.Wait()
	wgW.Wait()
}

func Example() {
	c := NewCounter()
	Register("money", c)
	c.Inc(17)

	// Threadsafe registration
	t := GetOrRegisterTimer("db.get.latency", nil)
	t.Time(func() { time.Sleep(10 * time.Millisecond) })
	t.Update(1)

	fmt.Println(c.Count())
	fmt.Println(t.Min())
	// Output: 17
	// 1
}
