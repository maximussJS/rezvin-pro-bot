package utils

import (
	"context"
)

func GetContextWithClientIp(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, "clientIp", ip)
}

func GetClientIpFromContext(ctx context.Context) string {
	return ctx.Value("clientIp").(string)
}
