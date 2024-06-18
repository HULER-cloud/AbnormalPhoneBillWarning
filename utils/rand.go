package utils

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func RandInt(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	number := r.Intn(9 * int(math.Pow10(length-1)))
	return fmt.Sprintf("%d", number+int(math.Pow10(length-1)))
}
