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

var TypeOfPaXStrings = map[TypeOfPaX]string{
	Type_OFP:   "OFP",
	Type_DXT1:  "DXT1",
	Type_DXT2:  "DXT2",
	Type_DXT3:  "DXT3",
	Type_DXT4:  "DXT4",
	Type_DXT5:  "DXT5",
	Type_RGBA4: "RGBA4",
	Type_RGBA5: "RGBA5",
	Type_RGBA8: "RGBA8",
	Type_Gray:  "Gray",
}
