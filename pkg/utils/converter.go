package utils

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
)

func NumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0.0
	}
	value := float64(n.Int.Int64()) * math.Pow10(int(n.Exp))
	return value
}

func Float64ToNumeric(f float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(fmt.Sprintf("%f", f))
	return n, err
}
