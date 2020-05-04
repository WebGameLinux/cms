package setter

type Setter interface {
		SetValue(interface{}, interface{})
}

func GetAnySetter(v interface{}) Setter {
		switch v.(type) {

		}
		panic("impl setter")
}
