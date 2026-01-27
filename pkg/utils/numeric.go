package utils

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func Int64ToNumeric(v int64) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   big.NewInt(v),
		Exp:   0,
		Valid: true,
	}
}

func NumericToInt64(n pgtype.Numeric) int64 {
	// Use Float64Value to handle exponent scaling (e.g. 400 * 10^-2 = 4)
	f, _ := n.Float64Value()
	return int64(f.Float64)
}

func Int64PtrToNumeric(v *int64) pgtype.Numeric {
	if v == nil {
		return pgtype.Numeric{Valid: false}
	}
	return Int64ToNumeric(*v)
}

func NumericToInt64Ptr(n pgtype.Numeric) *int64 {
	if !n.Valid {
		return nil
	}
	val := NumericToInt64(n)
	return &val
}
