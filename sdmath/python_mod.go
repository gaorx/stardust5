package sdmath

func Python2Mod(x, mod int64) int64 {
	r := x % mod
	if (mod < 0 && r > 0) || (mod > 0 && r < 0) {
		r += mod
	}
	return r
}
