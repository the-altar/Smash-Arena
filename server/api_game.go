package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func startGameHandler(c echo.Context) error {

	r := &startGameReq{}
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, 0)
	}
	fmt.Println(r)
	return c.JSON(http.StatusOK, 1)
}
