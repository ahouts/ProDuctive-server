go get -u github.com/jteeuwen/go-bindata/...
go-bindata -o migrations/bin-mig.go -pkg migrations migrations/*.sql

