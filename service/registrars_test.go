package service

import (
	"testing"

	"github.com/Comcast/webpa-common/service/servicemock"
	"github.com/stretchr/testify/assert"
)

func testRegistrars(t *testing.T, r Registrars, expectedInitialLen int) {
	assert := assert.New(t)

	assert.Equal(expectedInitialLen, r.Len())
	assert.NotPanics(func() { r.Register() })
	assert.NotPanics(func() { r.Deregister() })

	child := new(servicemock.Registrar)
	child.On("Register").Once()
	child.On("Deregister").Once()
	r.Add("child", child)

	assert.Equal(expectedInitialLen+1, r.Len())
	assert.NotPanics(func() { r.Register() })
	assert.NotPanics(func() { r.Deregister() })

	child.AssertExpectations(t)
}

func TestRegistrars(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		testRegistrars(t, nil, 0)
	})

	t.Run("Empty", func(t *testing.T) {
		testRegistrars(t, Registrars{}, 0)
	})
}
