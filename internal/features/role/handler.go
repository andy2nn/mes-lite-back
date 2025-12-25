package role

import (
	"encoding/json"
	"log/slog"
	"mes-lite-back/internal/features/permission"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
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

// CreateRequest запрос на создание роли
// @Description Запрос для создания новой роли с опциональным списком разрешений
type CreateRequest struct {
	Name          string  `json:"name" validate:"required" example:"Администратор"`
	PermissionIDs []int64 `json:"permission_ids,omitempty" example:"1,2,3,4,5"`
}

// UpdateRequest запрос на обновление роли
// @Description Запрос для обновления роли
type UpdateRequest struct {
	Name          string  `json:"name" validate:"required" example:"Модератор"`
	PermissionIDs []int64 `json:"permission_ids,omitempty" example:"1,2,3"`
}

// ErrorResponse стандартный ответ с ошибкой
// @Description Стандартный формат ответа при ошибке
type ErrorResponse struct {
	Error string `json:"error" example:"Описание ошибки"`
}

// SuccessResponse стандартный ответ об успехе
// @Description Стандартный формат ответа при успешной операции
type SuccessResponse struct {
	Message string `json:"message" example:"Операция выполнена успешно"`
}

// ListRoles godoc
// @Summary Получить список всех ролей
// @Description Возвращает список всех ролей с их разрешениями
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {array} Role "Список ролей"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.ListRoles()
	if err != nil {
		slog.Error("list roles failed", slog.Any("err", err))
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Не удалось получить список ролей",
		})
		return
	}

	respondJSON(w, http.StatusOK, roles)
}

// CreateRole godoc
// @Summary Создать новую роль
// @Description Создает новую роль с указанными разрешениями
// @Tags roles
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Данные для создания роли"
// @Success 201 {object} Role "Созданная роль"
// @Failure 400 {object} ErrorResponse "Некорректный запрос"
// @Failure 409 {object} ErrorResponse "Роль с таким именем уже существует"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Некорректный JSON",
		})
		return
	}

	if req.Name == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Название роли обязательно",
		})
		return
	}

	role := &Role{
		Name: req.Name,
	}

	// Используем метод Create с передачей permissionIDs
	if err := h.service.CreateRole(role, req.PermissionIDs); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			respondJSON(w, http.StatusConflict, ErrorResponse{
				Error: "Роль с таким названием уже существует",
			})
			return
		}

		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при создании роли",
		})
		return
	}

	// Получаем созданную роль с разрешениями
	createdRole, err := h.service.GetRole(role.ID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при получении созданной роли",
		})
		return
	}

	respondJSON(w, http.StatusCreated, createdRole)
}

// GetRole godoc
// @Summary Получить роль по ID
// @Description Возвращает роль по указанному ID со списком разрешений
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Success 200 {object} Role "Роль с разрешениями"
// @Failure 400 {object} ErrorResponse "Некорректный ID"
// @Failure 404 {object} ErrorResponse "Роль не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles/{id} [get]
func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)
	if id == 0 {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Некорректный ID роли",
		})
		return
	}

	role, err := h.service.GetRole(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "record not found") {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error: "Роль не найдена",
			})
			return
		}

		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при получении роли",
		})
		return
	}

	respondJSON(w, http.StatusOK, role)
}

// UpdateRole godoc
// @Summary Обновить роль
// @Description Обновляет информацию о роли и её разрешения
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Param request body UpdateRequest true "Данные для обновления роли"
// @Success 200 {object} Role "Обновленная роль"
// @Failure 400 {object} ErrorResponse "Некорректный запрос"
// @Failure 404 {object} ErrorResponse "Роль не найдена"
// @Failure 409 {object} ErrorResponse "Роль с таким именем уже существует"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)
	if id == 0 {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Некорректный ID роли",
		})
		return
	}

	// Проверяем существование роли
	existingRole, err := h.service.GetRole(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "record not found") {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error: "Роль не найдена",
			})
			return
		}

		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при проверке роли",
		})
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Некорректный JSON",
		})
		return
	}

	if req.Name == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Название роли обязательно",
		})
		return
	}

	// Обновляем роль
	existingRole.Name = req.Name
	if err := h.service.UpdateRole(existingRole); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			respondJSON(w, http.StatusConflict, ErrorResponse{
				Error: "Роль с таким названием уже существует",
			})
			return
		}

		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при обновлении роли",
		})
		return
	}

	// Обновляем разрешения, если они указаны
	if req.PermissionIDs != nil {
		if err := h.service.UpdatePermissions(id, req.PermissionIDs); err != nil {
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: "Ошибка при обновлении разрешений роли",
			})
			return
		}
	}

	// Получаем обновленную роль
	updatedRole, err := h.service.GetRole(id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при получении обновленной роли",
		})
		return
	}

	respondJSON(w, http.StatusOK, updatedRole)
}

// DeleteRole godoc
// @Summary Удалить роль
// @Description Удаляет роль по указанному ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Success 204 "Роль успешно удалена"
// @Failure 400 {object} ErrorResponse "Некорректный ID"
// @Failure 404 {object} ErrorResponse "Роль не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)
	if id == 0 {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Некорректный ID роли",
		})
		return
	}

	// Сначала получаем роль, чтобы передать в Delete
	role, err := h.service.GetRole(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "record not found") {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error: "Роль не найдена",
			})
			return
		}

		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при получении роли для удаления",
		})
		return
	}

	if err := h.service.DeleteRole(role.ID); err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при удалении роли",
		})
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
// @Success 200 {array} permission.Permission "Список разрешений роли"
// @Failure 400 {object} ErrorResponse "Некорректный ID"
// @Failure 404 {object} ErrorResponse "Роль не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles/{id}/permissions [get]
func (h *Handler) getRolePermissions(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)
	if id == 0 {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Некорректный ID роли",
		})
		return
	}

	role, err := h.service.GetRole(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "record not found") {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error: "Роль не найдена",
			})
			return
		}

		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при получении разрешений роли",
		})
		return
	}

	respondJSON(w, http.StatusOK, role.Permissions)
}

// UpdateRolePermissions godoc
// @Summary Обновить разрешения роли
// @Description Полностью заменяет список разрешений для указанной роли
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "ID роли"
// @Param request body permission.PermissionAssignRequest true "Список ID разрешений"
// @Success 200 {object} SuccessResponse "Разрешения успешно обновлены"
// @Failure 400 {object} ErrorResponse "Некорректный запрос"
// @Failure 404 {object} ErrorResponse "Роль или некоторые разрешения не найдены"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles/{id}/permissions [put]
func (h *Handler) updateRolePermissions(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)
	if id == 0 {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Некорректный ID роли",
		})
		return
	}

	var req permission.PermissionAssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Некорректный JSON",
		})
		return
	}

	req.RoleID = id

	if err := h.service.UpdatePermissions(req.RoleID, req.PermissionIDs); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Ошибка при обновлении разрешений",
		})
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Разрешения успешно обновлены",
	})
}

// TODO: Вынести из хендлеров
func respondJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func paramID(r *http.Request) int64 {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}
