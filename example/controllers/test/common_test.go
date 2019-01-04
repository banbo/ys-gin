package test

import (
	"net/http"

	ysGin "github.com/banbo/ys-gin"

	"github.com/banbo/ys-gin/example/router"
	"github.com/banbo/ys-gin/test"
)

var (
	lastUID string
)

// 登录cookie
var LOGIN_COOKIE = &http.Cookie{
	Name:  "mysession",
	Value: "MTUzODEyNzAzOXxOd3dBTkVzMFVrRldORmxLVFVjMVYwMU5OVU5VUmpkUVEwbEVSVWhhVWxGWVNGTlJORU5ZUjFsVlFVNVZRVkJYUXpSU1QxcFJVMEU9fNJWY6detApZ2ZL5MVVqm5wa4Nv9hKQlb1wCkXpDgrVy",
}

func getAPIClient() *test.APIClient {
	app := ysGin.NewApp("../../test.conf")
	router.Init(app.GinEngine)

	return test.NewAPIClient(app.GinEngine)
}
