package data

import "db/internal/txn"

type Manager interface {
	Close()

	Read(uid UID) (Item, bool, error)
	Insert(xid txn.TID, data []byte) (uint64, error)
}
