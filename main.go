package main

import "github.com/gin-gonic/gin"
import "net/http"



// Data structs

// TODO: support for multiple network adapters.
type Network_Info struct {
  Ip_addr string 
  Mac_addr string
  // TODO: add domain or FQDN
}

type Package_Manager struct {
  Package_manager_name string
  Package_repos []string
}

type Host struct {
  // TODO: support different linux distros
  Name string
  Network_info Network_Info
  Package_manager Package_Manager
  Current_packages []string
}

var hosts = []Host{
  {Name: "Host1", Network_info: Network_Info{Ip_addr: "192.168.1.1", Mac_addr: "AA:BB:CC:DD:EE:FF"}, Package_manager: Package_Manager{Package_manager_name: "pacman",Package_repos: []string{"Repo1", "Repo2", "Repo3"},},Current_packages: []string{"Package1", "package2", "Package3"},}, 
  {Name: "Host2", Network_info: Network_Info{Ip_addr: "192.168.1.2", Mac_addr: "AA:BB:CC:DD:EF:00"}, Package_manager: Package_Manager{Package_manager_name: "pacman",Package_repos: []string{"Repo1", "Repo2", "Repo3"},},Current_packages: []string{"Package1", "package2", "Package3"},},
  {Name: "Host3", Network_info: Network_Info{Ip_addr: "192.168.1.3", Mac_addr: "AA:BB:CC:DD:EF:01"}, Package_manager: Package_Manager{Package_manager_name: "pacman",Package_repos: []string{"Repo1", "Repo2", "Repo3"},},Current_packages: []string{"Package1", "package2", "Package3"},},
}

// Endpoints & Data Aggregation Functions
func getHosts(c *gin.Context) {
  c.IndentedJSON(http.StatusOK, hosts)
}



func main () {

  router := gin.Default()
  router.GET("/hosts", getHosts)

  router.Run("localhost:8080")

  

}
