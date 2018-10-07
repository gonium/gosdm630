package meters

type Measurement int

type Measurements map[Measurement]uint16

const (
	// phases
	VoltageL1 Measurement = iota
	VoltageL2
	VoltageL3
	CurrentL1
	CurrentL2
	CurrentL3
	PowerL1
	PowerL2
	PowerL3
	ActivePowerL1
	ReactivePowerL1
	ActivePowerL2
	ReactivePowerL2
	ActivePowerL3
	ReactivePowerL3
	ImportL1
	ImportL2
	ImportL3
	ExportL1
	ExportL2
	ExportL3
	PowerFactorL1
	PowerFactorL2
	PowerFactorL3
	CosphiL1
	CosphiL2
	CosphiL3
	THDL1
	THDL2
	THDL3

	// sum/avg
	Voltage
	Current
	Power
	ActivePower
	ReactivePower
	PowerFactor
	Cosphi
	THD
	Frequency

	// energy
	Net
	NetL1
	NetL2
	NetL3
	ActiveNet
	ActiveNetL1
	ActiveNetL2
	ActiveNetL3
	ReactiveNet
	ReactiveNetL1
	ReactiveNetL2
	ReactiveNetL3
	Import
	Export
	ActiveImportT1
	ActiveImportT2
	ReactiveImportT1
	ReactiveImportT2
	ActiveExportT1
	ActiveExportT2
	ReactiveExportT1
	ReactiveExportT2
)

var iec = map[Measurement]string{
	CurrentL1: "L1 Current (A)",
	CurrentL2: "L2 Current (A)",
	CurrentL3: "L3 Current (A)",
	CosphiL1:  "L1 Cosphi",
	CosphiL2:  "L2 Cosphi",
	CosphiL3:  "L3 Cosphi",
	Frequency: "Frequency of supply voltages",
	THD:       "Average voltage to neutral THD (%)",
	THDL1:     "L1 Voltage to neutral THD (%)",
	THDL2:     "L2 Voltage to neutral THD (%)",
	THDL3:     "L3 Voltage to neutral THD (%)",
	Export:    "Total Export (kWh)",
	ExportL1:  "L1 Export (kWh)",
	ExportL2:  "L2 Export (kWh)",
	ExportL3:  "L3 Export (kWh)",
	Import:    "Total Import (kWh)",
	ImportL1:  "L1 Import (kWh)",
	ImportL2:  "L2 Import (kWh)",
	ImportL3:  "L3 Import (kWh)",
	VoltageL1: "L1 Voltage (V)",
	VoltageL2: "L2 Voltage (V)",
	VoltageL3: "L3 Voltage (V)",
	PowerL1:   "L1 Power (W)",
	PowerL2:   "L2 Power (W)",
	PowerL3:   "L3 Power (W)",
}

func (m *Measurement) Description() string {
	if description, ok := iec[*m]; ok {
		return description
	}
	return m.String()
}
