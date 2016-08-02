# 

## dependencies
```
git submodule foreach git pull
go get github.com/LabNeuroCogDevel/ldapauth
```

## local devel notes

GOPATH="$GOPATH:$HOME/src/dbexperements/go/" go test
GOPATH="$GOPATH:$HOME/src/dbexperements/go/" go build && ./web

