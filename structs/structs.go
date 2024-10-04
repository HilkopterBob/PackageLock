// Structs
//
// The Structs Package privides needed structs.
// It is a Package Utility.
package structs

import (
	"time"

	"github.com/google/uuid"
)

type Package struct {
	ID             string `json:"id,omitempty"`
	PackageID      uuid.UUID
	PackageName    string
	PackageVersion string
	Updatable      bool
	CreationTime   time.Time
	UpdateTime     time.Time

type Network_Info struct {
	Interfaces []string
	Domain string
	DNSServers []string
}

type Interface struct {
	Name string
	IpAddr  string
	Gateway string
	Netmask string
	LinkSpeed int
	MacAddr string
}
type Package_Manager struct {
	ID                 string `json:"id,omitempty"`
	PackageManagerName string
	PackageRepos       []string // A Slice containing all Repository Links.
	CreationTime       time.Time
	UpdateTime         time.Time
}

type Host struct {
	ID             string `json:"id,omitempty"`
	Hostname       string // FQDN
	HostID         uuid.UUID
	FQDN           string
	NetworkInfo    map[string]string //	keys: InterfaceName, IPAddress/Net, MacAddress
	Distro         string
	Arch           string
	PackageManager Package_Manager
	Packages       []uuid.UUID
	CreationTime   time.Time
	UpdateTime     time.Time
}

type Agent struct {
	ID           string `json:"id,omitempty"`
	AgentName    string
	AgentSecret  string // a secret for encryption
	HostID       uuid.UUID
	AgentID      uuid.UUID
	CreationTime time.Time
	UpdateTime   time.Time
}

type ApiKey struct {
	ID               string `json:"id,omitempty"`
	KeyValue         string
	Description      string
	AccessSeperation bool     // true means fine grained access control
	AccessRights     []string // eg. read, write OR create, update, delete
	CreationTime     time.Time
	UpdateTime       time.Time
}

type User struct {
	ID           string `json:"id,omitempty"`
	UserID       uuid.UUID
	Username     string
	Password     string
	Groups       []string
	CreationTime time.Time
	UpdateTime   time.Time
	ApiKeys      []ApiKey

}
