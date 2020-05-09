package test

import (
		"fmt"
		"github.com/WebGameLinux/cms/dto/services"
		"testing"
)

func TestGetVerificationService(t *testing.T) {
		var id, b64, err = services.GetVerificationService().Image().Generate()
		id, b64, err = services.GetVerificationService().Mobile("15975798646").Generate()
		fmt.Println(id)
		fmt.Println(b64)
		fmt.Println(err)

}
