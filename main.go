package main

import (
	"log"
	"net/http"

	"github.com/bitCaskKV/global"
	"github.com/bitCaskKV/pkg/setting"
)

// A simple Router
func init() {
	_ = setupSetting()
}

func main() {
	serverMux := http.NewServeMux()
	//Get
	serverMux.HandleFunc("/v1/db", func(rw http.ResponseWriter, r *http.Request) {
		// r.Form.Get("Method")
		_, _ = rw.Write([]byte("HelloWorld"))
	})

	err := http.ListenAndServe(global.ServerSettingS.Addr+":"+global.ServerSettingS.HttpPort, serverMux)
	if err != nil {
		log.Fatalf("Run Server Err:%v", err)
	}
	return
}

func setupSetting() error {
	setting, _ := setting.NewSetting()
	_ = setting.ReadSection("Server", &global.ServerSettingS)
	return nil
}
