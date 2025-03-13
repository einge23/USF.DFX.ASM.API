package util

import (
	"fmt"
	"sync"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3/rpi"
)

var PrinterIdToGpio = NewPrinterToGpioMap()

type PrinterToGpioMap struct {
	Map       map[int]gpio.PinIO
	populated bool
}

//constructor
func NewPrinterToGpioMap() PrinterToGpioMap {
    return PrinterToGpioMap{
        Map:       make(map[int]gpio.PinIO),
        populated: false,
    }
}

type PinWithMutex struct {
	Pin gpio.PinIO
	Mutex sync.RWMutex
}

func populateGpioMap() {

	//maximum of 28 supported printers due to available GPIO pins
	//guide: PrinterIdToGpio.Map[printerId] = physical pin

	//3V3 power			   = rpi.P1_1
	//5V power 			   = rpi.P1_2
	PrinterIdToGpio.Map[1] = rpi.P1_3
	//5V power 			   = rpi.P1_4
	PrinterIdToGpio.Map[2] = rpi.P1_5
	//GND 				   = rpi.P1_6
	PrinterIdToGpio.Map[3] = rpi.P1_7
	PrinterIdToGpio.Map[4] = rpi.P1_8
	//GND				   = rpi.P1_9
	PrinterIdToGpio.Map[5] = rpi.P1_10
	PrinterIdToGpio.Map[6] = rpi.P1_11
	PrinterIdToGpio.Map[7] = rpi.P1_12
	PrinterIdToGpio.Map[8] = rpi.P1_13
	//GND				   = rpi.P1_14
	PrinterIdToGpio.Map[9] = rpi.P1_15
	PrinterIdToGpio.Map[10] = rpi.P1_16
	//3V3 power             = rpi.P1_17
	PrinterIdToGpio.Map[11] = rpi.P1_18
	PrinterIdToGpio.Map[12] = rpi.P1_19
	//GND				    = rpi.P1_20
	PrinterIdToGpio.Map[13] = rpi.P1_21
	PrinterIdToGpio.Map[14] = rpi.P1_22
	PrinterIdToGpio.Map[15] = rpi.P1_23
	PrinterIdToGpio.Map[16] = rpi.P1_24
	//GND				    = rpi.P1_25
	PrinterIdToGpio.Map[17] = rpi.P1_26
	PrinterIdToGpio.Map[18] = rpi.P1_27
	PrinterIdToGpio.Map[19] = rpi.P1_28
	PrinterIdToGpio.Map[20] = rpi.P1_29
	//GND		 		    = rpi.P1_30
	PrinterIdToGpio.Map[21] = rpi.P1_31
	PrinterIdToGpio.Map[22] = rpi.P1_32
	PrinterIdToGpio.Map[23] = rpi.P1_33
	//GND				    = rpi.P1_34
	PrinterIdToGpio.Map[24] = rpi.P1_35
	PrinterIdToGpio.Map[25] = rpi.P1_36
	PrinterIdToGpio.Map[26] = rpi.P1_37
	PrinterIdToGpio.Map[27] = rpi.P1_38
	//GND				    = rpi.P1_39
	PrinterIdToGpio.Map[28] = rpi.P1_40
}

func getPinFromPrinterId(printerId int) (gpio.PinIO, error) {

	if !PrinterIdToGpio.populated {
		populateGpioMap()
		PrinterIdToGpio.populated = true
	}

	value, exists := PrinterIdToGpio.Map[printerId]
	if !exists {
		return nil, fmt.Errorf("invalid Printer ID. Should be between 1 and 28 (inclusive)")
	}
	return value, nil
}

func TurnOnPrinter(printerId int) (bool, error) {

	pin, err := getPinFromPrinterId(printerId)
	if err != nil {
		return false, err
	}
	err = pin.Out(gpio.High)
	if err != nil {
		return false, fmt.Errorf("could not write HIGH to GPIO pin")
	}

	return true, nil
}

func TurnOffPrinter(printerId int) (bool, error) {

	pin, err := getPinFromPrinterId(printerId)
	if err != nil {
		return false, err
	}
	err = pin.Out(gpio.Low)
	if err != nil {
		return false, fmt.Errorf("could not write LOW to GPIO pin")
	}

	return true, nil
}