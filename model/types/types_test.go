package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeviceState(t *testing.T) {
	assert.Equal(t, "linked", DeviceStateLinked.String())
	assert.Equal(t, "unlinked", DeviceStateUnlinked.String())
}

func TestDeviceType(t *testing.T) {
	assert.Equal(t, "phone", DeviceTypePhone.String())
	assert.Equal(t, "chrome", DeviceTypeChrome.String())
	assert.Equal(t, "desktop", DeviceTypeDesktop.String())
}

func TestLinkedAccountType(t *testing.T) {
	assert.Equal(t, "google", LinkedAccountTypeGoogle.String())
}

func TestMessageStatus(t *testing.T) {
	assert.Equal(t, "started", MessageStatusStarted.String())
	assert.Equal(t, "sent", MessageStatusSent.String())
	assert.Equal(t, "delivered", MessageStatusDelivered.String())
	assert.Equal(t, "failed", MessageStatusFailed.String())
}
