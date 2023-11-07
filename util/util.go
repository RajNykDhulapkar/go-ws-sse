package util

func Round2Dec(val float64) float64 {
	return float64(int(val*100.0)) / 100.0
}
