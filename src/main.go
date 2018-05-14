/**
 * MIT License
 *
 * Copyright (c) 2018 CNES
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
/*package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var sinus_list []float64

	var t0 = time.Now().UnixNano() / int64(time.Millisecond)
	for i := 0; i < 50; i++ {
		//var t1 = time.Now().UnixNano() / int64(time.Millisecond)
		var randValue = time.Duration(rand.Int63n(4) + 1)
		time.Sleep(randValue * time.Millisecond / 5)
		var t2 = time.Now().UnixNano() / int64(time.Millisecond)
		var Delta = t2 - t0
		//var delta = t2 - t1
		//fmt.Println("Delta =", Delta, "; delta =", delta, "; y =", int(10*math.Sin(float64(Delta))))
		fmt.Println(randValue)

		// Append Delta and value in the attribute sinus_list
		sinus_list = append(sinus_list, float64(Delta), math.Sin(float64(Delta)))
	}

	fmt.Println(sinus_list)

	// Create the plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Sinus"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Y"

	var sinusPlot = make(plotter.XYs, len(sinus_list)/2)
	for i := 0; i < len(sinus_list)/2; i++ {
		sinusPlot[i].X = sinus_list[i*2]
		sinusPlot[i].Y = sinus_list[i*2+1]
	}

	err = plotutil.AddLinePoints(p, "First Sinus", sinusPlot)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err = p.Save(4*vg.Inch, 4*vg.Inch, "sinus.png"); err != nil {
		panic(err)
	}
}*/

package main

import "github.com/veandco/go-sdl2/sdl"

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	defer surface.Free()
	surface.FillRect(nil, 0)

	rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
	surface.FillRect(&rect, 0xffff0000)
	window.UpdateSurface()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.MouseButtonEvent:
				switch (*sdl.MouseButtonEvent).GetType(event.(*sdl.MouseButtonEvent)) {
				case sdl.MOUSEBUTTONDOWN:
					println("clic enfoncé")
				}
				break
			}
		}
	}
}
