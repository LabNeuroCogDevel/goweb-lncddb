# web wrapper for postgrest

Assume a postgresql database with login-only role `web` and select,update,delete role `lncd`:

```
postgrest -j abcd -m 300 postgres://postgres@localhost:5432/lncddb -a web
./web -secret abcd -dbrole lncd

# test
## should fail
curl localhost:3000/person |jq
## success?
token=$(curl  localhost:3001/login -d user=test -d pass=test )
curl localhost:3000/person -H "Range-Unit: items" -H "Range: 0-0" -H "Authorization: Bearer $token"|jq
```

## get (and dependencies)
```
git clone git@github.com:LabNeuroCogDevel/goweb-lncddb.git
git submodule update --init --recursive
git submodule foreach git pull
go get github.com/LabNeuroCogDevel/ldapauth
```

## usage
```
./web -help
```

> Usage of ./web:
>   -dbrole string
>        database user/role for token (default "lncd")
>   -login string
>     	where to send auth string (default "login")
>   -port int
>     	port to run proxy server on (default 3001)
>   -ppath string
>     	db route (default "db/")
>   -secret string
>     	hash for jwt (default "secret")
>   -static string
>     	filesystem path to static web assests (default "./elm-server-interface")
>   -uri string
>     	uri of server to be proxied (default "http://0.0.0.0:3000")

## usefulness

get a token to access database
```
curl  http://0.0.0.0:3001/login -d 'user=test' -d pass=test
```

## local devel notes

GOPATH="$GOPATH:$HOME/src/dbexperements/go/" go test
GOPATH="$GOPATH:$HOME/src/dbexperements/go/" go build && ./web

## related
* bitbucket.org/LabNeuroCogDevel/postgresql-export
* github.com/LabNeuroCogDevel/elm-server-interface
* github.com/LabNeuroCogDevel/ldapauth
