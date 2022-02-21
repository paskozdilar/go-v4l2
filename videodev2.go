package v4l2

import (
	"syscall"
	"unsafe"
)

const (
	VideoMaxFrames = 32
	VideoMaxPlanes = 8
)

func Fourcc(a, b, c, d uint32) PixFmt {
	return PixFmt(a | b<<8 | c<<16 | d<<24)
}

func FourccBe(a, b, c, d uint32) PixFmt {
	return Fourcc(a, b, c, d) | 1<<31
}

type Field uint32

const (
	Field_Any          Field = 0
	Field_None         Field = 1
	Field_Top          Field = 2
	Field_Bottom       Field = 3
	Field_Interlaced   Field = 4
	Field_SeqTb        Field = 5
	Field_SeqBt        Field = 6
	Field_Alternate    Field = 7
	Field_InterlacedTb Field = 8
	Field_InterlacedBt Field = 9
)

func (f Field) HasTop() bool {
	return (f == Field_Top ||
		f == Field_Interlaced ||
		f == Field_InterlacedTb ||
		f == Field_InterlacedBt ||
		f == Field_SeqTb ||
		f == Field_SeqBt)
}

func (f Field) HasBottom() bool {
	return (f == Field_Bottom ||
		f == Field_Interlaced ||
		f == Field_InterlacedTb ||
		f == Field_InterlacedBt ||
		f == Field_SeqTb ||
		f == Field_SeqBt)
}

func (f Field) HasBoth() bool {
	return (f == Field_Interlaced ||
		f == Field_InterlacedTb ||
		f == Field_InterlacedBt ||
		f == Field_SeqTb ||
		f == Field_SeqBt)
}

func (f Field) HasTOrB() bool {
	return (f == Field_Bottom ||
		f == Field_Top ||
		f == Field_Alternate)
}

func (f Field) IsInterlaced() bool {
	return (f == Field_Interlaced ||
		f == Field_InterlacedTb ||
		f == Field_InterlacedBt)
}

func (f Field) IsSequential() bool {
	return (f == Field_SeqTb || f == Field_SeqBt)
}

type BufType uint32

const (
	BufType_VideoCapture       BufType = 1
	BufType_VideoOutput        BufType = 2
	BufType_VideoOverlay       BufType = 3
	BufType_VbiCapture         BufType = 4
	BufType_VbiOutput          BufType = 5
	BufType_SlicedVbiCapture   BufType = 6
	BufType_SlicedVbiOutput    BufType = 7
	BufType_VideoOutputOverlay BufType = 8
	BufType_VideoCaptureMplane BufType = 9
	BufType_VideoOutputMplane  BufType = 10
	BufType_SdrCapture         BufType = 11
	BufType_SdrOutput          BufType = 12
	BufType_MetaCapture        BufType = 13
	BufType_MetaOutput         BufType = 14
	BufType_Private            BufType = 0x80
)

func (bt BufType) IsMultiplanar() bool {
	return (bt == BufType_VideoCaptureMplane ||
		bt == BufType_VideoOutputMplane)
}

func (bt BufType) IsOutput() bool {
	return (bt == BufType_VideoOutput ||
		bt == BufType_VideoOutputMplane ||
		bt == BufType_VideoOverlay ||
		bt == BufType_VideoOutputOverlay ||
		bt == BufType_VbiOutput ||
		bt == BufType_SlicedVbiOutput ||
		bt == BufType_SdrOutput ||
		bt == BufType_MetaOutput)
}

type TunerType int

const (
	Tuner_Radio     TunerType = 1
	Tuner_AnalogTv  TunerType = 2
	Tuner_DigitalTv TunerType = 3
	Tuner_Sdr       TunerType = 4
	Tuner_Rf        TunerType = 5
)

type Memory int

const (
	Memory_Mmap    Memory = 1
	Memory_Userptr Memory = 2
	Memory_Overlay Memory = 3
	Memory_Dmabuf  Memory = 4
)

type Colorspace uint32

const (
	Colorspace_Default     Colorspace = 0
	Colorspace_Smpte170M   Colorspace = 1
	Colorspace_Smpte240M   Colorspace = 2
	Colorspace_Rec709      Colorspace = 3
	Colorspace_Bt878       Colorspace = 4
	Colorspace_470SystemM  Colorspace = 5
	Colorspace_470SystemBg Colorspace = 6
	Colorspace_Jpeg        Colorspace = 7
	Colorspace_Srgb        Colorspace = 8
	Colorspace_Oprgb       Colorspace = 9
	Colorspace_Bt2020      Colorspace = 10
	Colorspace_Raw         Colorspace = 11
	Colorspace_DciP3       Colorspace = 12
)

func MapColorspaceDefault(isSdtv, isHdtv bool) Colorspace {
	if isSdtv {
		return Colorspace_Smpte170M
	} else if isHdtv {
		return Colorspace_Rec709
	} else {
		return Colorspace_Srgb
	}
}

type XferFunc uint32

const (
	XferFunc_Default   XferFunc = 0
	XferFunc_709       XferFunc = 1
	XferFunc_Srgb      XferFunc = 2
	XferFunc_Oprgb     XferFunc = 3
	XferFunc_Smpte240m XferFunc = 4
	XferFunc_None      XferFunc = 5
	XferFunc_DciP3     XferFunc = 6
	XferFunc_Smpte2084 XferFunc = 7
)

func MapXferFuncDefault(colsp Colorspace) XferFunc {
	switch colsp {
	case Colorspace_Oprgb:
		return XferFunc_Oprgb
	case Colorspace_Smpte240M:
		return XferFunc_Smpte240m
	case Colorspace_DciP3:
		return XferFunc_DciP3
	case Colorspace_Raw:
		return XferFunc_None
	case Colorspace_Srgb:
		return XferFunc_Srgb
	case Colorspace_Jpeg:
		return XferFunc_Srgb
	default:
		return XferFunc_709
	}
}

type YcbcrEncoding uint32

const (
	YcbcrEnc_Default        YcbcrEncoding = 0
	YcbcrEnc_601            YcbcrEncoding = 1
	YcbcrEnc_709            YcbcrEncoding = 2
	YcbcrEnc_Xv601          YcbcrEncoding = 3
	YcbcrEnc_Xv709          YcbcrEncoding = 4
	YcbcrEnc_Sycc           YcbcrEncoding = 5
	YcbcrEnc_Bt2020         YcbcrEncoding = 6
	YcbcrEnc_Bt2020ConstLum YcbcrEncoding = 7
	YcbcrEnc_Smpte240M      YcbcrEncoding = 8
)

type HsvEncoding uint32

const (
	HsvEnc_108 HsvEncoding = 128
	HsvEnc_256 HsvEncoding = 129
)

func MapYcbcrEncDefault(colsp Colorspace) YcbcrEncoding {
	switch colsp {
	case Colorspace_Rec709:
		return YcbcrEnc_709
	case Colorspace_DciP3:
		return YcbcrEnc_709
	case Colorspace_Bt2020:
		return YcbcrEnc_Bt2020
	case Colorspace_Smpte240M:
		return YcbcrEnc_Smpte240M
	default:
		return YcbcrEnc_601
	}
}

type Quantization uint32

const (
	Quantization_Default   Quantization = 0
	Quantization_FullRange Quantization = 1
	Quantization_LimRange  Quantization = 2
)

func MapQuantizationDefault(isRgbOrHsv bool, colsp Colorspace, ycbcrEnc YcbcrEncoding) Quantization {
	if isRgbOrHsv || colsp == Colorspace_Jpeg {
		return Quantization_FullRange
	} else {
		return Quantization_LimRange
	}
}

type Priority uint32

const (
	Priority_Unset       Priority = 0
	Priority_Background  Priority = 1
	Priority_Interactive Priority = 2
	Priority_Record      Priority = 3
	Priority_Default     Priority = Priority_Interactive
)

type Rect struct {
	Left   int32
	Top    int32
	Width  uint32
	Height uint32
}

type Fract struct {
	Numerator   uint32
	Denominator uint32
}

type Capability struct {
	Driver       [16]uint8
	Card         [32]uint8
	BusInfo      [32]uint8
	Version      uint32
	Capabilities Cap
	DeviceCaps   uint32
	Reserved     [3]uint32
}

type Cap uint32

const (
	Cap_VideoCapture       Cap = 0x00000001
	Cap_VideoOutput        Cap = 0x00000002
	Cap_VideoOverlay       Cap = 0x00000004
	Cap_VbiCapture         Cap = 0x00000010
	Cap_VbiOutput          Cap = 0x00000020
	Cap_SlicedVbiCapture   Cap = 0x00000040
	Cap_SlicedVbiOutput    Cap = 0x00000080
	Cap_RdsCapture         Cap = 0x00000100
	Cap_VideoOutputOverlay Cap = 0x00000200
	Cap_HwFreqSeek         Cap = 0x00000400
	Cap_RdsOutput          Cap = 0x00000800
	Cap_VideoCaptureMplane Cap = 0x00001000
	Cap_VideoOutputMplane  Cap = 0x00002000
	Cap_VideoM2MMplane     Cap = 0x00004000
	Cap_VideoM2M           Cap = 0x00008000
	Cap_Tuner              Cap = 0x00010000
	Cap_Audio              Cap = 0x00020000
	Cap_Radio              Cap = 0x00040000
	Cap_Modulator          Cap = 0x00080000
	Cap_SdrCapture         Cap = 0x00100000
	Cap_ExtPixFormat       Cap = 0x00200000
	Cap_SdrOutput          Cap = 0x00400000
	Cap_MetaCapture        Cap = 0x00800000
	Cap_Readwrite          Cap = 0x01000000
	Cap_Asyncio            Cap = 0x02000000
	Cap_Streaming          Cap = 0x04000000
	Cap_MetaOutput         Cap = 0x08000000
	Cap_Touch              Cap = 0x10000000
	Cap_DeviceCaps         Cap = 0x80000000
)

type PixFormat struct {
	Width            uint32
	Height           uint32
	PixelFormat      PixFmt
	Field            Field
	BytesPerLine     uint32
	SizeImage        uint32
	Colorspace       Colorspace
	Priv             uint32
	Flags            PixFmtFlag
	YcbcrEncOrHsvEnc uint32
	//  YcbcrEncOrHsvEnc union {
	//		YcbcrEncoding
	//		HsvEncoding
	//  }
	Quantization Quantization
	XferFunc     uint32
}

type PixFmt uint32

// var because golang doesn't have compile time functions
// TODO: convert to Fourcc literal manually or create code generator
var (
	PixFmtRgb332         PixFmt = Fourcc('R', 'G', 'B', '1')
	PixFmtRgb444         PixFmt = Fourcc('R', '4', '4', '4')
	PixFmtArgb444        PixFmt = Fourcc('A', 'R', '1', '2')
	PixFmtXrgb444        PixFmt = Fourcc('X', 'R', '1', '2')
	PixFmtRgba444        PixFmt = Fourcc('R', 'A', '1', '2')
	PixFmtRgbx444        PixFmt = Fourcc('R', 'X', '1', '2')
	PixFmtAbgr444        PixFmt = Fourcc('A', 'B', '1', '2')
	PixFmtXbgr444        PixFmt = Fourcc('X', 'B', '1', '2')
	PixFmtBgra444        PixFmt = Fourcc('G', 'A', '1', '2')
	PixFmtBgrx444        PixFmt = Fourcc('B', 'X', '1', '2')
	PixFmtRgb555         PixFmt = Fourcc('R', 'G', 'B', 'O')
	PixFmtArgb555        PixFmt = Fourcc('A', 'R', '1', '5')
	PixFmtXrgb555        PixFmt = Fourcc('X', 'R', '1', '5')
	PixFmtRgba555        PixFmt = Fourcc('R', 'A', '1', '5')
	PixFmtRgbx555        PixFmt = Fourcc('R', 'X', '1', '5')
	PixFmtAbgr555        PixFmt = Fourcc('A', 'B', '1', '5')
	PixFmtXbgr555        PixFmt = Fourcc('X', 'B', '1', '5')
	PixFmtBgra555        PixFmt = Fourcc('B', 'A', '1', '5')
	PixFmtBgrx555        PixFmt = Fourcc('B', 'X', '1', '5')
	PixFmtRgb565         PixFmt = Fourcc('R', 'G', 'B', 'P')
	PixFmtRgb555X        PixFmt = Fourcc('R', 'G', 'B', 'Q')
	PixFmtArgb555X       PixFmt = FourccBe('A', 'R', '1', '5')
	PixFmtXrgb555X       PixFmt = FourccBe('X', 'R', '1', '5')
	PixFmtRgb565X        PixFmt = Fourcc('R', 'G', 'B', 'R')
	PixFmtBgr666         PixFmt = Fourcc('B', 'G', 'R', 'H')
	PixFmtBgr24          PixFmt = Fourcc('B', 'G', 'R', '3')
	PixFmtRgb24          PixFmt = Fourcc('R', 'G', 'B', '3')
	PixFmtBgr32          PixFmt = Fourcc('B', 'G', 'R', '4')
	PixFmtAbgr32         PixFmt = Fourcc('A', 'R', '2', '4')
	PixFmtXbgr32         PixFmt = Fourcc('X', 'R', '2', '4')
	PixFmtBgra32         PixFmt = Fourcc('R', 'A', '2', '4')
	PixFmtBgrx32         PixFmt = Fourcc('R', 'X', '2', '4')
	PixFmtRgb32          PixFmt = Fourcc('R', 'G', 'B', '4')
	PixFmtRgba32         PixFmt = Fourcc('A', 'B', '2', '4')
	PixFmtRgbx32         PixFmt = Fourcc('X', 'B', '2', '4')
	PixFmtArgb32         PixFmt = Fourcc('B', 'A', '2', '4')
	PixFmtXrgb32         PixFmt = Fourcc('B', 'X', '2', '4')
	PixFmtGrey           PixFmt = Fourcc('G', 'R', 'E', 'Y')
	PixFmtY4             PixFmt = Fourcc('Y', '0', '4', ' ')
	PixFmtY6             PixFmt = Fourcc('Y', '0', '6', ' ')
	PixFmtY10            PixFmt = Fourcc('Y', '1', '0', ' ')
	PixFmtY12            PixFmt = Fourcc('Y', '1', '2', ' ')
	PixFmtY16            PixFmt = Fourcc('Y', '1', '6', ' ')
	PixFmtY16_Be         PixFmt = FourccBe('Y', '1', '6', ' ')
	PixFmtY10Bpack       PixFmt = Fourcc('Y', '1', '0', 'B')
	PixFmtY10P           PixFmt = Fourcc('Y', '1', '0', 'P')
	PixFmtPal8           PixFmt = Fourcc('P', 'A', 'L', '8')
	PixFmtUv8            PixFmt = Fourcc('U', 'V', '8', ' ')
	PixFmtYuyv           PixFmt = Fourcc('Y', 'U', 'Y', 'V')
	PixFmtYyuv           PixFmt = Fourcc('Y', 'Y', 'U', 'V')
	PixFmtYvyu           PixFmt = Fourcc('Y', 'V', 'Y', 'U')
	PixFmtUyvy           PixFmt = Fourcc('U', 'Y', 'V', 'Y')
	PixFmtVyuy           PixFmt = Fourcc('V', 'Y', 'U', 'Y')
	PixFmtY41P           PixFmt = Fourcc('Y', '4', '1', 'P')
	PixFmtYuv444         PixFmt = Fourcc('Y', '4', '4', '4')
	PixFmtYuv555         PixFmt = Fourcc('Y', 'U', 'V', 'O')
	PixFmtYuv565         PixFmt = Fourcc('Y', 'U', 'V', 'P')
	PixFmtYuv32          PixFmt = Fourcc('Y', 'U', 'V', '4')
	PixFmtAyuv32         PixFmt = Fourcc('A', 'Y', 'U', 'V')
	PixFmtXyuv32         PixFmt = Fourcc('X', 'Y', 'U', 'V')
	PixFmtVuya32         PixFmt = Fourcc('V', 'U', 'Y', 'A')
	PixFmtVuyx32         PixFmt = Fourcc('V', 'U', 'Y', 'X')
	PixFmtHi240          PixFmt = Fourcc('H', 'I', '2', '4')
	PixFmtHm12           PixFmt = Fourcc('H', 'M', '1', '2')
	PixFmtM420           PixFmt = Fourcc('M', '4', '2', '0')
	PixFmtNv12           PixFmt = Fourcc('N', 'V', '1', '2')
	PixFmtNv21           PixFmt = Fourcc('N', 'V', '2', '1')
	PixFmtNv16           PixFmt = Fourcc('N', 'V', '1', '6')
	PixFmtNv61           PixFmt = Fourcc('N', 'V', '6', '1')
	PixFmtNv24           PixFmt = Fourcc('N', 'V', '2', '4')
	PixFmtNv42           PixFmt = Fourcc('N', 'V', '4', '2')
	PixFmtNv12M          PixFmt = Fourcc('N', 'M', '1', '2')
	PixFmtNv21M          PixFmt = Fourcc('N', 'M', '2', '1')
	PixFmtNv16M          PixFmt = Fourcc('N', 'M', '1', '6')
	PixFmtNv61M          PixFmt = Fourcc('N', 'M', '6', '1')
	PixFmtNv12Mt         PixFmt = Fourcc('T', 'M', '1', '2')
	PixFmtNv12Mt16X16    PixFmt = Fourcc('V', 'M', '1', '2')
	PixFmtYuv410         PixFmt = Fourcc('Y', 'U', 'V', '9')
	PixFmtYvu410         PixFmt = Fourcc('Y', 'V', 'U', '9')
	PixFmtYuv411P        PixFmt = Fourcc('4', '1', '1', 'P')
	PixFmtYuv420         PixFmt = Fourcc('Y', 'U', '1', '2')
	PixFmtYvu420         PixFmt = Fourcc('Y', 'V', '1', '2')
	PixFmtYuv422P        PixFmt = Fourcc('4', '2', '2', 'P')
	PixFmtYuv420M        PixFmt = Fourcc('Y', 'M', '1', '2')
	PixFmtYvu420M        PixFmt = Fourcc('Y', 'M', '2', '1')
	PixFmtYuv422M        PixFmt = Fourcc('Y', 'M', '1', '6')
	PixFmtYvu422M        PixFmt = Fourcc('Y', 'M', '6', '1')
	PixFmtYuv444M        PixFmt = Fourcc('Y', 'M', '2', '4')
	PixFmtYvu444M        PixFmt = Fourcc('Y', 'M', '4', '2')
	PixFmtSbggr8         PixFmt = Fourcc('B', 'A', '8', '1')
	PixFmtSgbrg8         PixFmt = Fourcc('G', 'B', 'R', 'G')
	PixFmtSgrbg8         PixFmt = Fourcc('G', 'R', 'B', 'G')
	PixFmtSrggb8         PixFmt = Fourcc('R', 'G', 'G', 'B')
	PixFmtSbggr10        PixFmt = Fourcc('B', 'G', '1', '0')
	PixFmtSgbrg10        PixFmt = Fourcc('G', 'B', '1', '0')
	PixFmtSgrbg10        PixFmt = Fourcc('B', 'A', '1', '0')
	PixFmtSrggb10        PixFmt = Fourcc('R', 'G', '1', '0')
	PixFmtSbggr10P       PixFmt = Fourcc('p', 'B', 'A', 'A')
	PixFmtSgbrg10P       PixFmt = Fourcc('p', 'G', 'A', 'A')
	PixFmtSgrbg10P       PixFmt = Fourcc('p', 'g', 'A', 'A')
	PixFmtSrggb10P       PixFmt = Fourcc('p', 'R', 'A', 'A')
	PixFmtSbggr10Alaw8   PixFmt = Fourcc('a', 'B', 'A', '8')
	PixFmtSgbrg10Alaw8   PixFmt = Fourcc('a', 'G', 'A', '8')
	PixFmtSgrbg10Alaw8   PixFmt = Fourcc('a', 'g', 'A', '8')
	PixFmtSrggb10Alaw8   PixFmt = Fourcc('a', 'R', 'A', '8')
	PixFmtSbggr10Dpcm8   PixFmt = Fourcc('b', 'B', 'A', '8')
	PixFmtSgbrg10Dpcm8   PixFmt = Fourcc('b', 'G', 'A', '8')
	PixFmtSgrbg10Dpcm8   PixFmt = Fourcc('B', 'D', '1', '0')
	PixFmtSrggb10Dpcm8   PixFmt = Fourcc('b', 'R', 'A', '8')
	PixFmtSbggr12        PixFmt = Fourcc('B', 'G', '1', '2')
	PixFmtSgbrg12        PixFmt = Fourcc('G', 'B', '1', '2')
	PixFmtSgrbg12        PixFmt = Fourcc('B', 'A', '1', '2')
	PixFmtSrggb12        PixFmt = Fourcc('R', 'G', '1', '2')
	PixFmtSbggr12P       PixFmt = Fourcc('p', 'B', 'C', 'C')
	PixFmtSgbrg12P       PixFmt = Fourcc('p', 'G', 'C', 'C')
	PixFmtSgrbg12P       PixFmt = Fourcc('p', 'g', 'C', 'C')
	PixFmtSrggb12P       PixFmt = Fourcc('p', 'R', 'C', 'C')
	PixFmtSbggr14P       PixFmt = Fourcc('p', 'B', 'E', 'E')
	PixFmtSgbrg14P       PixFmt = Fourcc('p', 'G', 'E', 'E')
	PixFmtSgrbg14P       PixFmt = Fourcc('p', 'g', 'E', 'E')
	PixFmtSrggb14P       PixFmt = Fourcc('p', 'R', 'E', 'E')
	PixFmtSbggr16        PixFmt = Fourcc('B', 'Y', 'R', '2')
	PixFmtSgbrg16        PixFmt = Fourcc('G', 'B', '1', '6')
	PixFmtSgrbg16        PixFmt = Fourcc('G', 'R', '1', '6')
	PixFmtSrggb16        PixFmt = Fourcc('R', 'G', '1', '6')
	PixFmtHsv24          PixFmt = Fourcc('H', 'S', 'V', '3')
	PixFmtHsv32          PixFmt = Fourcc('H', 'S', 'V', '4')
	PixFmtMjpeg          PixFmt = Fourcc('M', 'J', 'P', 'G')
	PixFmtJpeg           PixFmt = Fourcc('J', 'P', 'E', 'G')
	PixFmtDv             PixFmt = Fourcc('d', 'v', 's', 'd')
	PixFmtMpeg           PixFmt = Fourcc('M', 'P', 'E', 'G')
	PixFmtH264           PixFmt = Fourcc('H', '2', '6', '4')
	PixFmtH264_NoSc      PixFmt = Fourcc('A', 'V', 'C', '1')
	PixFmtH264_Mvc       PixFmt = Fourcc('M', '2', '6', '4')
	PixFmtH263           PixFmt = Fourcc('H', '2', '6', '3')
	PixFmtMpeg1          PixFmt = Fourcc('M', 'P', 'G', '1')
	PixFmtMpeg2          PixFmt = Fourcc('M', 'P', 'G', '2')
	PixFmtMpeg2_Slice    PixFmt = Fourcc('M', 'G', '2', 'S')
	PixFmtMpeg4          PixFmt = Fourcc('M', 'P', 'G', '4')
	PixFmtXvid           PixFmt = Fourcc('X', 'V', 'I', 'D')
	PixFmtVc1_AnnexG     PixFmt = Fourcc('V', 'C', '1', 'G')
	PixFmtVc1_AnnexL     PixFmt = Fourcc('V', 'C', '1', 'L')
	PixFmtVp8            PixFmt = Fourcc('V', 'P', '8', '0')
	PixFmtVp9            PixFmt = Fourcc('V', 'P', '9', '0')
	PixFmtHevc           PixFmt = Fourcc('H', 'E', 'V', 'C')
	PixFmtFwht           PixFmt = Fourcc('F', 'W', 'H', 'T')
	PixFmtFwhtStateless  PixFmt = Fourcc('S', 'F', 'W', 'H')
	PixFmtCpia1          PixFmt = Fourcc('C', 'P', 'I', 'A')
	PixFmtWnva           PixFmt = Fourcc('W', 'N', 'V', 'A')
	PixFmtSn9C10X        PixFmt = Fourcc('S', '9', '1', '0')
	PixFmtSn9C20Xi420    PixFmt = Fourcc('S', '9', '2', '0')
	PixFmtPwc1           PixFmt = Fourcc('P', 'W', 'C', '1')
	PixFmtPwc2           PixFmt = Fourcc('P', 'W', 'C', '2')
	PixFmtEt61X251       PixFmt = Fourcc('E', '6', '2', '5')
	PixFmtSpca501        PixFmt = Fourcc('S', '5', '0', '1')
	PixFmtSpca505        PixFmt = Fourcc('S', '5', '0', '5')
	PixFmtSpca508        PixFmt = Fourcc('S', '5', '0', '8')
	PixFmtSpca561        PixFmt = Fourcc('S', '5', '6', '1')
	PixFmtPac207         PixFmt = Fourcc('P', '2', '0', '7')
	PixFmtMr97310A       PixFmt = Fourcc('M', '3', '1', '0')
	PixFmtJl2005Bcd      PixFmt = Fourcc('J', 'L', '2', '0')
	PixFmtSn9C2028       PixFmt = Fourcc('S', 'O', 'N', 'X')
	PixFmtSq905C         PixFmt = Fourcc('9', '0', '5', 'C')
	PixFmtPjpg           PixFmt = Fourcc('P', 'J', 'P', 'G')
	PixFmtOv511          PixFmt = Fourcc('O', '5', '1', '1')
	PixFmtOv518          PixFmt = Fourcc('O', '5', '1', '8')
	PixFmtStv0680        PixFmt = Fourcc('S', '6', '8', '0')
	PixFmtTm6000         PixFmt = Fourcc('T', 'M', '6', '0')
	PixFmtCitYyvyuy      PixFmt = Fourcc('C', 'I', 'T', 'V')
	PixFmtKonica420      PixFmt = Fourcc('K', 'O', 'N', 'I')
	PixFmtJpgl           PixFmt = Fourcc('J', 'P', 'G', 'L')
	PixFmtSe401          PixFmt = Fourcc('S', '4', '0', '1')
	PixFmtS5CuyvyJpg     PixFmt = Fourcc('S', '5', 'C', 'I')
	PixFmtY8I            PixFmt = Fourcc('Y', '8', 'I', ' ')
	PixFmtY12I           PixFmt = Fourcc('Y', '1', '2', 'I')
	PixFmtZ16            PixFmt = Fourcc('Z', '1', '6', ' ')
	PixFmtMt21C          PixFmt = Fourcc('M', 'T', '2', '1')
	PixFmtInzi           PixFmt = Fourcc('I', 'N', 'Z', 'I')
	PixFmtSunxiTiledNv12 PixFmt = Fourcc('S', 'T', '1', '2')
	PixFmtCnf4           PixFmt = Fourcc('C', 'N', 'F', '4')
	PixFmtIpu3_Sbggr10   PixFmt = Fourcc('i', 'p', '3', 'b')
	PixFmtIpu3_Sgbrg10   PixFmt = Fourcc('i', 'p', '3', 'g')
	PixFmtIpu3_Sgrbg10   PixFmt = Fourcc('i', 'p', '3', 'G')
	PixFmtIpu3_Srggb10   PixFmt = Fourcc('i', 'p', '3', 'r')
	SdrFmtCu8            PixFmt = Fourcc('C', 'U', '0', '8')
	SdrFmtCu16Le         PixFmt = Fourcc('C', 'U', '1', '6')
	SdrFmtCs8            PixFmt = Fourcc('C', 'S', '0', '8')
	SdrFmtCs14Le         PixFmt = Fourcc('C', 'S', '1', '4')
	SdrFmtRu12Le         PixFmt = Fourcc('R', 'U', '1', '2')
	SdrFmtPcu16Be        PixFmt = Fourcc('P', 'C', '1', '6')
	SdrFmtPcu18Be        PixFmt = Fourcc('P', 'C', '1', '8')
	SdrFmtPcu20Be        PixFmt = Fourcc('P', 'C', '2', '0')
	TchFmtDeltaTd16      PixFmt = Fourcc('T', 'D', '1', '6')
	TchFmtDeltaTd08      PixFmt = Fourcc('T', 'D', '0', '8')
	TchFmtTu16           PixFmt = Fourcc('T', 'U', '1', '6')
	TchFmtTu08           PixFmt = Fourcc('T', 'U', '0', '8')
	MetaFmtVsp1Hgo       PixFmt = Fourcc('V', 'S', 'P', 'H')
	MetaFmtVsp1Hgt       PixFmt = Fourcc('V', 'S', 'P', 'T')
	MetaFmtUvc           PixFmt = Fourcc('U', 'V', 'C', 'H')
	MetaFmtD4Xx          PixFmt = Fourcc('D', '4', 'X', 'X')
	PixFmtPrivMagic      PixFmt = 0xfeedcafe
)

type PixFmtFlag uint32

const (
	PixFmtFlag_PremulAlpha PixFmtFlag = 0x00000001
)

type FmtDesc struct {
	Index       uint32
	Type        BufType
	Flags       FmtFlag
	Description [32]uint8
	PixelFormat PixFmt
	Reserved    [4]uint32
}

type FmtFlag uint32

const (
	FmtFlag_Compressed           FmtFlag = 0x0001
	FmtFlag_Emulated             FmtFlag = 0x0002
	FmtFlag_ContinuousBytestream FmtFlag = 0x0004
	FmtFlag_DynResolution        FmtFlag = 0x0008
)

type FrmSizeType uint32

const (
	FrmSizeType_Discrete   FrmSizeType = 1
	FrmSizeType_Continuous FrmSizeType = 2
	FrmSizeType_Stepwise   FrmSizeType = 3
)

type FrmSizeDiscrete struct {
	Width  uint32
	Height uint32
}

type FrmSizeStepwise struct {
	MinWidth   uint32
	MaxWidth   uint32
	StepWidth  uint32
	MinHeight  uint32
	MaxHeight  uint32
	StepHeight uint32
}

type FrmSizeEnum struct {
	Index                     uint32
	PixelFormat               PixFmt
	Type                      FrmSizeType
	FrmSizeDiscreteOrStepwise [6]uint32
	//  FrmSizeDiscreteOrStepwise union {
	//		FrmSizeDiscrete
	//		FrmSizeStepwise
	//	}
	Reserved [2]uint32
}

type FrmIvalTypes uint32

const (
	FrmIvalType_Discrete   FrmIvalTypes = 1
	FrmIvalType_Continuous FrmIvalTypes = 2
	FrmIvalType_Stepwise   FrmIvalTypes = 3
)

type FrmIvalStepwise struct {
	Min  Fract
	Max  Fract
	Step Fract
}

type FrmIvalEnum struct {
	Index                  uint32
	PixelFormat            PixFmt
	Width                  uint32
	Height                 uint32
	Type                   uint32
	FractOrFrmIvalStepwise [3]Fract
	//  FractOrFmnIvalStepwise uniont {
	//		Fract
	//		FrmIvalStepwise
	//	}
	Reserved [2]uint32
}

type Timecode struct {
	Type     TcType
	Flags    TcFlag
	Frames   uint8
	Seconds  uint8
	Minutes  uint8
	Hours    uint8
	Userbits [4]uint8
}

type TcType uint32

const (
	TcType_24Fps TcType = 1
	TcType_25Fps TcType = 2
	TcType_30Fps TcType = 3
	TcType_50Fps TcType = 4
	TcType_60Fps TcType = 5
)

type TcFlag uint32

const (
	TcFlag_Dropframe  TcFlag = 0x0001
	TcFlag_Colorframe TcFlag = 0x0002

	TcUserbits_Field       = 0x000C
	TcUserbits_UserDefined = 0x0000
	TcUserbits_8BitChars   = 0x0008
)

type JpegCompression struct {
	Quality     int
	AppN        int
	AppLen      int
	AppData     [60]byte
	ComLen      int
	ComData     [60]byte
	JpegMarkers JpegMarker
}

type JpegMarker uint32

const (
	JpegMarker_Dht JpegMarker = 1 << 3
	JpegMarker_Dqt JpegMarker = 1 << 4
	JpegMarker_Dri JpegMarker = 1 << 5
	JpegMarker_Com JpegMarker = 1 << 6
	JpegMarker_App JpegMarker = 1 << 7
)

type RequestBuffers struct {
	Count        uint32
	Type         BufType
	Memory       Memory
	Capabilities BufCap
	Reserved     [1]uint32
}

type BufCap uint32

const (
	BufCap_SupportsMmap         BufCap = 1 << 0
	BufCap_SupportsUserptr      BufCap = 1 << 1
	BufCap_SupportsDmabuf       BufCap = 1 << 2
	BufCap_SupportsRequests     BufCap = 1 << 3
	BufCap_SupportsOrphanedBufs BufCap = 1 << 4
)

type Plane struct {
	BytesUsed              uint32
	Length                 uint32
	MemOffsetOrUserptrOrFd [8]byte // uint32 and uintptr and int32 union type
	//	MemOffsetOrUserptrOrFd union {
	//		MemOffset uint32
	//		Userptr uintptr
	//		Fd int32
	//	}
	DataOffset uint32
	Resertved  [11]uint32
}

type Buffer struct {
	Index     uint32
	Type      BufType
	BytesUsed uint32
	Flags     BufFlag
	Field     Field
	Timestamp [16]byte // Timestamp: [sys/time.h] sizeof(struct timeval) == 16
	// TODO: Write a timeval struct?
	Timecode                    Timecode
	Sequence                    uint32
	Memory                      uint32
	OffsetOrUserptrOrPlanesOrFd [8]byte
	//	OffsetOrUserptrOrPlanesOrFd union {
	//		Offset uint32
	//		Userptr uintptr
	//		Planes uintptr(Plane)
	//		Fd int32
	//	}
	Length              uint32
	Reserved2           uint32
	RequestFdOrReserved [4]byte
	//	RequestFdOrReserved union {
	//		RequestFd int32
	//		Reserved uint32
	//	}
}

// skipped v4l2_timeval_to_ns()

type BufFlag uint32

const (
	BufFlag_Mapped             BufFlag = 0x00000001
	BufFlag_Queued             BufFlag = 0x00000002
	BufFlag_Done               BufFlag = 0x00000004
	BufFlag_Keyframe           BufFlag = 0x00000008
	BufFlag_Pframe             BufFlag = 0x00000010
	BufFlag_Bframe             BufFlag = 0x00000020
	BufFlag_Error              BufFlag = 0x00000040
	BufFlag_InRequest          BufFlag = 0x00000080
	BufFlag_Timecode           BufFlag = 0x00000100
	BufFlag_Prepared           BufFlag = 0x00000400
	BufFlag_NoCacheInvalidate  BufFlag = 0x00000800
	BufFlag_NoCacheClean       BufFlag = 0x00001000
	BufFlag_TimestampMask      BufFlag = 0x0000e000
	BufFlag_TimestampUnknown   BufFlag = 0x00000000
	BufFlag_TimestampMonotonic BufFlag = 0x00002000
	BufFlag_TimestampCopy      BufFlag = 0x00004000
	BufFlag_TstampSrcMask      BufFlag = 0x00070000
	BufFlag_TstampSrcEof       BufFlag = 0x00000000
	BufFlag_TstampSrcSoe       BufFlag = 0x00010000
	BufFlag_Last               BufFlag = 0x00100000
	BufFlag_RequestFd          BufFlag = 0x00800000
)

type ExportBuffer struct {
	Type     BufType
	Index    uint32
	Plane    uint32
	Flags    uint32
	Fd       int32
	Reserved [11]uint32
}

type FrameBuffer struct {
	Capability FbufCap
	Flags      FbufFlag
	Base       uintptr
	Fmt        struct {
		Width        uint32
		Height       uint32
		PixelFormat  PixFmt
		Field        Field
		BytesPerLine uint32
		SizeImage    uint32
		Colorspace   Colorspace
		Priv         uint32
	}
}

type FbufCap uint32

const (
	FbufCap_Externoverlay  FbufCap = 0x0001
	FbufCap_Chromakey      FbufCap = 0x0002
	FbufCap_ListClipping   FbufCap = 0x0004
	FbufCap_BitmapClipping FbufCap = 0x0008
	FbufCap_LocalAlpha     FbufCap = 0x0010
	FbufCap_GlobalAlpha    FbufCap = 0x0020
	FbufCap_LocalInvAlpha  FbufCap = 0x0040
	FbufCap_SrcChromakey   FbufCap = 0x0080
)

type FbufFlag uint32

const (
	FbufFlag_Primary       FbufFlag = 0x0001
	FbufFlag_Overlay       FbufFlag = 0x0002
	FbufFlag_Chromakey     FbufFlag = 0x0004
	FbufFlag_LocalAlpha    FbufFlag = 0x0008
	FbufFlag_GlobalAlpha   FbufFlag = 0x0010
	FbufFlag_LocalInvAlpha FbufFlag = 0x0020
	FbufFlag_SrcChromakey  FbufFlag = 0x0040
)

type Clip struct {
	C    Rect
	Next uintptr // *Clip
}

type Window struct {
	W           Rect
	Field       Field
	Clips       uintptr // *Clip
	Clipcount   uint32
	Bitmap      uintptr
	GlobalAlpha uint8
}

type CaptureParm struct {
	Capability   CaptureFlag
	CaptureMode  CaptureFlag
	TimePerFrame Fract
	ExtendedMode uint32
	ReadBuffers  uint32
	Reserved     [4]uint32
}

type CaptureFlag uint32

const (
	CaptureFlag_ModeHighQuality CaptureFlag = 0x0001
	CaptureFlag_CapTimePerFrame CaptureFlag = 0x1000
)

type OutputParm struct {
	Capability   uint32
	OutputMode   uint32
	TimePerFrame Fract
	ExtendedMode uint32
	WriteBuffers uint32
	Reserved     [4]uint32
}

type CropCap struct {
	Type        BufType
	Bounds      Rect
	DefRect     Rect
	PixelAspect Rect
}

type Crop struct {
	Type BufType
	C    Rect
}

type Selection struct {
	Type     uint32
	Target   uint32
	Flags    uint32
	R        Rect
	Reserved [9]uint32
}

type StdId uint64

const (
	Std_PalB       StdId = 0x00000001
	Std_PalB1      StdId = 0x00000002
	Std_PalG       StdId = 0x00000004
	Std_PalH       StdId = 0x00000008
	Std_PalI       StdId = 0x00000010
	Std_PalD       StdId = 0x00000020
	Std_PalD1      StdId = 0x00000040
	Std_PalK       StdId = 0x00000080
	Std_PalM       StdId = 0x00000100
	Std_PalN       StdId = 0x00000200
	Std_PalNc      StdId = 0x00000400
	Std_Pal60      StdId = 0x00000800
	Std_NtscM      StdId = 0x00001000
	Std_NtscMJp    StdId = 0x00002000
	Std_Ntsc443    StdId = 0x00004000
	Std_NtscMKr    StdId = 0x00008000
	Std_SecamB     StdId = 0x00010000
	Std_SecamD     StdId = 0x00020000
	Std_SecamG     StdId = 0x00040000
	Std_SecamH     StdId = 0x00080000
	Std_SecamK     StdId = 0x00100000
	Std_SecamK1    StdId = 0x00200000
	Std_SecamL     StdId = 0x00400000
	Std_SecamLc    StdId = 0x00800000
	Std_Atsc8_Vsb  StdId = 0x01000000
	Std_Atsc16_Vsb StdId = 0x02000000
	Std_Ntsc       StdId = Std_NtscM | Std_NtscMJp | Std_NtscMKr
	Std_SecamDk    StdId = Std_SecamD | Std_SecamK | Std_SecamK1
	Std_Secam      StdId = Std_SecamB | Std_SecamG | Std_SecamH | Std_SecamDk | Std_SecamL | Std_SecamLc
	Std_PalBg      StdId = Std_PalB | Std_PalB1 | Std_PalG
	Std_PalDk      StdId = Std_PalD | Std_PalD1 | Std_PalK
	Std_Pal        StdId = Std_PalBg | Std_PalDk | Std_PalH | Std_PalI

	Std_B       StdId = Std_PalB | Std_PalB1 | Std_SecamB
	Std_G       StdId = Std_PalG | Std_SecamG
	Std_H       StdId = Std_PalH | Std_SecamH
	Std_L       StdId = Std_SecamL | Std_SecamLc
	Std_Gh      StdId = Std_G | Std_H
	Std_Dk      StdId = Std_PalDk | Std_SecamDk
	Std_Bg      StdId = Std_B | Std_G
	Std_Mn      StdId = Std_PalM | Std_PalN | Std_PalNc | Std_Ntsc
	Std_Mts     StdId = Std_NtscM | Std_PalM | Std_PalN | Std_PalNc
	Std_525_60  StdId = Std_PalM | Std_Pal60 | Std_Ntsc | Std_Ntsc443
	Std_625_50  StdId = Std_Pal | Std_PalN | Std_PalNc | Std_Secam
	Std_Atsc    StdId = Std_Atsc8_Vsb | Std_Atsc16_Vsb
	Std_Unknown StdId = 0
	Std_All     StdId = Std_525_60 | Std_625_50
)

type Standard struct {
	Index       uint32
	Id          StdId
	Name        [24]uint8
	FramePeriod Fract
	FrameLines  uint32
	Reserved    [4]uint32
}

type BtTimings struct {
	Width         uint32
	Height        uint32
	Interlaced    Dv
	Polarities    DvPol
	Pixelclock    uint64
	Hfrontporch   uint32
	Hsync         uint32
	Hbackporch    uint32
	Vfrontporch   uint32
	Vsync         uint32
	Vbackportch   uint32
	IlVfrontporch uint32
	IlVsync       uint32
	IlVbackporch  uint32
	Standards     DvBtStd
	Flags         uint32
	PictureAspect Fract
	Cea861Vic     uint8
	HdmiVic       uint8
	Reserved      [46]uint8
}

// TODO: __attribute__ ((packed));
// figure out what this means in videodev2.h?

func (bt *BtTimings) BlankingWidth() uint32 {
	return bt.Hfrontporch + bt.Hsync + bt.Hbackporch
}

func (bt *BtTimings) FrameWidth() uint32 {
	return bt.Width + bt.BlankingWidth()
}

func (bt *BtTimings) BlankingHeight() uint32 {
	return bt.Vfrontporch + bt.Vsync + bt.Vbackportch
}

func (bt *BtTimings) FrameHeight() uint32 {
	return bt.Height + bt.BlankingHeight()
}

type Dv uint32

const (
	Dv_Progressive Dv = 0
	Dv_Interlaced  Dv = 0
)

type DvPol uint32

const (
	DvPol_VsyncPos DvPol = 0x00000001
	DvPol_HsyncPos DvPol = 0x00000002
)

type DvBtStd uint32

const (
	DvBtStd_Cea861 DvBtStd = 1 << 0
	DvBtStd_Dmt    DvBtStd = 1 << 1
	DvBtStd_Cvt    DvBtStd = 1 << 2
	DvBtStd_Gtf    DvBtStd = 1 << 3
	DvBtStd_Sdi    DvBtStd = 1 << 4
)

type DvFl uint32

const (
	DvFl_ReducedBlanking     DvFl = 1 << 0
	DvFl_CanReduceFps        DvFl = 1 << 1
	DvFl_ReducedFps          DvFl = 1 << 2
	DvFl_HalfLine            DvFl = 1 << 3
	DvFl_IsCeVideo           DvFl = 1 << 4
	DvFl_FirstFieldExtraLine DvFl = 1 << 5
	DvFl_HasPictureAspect    DvFl = 1 << 6
	DvFl_HasCea861_Vic       DvFl = 1 << 7
	DvFl_HasHdmiVic          DvFl = 1 << 8
	DvFl_CanDetectReducedFps DvFl = 1 << 9
)

type DvTimings struct {
	Type DvType
	Bt   [128]byte
	//	Bt union {
	//		BtTimings
	//		[32]uint32
	//	}
}

// TODO: __attribute__((packed)) ? (probably unneeded)

type DvType uint32

const (
	Dv_Bt6561120 DvType = 0
)

type EnumDvTimings struct {
	Index    uint32
	Pad      uint32
	Reserved [2]uint32
	Timings  DvTimings
}

type BtTimingsCap struct {
	MinWidth      uint32
	MaxWidth      uint32
	MinHeight     uint32
	MaxHeight     uint32
	MinPixelclock uint64
	MaxPixelclock uint64
	Standards     uint32
	Capabilities  uint32
	Reserved      [16]uint32
}

type DvBtCap uint32

const (
	DvBtCap_Interlaced      DvBtCap = 1 << 0
	DvBtCap_Progressive     DvBtCap = 1 << 1
	DvBtCap_ReducedBlanking DvBtCap = 1 << 2
	DvBtCap_Custom          DvBtCap = 1 << 3
)

type DvTimingsCap struct {
	Type        uint32
	Pad         uint32
	Reserved    [2]uint32
	BtOrRawData [128]byte
	//  BtOrRawData union {
	//		BtTimingsCap
	//		uint32
	//	}
}

type Input struct {
	Index        uint32
	Name         [32]uint8
	Type         InputType
	Audioset     uint32
	Tuner        uint32
	Std          StdId
	Status       InSt
	Capabilities InCap
	Reserved     [3]uint32
}

type InputType uint32

const (
	InputType_Tuner  InputType = 1
	InputType_Camera InputType = 2
	InputType_Touch  InputType = 3
)

type InSt uint32

const (
	InSt_NoPower     InSt = 0x00000001
	InSt_NoSignal    InSt = 0x00000002
	InSt_NoColor     InSt = 0x00000004
	InSt_Hflip       InSt = 0x00000010
	InSt_Vflip       InSt = 0x00000020
	InSt_NoHLock     InSt = 0x00000100
	InSt_ColorKill   InSt = 0x00000200
	InSt_NoVLock     InSt = 0x00000400
	InSt_NoStdLock   InSt = 0x00000800
	InSt_NoSync      InSt = 0x00010000
	InSt_NoEqu       InSt = 0x00020000
	InSt_NoCarrier   InSt = 0x00040000
	InSt_Macrovision InSt = 0x01000000
	InSt_NoAccess    InSt = 0x02000000
	InSt_Vtr         InSt = 0x04000000
)

type InCap uint32

const (
	InCap_DvTimings     InCap = 0x00000002
	InCap_CustomTimings InCap = InCap_DvTimings
	InCap_Std           InCap = 0x00000004
	InCap_NativeSize    InCap = 0x00000008
)

type Output struct {
	Index        uint32
	Name         [32]uint8
	Type         OutputType
	Audioset     uint32
	Modulator    uint32
	Std          StdId
	Capabilities OutCap
	Reserved     [3]uint32
}

type OutputType uint32

const (
	OutputType_Modulator        OutputType = 1
	OutputType_Analog           OutputType = 2
	OutputType_AnalogVgaOverlay            = 3
)

type OutCap uint32

const (
	OutCap_DvTimings     OutCap = 0x00000002
	OutCap_CustomTimings OutCap = OutCap_DvTimings
	OutCap_Std           OutCap = 0x00000004
	OutCap_NativeSize    OutCap = 0x00000008
)

type Control struct {
	Id    uint32
	Value uint32
}

type ExtControl struct {
	Id        uint32
	Size      uint32
	Reserved2 [1]uint32
	Value     [8]byte
	//	Value nion {
	//  	int32
	//  	int64
	//  	uintptr(char [?])
	//  	uintptr(uint8)
	//  	uintptr(uint16)
	//  	uintptr(uint32)
	//  	uintptr(void)
	//	}
}

type ExtControls struct {
	CtrlClassOrWhich uint32 // uint32 union type
	Count            uint32
	ErrorIds         uint32
	RequestFd        int32
	Reserved         [1]uint32
	Controls         uintptr // *ExtControl
}

// TODO: type this?
const (
	CtrlIdMask          = 0x0fffffff
	CtrlMaxDims         = 4
	CtrlWhichCurVal     = 0
	CtrlWhichDefVal     = 0x0f000000
	CtrlWhichRequestVal = 0x0f010000
)

// TODO: find usage and figure out if correct type
func CtrlId2Class(id uint32) uint32 {
	return id & 0x0fff0000
}

func CtrlId2Which(id uint32) uint32 {
	return id & 0x0fff0000
}

func CtrlDriverPriv(id uint32) bool {
	return (id & 0xffff) >= 0x1000
}

type CtrlType uint32

const (
	CtrlType_Integer     CtrlType = 1
	CtrlType_Boolean     CtrlType = 2
	CtrlType_Menu        CtrlType = 3
	CtrlType_Button      CtrlType = 4
	CtrlType_Integer64   CtrlType = 5
	CtrlType_CtrlClass   CtrlType = 6
	CtrlType_String      CtrlType = 7
	CtrlType_Bitmask     CtrlType = 8
	CtrlType_IntegerMenu CtrlType = 9

	// Compound types
	CtrlType_CompoundTypes CtrlType = 0x0100
	CtrlType_U8            CtrlType = 0x0100
	CtrlType_U16           CtrlType = 0x0101
	CtrlType_U32           CtrlType = 0x0102
)

type QueryCtrl struct {
	Id           uint32
	Type         CtrlType
	Name         [32]uint8
	Minimum      int32
	Maximum      int32
	Step         int32
	DefaultValue int32
	Flags        uint32
	Reserved     [2]uint32
}

type QueryExtCtrl struct {
	Id           uint32
	Type         uint32
	Name         [32]uint8
	Minimum      int64
	Maximum      int64
	Step         uint64
	DefaultValue int64
	Flags        uint32
	ElemSize     uint32
	Elems        uint32
	NrOfDims     uint32
	Dims         [CtrlMaxDims]uint32
	Reserved     [32]uint32
}

type QueryMenu struct {
	Id          uint32
	Index       uint32
	NameOrValue [32]byte
	//	NameOrValue {
	//		[32]uint8
	//		int64
	//	}
	Reserved uint32
}

type CtrlFlag uint32

const (
	CtrlFlag_Disabled       CtrlFlag = 0x0001
	CtrlFlag_Grabbed        CtrlFlag = 0x0002
	CtrlFlag_ReadOnly       CtrlFlag = 0x0004
	CtrlFlag_Update         CtrlFlag = 0x0008
	CtrlFlag_Inactive       CtrlFlag = 0x0010
	CtrlFlag_Slider         CtrlFlag = 0x0020
	CtrlFlag_WriteOnly      CtrlFlag = 0x0040
	CtrlFlag_Volatile       CtrlFlag = 0x0080
	CtrlFlag_HasPayload     CtrlFlag = 0x0100
	CtrlFlag_ExecuteOnWrite CtrlFlag = 0x0200
	CtrlFlag_ModifyLayout   CtrlFlag = 0x0400
	CtrlFlag_NextCtrl       CtrlFlag = 0x80000000
	CtrlFlag_NextCompound   CtrlFlag = 0x40000000
)

// TODO: make these CtrlFlags?
const (
	CidMaxCtrls    = 1024
	CidPrivateBase = 0x08000000
)

type Tuner struct {
	Index      uint32
	Name       [32]uint8
	Type       TunerType
	Capability TunerCap
	RangeLow   uint32
	RangeHigh  uint32
	RxSubChans TunerSub
	AudMode    TunerMode
	Signal     int32
	Afc        int32
	Reserved   [4]uint32
}

type Modulator struct {
	Index      uint32
	Name       [32]uint8
	Capability TunerCap // TODO: check if this is right
	RangeLow   uint32
	RangeHigh  uint32
	RxSubChans TunerSub
	Type       TunerType
	Reserved   [3]uint32
}

type TunerCap uint32

const (
	TunerCap_Low           TunerCap = 0x0001
	TunerCap_Norm          TunerCap = 0x0002
	TunerCap_HwseekBounded TunerCap = 0x0004
	TunerCap_HwseekWrap    TunerCap = 0x0008
	TunerCap_Stereo        TunerCap = 0x0010
	TunerCap_Lang2         TunerCap = 0x0020
	TunerCap_Sap           TunerCap = 0x0020
	TunerCap_Lang1         TunerCap = 0x0040
	TunerCap_Rds           TunerCap = 0x0080
	TunerCap_RdsBlockIo    TunerCap = 0x0100
	TunerCap_RdsControls   TunerCap = 0x0200
	TunerCap_FreqBands     TunerCap = 0x0400
	TunerCap_HwseekProgLim TunerCap = 0x0800
	TunerCap_1Hz           TunerCap = 0x1000
)

type TunerSub uint32

const (
	TunerSub_Mono   TunerSub = 0x0001
	TunerSub_Stereo TunerSub = 0x0002
	TunerSub_Lang2  TunerSub = 0x0004
	TunerSub_Sap    TunerSub = 0x0004
	TunerSub_Lang1  TunerSub = 0x0008
	TunerSub_Rds    TunerSub = 0x0010
)

type TunerMode uint32

const (
	TunerMode_Mono       TunerMode = 0x0000
	TunerMode_Stereo     TunerMode = 0x0001
	TunerMode_Lang2      TunerMode = 0x0002
	TunerMode_Sap        TunerMode = 0x0002
	TunerMode_Lang1      TunerMode = 0x0003
	TunerMode_Lang1Lang2 TunerMode = 0x0004
)

type Frequency struct {
	Tuner     uint32
	Type      TunerType
	Frequency uint32
	Reserved  [8]uint32
}

// TODO: type BandModulation ?
const (
	BandModulation_Vsb = 1 << 1
	BandModulation_Fm  = 1 << 2
	BandModulation_Am  = 1 << 3
)

type FrequencyBand struct {
	Tuner      uint32
	Type       TunerType
	Index      uint32
	Capability uint32 // TODO: TunerCap?
	RangeLow   uint32
	RangeHigh  uint32
	Modulation uint32
	Reserved   [9]uint32
}

type HwFreqSeek struct {
	Tuner      uint32
	Type       TunerType
	SeekUpward uint32
	WrapAround uint32
	Spacing    uint32
	RangeLow   uint32
	RangeHigh  uint32
	Reserved   [5]uint32
}

type RdsData struct {
	Lsb   uint8
	Msb   uint8
	Block RdsBlock
}

type RdsBlock uint8

const (
	RdsBlock_Msk       RdsBlock = 0x7
	RdsBlock_A         RdsBlock = 0
	RdsBlock_B         RdsBlock = 1
	RdsBlock_C         RdsBlock = 2
	RdsBlock_D         RdsBlock = 3
	RdsBlock_CAlt      RdsBlock = 4
	RdsBlock_Invalid   RdsBlock = 7
	RdsBlock_Corrected RdsBlock = 0x40
	RdsBlock_Error     RdsBlock = 0x80
)

type Audio struct {
	Index      uint32
	Name       [32]uint8
	Capability AudCap
	Mode       AudMode
	Reserved   [2]uint32
}

type AudCap uint32

const (
	AudCap_Stereo AudCap = 0x00001
	AudCap_Avl    AudCap = 0x00002
)

type AudMode uint32

const (
	AudMode_Avl AudMode = 0x00001
)

type AudioOut struct {
	Index      uint32
	Name       [32]uint8
	Capability AudCap  // TODO: check
	Mode       AudMode // TODO: check
	Reserved   [2]uint32
}

// TODO: type EncIdxFrame ?
const (
	EncIdxFrame_I    = 0
	EncIdxFrame_P    = 1
	EncIdxFrame_B    = 2
	EncIdxFrame_Mask = 0xf
)

type EncIdxEntry struct {
	Offset   uint64
	Pts      uint64
	Length   uint32
	Flags    uint32
	Reserved [2]uint32
}

const EncIdxEntries = 64

type EncIdx struct {
	Entries    uint32
	EntriesCap uint32
	Reserved   [4]uint32
	Entry      [EncIdxEntries]EncIdxEntry
}

type EncCmd uint32

const (
	EncCmd_Start        EncCmd = 0
	EncCmd_Stop         EncCmd = 1
	EncCmd_Pause        EncCmd = 2
	EncCmd_Resume       EncCmd = 3
	EncCmd_StopAtGopEnd EncCmd = 1 << 0
)

type EncoderCmd struct {
	Cmd     EncCmd
	Flags   uint32
	RawData [8]uint32 // TODO: check if this fits union from videodev2.h?
}

type DecCmd uint32

const (
	DecCmd_Start           DecCmd = 0
	DecCmd_Stop            DecCmd = 1
	DecCmd_Pause           DecCmd = 2
	DecCmd_Resume          DecCmd = 3
	DecCmd_StartMuteAudio  DecCmd = 1 << 0
	DecCmd_PauseToBlack    DecCmd = 1 << 0
	DecCmd_StopToBlack     DecCmd = 1 << 0
	DecCmd_StopImmediately DecCmd = 1 << 1

	DecStartFmtNone DecCmd = 0
	DecStartFmtGop  DecCmd = 1
)

type DecoderCmd struct {
	Cmd          uint32
	Flags        uint32
	StopStartRaw [64]byte
	//	StopStartRaw union {
	//		uint64 (pts)
	//		struct {
	//			int32 (speed)
	//			uint32 (format)
	//		}
	//		[16]uint32 (data)
	//	}
}

type VbiFormat struct {
	SamplingRate   uint32 // Hz
	Offset         uint32
	SamplesPerLine uint32
	SampleFormat   PixFmt
	Start          [2]int32
	Count          [2]uint32
	Flags          Vbi
	Reserved       [2]uint32
}

type Vbi uint32

const (
	Vbi_Unsync        Vbi = 1 << 0
	Vbi_Interlaced    Vbi = 1 << 1
	Vbi_Itu525F1Start Vbi = 1
	Vbi_Itu525F2Start Vbi = 264
	Vbi_Itu625F1Start Vbi = 1
	Vbi_Itu625F2Start Vbi = 314
)

type SlicedVbiFormat struct {
	ServiceSet   Sliced
	ServiceLines [2][24]uint16
	IoSize       uint32
	Reserved     [2]uint32
}

type Sliced uint16

const (
	Sliced_TeletextB  Sliced = 0x0001
	Sliced_Vps        Sliced = 0x0400
	Sliced_Caption525 Sliced = 0x1000
	Sliced_Wss625     Sliced = 0x4000
	Sliced_Vbi525     Sliced = Sliced_Caption525
	Sliced_Vbi625     Sliced = Sliced_TeletextB | Sliced_Vps | Sliced_Wss625
)

type SlicedVbiCap struct {
	ServiceSet   Sliced
	ServiceLines [2][24]uint16
	Type         BufType
	Reserved     [3]uint32
}

type SlicedVbiData struct {
	Id       uint32
	Field    uint32
	Line     uint32
	Reserved uint32
	Data     [48]uint8
}

type MpegVbiIvtv uint8

const (
	MpegVbiIvtv_TeletextB  MpegVbiIvtv = 1
	MpegVbiIvtv_Caption525 MpegVbiIvtv = 4
	MpegVbiIvtv_Wss625     MpegVbiIvtv = 5
	MpegVbiIvtv_Vps        MpegVbiIvtv = 7
)

// Using var because arrays are mutable
// TODO: check if this char* hack works
var (
	MjpegVbiIvtvMagic0 [4]uint8 = [4]uint8{'i', 't', 'v', '0'}
	MjpegVbiIvtvMagic1 [4]uint8 = [4]uint8{'I', 'T', 'V', '0'}
)

type MpegVbiItv0Line struct {
	Id   MpegVbiIvtv
	Data [42]uint8
}

type MpegVbiItv0 struct {
	Linemask [2]uint32
	Line     [35]MpegVbiItv0Line
}

// ? wtf
type MpegVbiITV0 struct {
	Line [36]MpegVbiItv0Line
}

type MpegVbiFmtIvtv struct {
	Magic [4]uint8

	// TODO: find a smarter way to create union
	Itv0 [1548]byte
	//	Itv0 union {
	//		MpegVbiItv0
	//		MpegVbiITV0
	//	}
}

type PlanePixFormat struct {
	SizeImage    uint32
	BytesPerLine uint32
	Reserved     [6]uint16
}

type PixFormatMplane struct {
	Width       uint32
	Height      uint32
	PixelFormat uint32
	Field       uint32
	Colorspace  uint32

	PlaneFmt         [VideoMaxPlanes]PlanePixFormat
	NumPlanes        uint8
	Flags            uint8
	YcbcrEncOrHsvEnc uint8
	//	YcbcrEncOrHsvEnc union {
	//		YcbcrEncoding (!!!)
	//		HsvEncoding (!!!)
	//	}
	Quantization uint8 // Quantization (!!!)
	XferFunc     uint8 // XferFunc (!!!)
	Reserved     [7]uint8
}

type SdrFormat struct {
	PixelFormat PixFmt
	BufferSize  uint32
}

type MetaFormat struct {
	DataFormat PixFmt
	BufferSize uint32
}

type Format struct {
	Type uint32
	Fmt  [200]uint8
	//	Fmt union {
	//		PixFormat
	//		PixFormatMplane
	//		Window
	//		VbiFormat
	//		SlicedVbiFormat
	//		SdrFormat
	//		MetaFormat
	//		[200]byte
	//	}
}

type StreamParm struct {
	Type BufType
	Parm [200]byte
	//	Parm union {
	//		CaptureParm
	//		OutputParm
	//		[200]byte
	//	}
}

// TODO: type EventType ?
type EventType uint32

const (
	Event_All          EventType = 0
	Event_Vsync        EventType = 1
	Event_Eos          EventType = 2
	Event_Ctrl         EventType = 3
	Event_FrameSync    EventType = 4
	Event_SourceChange EventType = 5
	Event_MotionDet    EventType = 6
	Event_PrivateStart EventType = 0x08000000
)

type EventVsync struct {
	Field uint8 // Field (!!!)
}

// TODO: type EventCtrl ?
const (
	EventCtrlChValue = 1 << 0
	EventCtrlChFlags = 1 << 1
	EventCtrlChRange = 1 << 2
)

type EventCtrl struct {
	Changes uint32
	Type    uint32
	Value   int64
	//  Value union {
	//		int32
	//		int64
	//  }
	Flags        uint32
	Minimum      int32
	Maximum      int32
	Step         int32
	DefaultValue int32
}

type EventFrameSync struct {
	FrameSequence uint32
}

const EventSrcChResolution = 1 << 0

type EventSrcChange struct {
	changes uint32
}

const EventMdFlHaveFrameSeq = 1 << 0

type EventMotionDet struct {
	Flags         uint32
	FrameSequence uint32
	RegionMask    uint32
}

type Event struct {
	Type EventType
	U    [64]byte
	//	U union {
	//		EventVsync
	//		EventCtrl
	//		EventFrameSync
	//		EventSrcChange
	//		EventMotionDet
	//		[64]uint8
	//	}
	Pending   uint32
	Sequence  uint32
	Timestamp syscall.Timespec // TODO: check if same size as C timespec
	Id        uint32
	Reserved  [8]uint32
}

// TODO: type this?
const (
	EventSubFlSendInitial   = 1 << 0
	EventSubFlAllowFeedback = 1 << 1
)

type EventSubscription struct {
	Type     uint32
	Id       uint32
	Flags    uint32
	Reserved [5]uint32
}

//
// Note: "Advanced Debugging" section of videodev2.h not included
//

type CreateBuffers struct {
	Index        uint32
	Count        uint32
	Memory       uint32
	Format       Format
	Capabilities uint32
	Reserved     [7]uint32
}

// IOCTLS:

type Vidioc = uintptr // Video ioctl code ?

// var because golang doesn't have compile-time functions
var (
	Vidioc_QueryCap Vidioc = _IOR('V', 0, unsafe.Sizeof(Capability{}))
	Vidioc_EnumFmt  Vidioc = _IOWR('V', 2, unsafe.Sizeof(FmtDesc{}))

	VidiocGFmt      Vidioc = _IOWR('V', 4, unsafe.Sizeof(Format{}))
	VidiocSFmt      Vidioc = _IOWR('V', 5, unsafe.Sizeof(Format{}))
	VidiocReqbufs   Vidioc = _IOWR('V', 8, unsafe.Sizeof(RequestBuffers{}))
	VidiocQuerybuf  Vidioc = _IOWR('V', 9, unsafe.Sizeof(Buffer{}))
	VidiocGFbuf     Vidioc = _IOR('V', 10, unsafe.Sizeof(FrameBuffer{}))
	VidiocSFbuf     Vidioc = _IOW('V', 11, unsafe.Sizeof(FrameBuffer{}))
	VidiocOverlay   Vidioc = _IOW('V', 14, unsafe.Sizeof(int(0)))
	VidiocQbuf      Vidioc = _IOWR('V', 15, unsafe.Sizeof(Buffer{}))
	VidiocExpbuf    Vidioc = _IOWR('V', 16, unsafe.Sizeof(ExportBuffer{}))
	VidiocDqbuf     Vidioc = _IOWR('V', 17, unsafe.Sizeof(Buffer{}))
	VidiocStreamon  Vidioc = _IOW('V', 18, unsafe.Sizeof(int(0)))
	VidiocStreamoff Vidioc = _IOW('V', 19, unsafe.Sizeof(int(0)))
	VidiocGParm     Vidioc = _IOWR('V', 21, unsafe.Sizeof(StreamParm{}))
	VidiocSParm     Vidioc = _IOWR('V', 22, unsafe.Sizeof(StreamParm{}))
	VidiocGStd      Vidioc = _IOR('V', 23, unsafe.Sizeof(StdId(0)))
	VidiocSStd      Vidioc = _IOW('V', 24, unsafe.Sizeof(StdId(0)))
	VidiocEnumstd   Vidioc = _IOWR('V', 25, unsafe.Sizeof(Standard{}))
	VidiocEnuminput Vidioc = _IOWR('V', 26, unsafe.Sizeof(Input{}))
	VidiocGCtrl     Vidioc = _IOWR('V', 27, unsafe.Sizeof(Control{}))
	VidiocSCtrl     Vidioc = _IOWR('V', 28, unsafe.Sizeof(Control{}))
	VidiocGTuner    Vidioc = _IOWR('V', 29, unsafe.Sizeof(Tuner{}))
	VidiocSTuner    Vidioc = _IOW('V', 30, unsafe.Sizeof(Tuner{}))
	VidiocGAudio    Vidioc = _IOR('V', 33, unsafe.Sizeof(Audio{}))
	VidiocSAudio    Vidioc = _IOW('V', 34, unsafe.Sizeof(Audio{}))
	VidiocQueryctrl Vidioc = _IOWR('V', 36, unsafe.Sizeof(QueryCtrl{}))
	VidiocQuerymenu Vidioc = _IOWR('V', 37, unsafe.Sizeof(QueryMenu{}))
	VidiocGInput    Vidioc = _IOR('V', 38, unsafe.Sizeof(int(0)))
	VidiocSInput    Vidioc = _IOWR('V', 39, unsafe.Sizeof(int(0)))
	// TODO: get struct v4l2_edid from v4l2-common.h
	//VidiocGEdid              Vidioc = _IOWR('V', 40, unsafe.Sizeof(edid{}))
	//VidiocSEdid              Vidioc = _IOWR('V', 41, unsafe.Sizeof(edid{}))
	VidiocGOutput            Vidioc = _IOR('V', 46, unsafe.Sizeof(int(0)))
	VidiocSOutput            Vidioc = _IOWR('V', 47, unsafe.Sizeof(int(0)))
	VidiocEnumoutput         Vidioc = _IOWR('V', 48, unsafe.Sizeof(Output{}))
	VidiocGAudout            Vidioc = _IOR('V', 49, unsafe.Sizeof(AudioOut{}))
	VidiocSAudout            Vidioc = _IOW('V', 50, unsafe.Sizeof(AudioOut{}))
	VidiocGModulator         Vidioc = _IOWR('V', 54, unsafe.Sizeof(Modulator{}))
	VidiocSModulator         Vidioc = _IOW('V', 55, unsafe.Sizeof(Modulator{}))
	VidiocGFrequency         Vidioc = _IOWR('V', 56, unsafe.Sizeof(Frequency{}))
	VidiocSFrequency         Vidioc = _IOW('V', 57, unsafe.Sizeof(Frequency{}))
	VidiocCropcap            Vidioc = _IOWR('V', 58, unsafe.Sizeof(CropCap{}))
	VidiocGCrop              Vidioc = _IOWR('V', 59, unsafe.Sizeof(Crop{}))
	VidiocSCrop              Vidioc = _IOW('V', 60, unsafe.Sizeof(Crop{}))
	VidiocGJpegcomp          Vidioc = _IOR('V', 61, unsafe.Sizeof(JpegCompression{}))
	VidiocSJpegcomp          Vidioc = _IOW('V', 62, unsafe.Sizeof(JpegCompression{}))
	VidiocQuerystd           Vidioc = _IOR('V', 63, unsafe.Sizeof(StdId(0)))
	VidiocTryFmt             Vidioc = _IOWR('V', 64, unsafe.Sizeof(Format{}))
	VidiocEnumaudio          Vidioc = _IOWR('V', 65, unsafe.Sizeof(Audio{}))
	VidiocEnumaudout         Vidioc = _IOWR('V', 66, unsafe.Sizeof(AudioOut{}))
	VidiocGPriority          Vidioc = _IOR('V', 67, unsafe.Sizeof(Priority(0)))
	VidiocSPriority          Vidioc = _IOW('V', 68, unsafe.Sizeof(Priority(0)))
	VidiocGSlicedVbiCap      Vidioc = _IOWR('V', 69, unsafe.Sizeof(SlicedVbiCap{}))
	VidiocLogStatus          Vidioc = _IO('V', 70)
	VidiocGExtCtrls          Vidioc = _IOWR('V', 71, unsafe.Sizeof(ExtControls{}))
	VidiocSExtCtrls          Vidioc = _IOWR('V', 72, unsafe.Sizeof(ExtControls{}))
	VidiocTryExtCtrls        Vidioc = _IOWR('V', 73, unsafe.Sizeof(ExtControls{}))
	VidiocEnumFramesizes     Vidioc = _IOWR('V', 74, unsafe.Sizeof(FrmSizeEnum{}))
	VidiocEnumFrameintervals Vidioc = _IOWR('V', 75, unsafe.Sizeof(FrmIvalEnum{}))
	VidiocGEncIndex          Vidioc = _IOR('V', 76, unsafe.Sizeof(EncIdx{}))
	VidiocEncoderCmd         Vidioc = _IOWR('V', 77, unsafe.Sizeof(EncoderCmd{}))
	VidiocTryEncoderCmd      Vidioc = _IOWR('V', 78, unsafe.Sizeof(EncoderCmd{}))

//)
)

// ioctl.h (but a little more golangy)

func _IO(t rune, nr int) uintptr {
	return _IOC(_IOC_NONE, t, nr, 0)
}

func _IOR(t rune, nr int, size uintptr) uintptr {
	return _IOC(_IOC_READ, t, nr, size)
}

func _IOW(t rune, nr int, size uintptr) uintptr {
	return _IOC(_IOC_WRITE, t, nr, size)
}

func _IOWR(t rune, nr int, size uintptr) uintptr {
	return _IOC(_IOC_READ|_IOC_WRITE, t, nr, size)
}

func _IOC(dir int, t rune, nr int, size uintptr) uintptr {
	return ((uintptr(dir) << _IOC_DIRSHIFT) |
		(uintptr(t) << _IOC_TYPESHIFT) |
		(uintptr(nr) << _IOC_NRSHIFT) |
		(uintptr(size) << _IOC_SIZESHIFT))
}

const (
	_IOC_NONE  = 0
	_IOC_WRITE = 1
	_IOC_READ  = 2

	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_DIRBITS  = 2

	_IOC_NRSHIFT   = 0
	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS
)
