package main

import (
	"testing"
)

func TestFieldPosition(t *testing.T) {
	fp := FieldPosition{xPos: 4, yPos: 4}

	fp.move('^')
	if fp.yPos != 5 {
		t.Error("Move north failed")
	}
	fp.move('v')
	if fp.yPos != 4 {
		t.Error("Move south failed")
	}
	fp.move('<')
	if fp.xPos != 3 {
		t.Error("Move west failed")
	}
	fp.move('>')
	if fp.xPos != 4 {
		t.Error("Move east failed")
	}
}

func TestDeliverySet(t *testing.T) {
	channel := make(chan []string)
	go runDeliverySet(">>v", channel)
	results := <-channel

	if len(results) != 4 {
		t.Error("Not enough steps taken")
	}
	if results[0] != "3:3" && results[3] != "5:2" {
		t.Error("Delivery route incorrect")
	}
}

func TestFilteredData(t *testing.T) {
	data := "^v^v^v^v^v"

	if filterData(data, 0, 2) != "^^^^^" {
		t.Errorf("Filtered data is incorrect (0,2)")
	}
	if filterData(data, 1, 2) != "vvvvv" {
		t.Errorf("Filtered data is incorrect (1,2)")
	}
	if filterData(data, 0, 1) != data {
		t.Errorf("Filtered data is incorrect (0,1)")
	}
	if filterData(data, 0, 3) != "^v^v" {
		t.Errorf("Filtered data is incorrect (0,3)")
	}
}

func TestRunDeliveries(t *testing.T) {
	data := "^v^v^v^v^v"
	if len(runDeliveries(data, 1)) != 2 {
		t.Error("Single worker test failed")
	}
	if len(runDeliveries(data, 2)) != 11 {
		t.Error("Double worker test failed")
	}
}
