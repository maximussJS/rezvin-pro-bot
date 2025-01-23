package router

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"rezvin-pro-bot/utils"
	utils_context "rezvin-pro-bot/utils/context"
	"time"
)

func (router *router) authMiddleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")

		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)

			utils.MustWriteBytesToResponseWriter(w, []byte(fmt.Sprintf(`{"message": "Authorization token is required"}`)))

			return
		}

		if token != router.config.AuthToken() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			utils.MustWriteBytesToResponseWriter(w, []byte(fmt.Sprintf(`{"message": "Invalid authorization token"}`)))

			return
		}

		n(w, r, ps)
	}
}

func (router *router) logMiddleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()

		n(w, r, ps)

		duration := time.Since(start)

		clientIP := utils_context.GetClientIpFromContext(r.Context())

		router.logger.Log(fmt.Sprintf(
			"[%s] %s %s %s %.0fm%.0fs%dms%dns %s",
			start.Format(time.RFC3339),
			r.Method,
			r.RequestURI,
			r.Proto,
			duration.Minutes(),
			duration.Seconds(),
			duration.Milliseconds(),
			duration.Microseconds(),
			clientIP,
		))
	}
}

func (router *router) clientIpMiddleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		clientIP := r.RemoteAddr
		if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			clientIP = ip
		}

		ctx := utils_context.GetContextWithClientIp(r.Context(), clientIP)

		n(w, r.WithContext(ctx), ps)
	}
}

func (router *router) withMiddlewares(fn wrappedFn) httprouter.Handle {
	return router.clientIpMiddleware(router.logMiddleware(router.authMiddleware(router.wrapResponse(fn))))
}
