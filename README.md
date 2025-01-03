![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/HilkopterBob/PackageLock/.github%2Fworkflows%2Frun-tests.yml)
![GitHub Actions Workflow Status](https://github.com/hilkopterbob/packagelock/actions/workflows/release_tag.yml/badge.svg)
![GitHub Actions Workflow Status](https://github.com/hilkopterbob/packagelock/actions/workflows/golangci-lint.yml/badge.svg)
![GitHub Actions Workflow Status](https://github.com/hilkopterbob/packagelock/actions/workflows/build-docker-container.yml/badge.svg)
![GitHub Actions Workflow Status](https://github.com/hilkopterbob/packagelock/actions/workflows/unstable-build-docker-container.yml/badge.svg)

![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/HilkopterBob/PackageLock)


<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]




<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/HilkopterBob/PackageLock">
    <img src="README-Assets/logo.png" alt="Logo">
  </a>

<h3 align="center">PackageLock</h3>

  <p align="center">
    Leaning Project!
    PackageLock aims to be the distro-agnostic 🔋-included one-stop Solution for patchmanagement on linux systems.
    <br />
    <a href="https://github.com/HilkopterBob/PackageLock"><strong>Explore the docs (COMING SOON!) »</strong></a>
    <br />
    <br />
    <a href="https://github.com/HilkopterBob/PackageLock">View Demo (COMING SOON!)</a>
    ·
    <a href="https://github.com/HilkopterBob/PackageLock/issues/new?labels=bug&template=bug-report---.md">Report Bug (COMING SOON!)</a>
    ·
    <a href="https://github.com/HilkopterBob/PackageLock/issues/new?labels=enhancement&template=feature-request---.md">Request Feature (COMING SOON)</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

<!--- [![Product Name Screen Shot][product-screenshot]](https://example.com)
-->
I created PackageLock from the need to have a one-stop platform for package and software-management on my linux servers.
I wanted to design and create a system that allowed me to create....
<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

TODO: create getting started guide
TODO: general deployment?

### Prerequisites

Install docker. Thats it.



### Installation

#### Docker Run:
```bash
sudo docker run -p 8080:8080 hilkopterbob/packagelock
```

#### Docker Compose:

- download the docker-compose file:
`~/ $: wget https://github.com/HilkopterBob/PackageLock/blob/master/docker-compose.yml`

- get the default config & rename it to `config.yml`:
```bash
~/ $: wget https://github.com/HilkopterBob/PackageLock/blob/master/default-config.yml
~/ $: mv default-config.yml config.yml
```
- edit the config
- run docker compose:
`~/ $: docker-compose up -d`

The default-config:
```yaml
general:
  debug: true
  production: false
database:
  address: 127.0.0.1
  port: 8000
  username: root
  password: root
network:
  fqdn: 0.0.0.0
  port: 8080
  ssl: true
  ssl-config:
    allowselfsigned: true
    certificatepath: ./certs/testing.crt
    privatekeypath: ./certs/testing.key
    redirecthttp: true  
```




<!-- USAGE EXAMPLES -->
## Usage

TODO: explain usage


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [ ] backend-api to manage Agents & Hosts
- [ ] frontend to visualize backend data
- [ ] installable agent as background daemon
- [ ] agent CLI:
  - [ ] `packagelock id` -> returns agent id
- [ ] config management
- [ ] TLS Encryption
- [ ] Best Practice based Package Layout
- [ ] Check Vars and Func-Names for naming convention
- [ ] persistent storage
- [ ] implement interfaces for external functions for easier mocking in tests
- [ ] systemd service start/stop/enable/disable
- [ ] copy app file (.deb/rpm/binary) via SFTP to host and start stop
- [ ] binary self-Update
- [ ] agent can run docker/podman containers
- [ ] agent fetches running docker/podman containers, updates, restarts etc
- [ ] user management & SSH keys
- [ ] system definition in mpackagelock file for easy recovery & scaling
- [ ] CLI-Commands to add:
  - [ ] sync now|timestamp - force sync the server with the Agents
  - [ ] logs -s (severity) info|warning|error -d (date to start) 2024-08-23-10-00-00 (date-time)
  - [ ] backup - Creates a backup from server, server config, database
  - [ ] generate certs letsencrypt - lets encrypt certs
  - [ ] generate certs letsencrypt renew - renews
  - [ ] test - runs healthchecks on server
  - [ ] test agents - runs healthchecks on agents



TODO: create Issue template:
See the [open issues](https://github.com/HilkopterBob/PackageLock/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Top contributors:

<a href="https://github.com/HilkopterBob/PackageLock/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=HilkopterBob/PackageLock" alt="contrib.rocks image" />
</a>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Project Link: [https://github.com/HilkopterBob/PackageLock](https://github.com/HilkopterBob/PackageLock)

<p align="right">(<a href="#readme-top">back to top</a>)</p>




<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/HilkopterBob/PackageLock.svg?style=for-the-badge
[contributors-url]: https://github.com/HilkopterBob/PackageLock/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/HilkopterBob/PackageLock.svg?style=for-the-badge
[forks-url]: https://github.com/HilkopterBob/PackageLock/network/members
[stars-shield]: https://img.shields.io/github/stars/HilkopterBob/PackageLock.svg?style=for-the-badge
[stars-url]: https://github.com/HilkopterBob/PackageLock/stargazers
[issues-shield]: https://img.shields.io/github/issues/HilkopterBob/PackageLock.svg?style=for-the-badge
[issues-url]: https://github.com/HilkopterBob/PackageLock/issues
[license-shield]: https://img.shields.io/github/license/HilkopterBob/PackageLock.svg?style=for-the-badge
[license-url]: https://github.com/HilkopterBob/PackageLock/blob/master/LICENSE
[product-screenshot]: images/screenshot.png
[Next.js]: https://img.shields.io/badge/next.js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white
[Next-url]: https://nextjs.org/
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[Vue.js]: https://img.shields.io/badge/Vue.js-35495E?style=for-the-badge&logo=vuedotjs&logoColor=4FC08D
[Vue-url]: https://vuejs.org/
[Angular.io]: https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white
[Angular-url]: https://angular.io/
[Svelte.dev]: https://img.shields.io/badge/Svelte-4A4A55?style=for-the-badge&logo=svelte&logoColor=FF3E00
[Svelte-url]: https://svelte.dev/
[Laravel.com]: https://img.shields.io/badge/Laravel-FF2D20?style=for-the-badge&logo=laravel&logoColor=white
[Laravel-url]: https://laravel.com
[Bootstrap.com]: https://img.shields.io/badge/Bootstrap-563D7C?style=for-the-badge&logo=bootstrap&logoColor=white
[Bootstrap-url]: https://getbootstrap.com
[JQuery.com]: https://img.shields.io/badge/jQuery-0769AD?style=for-the-badge&logo=jquery&logoColor=white
[JQuery-url]: https://jquery.com 
