package http

import (
	"CompanySystemsMonitoring/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	resultService service.Result
}

func NewHandler(resultService service.Result) *Handler {
	return &Handler{resultService: resultService}
}

func (h *Handler) Init() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api", h.handleConnection).Methods("GET")
	return router
}

func (h *Handler) handleConnection(w http.ResponseWriter, r *http.Request) {
	err := error(nil)
	resultJSON := []byte{}
	result := h.resultService.GetResultData(r.Context())
	if resultJSON, err = json.Marshal(result); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	}
}
