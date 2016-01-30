package types

// DeviceState is a device state.
type DeviceState string

const (
	DeviceStateLinked   DeviceState = "linked"
	DeviceStateUnlinked DeviceState = "unlinked"
)

func (s DeviceState) String() string {
	return string(s)
}

// DeviceType is a type of device.
type DeviceType string

const (
	DeviceTypePhone   DeviceType = "phone"
	DeviceTypeChrome  DeviceType = "chrome"
	DeviceTypeDesktop DeviceType = "desktop"
)

func (t DeviceType) String() string {
	return string(t)
}

// A LinkedAccountType is a type of linked account.
type LinkedAccountType string

const (
	LinkedAccountTypeGoogle LinkedAccountType = "google"
)

func (t LinkedAccountType) String() string {
	return string(t)
}

type MessageStatus string

const (
	MessageStatusStarted   MessageStatus = "started"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusFailed    MessageStatus = "failed"
)

func (s MessageStatus) String() string {
	return string(s)
}
