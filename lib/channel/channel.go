package channel

// 許可されたチャンネルかを検証します
func IsAllowedChannelID(allowedChannels []string, channelID string) bool {
	for _, allowedChannel := range allowedChannels {
		if allowedChannel == channelID {
			return true
		}
	}

	return false
}
