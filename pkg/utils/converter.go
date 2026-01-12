package utils

import (
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func NullableUUIDToPointer(nu pgtype.UUID) *uuid.UUID {
	if !nu.Valid {
		return nil
	}
	uid := uuid.UUID(nu.Bytes)
	return &uid
}

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

func Int32Ptr(i int) *int32 {
	v := int32(i)
	return &v
}

func StringPtr(s string) *string {
	return &s
}
