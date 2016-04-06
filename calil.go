package calil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// AppKey はカーリル図書館API のアプリケーションキーです。API を呼ぶ前に有効な値を設定する必要があります。
var AppKey string

// Library は図書館の情報を表す構造体です。
type Library struct {
	SystemID        string `json:"systemid"`
	SystemName      string `json:"systemname"`
	LibKey          string `json:"libkey"`
	LibID           string `json:"libid"`
	ShortName       string `json:"short"`
	FormalName      string `json:"formal"`
	URLPc           string `json:"url_pc"`
	Address         string `json:"address"`
	Prefecture      string `json:"pref"`
	City            string `json:"city"`
	PostalNumber    string `json:"post"`
	TelephoneNumber string `json:"tel"`
	GeoCode         string `json:"geocode"`
	Category        string `json:"category"`
	Image           string `json:"image"`
	Distance        string `json:"distance"`
}

// SetAppKeyFromEnvironmentVariable は AppKey を環境変数「CALIL_APP_KEY」の値に設定します。
func SetAppKeyFromEnvironmentVariable() {
	AppKey = os.Getenv("CALIL_APP_KEY")
}

// SearchLibrary は指定した条件で図書館を検索します。
//     pref : 都道府県を指定します。例「青森県」
//     ciry : 市区町村を指定します。このパラメータはprefとセットで利用します。例「青森市」
//     systemid : 図書館のシステムIDを指定します。例「Aomori_Pref」
//     geocode : 緯度、経度を指定します。例「136.7163027,35.390516」
//     limit : 図書館の取得件数を指定します。
// pref, systemid, geocode のいずれかは必ず指定する必要があります。
func SearchLibrary(pref string, city string, systemid string, geocode string, limit int) ([]Library, error) {

	values := url.Values{}
	values.Add("appkey", AppKey)
	if pref != "" {
		values.Add("pref", pref)
	}
	if city != "" {
		values.Add("city", city)
	}
	if systemid != "" {
		values.Add("systemid", systemid)
	}
	if geocode != "" {
		values.Add("geocode", geocode)
	}
	if limit != 0 {
		values.Add("limit", strconv.Itoa(limit))
	}

	values.Add("format", "json")
	values.Add("callback", "")

	url := "http://api.calil.jp/library" + "?" + values.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var libs []Library

	if json.Unmarshal(data, &libs) != nil {
		return nil, err
	}

	return libs, nil
}
