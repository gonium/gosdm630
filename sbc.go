package sdm630

import (
	"math"
)

const (
	METERTYPE_SBC = "SBC"
)

type SBCProducer struct {
	ops MeasurementTypes
}

func NewSBCProducer() *SBCProducer {
	/**
	 * Opcodes for Saia Burgess ALE3
	 * http://datenblatt.stark-elektronik.de/saia_burgess/DE_DS_Energymeter-ALE3-with-Modbus.pdf
	 */
	ops := MeasurementTypes{
		Import: 28 - 1, // double, scaler 100
		Export: 32 - 1, // double, scaler 100
		// PartialImport: 30 - 1, // double, scaler 100
		// PartialExport: 34 - 1, // double, scaler 100

		VoltageL1:       36 - 1,
		CurrentL1:       37 - 1, // scaler 10
		PowerL1:         38 - 1, // scaler 100
		ReactivePowerL1: 39 - 1, // scaler 100
		CosphiL1:        40 - 1, // scaler 100

		VoltageL2:       41 - 1,
		CurrentL2:       42 - 1, // scaler 10
		PowerL2:         43 - 1, // scaler 100
		ReactivePowerL2: 44 - 1, // scaler 100
		CosphiL2:        45 - 1, // scaler 100

		VoltageL3:       46 - 1,
		CurrentL3:       47 - 1, // scaler 10
		PowerL3:         48 - 1, // scaler 100
		ReactivePowerL3: 49 - 1, // scaler 100
		CosphiL3:        50 - 1, // scaler 100

		Power:         51 - 1, // scaler 100
		ReactivePower: 52 - 1, // scaler 100
	}
	return &SBCProducer{
		ops: ops,
	}
}

func (p *SBCProducer) GetMeterType() string {
	return METERTYPE_SBC
}

func (p *SBCProducer) snip(devid uint8, iec MeasurementType, readlen uint16) QuerySnip {
	opcode := p.ops[iec]
	return QuerySnip{
		DeviceId: devid,
		FuncCode: ReadHoldingReg,
		OpCode:   opcode,
		ReadLen:  readlen,
		Value:    math.NaN(),
		IEC61850: iec,
	}
}

// snip16 creates modbus operation for single register
func (p *SBCProducer) snip16(devid uint8, iec MeasurementType, scaler ...float64) QuerySnip {
	snip := p.snip(devid, iec, 1)

	snip.Transform = RTU16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTU16ScaledIntToFloat64(scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *SBCProducer) snip32(devid uint8, iec MeasurementType, scaler ...float64) QuerySnip {
	snip := p.snip(devid, iec, 2)

	snip.Transform = RTU32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTU32ScaledIntToFloat64(scaler[0])
	}

	return snip
}

func (p *SBCProducer) Probe(devid uint8) QuerySnip {
	return p.snip16(devid, VoltageL1)
}

func (p *SBCProducer) Produce(devid uint8) (res []QuerySnip) {
	for _, op := range []MeasurementType{VoltageL1, VoltageL2, VoltageL1} {
		res = append(res, p.snip16(devid, op))
	}

	for _, op := range []MeasurementType{CurrentL1, CurrentL2, CurrentL1} {
		res = append(res, p.snip16(devid, op, 10))
	}

	for _, op := range []MeasurementType{
		PowerL1, PowerL2, PowerL1,
		CosphiL1, CosphiL2, CosphiL1,
	} {
		res = append(res, p.snip16(devid, op, 100))
	}

	res = append(res, p.snip32(devid, Import, 100))
	res = append(res, p.snip32(devid, Export, 100))

	return res
}
