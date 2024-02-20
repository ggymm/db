package str

func Serialize(str string) []byte {
	data := make([]byte, 4+len(str))

	// length
	l := len(str)
	data[0] = byte(l)
	data[1] = byte(l >> 8)
	data[2] = byte(l >> 16)
	data[3] = byte(l >> 24)

	// string
	copy(data[4:], str)
	return data
}

func Deserialize(data []byte) (string, int) {
	l := int(data[0]) |
		int(data[1])<<8 |
		int(data[2])<<16 |
		int(data[3])<<24
	return string(data[4 : 4+l]), l + 4
}
