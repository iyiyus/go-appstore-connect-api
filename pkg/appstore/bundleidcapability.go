package appstore

// BundleIdCapabilityAPI handles bundle ID capability-related operations
type BundleIdCapabilityAPI struct {
	client *Client
}

// NewBundleIdCapabilityAPI creates a new BundleIdCapability API client
func NewBundleIdCapabilityAPI(client *Client) *BundleIdCapabilityAPI {
	return &BundleIdCapabilityAPI{client: client}
}

// Enable enables a capability for a bundle ID
func (b *BundleIdCapabilityAPI) Enable(bId, capability string) (map[string]interface{}, error) {
	if err := b.client.EnsureAuth(); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "bundleIdCapabilities",
			"relationships": map[string]interface{}{
				"bundleId": map[string]interface{}{
					"data": map[string]string{
						"type": "bundleIds",
						"id":   bId,
					},
				},
			},
			"attributes": map[string]string{
				"capabilityType": capability,
			},
		},
	}

	return b.client.GetHTTPClient().PostJSON("/bundleIdCapabilities", data)
}

// Disable disables a bundle ID capability by ID
func (b *BundleIdCapabilityAPI) Disable(bcId string) (map[string]interface{}, error) {
	if err := b.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return b.client.GetHTTPClient().Delete("/bundleIdCapabilities/"+bcId, nil)
}
