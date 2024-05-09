package float_range

func convertFormat(format string) byte {
	switch format {
	case "binary":
		return 'b'
	case "decimal":
		return 'e'
	case "decimal_capitalize":
		return 'E'
	case "large":
		return 'g'
	case "large_capitalize":
		return 'G'
	case "hex":
		return 'x'
	case "hex_capitalize":
		return 'X'
	default:
		return 'f'
	}
}
