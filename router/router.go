package router

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"

	handler "github.com/PrimadonnaGit/seoulbitz-go/handler"
)

// Router function
func Router() *echo.Echo {
	// echo.New()를 사용하여 *Echo를 리턴 받는다
	e := echo.New()

	// debug 모드로 사용하기 위해서는 디버그 설정을 true로 변경
	e.Debug = false

	// echo middleware func
	e.Use(middleware.Logger())  //Setting logger
	e.Use(middleware.Recover()) //Recover from panics anywhere in the chain
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{ //CORS Middleware
	    AllowOrigins: []string{"*"},
	    AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Health check!
	e.GET("/healthy", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy!")
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Seoulbitz!")
	})

	// Router List

	crawlerRouterGroup := e.Group("/crawling")
	{
		crawlerRouterGroup.GET("/:searchKeyword", handler.ExecCrawling)
	}

	shopRouterGroup := e.Group("/shop")
	{
		shopRouterGroup.GET("/", handler.GetShop)
	}

	return e
}
