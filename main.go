package main

import (
	"image/color"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/text/unicode/norm"

	"github.com/PuerkitoBio/goquery"
	kml "github.com/twpayne/go-kml"
)

// Location は1件の店舗情報をもつ
type Location struct {
	Name       string  // 店舗名
	Address    string  // 住所
	Area       int     // エリア区分(1:北海道・東北, 2:関東, 3:東海, 4:北信越, 5:近畿, 6: 中国・四国, 7:九州・沖縄)
	Lat        float64 // 経度
	Long       float64 // 緯度
	ShopURL    string  // 店舗情報URL
	RankingURL string  // ランキングページURL
	Rank1st    string  // ランキング1位の人のpt
	Rank5th    string  // ランキング5位の人のpt
	Library    bool    // WonderlandLIBRARYの有無
}

var (
	locationURL  = "https://wonderland-wars.net/location_list.html"
	hostURL      = "https://wonderland-wars.net"
	gMapHostHead = "//maps.googleapis.com/maps/api/staticmap?center="
	gMapHostFoot = regexp.MustCompile("&markers=.*")
	locations1   = kml.Folder(kml.Name("北海道・東北"))
	locations2   = kml.Folder(kml.Name("関東"))
	locations3   = kml.Folder(kml.Name("東海"))
	locations4   = kml.Folder(kml.Name("北信越"))
	locations5   = kml.Folder(kml.Name("近畿"))
	locations6   = kml.Folder(kml.Name("中国・四国"))
	locations7   = kml.Folder(kml.Name("九州・沖縄"))
	libStyle     = "#icon-1664-0288D1/ "
	libSign      = "○"
	filePath     = "./result-" + strconv.FormatInt(time.Now().Unix(), 10) + ".kml"
	rankResult   = "0pt"
	rankNull     = "0pt"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {

	// 店舗情報一覧ページを取得
	locationPage, _ := goquery.NewDocument(locationURL)

	// 店舗数を取得
	shopSum := strconv.Itoa(locationPage.Find(".address_box").Length() - 1)

	// 全ての.address_box(店舗情報)に対して、以下の処理を繰り返す
	locationPage.Find(".address_box").Each(func(i int, s *goquery.Selection) {

		// pref=99の場合、NO DATAなので処理を飛ばす
		prefs, _ := s.Attr("pref")
		pref, _ := strconv.Atoi(prefs)
		if pref != 99 {

			// location構造体の定義
			l := Location{
				Name:       "",
				Address:    "",
				Area:       0,
				Lat:        0,
				Long:       0,
				ShopURL:    "",
				RankingURL: "",
				Rank1st:    rankNull,
				Rank5th:    rankNull,
				Library:    false,
			}

			// location_nameを取得して、Nameに追加(気持ち悪いので全半角の正規化を行う)
			l.Name = string(norm.NFKC.Bytes([]byte(s.Find(".location_name > a").Text())))

			// addressを取得して、Addressに追加(気持ち悪いので全半角の正規化を行う)
			l.Address = string(norm.NFKC.Bytes([]byte(s.Find(".address").Text())))

			// ページ分類ごとにエリア区分を区分け
			switch pref {
			case 40, 2, 6, 3, 42, 44, 39:
				l.Area = 1
			case 18, 28, 5, 20, 25, 26, 14, 46:
				l.Area = 2
			case 23, 1, 15, 41:
				l.Area = 3
			case 34, 32, 30, 4, 37:
				l.Area = 4
			case 22, 16, 33, 9, 47, 35:
				l.Area = 5
			case 29, 10, 36, 24, 45, 27, 12, 7, 19:
				l.Area = 6
			case 38, 8, 21, 31, 43, 17, 13, 11:
				l.Area = 7
			}

			// location_nameのURLを取得して、ShopURLに追加
			shopURL, locationExists := s.Find(".location_name > a").Attr("href")
			if locationExists {
				l.ShopURL = shopURL

				// ShopURLにアクセスし、ページ内のGoogleMapへのURLから緯度経度を取得してKMLにあうように転置させる
				// この際、非同期処理での一斉アクセスを避けるため、事前に配列番号秒分のsleepを入れて0.2rps程度になるように留める
				time.Sleep(2 * time.Second)
				shopPage, _ := goquery.NewDocument(l.ShopURL)
				gMapURL, mapExists := shopPage.Find(".access_map").Attr("src")
				if mapExists {

					// 緯度と経度部分のみ抜き出す
					gMapURL = strings.Replace(gMapURL, gMapHostHead, "", 1)
					gMapURL = gMapHostFoot.ReplaceAllString(gMapURL, "")

					// KMLに投入するため、逆転させてfloat64に変換させる
					longLat := strings.Split(gMapURL, ",")
					l.Long, _ = strconv.ParseFloat(longLat[1], 64)
					l.Lat, _ = strconv.ParseFloat(longLat[0], 64)
				}
			}

			// store_rankingのURLを取得して、ホスト名と結合してRankingURLに追加
			rankPath, storeRankingExists := s.Find(".store_ranking > a").Attr("href")
			if storeRankingExists {
				l.RankingURL = hostURL + strings.Trim(rankPath, ".")

				// ランキングURLページにもアクセスし、前月/今月ランキングの1位・5位を取得する
				rankPage, _ := goquery.NewDocument(l.RankingURL)
				rank1stNode := rankPage.Find(".block_rankig_special > .store_ranking_page")
				rank5thNode := rankPage.Find(".block_rankig_1st > .store_ranking_page").Eq(3)
				if rank1stNode.Length() != 0 {
					l.Rank1st = rank1stNode.Text()
				}
				if rank5thNode.Length() != 0 {
					l.Rank5th = rank5thNode.Text()
				}
			}

			// store_ranking配下のicon_terminalが存在する場合、Libraryをtrueに変更
			_, libraryExists := s.Find(".store_ranking > .icon_terminal > img").Attr("src")
			if libraryExists {
				l.Library = true
			} else {
				l.Library = false
			}

			// descとstyleはライブラリで内容に変化が出るので事前に作る
			desc := "所在地: " + l.Address + "<br>店舗URL: " + l.ShopURL + "<br>ランキング: " + l.RankingURL
			if l.Library {
				libSign = "○"
				libStyle = "#icon-1526-A52714"
			} else {
				libSign = "×"
				libStyle = "#icon-1598-0288D1"
			}
			desc += "<br>LIBRARY:" + libSign

			// ランキングも新規店舗の存在があるので多少作る
			if l.Rank5th == rankNull && l.Rank1st == rankNull {
				log.Warn("新店舗があるようです！: " + l.Name)
				rankResult = "ランキングなし"
				libStyle = "#icon-1881-0f9d58"
			} else {
				rankResult = l.Rank5th + " 〜 " + l.Rank1st
			}

			// PlaceMarkに全ての情報を結合して保管
			placemark := kml.Placemark(
				kml.Name(l.Name),
				kml.Description(desc),
				kml.ExtendedData(
					kml.SchemaData(
						"#extendInfomation",
						kml.SimpleData("住所", l.Address),
						// kml.SimpleData("店舗詳細情報", l.ShopURL),
						kml.SimpleData("ランキング", l.RankingURL),
						kml.SimpleData("ライブラリ設置", libSign),
						kml.SimpleData("ランキング結果(5〜1位)", rankResult),
					),
				),
				kml.Point(kml.Coordinates(kml.Coordinate{Lon: l.Long, Lat: l.Lat})),
				kml.StyleURL(libStyle),
			)

			// エリア別のフォルダに情報を保管
			switch l.Area {
			case 1:
				locations1.Add(placemark)
			case 2:
				locations2.Add(placemark)
			case 3:
				locations3.Add(placemark)
			case 4:
				locations4.Add(placemark)
			case 5:
				locations5.Add(placemark)
			case 6:
				locations6.Add(placemark)
			case 7:
				locations7.Add(placemark)
			}

			// カウント
			log.Info(strconv.Itoa(i+1) + "/" + shopSum + " done - " + l.Name)
		}
	})

	// フォルダ内のKMLを使って一気にKMLを作成
	result := kml.KML(
		kml.Document(
			kml.Name("WonderlandWars設置店舗"),
			kml.Description("公式のマップ情報を定期的に取得してプロットしています。<br>作者: @legnoh<br><br>図書館アイコン:ライブラリーあり<br>拠点アイコン:ライブラリーなし"),
			kml.SharedStyle(
				"icon-1526-A52714",
				kml.IconStyle(
					kml.Color(color.RGBA{R: 105, G: 27, B: 14, A: 0}),
					kml.Scale(1),
					kml.Icon(
						kml.Href("http://www.gstatic.com/mapspro/images/stock/503-wht-blank_maps.png"),
					),
				),
			),
			kml.SharedStyle(
				"icon-1598-0288D1",
				kml.IconStyle(
					kml.Color(color.RGBA{R: 2, G: 88, B: 209, A: 0}),
					kml.Scale(1),
					kml.Icon(
						kml.Href("http://www.gstatic.com/mapspro/images/stock/503-wht-blank_maps.png"),
					),
				),
			),
			kml.SharedStyle(
				"icon-1881-0f9d58",
				kml.IconStyle(
					kml.Color(color.RGBA{R: 15, G: 157, B: 58, A: 0}),
					kml.Scale(1),
					kml.Icon(
						kml.Href("http://www.gstatic.com/mapspro/images/stock/503-wht-blank_maps.png"),
					),
				),
			),
			kml.LabelStyle(
				kml.Scale(1),
			),
			kml.Schema(
				"extendInfomation",
				"extendInfomation",
				kml.SimpleField("住所", "string"),
				// kml.SimpleField("店舗詳細情報", "string"), googleMapのExtendDataのURLエスケープがバグってるのでこちらは出さないようにする
				kml.SimpleField("ランキング", "string"),
				kml.SimpleField("ランキング結果(5~1位)", "string"),
				kml.SimpleField("ライブラリ設置", "string"),
			),
			locations1,
			locations2,
			locations3,
			locations4,
			locations5,
			locations6,
			locations7,
		),
	)

	// ファイル書き込み
	file, err := os.Create(filePath)
	if err != nil {
		log.Error("KML output failed. Finish...")
	}
	result.WriteIndent(file, "", "  ")
	log.Info("KML output successfull.")
}
