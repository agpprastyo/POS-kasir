package escpos

// ASCII codes
const (
	NUL = 0x00
	LF  = 0x0a
	ESC = 0x1b
	GS  = 0x1d
)

// Commands
var (
	// Initialize printer
	Init = []byte{ESC, '@'}

	// Cut paper
	Cut     = []byte{GS, 'V', 66, 0}
	CutFull = []byte{GS, 'V', 65, 0}

	// Text format
	BoldOn         = []byte{ESC, 'E', 1}
	BoldOff        = []byte{ESC, 'E', 0}
	DoubleHeightOn = []byte{GS, '!', 0x10}
	DoubleWidthOn  = []byte{GS, '!', 0x20}
	DoubleSizeOn   = []byte{GS, '!', 0x30}
	NormalSize     = []byte{GS, '!', 0x00}

	// Alignment
	AlignLeft   = []byte{ESC, 'a', 0}
	AlignCenter = []byte{ESC, 'a', 1}
	AlignRight  = []byte{ESC, 'a', 2}
)
