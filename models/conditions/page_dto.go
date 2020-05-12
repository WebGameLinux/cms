package conditions

type PageCondition struct {
		Page       int                    `json:"page"`
		Count      int                    `json:"count"`
		Conditions map[string]interface{} `json:"conditions"`
}

func NewPageCondition(args ...interface{}) *PageCondition {
		condition := new(PageCondition)
		condition.Page = 1
		condition.Count = 10
		condition.Conditions = make(map[string]interface{})
		argc := len(args)
		if argc == 3 {
				if v, ok := args[0].(int); ok {
						condition.Page = v
				}
				if v, ok := args[1].(int); ok {
						condition.Count = v
				}
				if m, ok := args[2].(map[string]interface{}); ok {
						condition.Conditions = m
				}
				return condition
		}
		if argc > 0 {
				isOk := false
				for i, v := range args {
						if n, ok := v.(int); ok {
								if argc <= 2 && i <= 1 && !isOk {
										condition.Page = n
										isOk = true
								} else {
										condition.Count = n
								}
						}
						if m, ok := v.(map[string]interface{}); ok {
								condition.Conditions = m
						}
				}
		}
		return condition
}
