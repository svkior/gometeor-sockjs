package meteorMethods

type MeteorMethod struct {
	name string
	f    func(params interface{})
}

func Create(name string, f func(params interface{})) MeteorMethod {
	return MeteorMethod{name: name, f: f}
}

func (m *MeteorMethod) NameEquals(name string) bool {
	if m.name == name {
		return true
	} else {
		return false
	}
}

func (m *MeteorMethod) CallMethod(params interface{}) {
	m.f(params)
}
