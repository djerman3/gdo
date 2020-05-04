#!/bin/bash
go-bindata -fs  -pkg staticdata -o staticdata/staticdata.go  -prefix "static/" static/ 
go-bindata -debug -pkg data -o data/data.go  templates/
