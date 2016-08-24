package host

type ConfigFlagger struct {
	Data map[string]interface{}
}

func NewConfigFlagger(data map[string]interface{}) drivers.DriverOptions{
     return ConfigFlagger{Data: data}
}

func (this ConfigFlagger) String(key string) string {
	if value, ok := this.Data[key]; ok {
		return value.(string)
	}
	return ""
}

func (this ConfigFlagger) StringSlice(key string) []string {
	if value, ok := this.Data[key]; ok {
		return value.([]string)
	}
	return []string
}

func (this ConfigFlagger) Int(key string) int {
	if value, ok := this.Data[key]; ok {
		return value.(int)
	}
	return 0
}

func (this ConfigFlagger) Bool(key string) bool {
	if value, ok := this.Data[key]; ok {
		return value.(bool)
	}
	return false
}
