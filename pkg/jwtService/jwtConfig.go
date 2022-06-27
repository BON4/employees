package jwtService

import "time"

type jwtConfig struct {
	//Shoud be provided
	AccessKey []byte

	//Shoud be provided
	RefreshKey []byte

	AccessExpireTime  time.Duration
	RefreshExpireTime time.Duration
}

func NewJWTConfig(AccessKey, RefreshKey []byte, AccessExpireTime, RefreshExpireTime time.Duration) *jwtConfig {
	return &jwtConfig{
		AccessKey:         AccessKey,
		RefreshKey:        RefreshKey,
		AccessExpireTime:  AccessExpireTime,
		RefreshExpireTime: RefreshExpireTime,
	}
}
