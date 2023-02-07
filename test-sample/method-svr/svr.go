package main

import (
	"chreddal/pkgs/logger"
	"net/http"
)

var log = logger.StdLog.Specific("TestSrv")

func main() {
	http.ListenAndServe("127.0.0.1:8989", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Info("Method is ", req.Method)
		rw.WriteHeader(200)
	}))
}
