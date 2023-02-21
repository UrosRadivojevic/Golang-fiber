package hasher

import (
	"github.com/stretchr/testify/mock"
	"github.com/urosradivojevic/health/pkg/services/hasher"
)

type HasherMock struct {
	mock.Mock
}

var _ hasher.Interface = &HasherMock{}

func (h *HasherMock) Hash(password string) ([]byte, error) {
	args := h.Called(password)
	return args.Get(0).([]byte), args.Error(1)
}

func (h *HasherMock) Verify(password string, hashPassword []byte) bool {
	args := h.Called(password, hashPassword)
	return args.Bool(0)
}
