package util

import (
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ZERO = &pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true}
	NOW  = pgtype.Timestamp{Time: time.Now(), Valid: true}
)
