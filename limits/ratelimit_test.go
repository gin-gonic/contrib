package limits

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"text/template"
	"time"
)

var (
	params = struct {
		Size int
		Port int
	}{10, 9388}

	codeFile  = "/tmp/ratelimit_test_server.go"
	serverURL string
)

func init() {
	tmpl := template.Must(template.ParseFiles("test_server.tmpl"))
	fp, err := os.Create(codeFile)
	if err != nil {
		panic(fmt.Errorf("can't open %s - %s", codeFile, err))
	}
	err = tmpl.Execute(fp, params)
	if err != nil {
		panic(fmt.Errorf("can't create %s - %s", codeFile, err))
	}
	serverURL = fmt.Sprintf("http://localhost:%d", params.Port)
}

func waitForServer() error {
	timeout := 30 * time.Second
	ch := make(chan bool)
	go func() {
		for {
			_, err := http.Post(serverURL, "text/plain", nil)
			if err == nil {
				ch <- true
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	select {
	case <-ch:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("server did not reply after %v", timeout)
	}

}

func runServer() (*exec.Cmd, error) {
	cmd := exec.Command("go", "run", codeFile)
	cmd.Start()
	if err := waitForServer(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func doPost(val string) (*http.Response, error) {
	cmd, err := runServer()
	if err != nil {
		return nil, err
	}
	defer cmd.Process.Kill()

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "big=%s", val)
	return http.Post(serverURL, "application/x-www-form-urlencoded", &buf)
}

func TestRateLimiterOK(t *testing.T) {
	resp, err := doPost("abc")
	if err != nil {
		t.Fatalf("error posting - %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("bad status - %d", resp.StatusCode)
	}
}

func TestRateLimiterOver(t *testing.T) {
	resp, err := doPost("abcdefghijklmnop")
	if err != nil {
		t.Fatalf("error posting - %s", err)
	}
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("bad status - %d", resp.StatusCode)
	}
}
