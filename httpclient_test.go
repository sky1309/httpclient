package httpclient

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/sky1309/log"
)

type testPingResp struct {
	Now int64 `json:"now"`
}

type testUpdateReq struct {
	ID int `json:"id"`
}

type testUpdateResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`

	Data int
}

func testServe() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Info("ping")

		var resp testPingResp
		resp.Now = time.Now().Unix()

		b, _ := json.Marshal(resp)
		w.Write(b)

	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		log.Info("update, headers=%v", r.Header)

		var req testUpdateReq
		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &req)

		var resp testUpdateResp
		resp.Code = 1234
		resp.Msg = "ok"
		resp.Data = req.ID

		data, _ := json.Marshal(resp)
		w.Write(data)

	})

	address := ":8000"
	log.Info("listen and serve %s", address)
	http.ListenAndServe(address, nil)
}

func TestHttpClient(t *testing.T) {
	go testServe()

	var wg sync.WaitGroup
	wg.Add(1)
	time.AfterFunc(time.Second, func() {
		defer wg.Done()

		respPing, err := GetJson[testPingResp]("http://:8000/ping", nil)
		log.Info("respPing=%v, err=%v", respPing, err)

		opt := NewOptions().SetHeader("Content-Type", "application/json").SetHeader("Token", "foo")
		reqUpdate := &testUpdateReq{ID: rand.Intn(1024)}
		respUpdate, err := PostJson[testUpdateResp]("http://:8000/update", reqUpdate, opt)
		log.Info("reqUpdate=%v, respUpdate=%v, err=%v", reqUpdate, respUpdate, err)
	})

	wg.Wait()
	log.Info("quit")
}
