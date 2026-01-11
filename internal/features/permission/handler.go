package permission

import (
	"encoding/json"
	"mes-lite-back/pkg"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.create)
	r.Get("/{id}", h.getById)
	r.Get("/name/{name}", h.getByName)
	r.Get("/", h.list)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	return r
}

type CreatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

// @Summary Создать разрешение
// @Description Создать новое разрешение
// @Tags permissions
// @Accept json
// @Produce json
// @Param input body CreatePermissionRequest true "Данные разрешения"
// @Success 201 {object} PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /permissions [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{Message: "Неккоректный формат запроса"})
		return
	}

	if req.Name == "" {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{Message: "Название разрешения не может быть пустым"})
		return
	}

	permission := &Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.CreatePermission(permission); err != nil {
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Message: "Ошибка при создании разрешения"})
		return
	}

	resp := PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
	}
	pkg.RespondJSON(w, http.StatusCreated, resp)
}

// @Summary Получить разрешение по ID
// @Description Получить разрешение по идентификатору
// @Tags permissions
// @Produce json
// @Param id path int true "ID разрешения"
// @Success 200 {object} PermissionResponse
// @Failure 404 {object} ErrorResponse
// @Router /permissions/{id} [get]
func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {

	id := pkg.ParamID(r)

	permission, err := h.service.GetPermissionById(id)
	if err != nil {
		pkg.RespondJSON(w, http.StatusNotFound, ErrorResponse{Message: "Разрешение не найдено"})
		return
	}

	resp := PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
	}
	pkg.RespondJSON(w, http.StatusOK, resp)
}

// @Summary Получить разрешение по имени
// @Description Получить разрешение по имени
// @Tags permissions
// @Produce json
// @Param name path string true "Имя разрешения"
// @Success 200 {object} PermissionResponse
// @Failure 404 {object} ErrorResponse
// @Router /permissions/name/{name} [get]
func (h *Handler) getByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	permission, err := h.service.GetPermissionByName(name)
	if err != nil {
		pkg.RespondJSON(w, http.StatusNotFound, ErrorResponse{Message: "Разрешение не найдено"})
		return
	}

	resp := PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
	}
	pkg.RespondJSON(w, http.StatusOK, resp)
}

// @Summary Получить список разрешений
// @Description Получить все разрешения
// @Tags permissions
// @Produce json
// @Success 200 {array} PermissionResponse
// @Failure 500 {object} ErrorResponse
// @Router /permissions [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.service.ListPermissions()
	if err != nil {
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Message: "Ошибка при получении списка разрешений"})
		return
	}

	var resp []PermissionResponse
	for _, p := range permissions {
		resp = append(resp, PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
		})
	}
	pkg.RespondJSON(w, http.StatusOK, resp)
}

// @Summary Обновить разрешение
// @Description Обновить данные разрешения
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path int true "ID разрешения"
// @Param input body UpdatePermissionRequest true "Данные разрешения"
// @Success 200 {object} PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /permissions/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := pkg.ParamID(r)
	if id == 0 {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Некорректный ID разрешения",
		})
		return
	}

	var req UpdatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Некорректный JSON",
		})
		return
	}

	permission := &Permission{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.UpdatePermission(permission); err != nil {
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Message: "Ошибка при обновлении разрешения",
		})
		return
	}

	resp := PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
	}
	pkg.RespondJSON(w, http.StatusOK, resp)
}

// @Summary Удалить разрешение
// @Description Удалить разрешение по ID
// @Tags permissions
// @Param id path int true "ID разрешения"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /permissions/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := pkg.ParamID(r)
	if id == 0 {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Некорректный ID разрешения",
		})
		return
	}

	if err := h.service.DeletePermission(id); err != nil {
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Message: "Ошибка при удалении разрешения",
		})
		return
	}

	pkg.RespondJSON(w, http.StatusNoContent, nil)
}
