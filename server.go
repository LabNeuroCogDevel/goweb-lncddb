package main
// http://www.darul.io/post/2015-07-22_go-lang-simple-reverse-proxy

import (
  //"io"
  "net/http"
  "net/http/httputil"
  "net/url"
  "os"
  "log"
  "strconv"
)

var serveport string = "3001"
var proxyuri *url.URL

/* 
 * Proxy type and 'handle' for HTTP
 * */
type Proxy struct{
  proxy *httputil.ReverseProxy
}
func (p *Proxy)  handle (w http.ResponseWriter, r *http.Request) {
    //http.HandlerFunc 
    p.proxy.ServeHTTP(w,r)
}

func usage(){
 log.Fatalf("USAGE: %s proxyuri [serverport=%s]",os.Args[0],serveport)
 //os.Exit(1)
}


func main(){
  // must provide other server. can provide port to serve on
  switch len(os.Args) {
   case 2:
    // set proxyuri, always done not here
   case 3:
    sp,err  := strconv.Atoi(os.Args[2])
    if(err != nil) { usage() }
    serveport = strconv.Itoa(sp)

   default:
    usage()
  }

  proxyuri,err := url.Parse(os.Args[1])
  if(err != nil){
   usage()
  }

  // *httputil.ReverseProxy
  //
  p := &Proxy{ proxy: httputil.NewSingleHostReverseProxy(proxyuri) }
  //proxy := httputil.NewSingleHostReverseProxy(proxyuri)

  log.Printf("Serving %s on %s",proxyuri,serveport)

  http.HandleFunc("/db",p.handle)
  log.Fatal(http.ListenAndServe(":"+serveport,nil))


}
