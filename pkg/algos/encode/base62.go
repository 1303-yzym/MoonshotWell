package encode

var Base62 = Encoding{
	charset: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
}

type Encoding struct {
	charset string
}

func (e Encoding) Encode(num uint32) string {
	if num == 0 {
		return "0"
	}

	result := ""

	for num > 0 {
		remainder := num % 62
		result = string(e.charset[remainder]) + result
		num = num / 62
	}

	return result
}
