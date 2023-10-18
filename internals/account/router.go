package account

import (
	"encoding/json"
	"net/http"

	"github.com/flukis/expt/service/utils/response"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	r.Get("/{id}", findItemByIdHandler)
	r.Post("/", createItemHandler)
	r.Patch("/{id}", updateItemHandler)
	r.Post("/login", loginHandler)

	return r
}

type loginBodyRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var p loginBodyRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	resp, err := loginAccount(ctx, p.Email, p.Password)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	token, err := generateJwtToken(
		ctx,
		resp.Id,
		resp.DefaultBookId,
		resp.Email,
	)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(token); err != nil {
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
	resp, err := findAccountById(ctx, id)
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
	Email    string `json:"email"`
	Password string `json:"password"`
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

	resp, err := saveAccount(ctx, p.Email, p.Password)
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
	Email    string `json:"email"`
	Password string `json:"password"`
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

	resp, err := updateAccount(ctx, id, p.Email, p.Password)
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
