# tenant-utils-fake
A simple fake server which returns some fake responses for the tenant-utils tool

## Usage

Just build the project with `go build` and override the following environment variables if you don't want the fake
service running on the default `:12000` address:

* `HOST`. Example value: `http://localhost`
* `PORT`. Example value: `12000`
