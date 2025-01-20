package router

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
	"rezvin-pro-bot/types/responses"
	"rezvin-pro-bot/utils"
)

type wrappedFn func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) *responses.Response

func (router *router) wrapResponse(fn wrappedFn) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		resp := fn(w, r, p)

		if ch := r.Context().Err(); ch != nil {
			return
		}

		if resp.IsError() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.Status)

			utils.MustWriteBytesToResponseWriter(w, []byte(fmt.Sprintf(`{"message": "%s"}`, resp.Msg)))

			return
		}

		if resp.Data != nil {
			if reflect.ValueOf(resp.Data).Kind() == reflect.String {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(resp.Status)

				utils.MustWriteBytesToResponseWriter(w, []byte(resp.Data.(string)))

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.Status)

			utils.MustWriteJsonToResponseWriter(w, resp.Data)

			return
		}

		w.WriteHeader(resp.Status)
	}
}
