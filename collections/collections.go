package collections

type Collectioner interface {
	// Subscribe to collections change
	SubscribeChan() (s chan string)
	// Insert new document from Client
	Insert(doc interface{}) string
	// Get All Collections
	GetAllJSON() (s chan string)
}
