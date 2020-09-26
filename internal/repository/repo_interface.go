package repository

type Repository interface {
	InsertLink()
	GetLink()
	Start(maxConn, acquireTimeout int) error
	Close()
}
