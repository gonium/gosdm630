package sdm630_test

import (
	"testing"
	"time"

	. "github.com/gonium/gosdm630"
)

func TestQuerySnipMerge(t *testing.T) {
	r := Readings{
		Timestamp: time.Now(),
		Unix:      time.Now().Unix(),
		DeviceId:  1,
		UniqueId:  "Instrument1",
		Power: ThreePhaseReadings{
			L1: F2fp(1.0), L2: F2fp(2.0), L3: F2fp(3.0),
		},
		Voltage: ThreePhaseReadings{
			L1: F2fp(1.0), L2: F2fp(2.0), L3: F2fp(3.0),
		},
		Current: ThreePhaseReadings{
			L1: F2fp(4.0), L2: F2fp(5.0), L3: F2fp(6.0),
		},
		Cosphi: ThreePhaseReadings{
			L1: F2fp(7.0), L2: F2fp(8.0), L3: F2fp(9.0),
		},
		Import: ThreePhaseReadings{
			L1: F2fp(10.0), L2: F2fp(11.0), L3: F2fp(12.0),
		},
		Export: ThreePhaseReadings{
			L1: F2fp(13.0), L2: F2fp(14.0), L3: F2fp(15.0),
		},
	}

	setvalue := float64(230.0)
	var sniptests = []struct {
		snip  QuerySnip
		param func(Readings) float64
	}{
		{QuerySnip{DeviceId: 1, OpCode: VoltageL1,
			IEC61850: "VolLocPhsA", Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Voltage.L1) },
		},
		{QuerySnip{DeviceId: 1, OpCode: VoltageL2, IEC61850: "VolLocPhsB",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Voltage.L2) },
		},
		{QuerySnip{DeviceId: 1, OpCode: VoltageL3, IEC61850: "VolLocPhsC",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Voltage.L3) },
		},
		{QuerySnip{DeviceId: 1, OpCode: CurrentL1, IEC61850: "AmpLocPhsA",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Current.L1) },
		},
		{QuerySnip{DeviceId: 1, OpCode: CurrentL2, IEC61850: "AmpLocPhsB",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Current.L2) },
		},
		{QuerySnip{DeviceId: 1, OpCode: CurrentL3, IEC61850: "AmpLocPhsC",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Current.L3) },
		},
		{QuerySnip{DeviceId: 1, OpCode: PowerL1, IEC61850: "WLocPhsA",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Power.L1) },
		},
		{QuerySnip{DeviceId: 1, OpCode: PowerL2, IEC61850: "WLocPhsB",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Power.L2) },
		},
		{QuerySnip{DeviceId: 1, OpCode: PowerL3, IEC61850: "WLocPhsC",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Power.L3) },
		},
		{QuerySnip{DeviceId: 1, OpCode: CosphiL1, IEC61850: "AngLocPhsA",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Cosphi.L1) },
		},
		{QuerySnip{DeviceId: 1, OpCode: CosphiL2, IEC61850: "AngLocPhsB",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Cosphi.L2) },
		},
		{QuerySnip{DeviceId: 1, OpCode: CosphiL3, IEC61850: "AngLocPhsC",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Cosphi.L3) },
		},
		{QuerySnip{DeviceId: 1, OpCode: ImportL1, IEC61850: "TotkWhImportPhsA",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Import.L1) },
		},
		{QuerySnip{DeviceId: 1, OpCode: ImportL2, IEC61850: "TotkWhImportPhsB",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Import.L2) },
		},
		{QuerySnip{DeviceId: 1, OpCode: ImportL3, IEC61850: "TotkWhImportPhsC",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Import.L3) },
		},
		{QuerySnip{DeviceId: 1, OpCode: ExportL1, IEC61850: "TotkWhExportPhsA",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Export.L1) },
		},
		{QuerySnip{DeviceId: 1, OpCode: ExportL2, IEC61850: "TotkWhExportPhsB",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Export.L2) },
		},
		{QuerySnip{DeviceId: 1, OpCode: ExportL3, IEC61850: "TotkWhExportPhsC",
			Value: setvalue},
			func(r Readings) float64 { return Fp2f(r.Export.L3) },
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
