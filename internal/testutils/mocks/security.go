package mocks

type MockPasswordHasher struct {
	HashFunc          func(password string) (string, error)
	CheckPasswordFunc func(hashPassword, password string) error
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	if m.HashFunc != nil {
		return m.HashFunc(password)
	}
	return "", nil
}

func (m *MockPasswordHasher) CheckPassword(hashPassword, password string) error {
	if m.CheckPasswordFunc != nil {
		return m.CheckPasswordFunc(hashPassword, password)
	}
	return nil
}
