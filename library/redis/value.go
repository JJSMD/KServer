package redis

import (
	iface "KServer/library/kiface/iredis"
	"github.com/garyburd/redigo/redis"
)

type Value struct {
	Conn redis.Conn
}

func (v *Value) Get(key string) iface.IGetValue {
	return &GetValue{Conn: v.Conn, Key: key}
}
func (v *Value) Set(key string) iface.ISetValue {
	return &SetValue{Conn: v.Conn, Key: key}
}
func (v *Value) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	defer v.Conn.Close()
	return v.Conn.Do(commandName, args)
}
func (v *Value) Check(key string) bool {

	defer v.Conn.Close()
	_, err := v.Conn.Do("EXISTS", key)
	if err != nil {
		return false
	}
	return true
}

func (v *Value) Del(key string) (reply interface{}, err error) {
	defer v.Conn.Close()
	return v.Conn.Do("DEL", key, "SEX")
}
