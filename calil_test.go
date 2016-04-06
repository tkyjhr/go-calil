package calil

import "testing"

func TestSearchLibrary(t *testing.T) {
	SetAppKeyFromEnvironmentVariable()
	limit := 10
	libs, err := SearchLibrary("京都府", "京都市", "", "", limit);
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(libs) == 0 {
		t.Fatal("図書館がひとつも見つかりませんでした。")
	}
	if len(libs) != limit {
		t.Errorf("期待した数（%d）の図書館が得られませんでした", limit)
	}
	t.Log(libs)
}