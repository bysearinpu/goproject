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
ctr.s += "hellorrrrr"

fmt.Fprintf(c, "co2sddddddddd = %s\n", ctr.s)
}

func main() {
http.Handle("/counte", new(Counter))
log.Fatal("ListenAndServe: ", http.ListenAndServe(":80", nil))
}