package pipelines

import (
	"github.com/stretchr/testify/mock"
)

// FileWriterMock is a mock implementation of writer interface.
type FileWriterMock struct {
	mock.Mock
}

// Needed for make sure that FileWriterMock implements writer.
var _ FileWriter = &FileWriterMock{}

// NewServiceMock creates new ServiceMock.
func NewWriterMock(t mock.TestingT) *FileWriterMock {
	m := &FileWriterMock{}
	m.Test(t)

	return m
}

// WriteString mocks method of FileWriterMock.
func (mock *FileWriterMock) WriteString(s string) error {
	args := mock.Called(s)
	return args.Error(0)
}

// Close mocks method of FileWriterMock.
func (mock *FileWriterMock) Close() {}
