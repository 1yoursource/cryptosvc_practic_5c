package custom

type Ukrainian struct {
	*Custom
}

var alphabetUk = [38]string{
	"А", "Б", "В", "Г", "Ґ", "Д", "Е", "Є",
	"Ж", "З", "И", "І", "Ї", "Й", "К", "Л",
	"М", "Н", "О", "П", "Р", "С", "Т", "У",
	"Ф", "Х", "Ц", "Ч", "Ш", "Щ", "Ь", "Ю", "Я",
	space, dash, apostrophe, dot, comma}

func (c *Ukrainian) New() *Ukrainian {
	alp := make([]string, len(alphabetUk))

	for i, v := range alphabetUk {
		alp[i] = v
	}

	return &Ukrainian{
		Custom: newCustom(alp),
	}
}
