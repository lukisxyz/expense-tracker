package record

import (
	"encoding/json"
	"net/http"

	"github.com/flukis/expt/service/utils/response"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	r.Get("/", listItemHandler)
	r.Post("/", createItemHandler)
	return r
}

func listItemHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	queryParams := r.URL.Query()
	paramValue := queryParams.Get("id")
	id, err := ulid.Parse(paramValue)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ctx := r.Context()
	resp, err := listRecords(ctx, id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

type createItemBodyRequest struct {
	CategoryId ulid.ULID `json:"category_id"`
	BookId     ulid.ULID `json:"book_id"`
	Note       string    `json:"note"`
	Amount     float64   `json:"amount"`
	IsExpense  bool      `json:"is_expense"`
}

func createItemHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var p createItemBodyRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	resp, err := saveRecord(ctx, p.BookId, p.CategoryId, p.Note, p.Amount, p.IsExpense)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
