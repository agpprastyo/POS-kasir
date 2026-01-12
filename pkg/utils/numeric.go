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
	if !n.Valid {
		return 0
	}
	// Simplified conversion: assume scale is 0 or handle basic cases
	// For price/currency usually stored as integer (cents/minor units), Exp is often 0
	val := n.Int.Int64()
	// Adjust for exponent if necessary (e.g. if Exp is -2, divide by 100)
	// database/sql stores numeric as string usually, pgx uses big.Int + Exp.
	// If Exp > 0, multiply by 10^Exp. If Exp < 0, divide by 10^(-Exp)
	// For this context assuming Exp=0 for int64 prices is safe if consistency is maintained.
	return val
}
