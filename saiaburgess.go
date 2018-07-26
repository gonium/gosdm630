package sdm630

import (
	"math"
)

const (
	METERTYPE_SAIA = "SAIABURGESS"

	/***
	 * Opcodes for Saia Burgess ALE3
	 * http://datenblatt.stark-elektronik.de/saia_burgess/DE_DS_Energymeter-ALE3-with-Modbus.pdf
	 */
	OpCodeSaiaTotalImport   = 28 // double, scaler 100
	OpCodeSaiaPartialImport = 30 // double, scaler 100
	OpCodeSaiaTotalExport   = 32 // double, scaler 100
	OpCodeSaiaPartialExport = 34 // double, scaler 100

	OpCodeSaiaL1Voltage       = 36
	OpCodeSaiaL1Current       = 37 // scaler 10
	OpCodeSaiaL1Power         = 38 // scaler 100
	OpCodeSaiaL1ReactivePower = 39 // scaler 100
	OpCodeSaiaL1Cosphi        = 40 // scaler 100

	OpCodeSaiaL2Voltage       = 31
	OpCodeSaiaL2Current       = 32 // scaler 10
	OpCodeSaiaL2Power         = 43 // scaler 100
	OpCodeSaiaL2ReactivePower = 44 // scaler 100
	OpCodeSaiaL2Cosphi        = 45 // scaler 100

	OpCodeSaiaL3Voltage       = 46
	OpCodeSaiaL3Current       = 47 // scaler 10
	OpCodeSaiaL3Power         = 48 // scaler 100
	OpCodeSaiaL3ReactivePower = 49 // scaler 100
	OpCodeSaiaL3Cosphi        = 50 // scaler 100

	OpCodeSaiaTotalPower         = 51 // scaler 100
	OpCodeSaiaTotalReactivePower = 52 // scaler 100
)

type SaiaProducer struct {
}

func NewSaiaProducer() *SaiaProducer {
	return &SaiaProducer{}
}

func (p *SaiaProducer) snip(devid uint8, opcode uint16, iec string, readlen uint16) QuerySnip {
	return QuerySnip{
		DeviceId:    devid,
		FuncCode:    ReadHoldingReg,
		OpCode:      opcode,
		ReadLen:     readlen,
		Value:       math.NaN(),
		IEC61850:    iec,
		Description: GetIecDescription(iec),
	}
}

// snip16 creates modbus operation for single register
func (p *SaiaProducer) snip16(devid uint8, opcode uint16, iec string, scaler ...float64) QuerySnip {
	snip := p.snip(devid, opcode, iec, 1)

	snip.Transform = RTU16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTU16ScaledIntToFloat64(scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *SaiaProducer) snip32(devid uint8, opcode uint16, iec string, scaler ...float64) QuerySnip {
	snip := p.snip(devid, opcode, iec, 2)

	snip.Transform = RTU32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTU32ScaledIntToFloat64(scaler[0])
	}

	return snip
}

func (p *SaiaProducer) Probe(devid uint8) QuerySnip {
	return p.snip16(devid, OpCodeSaiaL1Voltage, "VolLocPhsA")
}

func (p *SaiaProducer) Produce(devid uint8) (res []QuerySnip) {
	res = append(res, p.snip16(devid, OpCodeSaiaL1Voltage, "VolLocPhsA"))
	res = append(res, p.snip16(devid, OpCodeSaiaL2Voltage, "VolLocPhsB"))
	res = append(res, p.snip16(devid, OpCodeSaiaL3Voltage, "VolLocPhsC"))

	res = append(res, p.snip16(devid, OpCodeSaiaL1Current, "AmpLocPhsA", 10))
	res = append(res, p.snip16(devid, OpCodeSaiaL2Current, "AmpLocPhsB", 10))
	res = append(res, p.snip16(devid, OpCodeSaiaL3Current, "AmpLocPhsC", 10))

	res = append(res, p.snip16(devid, OpCodeSaiaL1Power, "WLocPhsA", 100))
	res = append(res, p.snip16(devid, OpCodeSaiaL2Power, "WLocPhsB", 100))
	res = append(res, p.snip16(devid, OpCodeSaiaL3Power, "WLocPhsC", 100))

	res = append(res, p.snip16(devid, OpCodeSaiaL1Cosphi, "AngLocPhsA", 100))
	res = append(res, p.snip16(devid, OpCodeSaiaL2Cosphi, "AngLocPhsB", 100))
	res = append(res, p.snip16(devid, OpCodeSaiaL3Cosphi, "AngLocPhsC", 100))

	// res = append(res, p.snip16(devid, OpCodeSaiaTotalPower, "TotWLoc", 100))

	res = append(res, p.snip32(devid, OpCodeSaiaTotalImport, "TotkWhImport", 100))
	res = append(res, p.snip32(devid, OpCodeSaiaTotalExport, "TotkWhExport", 100))

	return res
}
