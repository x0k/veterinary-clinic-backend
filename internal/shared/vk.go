package shared

type VkUserId string

func NewVkUserId(id string) VkUserId {
	return VkUserId(id)
}

func (id VkUserId) String() string {
	return string(id)
}
