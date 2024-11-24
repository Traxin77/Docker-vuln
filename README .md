
# Docker-vuln

docker vuln is opensource cli-tool to scaan docker images.It scans the docker images by generating the sbom of the image and scanning the sbom for vulnerabilities 
## Dependencies
Install go installer through official go download page

verify installation
```bash
  go --version
```
Install syft

```bash
  scoop install syft
```
Install Cobra-cli tool

```bash
  go install github.com/spf13/cobra-cli@latest
```
verify  installation

```bash
  cobra-cli --version
```
Install osv-scanner

```bash
  go install github.com/google/osv-scanner/cmd/osv-scanner@latest

```
verify installation
```bash
  osv-scanner --version

```
## Installation

Install my project by

```bash
  git clone https://github.com/Traxin77/Docker-vuln
  cd docker-vuln
  go build
```
    
## Usage/Examples
List out the docker images to scan 
```Bash
.\docker-vuln.exe list
```
Dockerize a github repo
```Bash
.\docker-vuln.exe dock <Your_github_repo>
```
Scanning docker images
```Bash
.\docker-vuln.exe scan <test_image>
```

Scanning the docker images and storing the output
```Bash
.\docker-vuln.exe scan <test_image> -o <file_name>.json
```


## Authors

- [@Traxin](https://github.com/Traxin77)

