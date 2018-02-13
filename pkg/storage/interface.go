package storage

const (
	Create = iota
	Update
	Delete
	Error
)

type WatchOperation int

type Item interface{}

type GetOptions struct{}
type ListOptions struct{}
type WatchOptions struct{}

type Storage interface {
	Register(item Item) error
	Create(item Item) (Item, error)
	Update(item Item) (Item, error)
	Delete(item Item) error
	Get(item Item, getOptions *GetOptions) (Item, error)
	List(item Item, listOptions *ListOptions) ([]Item, error)
	Watch(item Item, watchOptions *WatchOptions, callback func(item Item, operation WatchOperation)) error
	WatchStop()
}
