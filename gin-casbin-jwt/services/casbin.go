package services

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func GetEnforcer(adapter *gormadapter.Adapter) (*casbin.Enforcer, error) {
	return casbin.NewEnforcer("conf/abac_model.conf", adapter)
}

func Enforce(sub, obj, act string, adapter *gormadapter.Adapter) (bool, error) {
	enforcer, err := GetEnforcer(adapter)
	if err != nil {
		panic(err)
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, err
	}
	ok, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func AddPolicy(username string, adapter *gormadapter.Adapter) {
	enforcer, err := GetEnforcer(adapter)
	if err != nil {
		panic(err)
	}
	enforcer.AddPolicy(username, "resource", "write")
}
