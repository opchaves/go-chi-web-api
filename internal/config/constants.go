package config

type ctxKey int

const (
	CtxClaims ctxKey = iota
	CtxRefreshToken
	CtxVersion
)
