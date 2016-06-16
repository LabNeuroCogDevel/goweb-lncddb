package main
// http://www.darul.io/post/2015-07-22_go-lang-simple-reverse-proxy

import (
  //"io"
  "net/http"
  "net/http/httputil"
  "net/url"
  "log"
  "flag"
  "fmt"
  "strings"
 // "github.com/gorilla/mux" // we need regexp routes
)

var serveport string = "3001"
var proxyuri *url.URL

//////////////////////
// Proxy type and 'handle' for HTTP
//////////////////////

type Proxy struct{
  proxy *httputil.ReverseProxy
  proxypath string
}
func (p *Proxy)  handle (w http.ResponseWriter, r *http.Request) {
    //http.HandlerFunc 
    r.URL.Path = strings.Replace(r.URL.Path,p.proxypath,"",1)
    log.Printf(r.URL.Path)
    p.proxy.ServeHTTP(w,r)
}

//////////////////////

func main(){

  serveport   := flag.Int   ("port"  ,3001,                 "port to run proxy server on")
  proxyuristr := flag.String("uri"   ,"http://0.0.0.0:3000","uri of server to be proxied")
  proxypath   := flag.String("ppath", "db/",                 "db route")
  staticpath  := flag.String("static","./static",           "filesystem path to static web assests")
  flag.Parse()

  proxyuri,err := url.Parse(*proxyuristr)
  if(err != nil){ log.Fatalf("%v",err) }

  // add ":" to port to get eg ":3001"
  port := fmt.Sprintf(":%d",*serveport)
  log.Printf("Serving %s on %s, proxy @ '/%s'",proxyuri,port,*proxypath)


  // tried to do without type+handler
  // *httputil.ReverseProxy
  //proxy := httputil.NewSingleHostReverseProxy(proxyuri)

  p := &Proxy{ proxy: httputil.NewSingleHostReverseProxy(proxyuri),
              proxypath: *proxypath }


  //r:=mux.NewRouter()
  http.HandleFunc("/"+*proxypath,p.handle)
  http.Handle("/",http.FileServer(http.Dir(*staticpath)))
  log.Fatal(http.ListenAndServe(port,nil))


}
