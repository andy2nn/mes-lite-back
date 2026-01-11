package role

import (
	"encoding/json"
	"log/slog"
	"mes-lite-back/pkg"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.getByID)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	r.Get("/{id}/permissions", h.getRolePermissions)
	r.Put("/{id}/permissions", h.updateRolePermissions)

	return r
}

type CreateRequest struct {
	Name          string  `json:"name" validate:"required" example:"Администратор"`
	PermissionIDs []int64 `json:"permission_ids,omitempty" example:"1,2,3,4,5"`
}

type UpdateRequest struct {
	RoleID        int64   `json:"role_id" example:"1"`
	Name          string  `json:"name" validate:"required" example:"Модератор"`
	PermissionIDs []int64 `json:"permission_ids,omitempty" example:"1,2,3"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Описание ошибки"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"Операция выполнена успешно"`
}

// ListRoles godoc
// @Summary Получить список всех ролей
// @Description Возвращает список всех ролей с их разрешениями
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {array} Role
// @Failure 500 {object} ErrorResponse
// @Router /roles [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.ListRoles()
	if err != nil {
		slog.Error("list roles failed", slog.Any("err", err))
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Не удалось получить список ролей"})
		return
	}
	pkg.RespondJSON(w, http.StatusOK, roles)
}

// CreateRole godoc
// @Summary Создать новую роль
// @Description Создает новую роль с указанными разрешениями
// @Tags roles
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Данные для создания роли"
// @Success 201 {object} Role
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Некорректный JSON"})
		return
	}

	if req.Name == "" {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Название роли обязательно"})
		return
	}

	role := &Role{Name: req.Name}

	if err := h.service.CreateRole(role, req.PermissionIDs); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			pkg.RespondJSON(w, http.StatusConflict, ErrorResponse{Error: "Роль с таким названием уже существует"})
			return
		}
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка при создании роли"})
		return
	}

	createdRole, err := h.service.GetRole(role.ID)
	if err != nil {
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка при получении созданной роли"})
		return
	}

	pkg.RespondJSON(w, http.StatusCreated, createdRole)
}

// GetRole godoc
// @Summary Получить роль по ID
// @Description Возвращает роль по указанному ID со списком разрешений
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Success 200 {object} Role
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles/{id} [get]
func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := pkg.ParamID(r)
	if id == 0 {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Некорректный ID роли"})
		return
	}

	role, err := h.service.GetRole(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			pkg.RespondJSON(w, http.StatusNotFound, ErrorResponse{Error: "Роль не найдена"})
			return
		}
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка при получении роли"})
		return
	}

	pkg.RespondJSON(w, http.StatusOK, role)
}

// UpdateRole godoc
// @Summary Обновить роль
// @Description Обновляет информацию о роли и её разрешения
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Param request body UpdateRequest true "Данные для обновления роли"
// @Success 200 {object} Role
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := pkg.ParamID(r)
	if id == 0 {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Некорректный ID роли"})
		return
	}

	existingRole, err := h.service.GetRole(id)
	if err != nil {
		pkg.RespondJSON(w, http.StatusNotFound, ErrorResponse{Error: "Роль не найдена"})
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Некорректный JSON"})
		return
	}

	existingRole.Name = req.Name
	if err := h.service.UpdateRole(existingRole); err != nil {
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка при обновлении роли"})
		return
	}

	if req.PermissionIDs != nil {
		if err := h.service.UpdatePermissions(id, req.PermissionIDs); err != nil {
			pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка при обновлении разрешений"})
			return
		}
	}

	updatedRole, _ := h.service.GetRole(id)
	pkg.RespondJSON(w, http.StatusOK, updatedRole)
}

// DeleteRole godoc
// @Summary Удалить роль
// @Description Удаляет роль по указанному ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := pkg.ParamID(r)
	if id == 0 {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Некорректный ID роли"})
		return
	}

	if err := h.service.DeleteRole(id); err != nil {
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка при удалении роли"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetRolePermissions godoc
// @Summary Получить разрешения роли
// @Description Возвращает список разрешений для указанной роли
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Success 200 {array} permission.Permission
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles/{id}/permissions [get]
func (h *Handler) getRolePermissions(w http.ResponseWriter, r *http.Request) {
	id := pkg.ParamID(r)
	role, err := h.service.GetRole(id)
	if err != nil {
		pkg.RespondJSON(w, http.StatusNotFound, ErrorResponse{Error: "Роль не найдена"})
		return
	}
	pkg.RespondJSON(w, http.StatusOK, role.Permissions)
}

// UpdateRolePermissions godoc
// @Summary Обновить разрешения роли
// @Description Полностью заменяет список разрешений для указанной роли
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Param request body UpdateRequest true "Список ID разрешений"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles/{id}/permissions [put]
func (h *Handler) updateRolePermissions(w http.ResponseWriter, r *http.Request) {
	id := pkg.ParamID(r)

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.RespondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Некорректный JSON"})
		return
	}

	if err := h.service.UpdatePermissions(id, req.PermissionIDs); err != nil {
		pkg.RespondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Ошибка при обновлении разрешений"})
		return
	}

	pkg.RespondJSON(w, http.StatusOK, SuccessResponse{Message: "Разрешения успешно обновлены"})
}
