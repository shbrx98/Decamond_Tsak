package handlers

import (
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

func NewUserHandler(get *appUser.GetUserUseCase, list *appUser.ListUsersUseCase, logger logger.Logger) *UserHandler {
    return &UserHandler{getUserUC: get, listUsersUC: list, logger: logger}
}

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
    c.JSON(http.StatusOK, res)
}

func parseInt(s string, def int) int {
    var i int
    _, err := fmt.Sscanf(s, "%d", &i)
    if err != nil || i <= 0 {
        return def
    }
    return i
}