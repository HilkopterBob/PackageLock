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
	PackageID      uuid.UUID
	PackageName    string
	PackageVersion string
	Updatable      bool
	CreationTime   time.Time
	UpdateTime     time.Time
}

type Package_Manager struct {
	PackageManagerName string
	PackageRepos       []string // A Slice containing all Repository Links.
	CreationTime       time.Time
	UpdateTime         time.Time
}

type Host struct {
	// TODO: support different linux distros
	HostName       string // FQDN
	CostID         int
	FQDN           string
	NetworkInfo    map[string]string //	keys: InterfaceName, IPAddress, MacAddress
	Distro         string
	Arch           string
	PackageManager Package_Manager
}

type Agent struct {
	AgentName   string
	AgentSecret string // a secret for encryption
	HostID      int
	AgentID     int
}

type ApiKey struct {
	KeyValue         string
	Description      string
	AccessSeperation bool     // true means fine grained access control
	AccessRights     []string // eg. read, write OR create, update, delete
	CreationTime     time.Time
	UpdateTime       time.Time
}

type User struct {
	UserID       uuid.UUID
	Username     string
	Password     string
	Groups       []string
	CreationTime time.Time
	UpdateTime   time.Time
	ApiKeys      []ApiKey
}
