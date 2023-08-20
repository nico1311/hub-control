## hub-control

A very simple HTTP server that allows to control a USB hub power state. I created this to be able control a USB powered lamp connected to a Raspberry Pi.

### Dependencies
* `uhubctl` - Available from the Debian repositories

### Usage
1. Build the binary: `go build`
2. Run the binary: `./hub-control` (needs root privileges)

### Routes
* GET `/ports/{port}`: Returns the power state of the given port
* POST `/ports/{port}`: Sets the power state of the given port. Expects a JSON body with a `status` key that can be `on`, `off`, `cycle` or `toggle`.

### Future improvements (maybe)
* Remove dependency on `uhubctl`
* Add authentication (not needed right now for my use case since I access behind a VPN)
* Add a web interface
* Dockerize
* Add tests
