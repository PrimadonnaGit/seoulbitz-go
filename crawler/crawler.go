package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PrimadonnaGit/seoulbitz-go/mysql"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	seleniumPath = "./crawler/chromedriver"
	searchURL    = "https://map.kakao.com/"
	port         = 5000
	KAKAO     = "KakaoAK cc116147fce20da7314166dce21f0305"
	KAKAO_URL = "https://dapi.kakao.com/v2/local/search/address.json"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func sleepSecond(second time.Duration) {
	time.Sleep(second * time.Second)
}

type KakaoResp struct {
	Documents []struct {
		Address struct {
			AddressName       string `json:"address_name"`
			BCode             string `json:"b_code"`
			HCode             string `json:"h_code"`
			MainAddressNo     string `json:"main_address_no"`
			MountainYn        string `json:"mountain_yn"`
			Region1DepthName  string `json:"region_1depth_name"`
			Region2DepthName  string `json:"region_2depth_name"`
			Region3DepthHName string `json:"region_3depth_h_name"`
			Region3DepthName  string `json:"region_3depth_name"`
			SubAddressNo      string `json:"sub_address_no"`
			X                 string `json:"x"`
		} `json:"address"`
		AddressName string `json:"address_name"`
		AddressType string `json:"address_type"`
		RoadAddress struct {
			AddressName      string `json:"address_name"`
			BuildingName     string `json:"building_name"`
			MainBuildingNo   string `json:"main_building_no"`
			Region1DepthName string `json:"region_1depth_name"`
			Region2DepthName string `json:"region_2depth_name"`
			Region3DepthName string `json:"region_3depth_name"`
			RoadName         string `json:"road_name"`
			NdYn             string `json:"nd_yn"`
			X                string `json:"x"`
			Y                string `json:"y"`
			ZoneNo           string `json:"zone_no"`
		} `json:"road_address"`
		X string `json:"x"`
		Y string `json:"y"`
	} `json:"documents"`
	Meta struct {
		IsEnd         bool `json:"is_end"`
		PageableCount int  `json:"pageable_count"`
		TotalCount    int  `json:"total_count"`
	} `json:"meta"`
}

func getKAKAOLatlng(address string) (string, string) {

	req, _ := http.NewRequest("GET", KAKAO_URL, nil)

	req.Header.Add("Authorization", KAKAO)
	query := req.URL.Query()
	query.Add("query", address)
	query.Add("size", "1")
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	checkErr(err)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var kakaoResp KakaoResp
	err = json.Unmarshal(body, &kakaoResp)
	checkErr(err)

	if len(kakaoResp.Documents) > 0 {
		return kakaoResp.Documents[0].Y, kakaoResp.Documents[0].X
	} else {
		return "0", "0"
	}

}

func loopPlaceElements(placeItems []selenium.WebElement) {
	for _, placeItem := range placeItems {
		placeTitleElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".head_item .tit_name .link_name")
		checkErr(err)

		placeSubCategoryElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".head_item .subcategory")
		checkErr(err)

		placeScoreElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".rating .score .num")
		checkErr(err)

		placeScoreCountElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".rating .score a")
		checkErr(err)

		placeReviewCountElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".rating a em")
		checkErr(err)

		placeAddressElement, err := placeItem.FindElement(selenium.ByCSSSelector, ".info_item .addr p[data-id='address']")
		checkErr(err)

		placeTitle, _ := placeTitleElement.Text()
		placeSubCategory, _ := placeSubCategoryElement.Text()
		placeScore, _ := placeScoreElement.Text()

		placeScoreCount, _ := placeScoreCountElement.Text()
		placeScoreCount = strings.Replace(placeScoreCount, ",", "", -1)
		placeScoreCount = strings.Trim(placeScoreCount, "건")

		placeReviewCount, _ := placeReviewCountElement.Text()
		placeReviewCount = strings.Replace(placeReviewCount, ",", "", -1)

		placeAddress, _ := placeAddressElement.Text()
		placeAddress = strings.Replace(placeAddress, ",", "", -1)
		placeAddress = strings.Join(strings.Split(placeAddress, " ")[:4], " ")

		lat, lng := getKAKAOLatlng(placeAddress)

		db := mysql.ConnectDB()
		defer db.Close()

		_, err = db.Exec("INSERT INTO kakao_map (title, xpoint, ypoint, tag, score, scoreCount, reviewCount, address) VALUES (?,?,?,?,?,?,?,?)", placeTitle, lat, lng, placeSubCategory, placeScore, placeScoreCount, placeReviewCount, placeAddress)
		checkErr(err)

	}
}

func KakaoCrawling(searchKeyword string) {

	// chromeDriver := webdriver.NewChromeDriver(seleniumPath)
	// defer chromeDriver.Stop()
	// err := chromeDriver.Start()
	// checkErr(err)

	// desired := webdriver.Capabilities{"Platform": "Windows"}

	// required := webdriver.Capabilities{}
	// session, err := chromeDriver.NewSession(desired, required)
	// defer session.Delete()
	// checkErr(err)
	
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(seleniumPath, port)
	checkErr(err)
	defer service.Stop()

	caps := selenium.Capabilities{}

	caps.AddChrome(chrome.Capabilities{
		Args: []string{"--headless"},
	})

	session, err := selenium.NewRemote(caps, "http://127.0.0.1:5000")
	checkErr(err)

	fmt.Println("sss")

	defer session.Quit()
	
	err = session.Get(searchURL)
	checkErr(err)
	
	// 검색 키워드 입력
	keywordInput, _ := session.FindElement(selenium.ByCSSSelector, ".box_searchbar > input.query")
	err = keywordInput.SendKeys(searchKeyword)
	checkErr(err)

	err = keywordInput.SendKeys(selenium.EnterKey) // Enter key
	checkErr(err)

	sleepSecond(1)

	// 더보기
	moreBtn, _ := session.FindElement(selenium.ByCSSSelector, ".places > .more")
	err = moreBtn.SendKeys(selenium.EnterKey)
	checkErr(err)

	sleepSecond(1)

	pageBtns, _ := session.FindElements(selenium.ByCSSSelector, ".keywordSearch .pages .pageWrap a")

	// 페이지 순회
	n := 0
	for n < 34 {
		nextBtn, _ := session.FindElement(selenium.ByCSSSelector, ".keywordSearch .pages .pageWrap .next")
		for _, pageBtn := range pageBtns {
			pageBtn.SendKeys(selenium.EnterKey)
			sleepSecond(1)
			placeItems, _ := session.FindElements(selenium.ByCSSSelector, ".PlaceItem")

			loopPlaceElements(placeItems)
		}

		classNames, _ := nextBtn.GetAttribute("class")
		if strings.Contains(classNames, "disabled") {
			break
		}
		err = nextBtn.SendKeys(selenium.EnterKey)
		checkErr(err)

		n++
	}

}
