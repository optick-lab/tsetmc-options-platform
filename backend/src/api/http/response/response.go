package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ResponseEnvelope is the standard API response structure.
type ResponseEnvelope struct {
	IsSuccessful bool        `json:"is_successful"`
	RequestUUID  string      `json:"request_uuid"`
	StatusCode   int         `json:"status_code"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
	Errors       interface{} `json:"errors"`
	ErrorCode    int         `json:"error_code"`
}

// Builder constructs and sends standard response envelopes.
type Builder struct {
	c            *gin.Context
	isSuccessful *bool
	statusCode   int
	message      string
	data         interface{}
	errors       interface{}
	errorCode    int
}

// Api initializes a response builder for the gin context.
func Api(c *gin.Context) *Builder {
	return &Builder{
		c:          c,
		statusCode: http.StatusOK,
		data:       struct{}{},
		errors:     nil,
		errorCode:  0,
	}
}

// SetStatusCode sets the HTTP status code.
func (b *Builder) SetStatusCode(code int) *Builder {
	b.statusCode = code
	return b
}

// SetMessage sets the response message.
func (b *Builder) SetMessage(msg string) *Builder {
	b.message = msg
	return b
}

// SetData sets the response data payload.
func (b *Builder) SetData(data interface{}) *Builder {
	if data == nil {
		b.data = struct{}{}
	} else {
		b.data = data
	}
	return b
}

// SetErrors sets error details.
func (b *Builder) SetErrors(errs interface{}) *Builder {
	b.errors = errs
	return b
}

// SetErrorCode sets an application-level error code.
func (b *Builder) SetErrorCode(code int) *Builder {
	b.errorCode = code
	return b
}

// SetSuccessful explicitly sets the is_successful flag.
func (b *Builder) SetSuccessful(isSuccessful bool) *Builder {
	b.isSuccessful = &isSuccessful
	return b
}

// Send renders the response envelope as JSON.
func (b *Builder) Send() {
	reqID := ""
	if b.c != nil {
		if id, exists := b.c.Get("request_uuid"); exists {
			if idStr, ok := id.(string); ok && idStr != "" {
				reqID = idStr
			}
		}
		if reqID == "" {
			reqID = b.c.GetHeader("X-Request-ID")
		}
	}
	if reqID == "" {
		reqID = uuid.NewString()
	}

	isSuccess := b.statusCode >= 200 && b.statusCode < 400
	if b.isSuccessful != nil {
		isSuccess = *b.isSuccessful
	}

	envelope := ResponseEnvelope{
		IsSuccessful: isSuccess,
		RequestUUID:  reqID,
		StatusCode:   b.statusCode,
		Message:      b.message,
		Data:         b.data,
		Errors:       b.errors,
		ErrorCode:    b.errorCode,
	}

	if b.c != nil {
		b.c.JSON(b.statusCode, envelope)
	}
}
