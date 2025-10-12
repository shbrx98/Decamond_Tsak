package handlers

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    appUser "github.com/shbrx98/Decamond_Tsak/internal/application/user"
    "github.com/shbrx98/Decamond_Tsak/internal/interfaces/http/dto"
    appErr "github.com/shbrx98/Decamond_Tsak/internal/pkg/errors"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"
)

type UserHandler struct {
    getUserUC   *appUser.GetUserUseCase
    listUsersUC *appUser.ListUsersUseCase
    logger      logger.Logger
}

func NewUserHandler(get *appUser.GetUserUseCase, list *appUser.ListUsersUseCase, l logger.Logger) *UserHandler {
    return &UserHandler{getUserUC: get, listUsersUC: list, logger: l}
}

// Get user godoc
// @Summary      Get user by ID
// @Description  Protected endpoint
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "VALIDATION_ERROR", Message: "invalid user id"})
        return
    }
    u, err := h.getUserUC.Execute(c.Request.Context(), id)
    if err != nil {
        switch err.(type) {
        case *appErr.NotFoundError:
            c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "NOT_FOUND", Message: "user not found"})
        default:
            h.logger.WithError(err).Error("get user error")
            c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "INTERNAL_ERROR", Message: "internal error"})
        }
        return
    }
    c.JSON(http.StatusOK, dto.UserFromDomain(u))
}

// List users godoc
// @Summary      List users
// @Description  Protected endpoint with pagination
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        query  query     string  false  "Search in phone"
// @Param        page   query     int     false  "Page"            default(1)   minimum(1)
// @Param        limit  query     int     false  "Items per page"  default(20)  minimum(1) maximum(100)
// @Success      200    {object}  dto.PaginatedUsersResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
    query := c.Query("query")
    page := parseInt(c.DefaultQuery("page", "1"), 1)
    limit := parseInt(c.DefaultQuery("limit", "20"), 20)

    res, err := h.listUsersUC.Execute(c.Request.Context(), query, page, limit)
    if err != nil {
        h.logger.WithError(err).Error("list users error")
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "INTERNAL_ERROR", Message: "internal error"})
        return
    }
    c.JSON(http.StatusOK, dto.PaginatedUsersFromDomain(res))
}

func parseInt(s string, def int) int {
    var i int
    if _, err := fmt.Sscanf(s, "%d", &i); err != nil || i <= 0 {
        return def
    }
    return i
}