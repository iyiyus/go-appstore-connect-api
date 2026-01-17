package appstore

// ProfilesAPI handles profile-related operations
type ProfilesAPI struct {
	client *Client
}

// NewProfilesAPI creates a new Profiles API client
func NewProfilesAPI(client *Client) *ProfilesAPI {
	return &ProfilesAPI{client: client}
}

// Query retrieves profiles with optional parameters
func (p *ProfilesAPI) Query(params map[string]string) (map[string]interface{}, error) {
	if err := p.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return p.client.GetHTTPClient().Get("/profiles", params)
}

// ProfileRelationship represents a relationship item
type ProfileRelationship struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Create creates a new profile
func (p *ProfilesAPI) Create(name, bId, profileType string, devices []string, certificates []string) (map[string]interface{}, error) {
	if err := p.client.EnsureAuth(); err != nil {
		return nil, err
	}

	// Prepare devices relationship
	devicesData := make([]ProfileRelationship, len(devices))
	for i, device := range devices {
		devicesData[i] = ProfileRelationship{
			Type: "devices",
			ID:   device,
		}
	}

	// Prepare certificates relationship
	certificatesData := make([]ProfileRelationship, len(certificates))
	for i, cert := range certificates {
		certificatesData[i] = ProfileRelationship{
			Type: "certificates",
			ID:   cert,
		}
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "profiles",
			"relationships": map[string]interface{}{
				"bundleId": map[string]interface{}{
					"data": ProfileRelationship{
						Type: "bundleIds",
						ID:   bId,
					},
				},
				"devices": map[string]interface{}{
					"data": devicesData,
				},
				"certificates": map[string]interface{}{
					"data": certificatesData,
				},
			},
			"attributes": map[string]string{
				"profileType": profileType,
				"name":        name,
			},
		},
	}

	return p.client.GetHTTPClient().PostJSON("/profiles", data)
}

// ListDevices lists devices for a profile
func (p *ProfilesAPI) ListDevices(pId string, params map[string]string) (map[string]interface{}, error) {
	if err := p.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return p.client.GetHTTPClient().Get("/profiles/"+pId+"/devices", params)
}

// ListCertificates lists certificates for a profile
func (p *ProfilesAPI) ListCertificates(pId string, params map[string]string) (map[string]interface{}, error) {
	if err := p.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return p.client.GetHTTPClient().Get("/profiles/"+pId+"/relationships/certificates", params)
}

// Delete deletes a profile by ID
func (p *ProfilesAPI) Delete(pId string) (map[string]interface{}, error) {
	if err := p.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return p.client.GetHTTPClient().Delete("/profiles/"+pId, nil)
}
