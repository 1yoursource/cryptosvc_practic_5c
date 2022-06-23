package custom

import "context"

type LikeCesarCrypt interface {
	Crypt(ctx context.Context, input Data, s Shift) (Data, Position)
}

// Custom - crypto algorithm like Cesar, but with shifting each symbol on diff position
type Custom struct {
	alphabetSize int8
	Alphabet     []string
}

func newCustom(alphabet []string) *Custom {
	return &Custom{
		Alphabet:     alphabet,
		alphabetSize: int8(len(alphabet)),
	}
}

func (c *Custom) Crypt(ctx context.Context, input Data, s Shift) (Data, Position) {
	result, position := c.shift(input, s)

	return Data(result), Position(position)
}

func (c *Custom) shift(symbol Data, s Shift) (string,int) {
	oldPosition := c.getPosition(symbol)

	shiftBy := int(c.getCorrectShift(s))

	newPosition := oldPosition + shiftBy

	newPosition = int(c.getShiftPositive(Shift(newPosition)))

	return c.Alphabet[newPosition-defaultIndex], shiftBy
}

func (c *Custom) getPosition(symbol Data) int {
	for index, v := range c.Alphabet {
		if Data(v) == symbol {
			return index + defaultIndex
		}
	}

	return defaultIndex
}

func (c *Custom) getCorrectShift(s Shift) Shift {
	shift := s

	if s.IsPositive() {
		shift = c.getShiftPositive(s)
	} else if s.IsNegative() {
		// get the positive shift for simplify position calculation
		shift = Shift(c.alphabetSize) + c.getShiftNegative(s)
	}

	// if shift is zero
	return shift
}

func (c *Custom) getShiftPositive(s Shift) Shift {
	if s <= Shift(c.alphabetSize) {
		return s
	}

	for s > Shift(c.alphabetSize) {
		s -= Shift(c.alphabetSize)
	}

	return s
}

func (c *Custom) getShiftNegative(s Shift) Shift {
	if s >= Shift(-c.alphabetSize) {
		return s
	}

	for s < Shift(-c.alphabetSize) {
		s += Shift(c.alphabetSize)
	}

	return s
}
