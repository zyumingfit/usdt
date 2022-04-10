package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/ebfe/keccak"
)

// Keccak hash.
func HashMethodName(method string) string {
	h := keccak.New256()
	h.Reset()
	h.Write([]byte(method))
	d := h.Sum(nil)
	return hex.EncodeToString(d[0:4])
}

// Get a certain number of 0 by num.
func GetStringZero(num int) string {
	s := "0000000000000000000000000000000000000000000000000000000000000000"
	return s[0:num]
}

// Decimal to hexadecimal
func DecHex(n int64) string {
	if n < 0 {
		fmt.Println("Decimal to hexadecimal error: the argument must be greater than zero.")
		return ""
	}
	if n == 0 {
		return "0"
	}
	hexTable := map[int64]int64{10: 65, 11: 66, 12: 67, 13: 68, 14: 69, 15: 70}
	s := ""
	for q := n; q > 0; q = q / 16 {
		m := q % 16
		if m > 9 && m < 16 {
			m = hexTable[m]
			s = fmt.Sprintf("%v%v", string(m), s)
			continue
		}
		s = fmt.Sprintf("%v%v", m, s)
	}
	return s
}
