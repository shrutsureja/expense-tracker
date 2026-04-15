package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"

	"github.com/go-chi/chi/v5"
)

type TagHandler struct {
	tagRepo *repository.TagRepository
}

func NewTagHandler(tagRepo *repository.TagRepository) *TagHandler {
	return &TagHandler{tagRepo: tagRepo}
}

func (h *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	var familyID int64
	if claims.FamilyID != nil {
		familyID = *claims.FamilyID
	}

	tags, err := h.tagRepo.GetAllForFamily(familyID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list tags"})
		return
	}
	if tags == nil {
		tags = []models.Tag{}
	}
	writeJSON(w, http.StatusOK, tags)
}

func (h *TagHandler) SuggestedTags(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())

	tags, err := h.tagRepo.GetSuggestedForUser(claims.UserID, 10)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get suggested tags"})
		return
	}
	if tags == nil {
		tags = []models.Tag{}
	}
	writeJSON(w, http.StatusOK, tags)
}

type createTagRequest struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	var req createTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tag name is required"})
		return
	}
	if len(req.Name) > 50 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tag name must be 50 characters or less"})
		return
	}

	tag := &models.Tag{
		Name:     req.Name,
		Icon:     req.Icon,
		FamilyID: claims.FamilyID,
	}

	if err := h.tagRepo.Create(tag); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to create tag (may already exist)"})
		return
	}

	writeJSON(w, http.StatusCreated, tag)
}

func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	tagID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid tag id"})
		return
	}

	if err := h.tagRepo.DeleteForFamily(tagID, *claims.FamilyID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "tag deleted"})
}
