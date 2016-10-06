package mysqlbind

type Item struct {
	Namespace string `migu:"default:,pk,size:64"`
	Key       string `migu:"default:,pk,size:64"`
	Value     string
}
