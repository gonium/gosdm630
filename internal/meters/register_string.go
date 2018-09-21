// Code generated by "stringer -type Register registers.go"; DO NOT EDIT.

package meters

import "strconv"

const _Register_name = "VoltageL1VoltageL2VoltageL3CurrentL1CurrentL2CurrentL3PowerL1PowerL2PowerL3ActivePowerL1ReactivePowerL1ActivePowerL2ReactivePowerL2ActivePowerL3ReactivePowerL3ImportL1ImportL2ImportL3ExportL1ExportL2ExportL3PowerFactorL1PowerFactorL2PowerFactorL3CosphiL1CosphiL2CosphiL3THDL1THDL2THDL3VoltageCurrentPowerActivePowerReactivePowerPowerFactorCosphiTHDFrequencyNetNetL1NetL2NetL3ActiveNetActiveNetL1ActiveNetL2ActiveNetL3ReactiveNetReactiveNetL1ReactiveNetL2ReactiveNetL3ImportExportActiveImportT1ActiveImportT2ReactiveImportT1ReactiveImportT2ActiveExportT1ActiveExportT2ReactiveExportT1ReactiveExportT2"

var _Register_index = [...]uint16{0, 9, 18, 27, 36, 45, 54, 61, 68, 75, 88, 103, 116, 131, 144, 159, 167, 175, 183, 191, 199, 207, 220, 233, 246, 254, 262, 270, 275, 280, 285, 292, 299, 304, 315, 328, 339, 345, 348, 357, 360, 365, 370, 375, 384, 395, 406, 417, 428, 441, 454, 467, 473, 479, 493, 507, 523, 539, 553, 567, 583, 599}

func (i Register) String() string {
	if i < 0 || i >= Register(len(_Register_index)-1) {
		return "Register(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Register_name[_Register_index[i]:_Register_index[i+1]]
}