package main

import (
	"net/http"
)

// func TestMain(m *testing.T) {
// 	os.Exit(m.Run())
// }

type myHandler struct {

}

func (mh *myHandler) ServeHTTP(res http.ResponseWriter,req *http.Request) {

}