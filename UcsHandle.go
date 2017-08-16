package goucs

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// debugPOST Trigger Debug messages for POST
const debugPOST = false

// debugXML Trigger Debug messages for XML structs
const debugXML = false

// UcsHandle UCS Handler Struct
type UcsHandle struct {
	cookie   string
	URL      string
	Username string
	Password string
}

// Login Login to UCS and get the Cookie
func (u *UcsHandle) Login() uint {
	var login aaaLogin
	login.InName = u.Username
	login.InPassword = u.Password

	var body = post(u.URL, login)

	var r aaaLogin
	xml.Unmarshal(body, &r)
	if debugXML {
		fmt.Printf("aaaLogin\n")
		spew.Dump(r)
	}

	u.cookie = r.OutCookie

	return r.ErrorCode
}

// Refresh Refresh the Cookie - Note this will generate a NEW Cookie
func (u *UcsHandle) Refresh() uint {
	var refresh aaaRefresh
	refresh.InName = u.Username
	refresh.InPassword = u.Password
	refresh.InCookie = u.cookie

	var body = post(u.URL, refresh)

	var r aaaRefresh
	xml.Unmarshal(body, &r)
	if debugXML {
		spew.Dump(r)
	}

	u.cookie = r.OutCookie

	return r.ErrorCode
}

// Logout Logout from the UCS session
func (u *UcsHandle) Logout() uint {
	var logout aaaLogout
	logout.InCookie = u.cookie

	body := post(u.URL, logout)

	var r aaaLogout
	xml.Unmarshal(body, &r)

	if debugXML {
		spew.Dump(r)
	}

	u.cookie = ""

	return r.ErrorCode
}

// ConfMo Configure a Managed Object
func (u *UcsHandle) ConfMo(mo managedObject) {
	var confMo configConfMo
	confMo.Cookie = u.cookie
	confMo.InConfig = append(confMo.InConfig, mo)

	if debugXML {
		xmlStr, xmlErr := xml.Marshal(confMo)
		spew.Dump(xmlStr, xmlErr)
	}

	post(u.URL, confMo)
}

// UCSObject convert a managedObject to a hiearchy of typed UCS objects
// func (u *UcsHandle) UCSObject(obj managedObject) interface{} {
// 	for _, element := range r.OutConfigs {
// 		e := element.(FaultInst)
//
// 	}
//
// 	return obj
// }

// ResolveClass Search by Class
func (u *UcsHandle) ResolveClass(class string, hierarchical string) ConfigResolveClass {
	var resolveClass ConfigResolveClass
	resolveClass.Cookie = u.cookie
	resolveClass.ClassID = class
	// Not handling the xs:choice field in XML, this may be causing the issue
	// TODO manage hierarchical objects in xsd
	resolveClass.InHierarchical = hierarchical

	if debugXML {
		xmlStr, xmlErr := xml.Marshal(resolveClass)
		fmt.Printf("%+v\n", xmlErr)
		fmt.Printf("%s\n", xmlStr)
	}

	body := post(u.URL, resolveClass)

	var r ConfigResolveClass
	xml.Unmarshal(body, &r)

	if debugXML {
		spew.Dump(r)
	}

	return r
}

func post(url string, xmlStruct interface{}) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var client = &http.Client{Transport: tr}

	xmlStr, xmlErr := xml.Marshal(xmlStruct)
	if debugXML {
		fmt.Printf("%+v\n", xmlErr)
	}

	resp, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(string(xmlStr)))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	t := time.Now()
	if debugPOST {
		fmt.Printf("%s - %s\n", t.Format(time.RFC3339), body)
	}

	return body
}
