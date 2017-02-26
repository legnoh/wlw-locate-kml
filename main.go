package main

import (
	//   log "github.com/Sirupsen/logrus"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

var (
	locationUrl = "https://wonderland-wars.net/location_list.html"
	hostUrl     = "https://wonderland-wars.net/"
)

type Location []struct {
	Name       string // 店舗名
	Address    string // 住所
	LatLong    string // 緯度経度
	ShopUrl    string // 店舗情報URL
	RankingUrl string // ランキングページURL
	Library    bool   // WonderlandLIBRARYの有無
}

type Locations []Location

func GetPage(url string) {
	doc, _ := goquery.NewDocument(url)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		fmt.Println(url)
	})
}

func GetLocationList(url string) {

	doc, _ := goquery.NewDocument(url)

	// 全ての.address_box(店舗情報)に対して、以下の処理を繰り返す
	doc.Find(".address_box").Each(func(i int, s *goquery.Selection) {

		// location構造体の定義
		// location := Location{}

		// location_nameを取得して、Nameに追加

		// addressを取得して、Addressに追加

		// ShopURLにアクセスし、ページ内のGoogleMapへのURLからLatLongを取得

		// location_nameのURLを取得して、ShopUrlに追加

		// store_rankingのURLを取得して、RankingUrlに追加

		// store_ranking配下のicon_terminalが存在する場合、Libraryをtrueに変更

		// s.Find("a").Each(func(_ int, t *goquery.Selection) {
		// 	url, _ := t.Attr("href")
		// 	fmt.Println(url)
		// })
	})
}

func main() {
	GetLocationList(locationUrl)
}
