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

	r.Get("/", listItemBookHandler)
	r.Get("/{id}", findItemByIdHandler)
	r.Get("/category", listItemCategoryHandler)
	r.Post("/", createItemHandler)
	r.Patch("/{id}", updateItemHandler)
	return r
}

func findItemByIdHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	resp, err := findRecordById(ctx, id)
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

func listItemCategoryHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	queryParams := r.URL.Query()
	idStr := queryParams.Get("category_id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ctx := r.Context()
	resp, err := listRecordByCategory(ctx, id)
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

func listItemBookHandler(
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
	resp, err := listRecordByBook(ctx, id)
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

type updateItemBodyRequest struct {
	CategoryId ulid.ULID `json:"category_id"`
	BookId     ulid.ULID `json:"book_id"`
	Note       string    `json:"note"`
	Amount     float64   `json:"amount"`
	IsExpense  bool      `json:"is_expense"`
}

func updateItemHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var p updateItemBodyRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	resp, err := updateRecord(ctx, id, p.BookId, p.CategoryId, p.Note, p.Amount, p.IsExpense)
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
