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
)

// debugPOST Trigger Debug messages for POST
const debugPOST = true

// debugXML Trigger Debug messages for XML structs
const debugXML = true

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
		fmt.Printf("%+v\n", r)
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
		fmt.Printf("%+v\n", r)
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
		fmt.Printf("%+v\n", r)
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
		fmt.Printf("%+v\n", xmlErr)
		fmt.Printf("%s\n", xmlStr)
	}

	post(u.URL, confMo)
}

// ResolveClass Search by Class
func (u *UcsHandle) ResolveClass(class string) {
	var resolveClass configResolveClass
	resolveClass.Cookie = u.cookie
	resolveClass.ClassID = class

	if debugXML {
		xmlStr, xmlErr := xml.Marshal(resolveClass)
		fmt.Printf("%+v\n", xmlErr)
		fmt.Printf("%s\n", xmlStr)
	}

	body := post(u.URL, resolveClass)

	var r configResolveClass
	xml.Unmarshal(body, &r)

	if debugXML {
		fmt.Printf("%+v\n", r)
	}
}

// <aaaLogin cookie="" response="yes" outCookie="1485838848/75d77788-475d-175d-8002-eebc39fc1d88" outRefreshPeriod="600" outPriv="admin" outSessionId="6" outVersion="2.0(9l)"> </aaaLogin>
// <aaaLogin cookie="" response="yes" errorCode="551" invocationResult="unidentified-fail" errorDescr="Authentication failed"> </aaaLogin>
/*
type aaaLogin struct {
	XMLName          xml.Name `xml:"aaaLogin"`
	InName           string   `xml:"inName,attr,omitempty"`
	InPassword       string   `xml:"inPassword,attr,omitempty"`
	Cookie           string   `xml:"cookie,attr,omitempty"`
	Response         string   `xml:"response,attr,omitempty"`
	OutCookie        string   `xml:"outCookie,attr,omitempty"`
	OutRefreshPeriod string   `xml:"outRefreshPeriod,attr,omitempty"`
	OutPriv          string   `xml:"outPriv,attr,omitempty"`
	OutSessionID     int      `xml:"outSessionId,attr,omitempty"`
	OutVersion       string   `xml:"outVersion,attr,omitempty"`
	ErrorCode        int      `xml:"errorCode,attr,omitempty"`
	InvocationResult string   `xml:"invocationResult,attr,omitempty"`
	ErrorDescr       string   `xml:"errorDescr,attr,omitempty"`
}
*/

// <aaaLogout cookie="" response="yes" outStatus="success"> </aaaLogout>
// <aaaLogout cookie="" response="yes" errorCode="555" invocationResult="unidentified-fail" errorDescr="Session not found"> </aaaLogout>
/*
type aaaLogout struct {
	XMLName          xml.Name `xml:"aaaLogout"`
	InCookie         string   `xml:"inCookie,attr,omitempty"`
	Cookie           string   `xml:"cookie,attr,omitempty"`
	Response         string   `xml:"response,attr,omitempty"`
	ErrorCode        int      `xml:"errorCode,attr,omitempty"`
	InvocationResult string   `xml:"invocationResult,attr,omitempty"`
	ErrorDescr       string   `xml:"errorDescr,attr,omitempty"`
	OutStatus        string   `xml:"outStatus,attr,omitempty"`
}
*/

// <aaaRefresh cookie="" response="yes" outCookie="1485969023/43f92d8e-477b-177b-8005-b523c4067ff0" outRefreshPeriod="600" outPriv="admin" outSessionId="6" outVersion="2.0(13i)"> </aaaRefresh>
// <error cookie="" response="yes" errorCode="ERR-xml-parse-error" invocationResult="594" errorDescr="XML PARSING ERROR: Required inCookie attribute missing in the xml request. " />
/*
type aaaRefresh struct {
	XMLName          xml.Name `xml:"aaaRefresh"`
	InName           string   `xml:"inName,attr,omitempty"`
	InPassword       string   `xml:"inPassword,attr,omitempty"`
	InCookie         string   `xml:"inCookie,attr,omitempty"`
	Cookie           string   `xml:"cookie,attr,omitempty"`
	Response         string   `xml:"response,attr,omitempty"`
	OutCookie        string   `xml:"outCookie,attr,omitempty"`
	OutRefreshPeriod string   `xml:"outRefreshPeriod,attr,omitempty"`
	OutPriv          string   `xml:"outPriv,attr,omitempty"`
	OutSessionID     int      `xml:"outSessionId,attr,omitempty"`
	OutVersion       string   `xml:"outVersion,attr,omitempty"`
	ErrorCode        int      `xml:"errorCode,attr,omitempty"`
	InvocationResult string   `xml:"invocationResult,attr,omitempty"`
	ErrorDescr       string   `xml:"errorDescr,attr,omitempty"`
}
*/

type outConfigs struct {
	XMLName    xml.Name      `xml:"outConfigs"`
	OutConfigs []interface{} `xml:"outConfigs"`
}

/*
type configConfMo struct {
	XMLName        xml.Name    `xml:"configConfMo"`
	Cookie         string      `xml:"cookie,attr"`
	Dn             string      `xml:"dn,attr,omitempty"`
	InHierarchical string      `xml:"inHierarchical,attr,omitempty"`
	InConfig       interface{} `xml:"inConfig>InConfig,omitempty"`
	OutConfig      interface{} `xml:"inConfig>OutConfig,omitempty"`
}
*/

// <configResolveClass cookie="1313086522/c7c08988-aa3e-1a3e-8005-5e61c2e14388" inHierarchical="false" classId="firmwareRunning"/>
/*
type configResolveClass struct {
	XMLName        xml.Name `xml:"configResolveClass"`
	Cookie         string   `xml:"cookie,attr"`
	ClassID        string   `xml:"classId,attr,omitempty"`
	InHierarchical string   `xml:"inHierarchical,attr,omitempty"`
	OutConfigs     *outConfigs
}
*/

func post(url string, xmlStruct interface{}) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var client = &http.Client{Transport: tr}

	xmlStr, xmlErr := xml.Marshal(xmlStruct)
	fmt.Printf("%+v\n", xmlErr)

	resp, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(string(xmlStr)))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%+v\n", err)

	t := time.Now()
	if debugPOST {
		fmt.Printf("%s - %s\n", t.Format(time.RFC3339), body)
	}

	return body
}
