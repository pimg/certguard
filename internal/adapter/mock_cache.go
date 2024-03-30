package adapter

type MockCache struct {
	filename    string
	fileContent []byte
}

func (m *MockCache) Write(filename string, fileContent []byte) error {
	m.filename = filename
	m.fileContent = fileContent

	return nil
}

func (m *MockCache) Read(filename string) ([]byte, error) {
	return m.fileContent, nil
}
