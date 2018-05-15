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
package workonvalues

import (
	"math"
	"math/rand"
	"os/exec"
	"sync"
	"time"

	. "github.com/etiennelndr/archiveservice/archive/service"
	_ "github.com/go-sql-driver/mysql"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"

	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/tests/data"
)

// Show :
func Show() {
	cmd := exec.Command("gnuplot", "liveplot.gnu")

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

// Store :
func Store() {
	// Create the wait group to synchronize the go-routines
	var wg = new(sync.WaitGroup)
	wg.Add(2)

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the archive service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Launch the provider
	go storeProvider(wg, archiveService)

	go func() {
		var t0 = time.Now().UnixNano() / int64(time.Millisecond)
		for i := 0; i < 50; i++ {
			//var t1 = time.Now().UnixNano() / int64(time.Millisecond)
			var randValue = time.Duration(rand.Int63n(4) + 1)
			time.Sleep(randValue * time.Millisecond / 5)
			var t2 = time.Now().UnixNano() / int64(time.Millisecond)
			var Delta = t2 - t0
			//var delta = t2 - t1
			//fmt.Println("Delta =", Delta, "; delta =", delta, "; y =", int(10*math.Sin(float64(Delta))))
			// Append Delta and value in the attribute sinus_list
			//sinusList = append(sinusList, float64(Delta), math.Sin(float64(Delta)))
			store(archiveService, float32(math.Sin(float64(Delta))))
		}
		wg.Done()
	}()
}

func storeProvider(wg *sync.WaitGroup, archiveService *ArchiveService) {
	archiveService.StartProviders("maltcp://127.0.0.1:12400")
	wg.Done()
}

func store(archiveService *ArchiveService, valueOfSine float32) {
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = *NewLong(81)
	// Variables for ArchiveDetailsList
	var objectKey = ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = ObjectId{
		Type: &objectType,
		Key:  &objectKey,
	}
	var objectDetails = ObjectDetails{
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("network")
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	archiveService.Store("maltcp://127.0.0.1:12500", "maltcp://127.0.0.1:12400", boolean, objectType, identifierList, archiveDetailsList, elementList)
}

// Retrieve :
func Retrieve() {

}
