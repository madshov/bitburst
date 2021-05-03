package http

import (
	"encoding/json"
	"log"
	"net/http"

	service "github.com/madshov/bitburst/app/service"
)

type reqBody struct {
	ObjIDs []int `json:"object_ids"`
}

type objectHandler struct {
	logger *log.Logger
	objSvc service.ObjectService
}

func NewObjectHandler(logger *log.Logger, svc service.ObjectService) http.Handler {
	h := &objectHandler{
		logger: logger,
		objSvc: svc,
	}

	return h
}

func (o *objectHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodPost:
		if r.Body == nil {
			http.Error(rw, "empty body", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		var req reqBody
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		// starting in a new routine, to avoid timeout from caller
		go o.objSvc.StoreObjects(ctx, req.ObjIDs)
		go o.objSvc.DeleteObjects(ctx)

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
		return

	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}
