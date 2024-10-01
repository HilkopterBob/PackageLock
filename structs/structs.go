// Structs
//
// The Structs Package privides needed structs.
// It is a Package Utility.
package structs

type Network_Info struct {
	Ip_addr  string
	Gateway string
	DNS_servers []string
	Netmask string
	Link_speed int
	Mac_addr string
	Domain string
}

type Package_Manager struct {
	Package_manager_name string
	Package_repos        []string // A Slice containing all Repository Links.
}

type Host struct {
	Name             string // FQDN
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
