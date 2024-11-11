package handlers

import (
	"cache-demo/internal/errors"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const(
	EMPTY = ""
	HTTP_REQUEST_TIMEOUT = 5*time.Second
	OBJECT_ID = "objectId"
	EXPIRES_AT = "expires_at"
	DEFAULT_EXPIRES = 86400
)

type CacheServiceInterface interface {
	GetObjectById(objectId string) (any, error)
	PutObjectById(objectId string, object any, expired uint32) bool
}

type CacheHandler struct {
	CacheService CacheServiceInterface
}

func NewCacheHandler(service CacheServiceInterface) *CacheHandler {
	return &CacheHandler{CacheService: service}
}

func(h *CacheHandler) GetObjectById() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), HTTP_REQUEST_TIMEOUT)
		defer cancel()

		objectId := c.Param(OBJECT_ID)

		object, err := h.CacheService.GetObjectById(objectId)
		if object == nil {
			slog.Error(
				"Object not found", 
				slog.String("Object ID", objectId),
			)
			c.JSON(
				http.StatusNotFound,
				err.Error(),
			)
			return
		}
		if err != nil {
			slog.Error(
				"Exception during the object finding",
				slog.String("Error", err.Error()))
			c.JSON(
				http.StatusInternalServerError,
				err.Error(),
			)
			return
		}
		c.JSON(
			http.StatusOK,
			object,
		)
		slog.Info(
			"Object found",
			slog.String("Object Id", objectId),
		)
		slog.Debug(
			"Object found",
			slog.String("Object Id", objectId),
			slog.Any("Object", object),
		)
	}
}

func(h *CacheHandler) PutObjectById() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), HTTP_REQUEST_TIMEOUT)
		defer cancel()

		var expires uint32
		var object map[string]interface{}

		objectId := c.Param(OBJECT_ID)
		body, _ := io.ReadAll(c.Request.Body)
		err := json.Unmarshal(body, &object)
		if err != nil {
			slog.Error(
				"Invalid request body. Body should be in valid JSON format",
				slog.String("Error", err.Error()),
			)
			c.JSON(
				http.StatusBadRequest,
				err.Error(),
			)
		}
		expires_at := c.GetHeader(EXPIRES_AT)
		expires, err = calculateExpirationDate(expires_at)
		if err != nil {
			slog.Error(
				"Incorrect Date/Time format in request",
				slog.String("Error", err.Error()))
			c.JSON(
				http.StatusBadRequest,
				err.Error(),
			)
		}
		ok := h.CacheService.PutObjectById(objectId, object, expires)
		if ok {
			c.JSON(
				http.StatusOK,
				nil,
			)
			slog.Info(
				"Object successfully updated",
				slog.String("Object Id", objectId),
			)
		} else {
			c.JSON(
				http.StatusCreated,
				nil,
			)
			slog.Info(
				"Object successfully created",
				slog.String("Object Id", objectId),
			)
		}
	}
}

func calculateExpirationDate(rawDate string) (uint32, error) {
	if rawDate != EMPTY {
		layout := "2006-01-02 15:04:05"
		date, err := time.Parse(layout, rawDate) 
		if err != nil {
			return 0, errors.IncorrectDateTimeFormatException()
		}
		return uint32(date.Unix()), nil
	}
	return uint32(time.Now().Local().Unix()) + DEFAULT_EXPIRES, nil
}