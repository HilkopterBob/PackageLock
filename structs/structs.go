package structs

type Network_Info struct {
	Ip_addr  string
	Mac_addr string
	// TODO: add domain or FQDN
}

type Package_Manager struct {
	Package_manager_name string
	Package_repos        []string
}

type Host struct {
	// TODO: support different linux distros
	ID               int
	Name             string
	Current_packages []string
	Network_info     Network_Info
	Package_manager  Package_Manager
}

type Agent struct {
	Agent_name   string
	Agent_secret string
	Host_ID      int
	Agent_ID     int
}
