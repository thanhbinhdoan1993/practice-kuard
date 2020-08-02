// Package sitedate provides static data compiled in via go-bindata.
package sitedata

//go:generate go-bindata -pkg $GOPACKAGE -prefix ../../sitedata/ -o bindata.go ../../sitedata/...
