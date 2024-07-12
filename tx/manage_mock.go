package tx

type mockManager struct {
}

func NewMockManage() Manage {
	return new(mockManager)
}

func (m *mockManager) Close() {
}

func (m *mockManager) Begin() uint64 {
	return Super
}

func (m *mockManager) Commit(_ uint64) {
}

func (m *mockManager) Rollback(_ uint64) {
}

func (m *mockManager) IsActive(_ uint64) bool {
	return false
}

func (m *mockManager) IsCommitted(_ uint64) bool {
	return false
}

func (m *mockManager) IsRolledBack(_ uint64) bool {
	return false
}
