package testhelper

import (
	"github.com/stretchr/testify/suite"
)

type MockTestSuite[T any] struct {
	TestObjectGenerator func() T
	suite.Suite
	TestO T
}

func (s *MockTestSuite[T]) BeforeTest(suiteName, testName string) {
	s.TestO = s.TestObjectGenerator()
}
