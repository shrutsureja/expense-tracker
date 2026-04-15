package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"expense-tracker/internal/models"
	"expense-tracker/internal/service"

	"github.com/go-chi/chi/v5"
)

type FamilyHandler struct {
	familyService *service.FamilyService
}

func NewFamilyHandler(familyService *service.FamilyService) *FamilyHandler {
	return &FamilyHandler{familyService: familyService}
}

// Admin endpoints

type createFamilyRequest struct {
	Name             string `json:"name"              validate:"required,min=1,max=100"`
	OwnerUsername    string `json:"owner_username"    validate:"required,min=2,max=50,alphanum"`
	OwnerPassword    string `json:"owner_password"    validate:"required,len=6,numeric"`
	OwnerDisplayName string `json:"owner_display_name" validate:"required,min=1,max=100"`
}

func (h *FamilyHandler) CreateFamily(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())

	var req createFamilyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if !validateRequest(w, req) {
		return
	}

	family, owner, err := h.familyService.CreateFamily(
		req.Name, claims.UserID, req.OwnerUsername, req.OwnerPassword, req.OwnerDisplayName,
	)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"family": family,
		"owner":  owner,
	})
}

func (h *FamilyHandler) ListFamilies(w http.ResponseWriter, r *http.Request) {
	families, err := h.familyService.ListFamilies()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list families"})
		return
	}
	if families == nil {
		families = []models.FamilyWithOwner{}
	}
	writeJSON(w, http.StatusOK, families)
}

func (h *FamilyHandler) DeleteFamily(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid family id"})
		return
	}

	if err := h.familyService.DeleteFamily(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete family"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "family deleted"})
}

// Family owner/member endpoints

func (h *FamilyHandler) GetFamily(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	family, err := h.familyService.GetFamily(*claims.FamilyID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get family"})
		return
	}

	writeJSON(w, http.StatusOK, family)
}

func (h *FamilyHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	members, err := h.familyService.ListMembers(*claims.FamilyID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list members"})
		return
	}
	if members == nil {
		members = []models.User{}
	}
	writeJSON(w, http.StatusOK, members)
}

type addMemberRequest struct {
	Username    string `json:"username"     validate:"required,min=2,max=50,alphanum"`
	Pin         string `json:"pin"          validate:"required,len=6,numeric"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
}

func (h *FamilyHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	var req addMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if !validateRequest(w, req) {
		return
	}

	member, err := h.familyService.AddMember(*claims.FamilyID, req.Username, req.Pin, req.DisplayName)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, member)
}

type updateMemberRequest struct {
	DisplayName string `json:"display_name" validate:"omitempty,min=1,max=100"`
	Pin         string `json:"pin"          validate:"omitempty,len=6,numeric"`
}

func (h *FamilyHandler) UpdateMember(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid member id"})
		return
	}

	var req updateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if !validateRequest(w, req) {
		return
	}

	if err := h.familyService.UpdateMember(id, req.DisplayName, req.Pin, *claims.FamilyID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "member updated"})
}

func (h *FamilyHandler) DeactivateMember(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid member id"})
		return
	}

	if err := h.familyService.DeactivateMember(id, *claims.FamilyID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "member deactivated"})
}

func (h *FamilyHandler) ReactivateMember(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid member id"})
		return
	}

	if err := h.familyService.ReactivateMember(id, *claims.FamilyID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "member reactivated"})
}
