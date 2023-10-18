package category

import (
	"encoding/json"
	"net/http"

	customMiddleware "github.com/flukis/expt/service/internals/middleware"
	"github.com/flukis/expt/service/utils/response"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	r.Post("/", createItemHandler)
	r.Get("/{id}", findItemByIdHandler)
	r.Patch("/{id}", updateItemHandler)
	r.Get("/", listItemHandler)
	return r
}

func listItemHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx := r.Context()
	res := ctx.Value(customMiddleware.ClaimJWTKey)
	id := res.(*customMiddleware.MapClaimResponse).Id
	resp, err := listCategories(ctx, id)
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
	resp, err := findCategoryById(ctx, id)
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
	Name string `json:"name"`
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
	res := ctx.Value(customMiddleware.ClaimJWTKey)
	ownerId := res.(*customMiddleware.MapClaimResponse).Id

	resp, err := saveCategory(ctx, ownerId, p.Name)
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
	Name string `json:"name"`
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

	ctx := r.Context()
	res := ctx.Value(customMiddleware.ClaimJWTKey)
	ownerId := res.(*customMiddleware.MapClaimResponse).Id

	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := updateBook(ctx, id, ownerId, p.Name)
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
