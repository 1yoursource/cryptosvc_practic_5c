package custom

type English struct {
	*Custom
}

var alphabetEn = [31]string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I",
	"J", "K", "L", "M", "N", "O", "P", "Q", "R",
	"S", "T", "U", "V", "W", "X", "Y", "Z",
	space, dash, apostrophe, dot, comma}

func (c *English) New() *Custom {
	alp := make([]string, len(alphabetEn))

	for _, v := range alphabetEn {
		alp = append(alp, v)
	}

	return newCustom(alp)
}
