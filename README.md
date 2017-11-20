# ProDuctive Server
## Install
choose a directory to install
install oracle OCI-8 BS, you will have to figure that one out

get and build the program
```bash
go get github.com/ahouts/ProDuctive-server
cd $GOPATH/src/github.com/ahouts/ProDuctive-server
go build .
```
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
