package webapp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/domains"
)

type apiEnvelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
}

func TestWebAppSmoke(t *testing.T) {
	tmp := t.TempDir()
	t.Chdir(tmp)

	if err := os.MkdirAll("config", 0o755); err != nil {
		t.Fatalf("mkdir config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "logo.png"), []byte("png"), 0o644); err != nil {
		t.Fatalf("write logo: %v", err)
	}

	db.InitDB("config/database")
	t.Cleanup(db.CloseDB)
	db.AddLegitDomainToDB(domains.LegitDomain{Domain: domains.Domain{Name: "example.com", Status: constants.INACTIVE}})
	db.AddTypoDomainToDB(domains.NewTypoDomain("examp1e.com", "example.com", "CO"))

	srv := NewServer()
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)

	assertStatus(t, get(t, ts.URL+"/"), http.StatusOK)
	assertStatus(t, get(t, ts.URL+"/logo.png"), http.StatusOK)

	resp := postJSON(t, ts.URL+"/api/domain/import", map[string]string{
		"domain": "example.com",
		"typos":  "",
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty import, got %d", resp.StatusCode)
	}

	assertOK(t, postEmpty(t, ts.URL+"/api/background/start"))
	assertOK(t, postEmpty(t, ts.URL+"/api/background/stop"))
	assertOK(t, postEmpty(t, ts.URL+"/api/background/run"))
	assertOK(t, postEmpty(t, ts.URL+"/api/check/all"))
	assertOK(t, get(t, ts.URL+"/api/overview"))
	assertOK(t, get(t, ts.URL+"/api/config"))
	assertOK(t, get(t, ts.URL+"/api/stats"))
	assertOK(t, get(t, ts.URL+"/api/expiration"))
	assertOK(t, get(t, ts.URL+"/api/reliable"))
	assertOK(t, get(t, ts.URL+"/api/domain/example.com/typos"))
	assertOK(t, get(t, ts.URL+"/api/domain/example.com/typos/page?page=1&size=20&q=exa"))
	assertOK(t, get(t, ts.URL+"/api/domain/example.com/typos/page?page=1&size=20&name=exa&algorithm=co&sort_by=date_check&sort_dir=desc"))
	assertOK(t, get(t, ts.URL+"/api/typodomain/examp1e.com/history"))
	assertOK(t, postEmpty(t, ts.URL+"/api/typodomain/examp1e.com/delete"))
	assertOK(t, postJSON(t, ts.URL+"/api/domain/import", map[string]interface{}{
		"domain":              "example.com",
		"typos":               "example-typo-one.com\n",
		"include_typodomains": false,
	}))
	assertOK(t, postJSON(t, ts.URL+"/api/presence", map[string][]string{
		"domains": []string{"example.com", "missing.local"},
	}))
}

func assertStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	defer resp.Body.Close()
	if resp.StatusCode != want {
		t.Fatalf("status: got %d want %d", resp.StatusCode, want)
	}
}

func assertOK(t *testing.T, resp *http.Response) {
	t.Helper()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d want 200", resp.StatusCode)
	}
	var env apiEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if !env.OK {
		t.Fatalf("expected ok=true, err=%s", env.Error)
	}
}

func get(t *testing.T, url string) *http.Response {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	return resp
}

func postEmpty(t *testing.T, url string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	return resp
}

func postJSON(t *testing.T, url string, payload interface{}) *http.Response {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	return resp
}
