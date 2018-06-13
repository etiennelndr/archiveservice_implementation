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
package main

import (
	"sync"
	"time"

	work "github.com/etiennelndr/archiveservice_implementation/src/workonvalues"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)

	// Go-routine to show all the values
	go func() {
		// Wait a little for the two consumers to start
		time.Sleep(3 * time.Second)
		work.Show()
		wg.Done()
	}()

	// Go-routine to retrieve all the values
	go func() {
		// Wait a little for the provider to start
		time.Sleep(2 * time.Second)
		work.Retrieve()
		wg.Done()
	}()

	// Go-routine to store all the values
	go func() {
		// Wait a little for the provider to start
		time.Sleep(1 * time.Second)
		work.Store()
		wg.Done()
	}()

	// Go-routine to start a provider
	go func() {
		work.Provider()
		wg.Done()
	}()

	// Wait
	wg.Wait()

	return
}
