package util

import (
	"fmt"
	"reflect"
	"periph.io/x/conn/v3/gpio"
    "periph.io/x/host/v3"
    "periph.io/x/host/v3/rpi"
)

var PrinterIdToGpio PrinterToGpioMap

type PrinterToGpioMap struct {
	Map       map[int]int
	populated bool
}

func populateGpioMap() {

	//maximum of 28 supported printers due to available GPIO pins
	//guide: PrinterIdToGpio.Map[printerId] = physical pin Number

	//3V3 power			    = 1
	//5V power 			    = 2
	PrinterIdToGpio.Map[1] = 3
	//5V power 			    = 4
	PrinterIdToGpio.Map[2] = 5
	//GND 				    = 6
	PrinterIdToGpio.Map[3] = 7
	PrinterIdToGpio.Map[4] = 8
	//GND				    = 9
	PrinterIdToGpio.Map[5] = 10
	PrinterIdToGpio.Map[6] = 11
	PrinterIdToGpio.Map[7] = 12
	PrinterIdToGpio.Map[8] = 13
	//GND				    = 14
	PrinterIdToGpio.Map[9] = 15
	PrinterIdToGpio.Map[10] = 16
	//3V3 power             = 17
	PrinterIdToGpio.Map[11] = 18
	PrinterIdToGpio.Map[12] = 19
	//GND				    = 20
	PrinterIdToGpio.Map[13] = 21
	PrinterIdToGpio.Map[14] = 22
	PrinterIdToGpio.Map[15] = 23
	PrinterIdToGpio.Map[16] = 24
	//GND				    = 25
	PrinterIdToGpio.Map[17] = 26
	PrinterIdToGpio.Map[18] = 27
	PrinterIdToGpio.Map[19] = 28
	PrinterIdToGpio.Map[20] = 29
	//GND		 		    = 30
	PrinterIdToGpio.Map[21] = 31
	PrinterIdToGpio.Map[22] = 32
	PrinterIdToGpio.Map[23] = 33
	//GND				    = 34
	PrinterIdToGpio.Map[24] = 35
	PrinterIdToGpio.Map[25] = 36
	PrinterIdToGpio.Map[26] = 37
	PrinterIdToGpio.Map[27] = 38
	//GND				    = 39
	PrinterIdToGpio.Map[28] = 40
}

func getPinNumberFromPrinterId(printerId int) (int, error) {

	if !PrinterIdToGpio.populated {
		populateGpioMap()
		PrinterIdToGpio.populated = true
	}

	value, exists := PrinterIdToGpio.Map[printerId]
	if !exists {
		return -1, fmt.Errorf("invalid Printer ID. Should be between 1 and 28")
	}
	return value, nil
}

func TogglePrinterPower(printerId int) (bool, error) {
	pinNumber, err := getPinNumberFromPrinterId(printerId)
	if err != nil {
		return false, err
	}
	
}
