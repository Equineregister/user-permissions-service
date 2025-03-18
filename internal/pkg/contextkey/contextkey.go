package contextkey

type CtxKey int

const (
	CtxKeyUserID CtxKey = iota
	CtxKeyTenantID
)
