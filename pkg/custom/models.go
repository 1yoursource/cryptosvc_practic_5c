package custom

type Data string

func (s Data) String() string {
	return string(s)
}

type Shift int8

func (s Shift) Int8() int8       { return int8(s) }
func (s Shift) IsZero() bool     { return s == 0 }
func (s Shift) IsPositive() bool { return s > 0 }
func (s Shift) IsNegative() bool { return s < 0 }

type Position int

func (p Position) Int() int { return int(p) }

const (
	defaultIndex = 1

	space      = "_"
	dash       = "-"
	apostrophe = "'"
	dot        = "."
	comma      = ","
)
