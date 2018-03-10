package main

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"net/http"
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

// 1地方ごとの店舗情報データ
type Area struct {
	Name string
	Pref []Prefacture
}

// 1県ごとの店舗情報データ
type Prefacture struct {
	Name  string
	Store []Store
}

// 1店舗ごとの店舗情報データ
type Store struct {
	Id   int
	Name string
	Add  string
	Lib  bool
}

// 1ストアのランキング情報
type StoreScore struct {
	Id      int
	Ranking []StoreScoreMonthly
}

// 1ストアの3ヶ月分のランキング情報
type StoreScoreMonthly struct {
	Name    string
	Updtime string
	Data    []Ranker
}

// ランカー情報
type Ranker struct {
	Rank  int
	Upd   int
	Name  string
	Cast  string
	Honor string
	Opera string
	Score int
}

// 1件の店舗情報に諸情報を加えた完成形データ
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
	locationURL     = "https://wonder.sega.jp/json/store-list.json"
	scoreRankingURL = "https://wonder.sega.jp/json/store-ranking-"
	shopURL         = "https://location.am-all.net/alm/shop?gm=43&sid="
	rankingURL      = "https://wonder.sega.jp/ranking/store/#!/store:"
	hostURL         = "https://wonderland-wars.net"
	gMapHostHead    = "//maps.googleapis.com/maps/api/staticmap?center="
	gMapHostFoot    = regexp.MustCompile("&markers=.*")
	iconImage       = "http://www.gstatic.com/mapspro/images/stock/503-wht-blank_maps.png"
	locations0      = kml.Folder(kml.Name("北海道・東北"))
	locations1      = kml.Folder(kml.Name("関東"))
	locations2      = kml.Folder(kml.Name("東海"))
	locations3      = kml.Folder(kml.Name("北信越"))
	locations4      = kml.Folder(kml.Name("近畿"))
	locations5      = kml.Folder(kml.Name("中国・四国"))
	locations6      = kml.Folder(kml.Name("九州・沖縄"))
	long            float64
	lat             float64
	libStyle        = "#icon-1664-0288D1/ "
	libSign         = "○"
	filePath        = "./result-" + strconv.FormatInt(time.Now().Unix(), 10) + ".kml"
	rankResult      = "0pt"
	rankNull        = "0pt"
	rank1st         = 0
	rank5th         = 0
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {

	// 4〜7時はメンテナンス中なので実行しない
	hour := time.Now().Hour()
	if hour == 4 || hour == 5 || hour == 6 {
		log.Error("現在メンテナンス中なので実行できません。7時以降に再度実行してください。")
		os.Exit(1)
	}

	// 店舗情報一覧を取得して構造体に分解
	resp, _ := http.Get(locationURL)
	defer resp.Body.Close()
	jsonBlob, _ := ioutil.ReadAll(resp.Body)
	var areas []Area
	json.Unmarshal(jsonBlob, &areas)

	// ストア情報を保持する配列を定義
	locations := make(map[int]*kml.CompoundElement)

	// すべてのAreaに対して以下の処理を行う
	for i := 0; i < len(areas); i++ {

		// エリア名とKML指定
		log.Info("# " + areas[i].Name)
		locations[i] = kml.Folder(kml.Name(areas[i].Name))

		// ランキング情報を取得する
		resp, _ := http.Get(scoreRankingURL + strconv.Itoa(i) + ".json")
		defer resp.Body.Close()
		jsonBlob, _ := ioutil.ReadAll(resp.Body)
		var storeScoresData []StoreScore
		json.Unmarshal(jsonBlob, &storeScoresData)

		storeScores := make(map[int]StoreScoreMonthly)

		// 配列に組み替え
		for x := 0; x < len(storeScoresData); x++ {
			storeScores[storeScoresData[x].Id] = storeScoresData[x].Ranking[0]
		}

		// すべての県に対して以下の処理を繰り返す
		for j := 0; j < len(areas[i].Pref); j++ {

			log.Info("## " + areas[i].Pref[j].Name)

			// すべての店舗に対して以下の処理を繰り返す
			for k := 0; k < len(areas[i].Pref[j].Store); k++ {

				log.Info("- " + areas[i].Pref[j].Store[k].Name)

				// 店舗IDから位置情報取得
				// ShopURLにアクセスしGoogleMapへのURLから緯度経度を取得する
				// 一斉アクセスを避けるため、sleepを入れて0.2rps程度になるように留める
				time.Sleep(1 * time.Second)
				shopPage, _ := goquery.NewDocument(shopURL + strconv.Itoa(areas[i].Pref[j].Store[k].Id))
				gMapURL, mapExists := shopPage.Find(".access_map").Attr("src")
				if mapExists {

					// 緯度と経度部分のみ抜き出す
					gMapURL = strings.Replace(gMapURL, gMapHostHead, "", 1)
					gMapURL = gMapHostFoot.ReplaceAllString(gMapURL, "")

					// KMLに投入するため、逆転させてfloat64に変換させる
					longLat := strings.Split(gMapURL, ",")
					long, _ = strconv.ParseFloat(longLat[1], 64)
					lat, _ = strconv.ParseFloat(longLat[0], 64)
				}

				// ランキングが空の場合、もしくは5位までない場合はダミーを入れる
				if storeScores[areas[i].Pref[j].Store[k].Id].Data == nil {
					rank1st = 0
					rank5th = 0
				} else if len(storeScores[areas[i].Pref[j].Store[k].Id].Data) < 5 {
					rank1st = storeScores[areas[i].Pref[j].Store[k].Id].Data[0].Score
					rank5th = 0
				} else {
					rank1st = storeScores[areas[i].Pref[j].Store[k].Id].Data[0].Score
					rank5th = storeScores[areas[i].Pref[j].Store[k].Id].Data[4].Score
				}

				// location構造体の定義
				l := Location{
					Name:       string(norm.NFKC.Bytes([]byte(areas[i].Pref[j].Store[k].Name))),
					Address:    string(norm.NFKC.Bytes([]byte(areas[i].Pref[j].Store[k].Add))),
					Area:       i,
					Lat:        lat,
					Long:       long,
					ShopURL:    shopURL + strconv.Itoa(areas[i].Pref[j].Store[k].Id),
					RankingURL: rankingURL + strconv.Itoa(areas[i].Pref[j].Store[k].Id),
					Rank1st:    strconv.Itoa(rank1st) + "pt",
					Rank5th:    strconv.Itoa(rank5th) + "pt",
					Library:    areas[i].Pref[j].Store[k].Lib,
				}

				// ライブラリの有無によってアイコンを変更
				if l.Library {
					libStyle = "#icon-1526-A52714"
					libSign = "○"
				} else {
					libStyle = "#icon-1598-0288D1"
					libSign = "×"
				}

				// 新店舗のみ個別のアイコンに変更し、ランキング情報を変更
				if l.Rank5th == rankNull && l.Rank1st == rankNull {
					log.Warn("新店舗があるようです！: " + l.Name)
					rankResult = "ランキングなし"
					libStyle = "#icon-1881-0f9d58"
				} else {
					rankResult = l.Rank5th + " 〜 " + l.Rank1st
				}

				// PlaceMarkに全ての情報を結合して配列にKMLに追加
				placemark := kml.Placemark(
					kml.Name(l.Name),
					kml.Description("店舗情報: "+l.ShopURL),
					kml.ExtendedData(
						kml.SchemaData(
							"#extendInfomation",
							kml.SimpleData("ライブラリ設置", libSign),
							kml.SimpleData("住所", l.Address),
							kml.SimpleData("ランキング", l.RankingURL),
							kml.SimpleData("ランキング結果(5〜1位)", rankResult),
						),
					),
					kml.StyleURL(libStyle),
					kml.Point(kml.Coordinates(kml.Coordinate{Lon: l.Long, Lat: l.Lat})),
				)
				locations[i].Add(placemark)
			}
		}
	}

	// フォルダ内のKMLを使って一気にKMLを作成
	result := kml.KML(
		kml.Document(
			kml.SharedStyle(
				"icon-1526-A52714",
				kml.IconStyle(
					kml.Color(color.RGBA{R: 105, G: 27, B: 14, A: 0}),
					kml.Scale(1),
					kml.Icon(
						kml.Href(iconImage),
					),
				),
			),
			kml.SharedStyle(
				"icon-1598-0288D1",
				kml.IconStyle(
					kml.Color(color.RGBA{R: 2, G: 88, B: 209, A: 0}),
					kml.Scale(1),
					kml.Icon(
						kml.Href(iconImage),
					),
				),
			),
			kml.SharedStyle(
				"icon-1881-0f9d58",
				kml.IconStyle(
					kml.Color(color.RGBA{R: 15, G: 157, B: 58, A: 0}),
					kml.Scale(1),
					kml.Icon(
						kml.Href(iconImage),
					),
				),
			),
			kml.LabelStyle(
				kml.Scale(1),
			),
			kml.Schema(
				"extendInfomation",
				"extendInfomation",
				kml.SimpleField("ライブラリ設置", "string"),
				kml.SimpleField("住所", "string"),
				kml.SimpleField("ランキング", "string"),
				kml.SimpleField("ランキング結果(5~1位)", "string"),
			),
			locations[0],
			locations[1],
			locations[2],
			locations[3],
			locations[4],
			locations[5],
			locations[6],
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
