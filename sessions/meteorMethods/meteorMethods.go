package meteorMethods

type MeteorMethod struct {
	name string
	f    func(params interface{}) string
}

func Create(name string, f func(params interface{}) string) MeteorMethod {
	return MeteorMethod{name: name, f: f}
}

func (m *MeteorMethod) NameEquals(name string) bool {
	if m.name == name {
		return true
	} else {
		return false
	}
}

func (m *MeteorMethod) CallMethod(params interface{}) string {
	return m.f(params)
}
