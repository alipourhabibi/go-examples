package settings

import (
        "gopkg.in/gcfg.v1"
)

type App struct {
        Items AppItems `gcfg:"app"`
}
type AppItems struct {
        JwtAccess string `gcfg:"AccessSecret"`
        JwtRefresh string `gcfg:"RefreshSecret"`
}

type Postgresql struct {
        Items PostgresqlItems `gcfg:"postgresql"`
}
type PostgresqlItems struct {
        Host string `gcfg:"Host"`
        Port int `gcfg:"Port"`
        User string `gcfg:"User"`
        Password string `gcfg:"Password"`
        DBName string `gcfg:"DBName"`
}

type Server struct {
        Items ServerItems `gcfg:"server"`
}
type ServerItems struct {
        RunMode string `gcfg:"RunMode"`
        Port string `gcfg:"Port"`
}

var AppSettings App
var PostgresqlSettings Postgresql
var ServerSetting Server

func SetUp() {
        gcfg.ReadFileInto(&AppSettings, "conf/app.ini")
        gcfg.ReadFileInto(&PostgresqlSettings, "conf/app.ini")
        gcfg.ReadFileInto(&ServerSetting, "conf/app.ini")
}

