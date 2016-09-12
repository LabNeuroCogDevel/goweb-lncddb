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
  //"ldapauth"
  "github.com/LabNeuroCogDevel/ldapauth"
  //"github.com/LabNeuroCogDevel/lncdcal"
  "github.com/dgrijalva/jwt-go"

)

var serveport string = "3001"
var proxyuri *url.URL

//////////////////////
// Proxy type and 'handle' for HTTP
//////////////////////

type Proxy struct{
  proxy *httputil.ReverseProxy
  proxyroute string
}
func (p *Proxy)  handle (w http.ResponseWriter, r *http.Request) {
    //http.HandlerFunc 
    r.URL.Path = strings.Replace(r.URL.Path,p.proxyroute,"",1)
    //TODO: if we are adding or removing a visit
    // modify the google calendar
    log.Printf(r.URL.Path)
    p.proxy.ServeHTTP(w,r)
}

func login (w http.ResponseWriter, r *http.Request, sstr string,dbrole string) {
    r.ParseForm()

    if ( len(r.Form["user"])==0 || len(r.Form["pass"])==0 ) {
       http.Error(w, "\"bad input\"", http.StatusNotAcceptable)
       return
    }

    // check ldap and hardcoded test
    isauth := false
    if ( r.Form["user"][0] == "test" &&  r.Form["pass"][0] == "test" ){
     isauth = true
    } else {
     isauth = ldapauth.IsAuth(r.Form["user"][0],r.Form["pass"][0])
    }

    // return authentication
    if ! isauth {
       http.Error(w, "\"Not authorized\"", http.StatusUnauthorized)
       return
    } else {
       token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
             "user": r.Form["user"][0],
             "role": dbrole })
       tokenString, err := token.SignedString([]byte(sstr))
       if err == nil {
         fmt.Fprintf(w, "\"%s\"", tokenString)
       }
       fmt.Println(err)
    }
}

//////////////////////

func main(){

  serveport   := flag.Int   ("port"  ,3001,                 "port to run proxy server on")
  proxyuristr := flag.String("uri"   ,"http://0.0.0.0:3000","uri of server to be proxied")
  proxyroute  := flag.String("ppath", "db/",                "db route")
  staticpath  := flag.String("static","./elm-server-interface","filesystem path to static web assests")
  loginroute  := flag.String("login","login",               "where to send auth string")
  sstr        := flag.String("secret","secret",             "hash for jwt")
  dbrole      := flag.String("dbrole","lncd",               "database user/role for token")
  flag.Parse()

  proxyuri,err := url.Parse(*proxyuristr)
  if(err != nil){ log.Fatalf("%v",err) }

  // add ":" to port to get eg ":3001"
  port := fmt.Sprintf(":%d",*serveport)
  log.Printf("Serving %s on %s, proxy @ '/%s', login @ '/%s'",proxyuri,port,*proxyroute, *loginroute)


  // tried to do without type+handler
  // *httputil.ReverseProxy
  //proxy := httputil.NewSingleHostReverseProxy(proxyuri)

  p := &Proxy{ proxy: httputil.NewSingleHostReverseProxy(proxyuri),
              proxyroute: *proxyroute }



  // login
  // use the provided -login path to make a token with the user and -dbuser for postgrest
  loginsstr:=func(w http.ResponseWriter, r *http.Request) { login(w,r,*sstr,*dbrole) }
  http.HandleFunc("/"+*loginroute,loginsstr)

  // database proxy
  http.HandleFunc("/"+*proxyroute,p.handle)
  // static assests
  http.Handle("/",http.FileServer(http.Dir(*staticpath)))

  // all done
  log.Fatal(http.ListenAndServe(port,nil))


}
