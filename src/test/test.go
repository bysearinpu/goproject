package main

import (

    "net/http"
    "fmt"
    "log"
)

type Counter struct {
n int
s string

}

func (ctr *Counter) ServeHTTP(c http.ResponseWriter, req *http.Request) {
fmt.Fprintf(c, "%08x\n", ctr)
ctr.n++
ctr.s += "hello"

fmt.Fprintf(c, "coun2222 = %s\n", ctr.s)
}

func main() {
http.Handle("/counter", new(Counter))
log.Fatal("ListenAndServe: ", http.ListenAndServe(":80", nil))
}