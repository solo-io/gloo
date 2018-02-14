package storage

const (
	// Create is WatchOperation passed to the watch callback
	CreateOp = iota
	// Update is WatchOperation passed to the watch callback
	UpdateOp
	// Delete is WatchOperation passed to the watch callback
	DeleteOp
	// Error is WatchOperation passed to the watch callback
	ErrorOp
)

// WatchOperation is an operation that was detected by watcher and passed to the callbacl method
type WatchOperation int

// Item is used to represent all objects to be stored. Currently Upstream and VirtualHost are supported
type Item interface{}

// GetOptions is options object passed to Get method. Currently unused
type GetOptions struct{}

// ListOptions is options object passed to List method. Currently unused
type ListOptions struct{}

// WatchOptions is options object passed to Watch method. Currently unused
type WatchOptions struct{}

// Storage is interface to the storage backend
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
