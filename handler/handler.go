package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	echo "github.com/labstack/echo/v4"
)

type Shop struct {
	index   int
	title   string
	y       float32
	x       float32
	tag     string
	addr    string
	like    int
	insta   string
	imgList string
}

func GetShop(c echo.Context) error {
	// id := c.Param("id")

	var shops []Shop

	data, err := os.Open("foodie.json")
	defer data.Close()
	if err != nil {
		fmt.Println(err)
	}

	byteValue, err := ioutil.ReadAll(data)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(byteValue, &shops)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(shops)

	return c.JSON(http.StatusOK, shops)
}
