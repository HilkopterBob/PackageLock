// Structs
//
// The Structs Package privides needed structs.
// It is a Package Utility.
package structs

type Network_Info struct {
	Ip_addr  string
	Mac_addr string
	// TODO: add domain or FQDN
}

type Package_Manager struct {
	Package_manager_name string
	Package_repos        []string // A Slice containing all Repository Links.
}

type Host struct {
	// TODO: support different linux distros
	Name             string // FQDN
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
