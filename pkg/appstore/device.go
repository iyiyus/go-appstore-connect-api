package appstore

import (
	"fmt"
	"strings"
)

// DeviceAPI handles device-related operations
type DeviceAPI struct {
	client *Client
}

// NewDeviceAPI creates a new Device API client
func NewDeviceAPI(client *Client) *DeviceAPI {
	return &DeviceAPI{client: client}
}

// All retrieves all devices
func (d *DeviceAPI) All(params map[string]string) (map[string]interface{}, error) {
	if err := d.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return d.client.GetHTTPClient().Get("/devices", params)
}

// Register registers a new device
func (d *DeviceAPI) Register(name, platform, udid string) (map[string]interface{}, error) {
	if err := d.client.EnsureAuth(); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "devices",
			"attributes": map[string]string{
				"name":     name,
				"platform": strings.ToUpper(platform),
				"udid":     udid,
			},
		},
	}

	return d.client.GetHTTPClient().PostJSON("/devices", data)
}

// DeviceType represents device type information
type DeviceType struct {
	Success    bool   `json:"success"`
	DeviceClass string `json:"deviceClass"`
	Model      string `json:"model"`
	Platform   string `json:"platform"`
	Status     string `json:"status"`
	IsIPhone   bool   `json:"isIPhone"`
	IsIPad     bool   `json:"isIPad"`
	IsMac      bool   `json:"isMac"`
	Error      string `json:"error,omitempty"`
}

// GetDeviceType retrieves device type information for a given UDID
func (d *DeviceAPI) GetDeviceType(udid string) (DeviceType, error) {
	params := map[string]string{
		"filter[udid]":      udid,
		"fields[devices]":   "deviceClass,model,platform,status",
	}

	response, err := d.All(params)
	if err != nil {
		return DeviceType{Success: false, Error: err.Error()}, nil
	}

	// Check for API errors
	if errors, ok := response["errors"].([]interface{}); ok && len(errors) > 0 {
		if errorDetail, ok := errors[0].(map[string]interface{})["detail"].(string); ok {
			return DeviceType{Success: false, Error: errorDetail}, nil
		}
	}

	// Check if device data exists
	data, ok := response["data"].([]interface{})
	if !ok || len(data) == 0 {
		return DeviceType{Success: false, Error: "Device not found"}, nil
	}

	// Parse device data
	device, ok := data[0].(map[string]interface{})
	if !ok {
		return DeviceType{Success: false, Error: "Invalid device data"}, nil
	}

	attributes, ok := device["attributes"].(map[string]interface{})
	if !ok {
		return DeviceType{Success: false, Error: "Invalid device attributes"}, nil
	}

	deviceClass := ""
	if v, ok := attributes["deviceClass"].(string); ok {
		deviceClass = v
	}

	model := ""
	if v, ok := attributes["model"].(string); ok {
		model = v
	}

	platform := ""
	if v, ok := attributes["platform"].(string); ok {
		platform = v
	}

	status := ""
	if v, ok := attributes["status"].(string); ok {
		status = v
	}

	return DeviceType{
		Success:     true,
		DeviceClass: deviceClass,
		Model:       model,
		Platform:    platform,
		Status:      status,
		IsIPhone:    deviceClass == "IPHONE",
		IsIPad:      deviceClass == "IPAD",
		IsMac:       deviceClass == "MAC",
	}, nil
}

// RegisterAndGetType attempts to register a device and returns device type
// If device already exists, it queries existing device information
func (d *DeviceAPI) RegisterAndGetType(name, platform, udid string) (DeviceType, error) {
	// Try to register device first
	registration, err := d.Register(name, platform, udid)
	if err != nil {
		return DeviceType{Success: false, Error: err.Error()}, nil
	}

	// Check for errors
	if errors, ok := registration["errors"].([]interface{}); ok && len(errors) > 0 {
		if errorDetail, ok := errors[0].(map[string]interface{})["detail"].(string); ok {
			// If device already exists, query existing device information
			if strings.Contains(errorDetail, "already exists on this team") {
				return d.GetDeviceType(udid)
			}
			return DeviceType{Success: false, Error: errorDetail}, nil
		}
	}

	// Registration successful, return device information
	data, ok := registration["data"].(map[string]interface{})
	if !ok {
		return DeviceType{Success: false, Error: "Invalid registration data"}, nil
	}

	attributes, ok := data["attributes"].(map[string]interface{})
	if !ok {
		return DeviceType{Success: false, Error: "Invalid device attributes"}, nil
	}

	deviceClass := ""
	if v, ok := attributes["deviceClass"].(string); ok {
		deviceClass = v
	}

	model := ""
	if v, ok := attributes["model"].(string); ok {
		model = v
	}

	platformResult := ""
	if v, ok := attributes["platform"].(string); ok {
		platformResult = v
	}

	status := ""
	if v, ok := attributes["status"].(string); ok {
		status = v
	}

	deviceID := ""
	if v, ok := data["id"].(string); ok {
		deviceID = v
	}

	return DeviceType{
		Success:     true,
		DeviceClass: deviceClass,
		Model:       model,
		Platform:    platformResult,
		Status:      status,
		IsIPhone:    deviceClass == "IPHONE",
		IsIPad:      deviceClass == "IPAD",
		IsMac:       deviceClass == "MAC",
	}, nil
}

// DeviceSortResult represents the result of device sorting
type DeviceSortResult struct {
	Code   int      `json:"code"`
	Msg    string   `json:"msg"`
	Data   DataInfo `json:"data"`
}

// DataInfo contains device count information
type DataInfo struct {
	IPHONE int    `json:"IPHONE"`
	IPAD   int    `json:"IPAD"`
	MAC    int    `json:"MAC"`
	Email  string `json:"email"`
}

// DeviceSort counts devices by type and returns available slots
func (d *DeviceAPI) DeviceSort() (DeviceSortResult, error) {
	result := DeviceSortResult{}

	// Query iOS devices
	iOSParams := map[string]string{
		"filter[platform]":  "IOS",
		"fields[devices]":   "deviceClass",
		"limit":             "200",
	}

	iOSData, err := d.All(iOSParams)
	if err != nil {
		return result, err
	}

	iPhone := 0
	iPad := 0

	if data, ok := iOSData["data"].([]interface{}); ok {
		for _, item := range data {
			if device, ok := item.(map[string]interface{}); ok {
				if attributes, ok := device["attributes"].(map[string]interface{}); ok {
					if deviceClass, ok := attributes["deviceClass"].(string); ok {
						if deviceClass == "IPHONE" {
							iPhone++
						} else if deviceClass == "IPAD" {
							iPad++
						}
					}
				}
			}
		}
	}

	// Query Mac devices
	macParams := map[string]string{
		"filter[platform]":  "MAC_OS",
		"fields[devices]":   "deviceClass",
	}

	macData, err := d.All(macParams)
	if err != nil {
		return DeviceSortResult{
			Code: 1001,
			Msg:  "Failed to query Mac devices",
		}, nil
	}

	// Check for errors
	if errors, ok := macData["errors"].([]interface{}); ok && len(errors) > 0 {
		if errorDetail, ok := errors[0].(map[string]interface{})["detail"].(string); ok {
			return DeviceSortResult{
				Code: 1001,
				Msg:  errorDetail,
			}, nil
		}
	}

	mac := 0
	if meta, ok := macData["meta"].(map[string]interface{}); ok {
		if paging, ok := meta["paging"].(map[string]interface{}); ok {
			if total, ok := paging["total"].(float64); ok {
				mac = int(total)
			}
		}
	}

	// Query user email
	userData, err := d.client.GetHTTPClient().Get("/users", nil)
	if err != nil {
		return result, err
	}

	email := ""
	if data, ok := userData["data"].([]interface{}); ok && len(data) > 0 {
		if user, ok := data[0].(map[string]interface{}); ok {
			if attributes, ok := user["attributes"].(map[string]interface{}); ok {
				if username, ok := attributes["username"].(string); ok {
					email = username
				}
			}
		}
	}

	return DeviceSortResult{
		Code: 1,
		Msg:  "ok",
		Data: DataInfo{
			IPHONE: 100 - iPhone,
			IPAD:   100 - iPad,
			MAC:    100 - mac,
			Email:  email,
		},
	}, nil
}
