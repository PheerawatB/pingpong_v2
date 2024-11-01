package server

import "golang.org/x/exp/rand"

// Random for return 70-95% power
func BallPowerTo(power uint, name string) uint {
	lowPowerPercentage := 70 + rand.Intn(26)
	lowPower := uint(float64(power) * float64(lowPowerPercentage) / 100.0)
	return lowPower
}
