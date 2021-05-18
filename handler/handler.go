package handler

import (
	"fmt"
	"net/http"

	crawler "github.com/PrimadonnaGit/seoulbitz-go/crawler"
	"github.com/PrimadonnaGit/seoulbitz-go/mysql"
	echo "github.com/labstack/echo/v4"
)

type Shop struct {
	Index   int `db:"idx"`
	Title   string `db:"title"`
	X       float32 `db:"xpoint"`
	Y       float32 `db:"ypoint"`
	Tag     string `db:"tag"`
	Score    float64 `db:"score"`
	ScoreCount   int `db:"scoreCount"`
	ReviewCount int `db:"reviewCount"`
	Address  string `db:"address"`
}

func GetShop(c echo.Context) error {

	db := mysql.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM kakao_map Where Score > 4")
	if err != nil {
		fmt.Println(err)	
	}

	defer rows.Close()

	var shop Shop
	var shops []Shop
	for rows.Next() {
		err := rows.Scan(&shop.Index, &shop.Title, &shop.X, &shop.Y, &shop.Tag, &shop.Score, &shop.ScoreCount, &shop.ReviewCount, &shop.Address)
		if err != nil {
			fmt.Println(err)
		}
		shops = append(shops, shop)
	}

	return c.JSON(http.StatusOK, shops)
}

func ExecCrawling(c echo.Context) error {

	fmt.Println("here")

	searchKeyword := c.Param("searchKeyword")
	
	go func() {
		crawler.KakaoCrawling(searchKeyword)
	}()

	return c.String(http.StatusOK, "")
}
