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
	for i := 0; i < len(ms.subs); i++ {
		if ms.subs[i] == name {
			return true
		}
	}
	return false
}

func (ms *MeteorSession) Subscribe(name string) {
	if ms.IsSubcribed(name) == false {
		ms.subs = append(ms.subs, name)
	}

}
