package handler

import (
	"TaskForge/internal/interfaces/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	usecase auth.UseCaseAuth
}

func NewAuthHandler(u auth.UseCaseAuth) *AuthHandler {
	return &AuthHandler{usecase: u}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} auth.ResponseAuth
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.usecase.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя в систему
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Данные для входа"
// @Success 200 {object} auth.ResponseAuth
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.usecase.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
