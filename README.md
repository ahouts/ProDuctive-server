# ProDuctive Server
## Install
access to a database on a MySql server is required

get and build the program
```bash
go get github.com/ahouts/ProDuctive-server
cd $GOPATH/src/github.com/ahouts/ProDuctive-server
go build .
```
choose a directory to serve from

move swagger-dist, ProDuctive-server, and config-example.json binary to that directory
```bash
cp -r ./swagger-dist ./ProDuctive-server <install dir>
cp ./config-example.json <install dir>/config.json
```
cd to that directory and edit the configuration
```bash
cd <install dir>
vi config.json
```
get help from the executable
```bash
./ProDuctive-server help
```
start the server
```bash
./ProDuctive-server serve
```
