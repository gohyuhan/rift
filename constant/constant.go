package constant

const (
	APPNAME         = "rift"
	APPDBNAME       = "rift.db"
	APPSETTINGSNAME = "rift_settings.json"
)

// this will be injected during build
// exmaple) go build -ldflags "-X rift/constant.APPVERSION=v1.x.x" -o main
var APPVERSION = "v0.1.0"
