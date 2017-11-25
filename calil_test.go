package calil

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	SetAppKeyFromEnvironmentVariable()
	os.Exit(m.Run())
}

func TestSearchLibrary(t *testing.T) {
	limit := 10
	libs, err := SearchLibrary("京都府", "京都市", "", "", limit, http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(libs) == 0 {
		t.Fatal("図書館がひとつも見つかりませんでした。")
	}
	if len(libs) != limit {
		t.Errorf("期待した数（%d）の図書館が得られませんでした", limit)
	}
	t.Log(libs)
}

func TestCheckBooks(t *testing.T) {
	isbn := "4834000826"
	pref := "Kyoto_Kyoto"
	checkBooksResult, err := CheckBooks(isbn, pref, true, http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	if checkBooksResult.Books[isbn][pref].ReserveURL == "" {
		t.Fatalf("ReserveURL が取得できませんでした。")
	}
	t.Log(checkBooksResult)
}

func TestContinueCheckBooks(t *testing.T) {
	isbn := "4834000826"
	pref := "Aomori_Pref"
	checkBooksResult, err := CheckBooks(isbn, pref, true, http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	for {
		time.Sleep(time.Second * 2)
		// checkBooksResult.Continue をチェックせず ContinueCheckBooks を呼び出す。
		// 既に完了している session に対しても蔵書検索は有効。
		checkBooksResult, err = ContinueCheckBooks(checkBooksResult.Session, http.DefaultClient)
		if err != nil {
			t.Fatal(err)
		}
		if checkBooksResult.Continue == 0 {
			break
		}
	}
	if checkBooksResult.Books[isbn][pref].ReserveURL == "" {
		t.Fatalf("ReserveURL が取得できませんでした。")
	}
	t.Log(checkBooksResult)
}
