// Code generated by gonectl. DO NOT EDIT.
package router

import "github.com/gone-io/gone/v2"

func init() {
	gone.
		Load(&authRouter{}).
		Load(&tokenParser{}).
		Load(&pubRouter{})
}
