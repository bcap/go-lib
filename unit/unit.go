package unit

//
// 10 based units
//

const One = 1
const Thousand = 1000 * One
const Million = 1000 * Thousand
const Billion = 1000 * Million
const Trillion = 1000 * Billion
const Quadrillion = 1000 * Trillion

const OneF = float64(One)
const ThousandF = float64(Thousand)
const MillionF = float64(Million)
const BillionF = float64(Billion)
const TrillionF = float64(Trillion)
const QuadrillionF = float64(Quadrillion)

//
// 2 based units
//

const B int64 = 1
const KiB int64 = 1024 * B
const MiB int64 = 1024 * KiB
const GiB int64 = 1024 * MiB
const TiB int64 = 1024 * GiB
const PiB int64 = 1024 * TiB

const Bf = float64(B)
const KiBf = float64(KiB)
const MiBf = float64(MiB)
const GiBf = float64(GiB)
const TiBf = float64(TiB)
const PiBf = float64(PiB)

const Byte = B
const KibiByte = KiB
const MebiByte = MiB
const GibiByte = GiB
const TebiByte = TiB

const ByteF = Bf
const KibiByteF = KiBf
const MebiByteF = MiBf
const GibiByteF = GiBf
const TebiByteF = TiBf
