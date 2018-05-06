package sdm630_test

import (
	"testing"
	"time"

	"github.com/gonium/gosdm630"
)

func TestQuerySnipMerge(t *testing.T) {
	r := sdm630.Readings{
		Timestamp:      time.Now(),
		Unix:           time.Now().Unix(),
		ModbusDeviceID: 1,
		UniqueID:       "Instrument1",
		Power: sdm630.ThreePhaseReadings{
			L1: sdm630.F2fp(1.0), L2: sdm630.F2fp(2.0), L3: sdm630.F2fp(3.0),
		},
		Voltage: sdm630.ThreePhaseReadings{
			L1: sdm630.F2fp(1.0), L2: sdm630.F2fp(2.0), L3: sdm630.F2fp(3.0),
		},
		Current: sdm630.ThreePhaseReadings{
			L1: sdm630.F2fp(4.0), L2: sdm630.F2fp(5.0), L3: sdm630.F2fp(6.0),
		},
		Cosphi: sdm630.ThreePhaseReadings{
			L1: sdm630.F2fp(7.0), L2: sdm630.F2fp(8.0), L3: sdm630.F2fp(9.0),
		},
		Import: sdm630.ThreePhaseReadings{
			L1: sdm630.F2fp(10.0), L2: sdm630.F2fp(11.0), L3: sdm630.F2fp(12.0),
		},
		Export: sdm630.ThreePhaseReadings{
			L1: sdm630.F2fp(13.0), L2: sdm630.F2fp(14.0), L3: sdm630.F2fp(15.0),
		},
	}

	setvalue := float64(230.0)
	var sniptests = []struct {
		snip  sdm630.QuerySnip
		param func(sdm630.Readings) float64
	}{
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML1Voltage,
			IEC61850: "VolLocPhsA", Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Voltage.L1) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML2Voltage, IEC61850: "VolLocPhsB",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Voltage.L2) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML3Voltage, IEC61850: "VolLocPhsC",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Voltage.L3) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML1Current, IEC61850: "AmpLocPhsA",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Current.L1) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML2Current, IEC61850: "AmpLocPhsB",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Current.L2) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML3Current, IEC61850: "AmpLocPhsC",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Current.L3) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML1Power, IEC61850: "WLocPhsA",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Power.L1) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML2Power, IEC61850: "WLocPhsB",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Power.L2) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML3Power, IEC61850: "WLocPhsC",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Power.L3) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML1Cosphi, IEC61850: "AngLocPhsA",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Cosphi.L1) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML2Cosphi, IEC61850: "AngLocPhsB",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Cosphi.L2) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML3Cosphi, IEC61850: "AngLocPhsC",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Cosphi.L3) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML1Import, IEC61850: "TotkWhImportPhsA",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Import.L1) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML2Import, IEC61850: "TotkWhImportPhsB",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Import.L2) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML3Import, IEC61850: "TotkWhImportPhsC",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Import.L3) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML1Export, IEC61850: "TotkWhExportPhsA",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Export.L1) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML2Export, IEC61850: "TotkWhExportPhsB",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Export.L2) },
		},
		{sdm630.QuerySnip{DeviceID: 1, OpCode: sdm630.OpCodeSDML3Export, IEC61850: "TotkWhExportPhsC",
			Value: setvalue},
			func(r sdm630.Readings) float64 { return sdm630.Fp2f(r.Export.L3) },
		},
	}

	for _, test := range sniptests {
		r.MergeSnip(test.snip)
		if test.param(r) != setvalue {
			t.Errorf("Merge of querysnip failed: Expected %.2f, got %.2f",
				setvalue, test.param(r))
		}
	}
}
