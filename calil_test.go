package calil

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetAppKeyFromEnvironmentVariable()
	os.Exit(m.Run())
}

func TestSearchLibrary(t *testing.T) {
	limit := 10
	libs, err := SearchLibrary("京都府", "京都市", "", "", limit)
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
	pref := "Aomori_Pref"
	books, err := CheckBooks(isbn, pref)
	if err != nil {
		t.Fatal(err)
	}
	if books[isbn][pref].ReserveUrl == "" {
		t.Fatalf("ReserveURL が取得できませんでした。")
	}
	t.Log(books)
}
