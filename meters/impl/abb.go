package impl

import . "github.com/gonium/gosdm630/meters"

func init() {
	Register(NewABBProducer)
}

const (
	METERTYPE_ABB = "ABB"
)

type ABBProducer struct {
	RS485Core
	Opcodes
}

func NewABBProducer() Producer {
	/***
	 * http://datenblatt.stark-elektronik.de/Energiezaehler_B-Serie_Handbuch.pdf
	 */
	ops := Opcodes{
		VoltageL1: 0x5B00,
		VoltageL2: 0x5B02,
		VoltageL3: 0x5B04,

		CurrentL1: 0x5B0C,
		CurrentL2: 0x5B0E,
		CurrentL3: 0x5B10,

		Power:   0x5B14,
		PowerL1: 0x5B16,
		PowerL2: 0x5B18,
		PowerL3: 0x5B1A,
		
		ImportL1:  0x5460,
		ImportL2:  0x5464,
		ImportL3:  0x5468,
		Import:    0x5000,
		
		ExportL1:  0x546C,
		ExportL2:  0x5470,
		ExportL3:  0x5474,
		Export:    0x5004,

		Cosphi:   0x5B3A,
		CosphiL1: 0x5B3B,
		CosphiL2: 0x5B3C,
		CosphiL3: 0x5B3D,

		Frequency: 0x5B2C,
	}
	return &ABBProducer{Opcodes: ops}
}

func (p *ABBProducer) Type() string {
	return METERTYPE_ABB
}

func (p *ABBProducer) Description() string {
	return "ABB A/B-Series meters"
}

func (p *ABBProducer) snip(iec Measurement, readlen uint16) Operation {
	opcode := p.Opcodes[iec]
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   opcode,
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip16 creates modbus operation for single register
func (p *ABBProducer) snip16(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *ABBProducer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip64 creates modbus operation for double register
func (p *ABBProducer) snip64(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 4)

	snip.Transform = RTUUint64ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *ABBProducer) Probe() Operation {
	return p.snip16(VoltageL1)
}

func (p *ABBProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
		CurrentL1, CurrentL2, CurrentL3,
		Power, PowerL1, PowerL2, PowerL3,
		Cosphi, CosphiL1, CosphiL2, CosphiL3,
	} {
		res = append(res, p.snip32(op, 10))
	}

	for _, op := range []Measurement{
		Frequency,
	} {
		res = append(res, p.snip16(op,100))
	}

	for _, op := range []Measurement{
		Import, ImportL1, ImportL2, ImportL3,
		Export, ExportL1, ExportL2, ExportL3,
	} {
		res = append(res, p.snip64(op,100))
	}

	return res
}
