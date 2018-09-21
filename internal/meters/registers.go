package meters

type Register int

type RegisterDefinitions map[Register]uint16

const (
	// phases
	VoltageL1 Register = iota
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
