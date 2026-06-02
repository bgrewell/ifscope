package model

// PCIDevice is a PCIe network device correlated to its Linux interface name.
//
// Vendor/device IDs, NUMA node, and link speed/width are optional and depend
// on what lspci and sysfs expose.
type PCIDevice struct {
	Bus               string `json:"bus"`
	Interface         string `json:"interface,omitempty"`
	Driver            string `json:"driver,omitempty"`
	VendorID          string `json:"vendor_id,omitempty"`
	DeviceID          string `json:"device_id,omitempty"`
	SubsystemVendorID string `json:"subsystem_vendor_id,omitempty"`
	SubsystemDeviceID string `json:"subsystem_device_id,omitempty"`
	Description       string `json:"description,omitempty"`
	NUMANode          *int   `json:"numa_node,omitempty"`
	LinkSpeed         string `json:"link_speed,omitempty"`
	LinkWidth         string `json:"link_width,omitempty"`
}
