package sdm630

import (
	"math"
)

const (
	METERTYPE_DZG = "DZG"
)

type DZGProducer struct {
	ops MeasurementTypes
}

func NewDZGProducer() *DZGProducer {
	/**
	 * Opcodes for DZG DVH4014.
	 * https://www.dzg.de/fileadmin/dzg/content/downloads/produkte-zaehler/dvh4013/Communication-Protocol_DVH4013.pdf
	 */
	ops := MeasurementTypes{
		ActivePower:   0x0000, // 0x0 instant values and parameters
		ReactivePower: 0x0002,
		VoltageL1:     0x0004,
		VoltageL2:     0x0006,
		VoltageL3:     0x0008,
		CurrentL1:     0x000A,
		CurrentL2:     0x000C,
		CurrentL3:     0x000E,
		Cosphi:        0x0010, // DVH4013
		Frequency:     0x0012, // DVH4013
		Import:        0x0014, // DVH4013
		Export:        0x0016, // DVH4013
		Import:        0x4000, // 0x4 energy
		ImportL1:      0x4020,
		ImportL2:      0x4040,
		ImportL3:      0x4060,
		Export:        0x4100,
		ExportL1:      0x4120,
		ExportL2:      0x4140,
		ExportL3:      0x4160,
		// 0x8 max demand
	}
	return &DZGProducer{
		ops: ops,
	}
}

func (p *DZGProducer) GetMeterType() string {
	return METERTYPE_DZG
}

func (p *DZGProducer) snip(devid uint8, iec MeasurementType, scaler ...float64) QuerySnip {
	opcode := p.ops[iec]

	transform := RTU32ToFloat64 // default conversion
	if len(scaler) > 0 {
		transform = MakeRTU32ScaledIntToFloat64(scaler[0])
	}

	snip := QuerySnip{
		DeviceId:  devid,
		FuncCode:  ReadHoldingReg,
		OpCode:    opcode,
		ReadLen:   2,
		Value:     math.NaN(),
		IEC61850:  iec,
		Transform: transform,
	}
	return snip
}

func (p *DZGProducer) Probe(devid uint8) QuerySnip {
	return p.snip(devid, VoltageL1, 100)
}

func (p *DZGProducer) Produce(devid uint8) (res []QuerySnip) {
	for _, op := range []MeasurementType{VoltageL1, VoltageL2, VoltageL1} {
		res = append(res, p.snip(devid, op, 100))
	}

	for _, op := range []MeasurementType{
		CurrentL1, CurrentL2, CurrentL1,
		Import, Export, Cosphi,
	} {
		res = append(res, p.snip(devid, op, 1000))
	}

	for _, op := range []MeasurementType{
		ImportL1, ImportL2, ImportL1,
		ExportL1, ExportL2, ExportL1,
	} {
		res = append(res, p.snip(devid, op))
	}

	return res
}
