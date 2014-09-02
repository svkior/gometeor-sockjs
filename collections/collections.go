package collections

type Collectioner interface {
	GetAllJSON() (s chan string)
}
