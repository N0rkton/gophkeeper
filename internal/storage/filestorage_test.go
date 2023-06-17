package storage

import (
	"gophkeeper/internal/datamodels"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_Login(t *testing.T) {
	s := NewMemoryStorage()
	Init()
	id, err := s.Login("final", "1")
	assert.NoError(t, err)
	assert.NotNil(t, id)
}
func TestMemoryStorage_AddData(t *testing.T) {
	s := NewMemoryStorage()
	Init()
	err := s.AddData(datamodels.Data{DataID: "new", Data: "test", Metadata: "test"})
	assert.NoError(t, err)

}
func TestMemoryStorage_DelData(t *testing.T) {
	s := NewMemoryStorage()
	Init()
	err := s.DelData("new", 0)
	assert.NoError(t, err)
}
func TestMemoryStorage_Get(t *testing.T) {
	s := NewMemoryStorage()
	Init()
	err := s.AddData(datamodels.Data{DataID: "new", Data: "test", Metadata: "test"})
	assert.NoError(t, err)
	data, err := s.GetData("new", 0)
	assert.NoError(t, err)
	assert.NotNil(t, data)
}
