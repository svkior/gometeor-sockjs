package collections

type Collection interface {
	GetAllJSON() (s chan string)
}
