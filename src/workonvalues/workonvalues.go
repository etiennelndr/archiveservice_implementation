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
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"time"

	_ "github.com/go-sql-driver/mysql"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"

	. "github.com/etiennelndr/archiveservice/archive/service"
	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/data/implementation"
	. "github.com/etiennelndr/archiveservice/errors"
)

const (
	providerURL         = "maltcp://127.0.0.1:12400"
	consumerStoreURL    = "maltcp://127.0.0.1:12500"
	consumerRetrieveURL = "maltcp://127.0.0.1:12600"
	consumerCountURL    = "maltcp://127.0.0.1:12700"
)

var (
	Delta int64
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
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the archive service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Now, we can store the values
	for i := 0; i < 1000; i++ {
		var t1 = time.Now().UnixNano() / int64(time.Millisecond)
		// Create a random value
		var randValue = time.Duration(rand.Int63n(4) + 1)
		// Wait during a random time
		time.Sleep(randValue * time.Millisecond / 5)
		// Then, retrieve the elapsed time
		var t2 = time.Now().UnixNano() / int64(time.Millisecond)
		// Calculate the difference between the two times
		Delta = t2 - t1 + Delta
		// Finally, store this value in the archive
		store(archiveService, float32(math.Sin(float64(Delta))), Delta)
	}
}

// Provider :
func Provider() {
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the archive service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	archiveService.StartProvider(providerURL)
}

func store(archiveService *ArchiveService, valueOfSine float32, t int64) {
	var elementList = NewSineList(1)
	(*elementList)[0] = NewSine(Long(t), Float(valueOfSine))
	var boolean = NewBoolean(false)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("implementation")})
	// Object instance identifier
	var objectInstanceIdentifier = *NewLong(0)
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
		Related: NewLong(0),
		Source:  &objectID,
	}
	var network = NewIdentifier("network")
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	archiveService.Store(consumerStoreURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
}

// Retrieve :
func Retrieve() {
	var nbrOfEelements = NewLong(0)
	for true {
		// Count the number of objets stored in the archive
		nbrOfElementsInDB, err := countInDB()
		if err != nil {
			panic(err)
		}

		if *nbrOfElementsInDB > 0 && *nbrOfElementsInDB != *nbrOfEelements {
			println("COUNT:", *nbrOfElementsInDB-*nbrOfEelements)
			// Retrieve in database
			t, y, err := retrieveInDB(*nbrOfEelements)
			if err != nil {
				panic(err)
			}

			// Write in plot file
			err = writeInPlot("plot.dat", t, y)
			if err != nil {
				panic(err)
			}

			*nbrOfEelements = *nbrOfElementsInDB
		}

		time.Sleep(2 * time.Second)
	}
}

func countInDB() (*Long, error) {
	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	var err error

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the archive service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = &ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("implementation")})
	archiveQuery := &ArchiveQuery{
		Domain:    &domain,
		Related:   Long(0),
		SortOrder: NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *CompositeFilterSetList

	// Variable to retrieve the return of this function
	var longList *LongList
	// Start the consumer
	longList, errorsList, err = archiveService.Count(consumerCountURL, providerURL, objectType, archiveQueryList, queryFilterList)

	if errorsList != nil {
		return nil, errors.New(string(*errorsList.ErrorComment))
	} else if err != nil {
		return nil, err
	}

	return (*longList)[0], nil
}

func retrieveInDB(nbrOfElements Long) ([]string, []string, error) {
	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	var err error

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the archive service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("implementation")})
	archiveQuery := &ArchiveQuery{
		Domain:    &domain,
		Related:   Long(0),
		SortOrder: NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList = NewCompositeFilterSetList(0)
	filters := NewCompositeFilterList(0)
	filter := NewCompositeFilter(String("id"), COM_EXPRESSIONOPERATOR_GREATER, &nbrOfElements)
	filters.AppendElement(filter)
	compositeFilterSet := &CompositeFilterSet{
		Filters: filters,
	}
	queryFilterList.AppendElement(compositeFilterSet)

	// Variable to retrieve the responses
	var responses []interface{}

	// Start the consumer
	responses, errorsList, err = archiveService.Query(consumerRetrieveURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	if errorsList != nil {
		return nil, nil, errors.New(string(*errorsList.ErrorComment))
	} else if err != nil {
		return nil, nil, err
	}

	// Verify the response
	if len(responses) != 4 {
		return nil, nil, errors.New("Bad response")
	}

	var liste = responses[3].(*SineList)
	var t []string
	var y []string
	for i := 0; i < liste.Size(); i++ {
		sine := liste.GetElementAt(i).(*Sine)
		t = append(t, fmt.Sprintf("%d", sine.T))
		y = append(y, fmt.Sprintf("%f", sine.Y))
	}

	return t, y, nil
}

func writeInPlot(filename string, x []string, y []string) error {
	if len(x) != len(y) {
		return errors.New("lists don't have the same size")
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < len(x); i++ {
		if _, err = file.WriteString(x[i] + " " + y[i] + "\n"); err != nil {
			return err
		}
	}

	return nil
}
