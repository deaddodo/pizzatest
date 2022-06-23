package main

import (
	"fmt"
	"os"
	"strconv"
)

func getData(filename string) string {
	var data, err = os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	return string(data)
}

type FieldPosition struct {
	xPos, yPos int32
}

func (fp *FieldPosition) move(direction rune) {
	if direction == '<' {
		fp.xPos--
	}
	if direction == '>' {
		fp.xPos++
	}
	if direction == '^' {
		fp.yPos++
	}
	if direction == 'v' {
		fp.yPos--
	}
}

func runDeliverySet(data string, channel chan<- []string) {
	// we'll start in the middle to make things easier to track
	// since we're only storing unique positions, the basepoint is moot
	position := FieldPosition{
		xPos: int32(len(data)),
		yPos: int32(len(data)),
	}
	deliveries := make([]string, 0)

	// our basepoint starts with a delivered pizza
	deliveries = append(deliveries, fmt.Sprintf("%d:%d", position.xPos, position.yPos))

	for _, ch := range data {
		position.move(ch)
		deliveries = append(deliveries, fmt.Sprintf("%d:%d", position.xPos, position.yPos))
	}

	channel <- deliveries
}

func filterData(data string, start int, steps int) string {
	filteredData := ""

	stepCounter := steps
	for idx, ch := range data {
		if idx < start {
			continue
		}

		if stepCounter == steps {
			filteredData = filteredData + string(ch)
			stepCounter = 1
		} else {
			stepCounter++
		}
	}

	return filteredData
}

func runDeliveries(data string, workerCount int) map[string]struct{} {
	locationMap := make(map[string]struct{})

	channel := make(chan []string)
	for idx := 0; idx < workerCount; idx++ {
		go runDeliverySet(filterData(data, idx, workerCount), channel)
	}

	for idx := 0; idx < workerCount; idx++ {
		for _, address := range <-channel {
			if _, ok := locationMap[address]; !ok {
				locationMap[address] = struct{}{}
			}
		}
	}

	return locationMap
}

func main() {
	if len(os.Args) < 2 {
		panic("File location required!")
	}

	workerCount := 1
	if len(os.Args) > 2 {
		wcArgument, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err.Error())
		}
		workerCount = wcArgument
	}

	data := getData(os.Args[1])
	locationMap := runDeliveries(data, workerCount)
	fmt.Printf("Steps: %d, Deliveries: %d\n", len(data), len(locationMap))
}
