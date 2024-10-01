// Structs
//
// The Structs Package privides needed structs.
// It is a Package Utility.
package structs

type Network_Info struct {
	Interfaces []string
	Domain string
	DNS_Servers []string
}

type Interface struct {
	Name string
	Ip_Addr  string
	Gateway string
	Netmask string
	Link_Speed int
	Mac_Addr string
}
type Package_Manager struct {
	Package_manager_name string
	Package_repos        []string // A Slice containing all Repository Links.
}

type Host struct {
	Name             string // Hostname
	Uname			 string //Linux Distro
	ID               int
	Current_packages []string // A Slice with all currently installed Packages.
	Network_info     Network_Info
	Package_manager  Package_Manager
}

type Agent struct {
	Agent_name   string
	Agent_secret string // a secret for encryption
	Host_ID      int
	Agent_ID     int
}

type User struct {
	Username string
	Group string
	Password string
	UserID   string
	APIToken []string
}
