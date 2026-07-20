package paa

type TypeOfPaX uint16

const (
	Type_OFP   TypeOfPaX = 0x0
	Type_DXT1  TypeOfPaX = 0xFF01
	Type_DXT2  TypeOfPaX = 0xFF02
	Type_DXT3  TypeOfPaX = 0xFF03
	Type_DXT4  TypeOfPaX = 0xFF04
	Type_DXT5  TypeOfPaX = 0xFF05
	Type_RGBA4 TypeOfPaX = 0x4444
	Type_RGBA5 TypeOfPaX = 0x1555
	Type_RGBA8 TypeOfPaX = 0x8888
	Type_Gray  TypeOfPaX = 0x8080
)
