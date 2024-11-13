package handlers

import (
	"bytes"
	"cache-demo/internal/errors"
	"cache-demo/mocks"
	"encoding/json"
	"path"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
) 

func TestGetObjectByIdHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)

    cacheServiceMock := mocks.NewServiceMock()
    cacheHandler := NewCacheHandler(cacheServiceMock)

    testCases := []struct {
        name         string
        objectId       string
        mockResponse any
        mockError    error
        expectedCode int
        expectedBody string
    }{
        {
            name:         "Object Found",
            objectId:      "123",
            mockResponse: map[string]interface{}{"name":"John Doe"},
            mockError:    nil,
            expectedCode: http.StatusOK,
            expectedBody: `{"name":"John Doe"}`,
        },
        {
            name:         "Object Not Found",
            objectId:      "999",
            mockResponse: nil,
            mockError:    errors.ObjectNotFoundException("999"),
            expectedCode: http.StatusNotFound,
            expectedBody: "\"Object with id: 999 is not found\"",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {

            cacheServiceMock.On("GetObjectById", tc.objectId).Return(tc.mockResponse, tc.mockError)

            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)

            c.Params = append(c.Params, gin.Param{Key: OBJECT_ID, Value: tc.objectId})

            handlerFunc := cacheHandler.GetObjectById()
            handlerFunc(c)

            assert.Equal(t, tc.expectedCode, w.Code)
            assert.Equal(t, tc.expectedBody, w.Body.String())

            cacheServiceMock.AssertExpectations(t)
        })
    }
}

func TestPutObjectByIdHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)

    cacheServiceMock := mocks.NewServiceMock()
    cacheHandler := NewCacheHandler(cacheServiceMock)

    testCases := []struct {
        name         string
        objectId       string
        object map[string]interface{}
        expires uint32
        mockResponse bool
        expectedCode int
        expectedBody string
    }{
        {
            name:         "Object Created",
            objectId:      "123",
            object: map[string]interface{}{"name":"John Doe"},
            expires: 1728555010,
            mockResponse: false,
            expectedCode: http.StatusCreated,
            expectedBody: "",
        },
        {
            name:         "Object Updated",
            objectId:      "999",
            object: map[string]interface{}{"name":"John Doe"},
            expires: 1728555010,
            mockResponse: true,
            expectedCode: http.StatusOK,
            expectedBody: "",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {

            cacheServiceMock.On("PutObjectById", tc.objectId, tc.object, tc.expires).Return(tc.mockResponse)

            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)

            c.Params = append(c.Params, gin.Param{Key: OBJECT_ID, Value: tc.objectId})

            obj, _ := json.Marshal(tc.object)
            reqBody := bytes.NewBuffer(obj)
            c.Request, _ = http.NewRequest(http.MethodPut, path.Join(OBJECT_ID, tc.objectId) , reqBody)
            c.Request.Header.Set("Content-Type", "application/json")
            c.Request.Header.Set(EXPIRES_AT, "2024-10-10 10:10:10")
            
            handlerFunc := cacheHandler.PutObjectById()
            handlerFunc(c)

            assert.Equal(t, tc.expectedCode, w.Code)

            cacheServiceMock.AssertExpectations(t)
        })
    }
}