package gmongo

import (
    "fmt"
    "net/url"
    "strings"
)

type MongoConfig struct {
    User           string   `json:"User"     yaml:"User"`
    Password       string   `json:"Password" yaml:"Password"`
    Addr           []string `json:"Addr"     yaml:"Addr"`
    Database       string   `json:"Database" yaml:"Database"`
    ConnectTimeout int      `json:"ConnectTimeout" yaml:"ConnectTimeout"`
}

// URI mongo connect with String
func (mc *MongoConfig) URI() string {
    return fmt.Sprintf("mongodb://%v:%v@%v/%v", mc.User, url.QueryEscape(mc.Password), strings.Join(mc.Addr, ","), mc.Database)
}
