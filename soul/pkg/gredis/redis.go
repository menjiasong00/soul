package gredis

import (
	"encoding/json"
	"time"
	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

var GRedis RediGo

type RediGo struct {
	RedisConn *redis.Pool
	DbName string
	ExpireTime
}

func init(){
	Setup()
}

// Setup 初始化连接池
func Setup() error {
	RediGo.RedisConn = &redis.Pool{
		MaxIdle:     "setting.RedisSetting.MaxIdle",
		MaxActive:   "setting.RedisSetting.MaxActive",
		IdleTimeout: "setting.RedisSetting.IdleTimeout",
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "setting.RedisSetting.Host")
			if err != nil {
				return nil, err
			}
			if "setting.RedisSetting.Password" != "" {
				if _, err := c.Do("AUTH", "setting.RedisSetting.Password"); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// New 返回gredis实例
func New() *RediGo {
	return &GRedis
}

// Db 选择库
func (h *RediGo) Db(name string) *RediGo {
	h.DbName = name
	return h
}

// SetExpireTime 设置过期时间
func (h *RediGo) SetExpireTime(eTime int) *RediGo {
	h.ExpireTime = eTime
	return h
}


// Set 设置缓存
func (h *RediGo) Set(key string, data interface{}) error {
	conn := h.RedisConn.Get()

	_,err := conn.Do("SELECT",h.DbName)
	if err != nil {
		return err
	}
	
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, h.ExpireTime)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func (h *RediGo) Exists(key string) bool {
	conn := h.RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func (h *RediGo) Get(key string) ([]byte, error) {
	conn := h.RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := h.RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func (h *RediGo) LikeDeletes(key string) error {
	conn := h.RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

//err := gredis.New().Db("xx").SetByFunc("KEYNAME",&user,func(interface{}) {
//  	db.Find(&user)
//})

//Callback
type CacheFunc func(interface{})

//SetByFunc 
func (h *RediGo) SetByFunc(key string, scan interface{},selectFunc CacheFunc) error {
	val,err:= h.Get(key)
	if len(val)>2 {
		//有缓存走缓存
		err = json.Unmarshal([]byte(val),&scan)
		if err != nil {
			return err
		}
	} else {
		//没缓存，执行selectFunc，然后把结果生成缓存
		selectFunc(&scan)
		err = h.Set(key,scan)
		if err != nil {
			return err
		}
	}
	return nil 
}



