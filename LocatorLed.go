package goucs

// <equipmentLocatorLed adminState='on' dn='sys/rack-unit-1/locator-led'></equipmentLocatorLed>
/*
type equipmentLocatorLed struct {
	XMLName    xml.Name `xml:"equipmentLocatorLed"`
	Dn         string   `xml:"dn,attr,omitempty"`
	AdminState string   `xml:"adminState,attr,omitempty"`
	OperState  string   `xml:"operState,attr,omitempty"`
}
*/

// LocatorLedOn Enables Locator LED
func LocatorLedOn(u *UcsHandle, dn string) {
	var loc equipmentLocatorLed
	loc.Dn = dn
	loc.AdminState = "on"

	u.ConfMo(loc)
}

// LocatorLedOff Disables Locator LED
func LocatorLedOff(u *UcsHandle, dn string) {
	var loc equipmentLocatorLed
	loc.Dn = dn
	loc.AdminState = "off"

	u.ConfMo(loc)
}
