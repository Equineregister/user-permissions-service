package contextkey

import "context"

type CtxKey int

const (
	CtxKeyUserID CtxKey = iota
	CtxKeyTenantID
)

func TenantID(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(CtxKeyTenantID).(string)
	if ok {
		return tenantID, true
	}
	return "", false
}

func UserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(CtxKeyUserID).(string)
	if ok {
		return userID, true
	}
	return "", false
}
