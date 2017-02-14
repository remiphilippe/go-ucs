package goucs

// <firmwareBootUnit dn='sys/rack-unit-1/mgmt/fw-boot-def/bootunit-combined' adminState='trigger' image='backup' resetOnActivate='yes' />
/*
type firmwareBootUnit struct {
	XMLName    xml.Name `xml:"firmwareBootUnit"`
	Dn         string   `xml:"dn,attr,omitempty"`
	AdminState string   `xml:"adminState,attr,omitempty"`
}
*/

// Updating BIOS Firmware
/*
<configConfMo cookie="0000020175/b2909140-0004-1004-8002-cdac38e14388"
inHierarchical="true" dn="sys/rack-unit-1/bios/fw-updatable">
 <inConfig>
   <firmwareUpdatable adminState='trigger' dn='sys/rack-unit-1/bios/fw-updatable'
    protocol='tftp' remoteServer='10.xxx.196.xxx' remotePath='HP-SL2.cap'
    type='blade-bios' />
 </inConfig>
</configConfMo>
*/
