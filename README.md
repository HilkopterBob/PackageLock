<div align="center">
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/HilkopterBob/PackageLock/.github%2Fworkflows%2Frun-tests.yml)
![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/HilkopterBob/PackageLock)
[CodeQL](https://github.com/HilkopterBob/Packagelock/.github/workflows/codeql.yml/badge.svg
?branch=main&event=push)]
(https://github.com/HilkopterBob/PackageLock/.github/workflows/codeql.yml) 

</div>

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/HilkopterBob/PackageLock/.github%2Fworkflows%2Frun-tests.yml)
![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/HilkopterBob/PackageLock)

[CodeQL](https://github.com/HilkopterBob/Packagelock/.github/workflows/codeql.yml/badge.svg?branch=main&event=push)

Packagelock is an interactive Serversoftware that shows all packages on your servers collected by agents.  

## Backend
the go based Backend provides a JSON-REST API for the frontend.  

## Frontend 
the TypeScript based frontend is a Single Page Application that uses the REST backend to display the data. it uses patternfly components (RedHat UI tool-Kit).


feature creep:
- timed pooling from agents to '/pool/' to get commands like rescans or updates

