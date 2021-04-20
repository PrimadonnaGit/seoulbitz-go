package router

import (
	"github.com/labstack/echo/v4"
	echo "github.com/labstack/echo/v4/middleware"
	middleware "net/http"
	// auth "[auth 폴더 위치]"
	// handler "[handler 폴더 위치]"
)

// Router function
func Router() *echo.Echo {
	// echo.New()를 사용하여 *Echo를 리턴 받는다
	e := echo.New()

	// debug 모드로 사용하기 위해서는 디버그 설정을 true로 변경
	e.Debug = true

	// echo middleware func
	e.Use(middleware.Logger())  //Setting logger
	e.Use(middleware.Recover()) //Recover from panics anywhere in the chain
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{ //CORS Middleware
	//     AllowOrigins: []string{"*"},
	//     AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	// }))

	// Health check!
	e.GET("/healthy", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy!")
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Seoulbitz!")
	})

	// Router List
	// getList := e.Group("/get")
	// {
	// getList.GET("[path]", handler.[요청함수])
	// getList.GET("[path][:pathParameter]", handler.[요청함수])
	// }

	// admin := e.Group("/admin")
	// {
	// admin.GET("[path]", handler.[요청함수])
	// admin.GET("[path]", handler.[요청함수], auth.[로그인체크함수], auth.[어드민체크함수])
	// }

	// login := e.Group("/login")
	// {
	// login.POST("", auth.auth.[로그인함수])
	// }

	return e
}
