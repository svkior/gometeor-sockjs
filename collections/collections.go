package collections

type Collectioner interface {
	// Subscribe to collections change
	SubscribeChan() (s chan string)
	// Insert new document from Client
	Insert(doc interface{}) string
	// Delete old document from Client
	Remove(params interface{}) string
	// Update document from Client
	Update(params interface{}) string
	// Get All Collections
	GetAllJSON() (s chan string)
}
