// Package calil はカーリル図書館API (https://calil.jp/doc/api_ref.html) が提供する機能を利用するための API を提供します。
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

// BookStatus は本の検索状態を表す構造体です。
//    status : システムIDに対して、検索状態を示します。OK, Cache, Running, Errorの4つをとります。Cache は OK と同じですが内部的にキャッシュ結果が利用されています。
//    reserveurl : 予約ページのアドレスです。
//    libkey : システムIDに紐尽く図書館のキーの配列です。図書館館のキー毎に貸出状況（「貸出中」、「貸出可」など）を値として持ちます。蔵書がない場合は、図書館キー自体が配列に含まれず空になります。
// libkey の値となる貸出状況は、貸出可、 蔵書あり、 館内のみ、 貸出中、 予約中、 準備中、 休館中、 蔵書なし、の8つに分類されます。
type BookStatus struct {
	Status     string            `json:"status"`
	ReserveUrl string            `json:"reserveurl"`
	LibKey     map[string]string `json:"libkey"`
}

// CheckBooksResult は蔵書情報の問い合わせ結果を受け取るための構造体です。
//     session : 検索に時間がかかる場合に、セッションが文字列として返ります。
//     books : 本毎にどの図書館に蔵書があるか返ります
//     continued : 0（偽）または1（真）が返ります。1の場合は、まだすべての取得が完了していないことを示します。
// books は isbn を1つめのキー、systemid を2つめのキーにとる map の map です。
type CheckBooksResult struct {
	Session  string `json:"session"`
	Books    map[string]map[string]BookStatus
	Continue int `json:"continue"`
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

	url := "https://api.calil.jp/library" + "?" + values.Encode()

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

// CheckBooks は指定した図書館システムに対して蔵書の有無と貸出状況を問い合わせます。
//     isbn : 書籍のISBNを指定します。カンマ区切りで複数指定できます。例「4834000826」
//     systemid : 図書館のシステムIDを指定します。カンマ区切りで複数指定できます。例「Aomori_Pref」
func CheckBooks(isbn string, systemid string) (map[string]map[string]BookStatus, error) {
	values := url.Values{}
	values.Add("isbn", isbn)
	values.Add("systemid", systemid)
	values.Add("format", "json")
	values.Add("callback", "")

	url := "https://api.calil.jp/check" + "?" + values.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// ({"session": "..." ...., "continue":0}); のような形で返ってくるので、json.Unmarshal の妨げとなる先頭の ( と末尾の ); を削除
	data = data[1 : len(data) - 2]

	var checkBooksResult CheckBooksResult

	if json.Unmarshal(data, &checkBooksResult) != nil {
		return nil, err
	}

	// TODO: continue = 1 だった場合の対応

	return checkBooksResult.Books, nil

}
