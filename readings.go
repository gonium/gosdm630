package sdm630

import (
	"fmt"
	"time"
)

type MeterReadings struct {
	Lastminutereadings ReadingSlice
	Lastreading        Readings
}

func NewMeterReadings(devid uint8, secondsToStore time.Duration) (retval *MeterReadings) {
	reading := Readings{
		UniqueId: fmt.Sprintf(UniqueIdFormat, devid),
		DeviceId: devid,
	}
	retval = &MeterReadings{
		Lastminutereadings: ReadingSlice{},
		Lastreading:        reading,
	}
	go func() {
		for {
			time.Sleep(secondsToStore)
			//before := len(retval.lastminutereadings)
			retval.Lastminutereadings =
				retval.Lastminutereadings.NotOlderThan(time.Now().Add(-1 *
					secondsToStore))
			//after := len(retval.lastminutereadings)
			//fmt.Printf("Cache cleanup: Before %d, after %d\r\n", before, after)
		}
	}()
	return retval
}

func (mr *MeterReadings) Purge(devid uint8) {
	mr.Lastminutereadings = ReadingSlice{}
	mr.Lastreading = Readings{
		UniqueId: fmt.Sprintf(UniqueIdFormat, devid),
		DeviceId: devid,
	}
}

func (mr *MeterReadings) AddSnip(snip QuerySnip) {
	// 1. Merge the snip to the last values.
	reading := mr.Lastreading
	reading.MergeSnip(snip)
	// 2. store it
	mr.Lastreading = reading
	mr.Lastminutereadings = append(mr.Lastminutereadings, reading)
}
