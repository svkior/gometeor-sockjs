package meteorSession

// Сессия с клиентом Meteor
type MeteorSession struct {
	id   string
	subs []string
}

func Create(id string) MeteorSession {
	return MeteorSession{id: id}
}

func (s *MeteorSession) GetId() string {
	return s.id
}

func (ms *MeteorSession) IsSubcribed(name string) bool {
	//fmt.Println("Hi")
	for _, c := range ms.subs {
		//		fmt.Println(c)
		if c == name {
			return true
		}
	}
	//fmt.Println("Return false")
	return false
}

func (ms *MeteorSession) Subscribe(name string) {
	if ms.IsSubcribed(name) == false {
		ms.subs = append(ms.subs, name)
	}

}
