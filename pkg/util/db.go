package util

import (
	"log/slog"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ZERO = &pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true}
	NOW  = pgtype.Timestamp{Time: time.Now(), Valid: true}
)

func ParseNumeric(s string) *pgtype.Numeric {
	n := &pgtype.Numeric{}
	err := n.Scan(s)
	if err != nil {
		slog.Default().Warn("error parsing numeric: %v", err)
		return nil
	}
	return n
}

func ParseTimestamp(s string) *pgtype.Timestamp {
	dt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		slog.Default().Warn("error parsing date: %v", err)
		return nil
	}

	t := &pgtype.Timestamp{}
	if err := t.Scan(dt); err != nil {
		slog.Default().Warn("error parsing timestamp: %v", err)
		return nil
	}

	return t
}
