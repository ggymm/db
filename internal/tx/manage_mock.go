package tx

type mockManager struct {
}

func (m *mockManager) Close() {
	//TODO implement me
	panic("implement me")
}

func (m *mockManager) Begin() uint64 {
	//TODO implement me
	panic("implement me")
}

func (m *mockManager) Abort(tid uint64) {
	//TODO implement me
	panic("implement me")
}

func (m *mockManager) Commit(tid uint64) {
	//TODO implement me
	panic("implement me")
}

func (m *mockManager) IsActive(tid uint64) bool {
	//TODO implement me
	panic("implement me")
}

func (m *mockManager) IsCommitted(tid uint64) bool {
	//TODO implement me
	panic("implement me")
}

func (m *mockManager) IsAborted(tid uint64) bool {
	//TODO implement me
	panic("implement me")
}
