package sdm630

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type MeterType string
type MeterState uint8

const (
	METERTYPE_JANITZA = "JANITZA"
	METERTYPE_SDM     = "SDM"
	METERTYPE_DZG     = "DZG"
)

const (
	METERSTATE_AVAILABLE   = iota // The device responds (initial state)
	METERSTATE_UNAVAILABLE        // The device does not respond
)

type Meter struct {
	Type          MeterType
	DeviceID      uint8
	Scheduler     Scheduler
	MeterReadings *MeterReadings
	state         MeterState
	mux           sync.Mutex // syncs the meter state variable
}

func NewMeter(
	typeid MeterType,
	devid uint8,
	scheduler Scheduler,
	timeToCacheReadings time.Duration,
) *Meter {
	r := NewMeterReadings(devid, timeToCacheReadings)
	return &Meter{
		Type:          typeid,
		Scheduler:     scheduler,
		DeviceID:      devid,
		MeterReadings: r,
		state:         METERSTATE_AVAILABLE,
	}
}

func (m *Meter) UpdateState(newstate MeterState) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.state = newstate
	if newstate == METERSTATE_UNAVAILABLE {
		m.MeterReadings.Purge(m.DeviceID)
	}
}

func (m *Meter) GetState() MeterState {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.state
}

func (m *Meter) GetReadableState() string {
	var retval string
	switch m.GetState() {
	case METERSTATE_AVAILABLE:
		retval = "available"
	case METERSTATE_UNAVAILABLE:
		retval = "unavailable"
	default:
		log.Fatal("Unknown meter state, aborting.")
	}
	return retval
}

func (m *Meter) GetMeterType() MeterType {
	return m.Type
}

func (m *Meter) AddSnip(snip QuerySnip) {
	m.MeterReadings.AddSnip(snip)
}

type MeterReadings struct {
	LastMinuteReadings ReadingSlice
	LastReading        Readings
}

func NewMeterReadings(devid uint8, secondsToStore time.Duration) (retval *MeterReadings) {
	reading := Readings{
		UniqueID:       fmt.Sprintf(UniqueIDFormat, devid),
		ModbusDeviceID: devid,
	}
	retval = &MeterReadings{
		LastMinuteReadings: ReadingSlice{},
		LastReading:        reading,
	}
	go func() {
		for {
			time.Sleep(secondsToStore)
			//before := len(retval.LastMinuteReadings)
			retval.LastMinuteReadings =
				retval.LastMinuteReadings.NotOlderThan(time.Now().Add(-1 *
					secondsToStore))
			//after := len(retval.LastMinuteReadings)
			//fmt.Printf("Cache cleanup: Before %d, after %d\r\n", before, after)
		}
	}()
	return retval
}

func (mr *MeterReadings) Purge(devid uint8) {
	mr.LastMinuteReadings = ReadingSlice{}
	mr.LastReading = Readings{
		UniqueID:       fmt.Sprintf(UniqueIDFormat, devid),
		ModbusDeviceID: devid,
	}
}

func (mr *MeterReadings) AddSnip(snip QuerySnip) {
	// 1. Merge the snip to the last values.
	reading := mr.LastReading
	reading.MergeSnip(snip)
	// 2. store it
	mr.LastReading = reading
	mr.LastMinuteReadings = append(mr.LastMinuteReadings, reading)
}
