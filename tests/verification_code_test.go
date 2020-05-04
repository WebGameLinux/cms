package test

import (
		"fmt"
		"github.com/WebGameLinux/cms/dto/repositories"
		"github.com/WebGameLinux/cms/dto/services"
		"testing"
		"time"
)

func TestGetVerificationService(t *testing.T) {
		var id, b64, err = services.GetVerificationService().Image().Generate()
		id, b64, err = services.GetVerificationService().Mobile("15975798646").Generate()
		fmt.Println(id)
		fmt.Println(b64)
		fmt.Println(err)
		repositories.GetCacheManager().Store(repositories.DriverStoreRedis).Put("mysql", time.Now().Unix(), 3*time.Minute)
		repositories.GetCacheManager().Store(repositories.DriverStoreFile).Put("file123", time.Now().Unix(), 3*time.Minute)

}
