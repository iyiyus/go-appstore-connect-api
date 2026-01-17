package appstore

// BundleIdAPI handles bundle ID-related operations
type BundleIdAPI struct {
	client *Client
}

// NewBundleIdAPI creates a new BundleId API client
func NewBundleIdAPI(client *Client) *BundleIdAPI {
	return &BundleIdAPI{client: client}
}

// All retrieves all bundle IDs
func (b *BundleIdAPI) All(params map[string]string) (map[string]interface{}, error) {
	if err := b.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return b.client.GetHTTPClient().Get("/bundleIds", params)
}

// Register registers a new bundle ID
func (b *BundleIdAPI) Register(name, platform, bundleId string) (map[string]interface{}, error) {
	if err := b.client.EnsureAuth(); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "bundleIds",
			"attributes": map[string]string{
				"identifier": bundleId,
				"name":       name,
				"platform":   platform,
			},
		},
	}

	return b.client.GetHTTPClient().PostJSON("/bundleIds", data)
}

// Delete deletes a bundle ID by ID
func (b *BundleIdAPI) Delete(bId string) (map[string]interface{}, error) {
	if err := b.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return b.client.GetHTTPClient().Delete("/bundleIds/"+bId, nil)
}

// Query queries bundle ID capabilities for a specific bundle ID
func (b *BundleIdAPI) Query(bId string, params map[string]string) (map[string]interface{}, error) {
	if err := b.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return b.client.GetHTTPClient().Get("/bundleIds/"+bId+"/bundleIdCapabilities", params)
}
