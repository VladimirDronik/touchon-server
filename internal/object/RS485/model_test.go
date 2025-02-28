package RS485

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"touchon-server/internal/g"
)

func init() {
	g.Logger = logrus.New()
}

func setUp(t *testing.T) (*RS485Impl, *MockClient) {
	o, err := MakeModel(true)
	require.NotNil(t, o)
	require.NoError(t, err)

	mb, ok := o.(*RS485Impl)
	require.True(t, ok)

	client := new(MockClient)
	mb.client = client

	require.NoError(t, mb.GetProps().Set("tries", 1))
	require.NoError(t, mb.Start())

	return mb, client
}

func TestRS485Model(t *testing.T) {
	mb, client := setUp(t)

	client.EXPECT().Open().Return(nil)
	client.EXPECT().SetUnitId(uint8(0x0001)).Return(nil)
	for i := 0; i < 10; i++ {
		client.EXPECT().WriteCoil(uint16(i), true).Return(nil)
	}
	client.EXPECT().Close().Return(nil)

	for i := 0; i < 10; i++ {
		go func(i int) {
			action := func(client Client) (interface{}, error) {
				return nil, client.WriteCoil(uint16(i), true)
			}
			requestHandler := func(result interface{}, err error) {}

			_ = mb.DoAction(0x0001, action, 4, requestHandler, 2)
		}(i)
	}

	time.Sleep(500 * time.Millisecond)

	client.AssertExpectations(t)
}
