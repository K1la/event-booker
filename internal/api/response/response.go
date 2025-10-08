package response

import (
	"github.com/wb-go/wbf/ginext"
	"net/http"
)

type Success struct {
	Result interface{} `json:"result"`
}

type Error struct {
	Message string `json:"error"`
}

func JSON(c *ginext.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func Created(c *ginext.Context, result interface{}) {
	JSON(c, http.StatusCreated, Success{Result: result})
}

func OK(c *ginext.Context, result interface{}) {
	JSON(c, http.StatusOK, Success{Result: result})
}

func Internal(c *ginext.Context, err error) {
	JSON(c, http.StatusInternalServerError, Error{Message: err.Error()})
}

func BadRequest(c *ginext.Context, err error) {
	JSON(c, http.StatusBadRequest, Error{Message: err.Error()})
}

func Fail(c *ginext.Context, status int, err error) {
	JSON(c, status, Error{Message: err.Error()})
}
