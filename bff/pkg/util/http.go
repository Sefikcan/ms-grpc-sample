package util

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/logger"
)

func GetRequestId(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), "RequestCtx", GetRequestId(c))
}

func PrepareLogging(ctx echo.Context, logger logger.Logger, err error) {
	logger.Errorf("Error, RequestId: %s, IPAddress: %s, Error: %s", GetRequestId(ctx), GetIPAddress(ctx), err)
}
