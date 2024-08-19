package handler

import "packagelock/structs"

var Hosts = []structs.Host{
	{ID: 1, Name: "Host1", Network_info: structs.Network_Info{Ip_addr: "192.168.1.1", Mac_addr: "AA:BB:CC:DD:EE:FF"}, Package_manager: structs.Package_Manager{Package_manager_name: "pacman", Package_repos: []string{"Repo1", "Repo2", "Repo3"}}, Current_packages: []string{"Package1", "package2", "Package3"}},
	{ID: 2, Name: "Host2", Network_info: structs.Network_Info{Ip_addr: "192.168.1.2", Mac_addr: "AA:BB:CC:DD:EF:00"}, Package_manager: structs.Package_Manager{Package_manager_name: "pacman", Package_repos: []string{"Repo1", "Repo2", "Repo3"}}, Current_packages: []string{"Package1", "package2", "Package3"}},
	{ID: 3, Name: "Host3", Network_info: structs.Network_Info{Ip_addr: "192.168.1.3", Mac_addr: "AA:BB:CC:DD:EF:01"}, Package_manager: structs.Package_Manager{Package_manager_name: "pacman", Package_repos: []string{"Repo1", "Repo2", "Repo3"}}, Current_packages: []string{"Package1", "package2", "Package3"}},
}

var Agents = []structs.Agent{
	{Agent_name: "Agent Host1", Agent_secret: "11:11:11:11", Host_ID: 1, Agent_ID: 1},
	{Agent_name: "Agent Host2", Agent_secret: "11:11:11:12", Host_ID: 2, Agent_ID: 2},
	{Agent_name: "Agent Host3", Agent_secret: "11:11:11:13", Host_ID: 3, Agent_ID: 3},
}

var Users = []structs.User{
	{Username: "JohnDoe", Password: "password123", UserID: "12345", APIToken: []string{"token1", "token2", "token3"}},
	{Username: "JaneDoe", Password: "password456", UserID: "67890", APIToken: []string{"token4", "token5", "token6"}},
}
