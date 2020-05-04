package beego

func Onerror(err error) {
		if err != nil {
				GetLogger().Error(err.Error())
		}
}
