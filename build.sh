#!/bin/bash
# BUILD gdo
#

go-bindata -fs  -pkg staticdata -o staticdata/staticdata.go  -prefix "static/" static/ 
#go-bindata -debug -pkg data -o data/data.go  templates/
go-bindata -pkg data -o data/data.go  templates/

go mod tidy 
go build -o ~/tmp/gdoserver cmd/main.go 
sudo killall gdoserver 
sudo cp ~/tmp/gdoserver ~/bin/gdoserver 
#sudo ~/bin/gdoserver &
sudo rm ~/tmp/gdoserver
