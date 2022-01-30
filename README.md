<h1 align="center">huekit</h1>

<p align="center">
  Automatically bridge third party hue lights to HomeKit.
  <br><br>
  Since philips hue does not enable HomeKit for third party ZigBee devices, this project aims to close the gap.
  It automatically detects, whether a light is an original device or a third party device and then bridges it
  to HomeKit.
  <br><br>
  <a href="https://cloud.drone.io/dj95/huekit">
    <img alt="BuildStatus" src="https://cloud.drone.io/api/badges/dj95/huekit/status.svg" />
  </a>
  <a href="https://github.com/dj95/huekit/actions?query=workflow%3AGo">
    <img alt="GoActions" src="https://github.com/dj95/huekit/workflows/Go/badge.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/dj95/huekit">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/dj95/huekit" />
  </a>
  <a href="https://github.com/dj95/huekit/releases">
    <img alt="latest version" src="https://img.shields.io/github/tag/dj95/huekit.svg" />
  </a>
</p>


## 📦 Requirements

- Golang (>=1.14.3)
- Make
- golint (✅ tests/linting)
- staticcheck (✅ tests/linting)
- gosec (✅ tests/linting)


## 🔧 Usage/Installation

- Download the released package for your operating system or follow the build instructions
- Retrieve the ip address of the bridge in the hue app (Settings > Hue Bridges > i near the [Bridge Name])
- Place the config.yml, that is contained in the release package or from [./configs/config.yml](./configs/config.yml), near the retrieved binary.
- Open the config.yml and insert the ip address in the `''` behind the `bridge_address` key
- Change the `homekit_pin` to a random 8-digit pin
- Run `./huekit`
- Check, if it says, that you need to press the link button. If so, press the button to authenticate huekit at your hue bridge
- Open the Home app on your apple device and click on the top right `+` icon to add an accessory.
- Tap on "I Don't Have A Code or Cannot Scan"
- Select `HueKit` and insert the code from the `config.yml`, that is configured after the `homekit_pin` key


**Hint** In order to reset the huekit, remove the `huekit_data` directory near the binary.


## 🏗 Build

In order to build the binary, just run `make build`. The binary will be placed in the `./bin` directory.


## 🐳 Docker

This project features a Dockerfile, docker-compose file and the way to read the configuration from the environment.
The Dockerfile can be found in `build/package/docker/huekit` and the docker-compose.yml in `deployments/docker`.
For building an image for different OS and CPU architextures, the Dockerfile features build arguments, that allow cross compiling the binary.
The docker-compose.yml features a configuration for linux/amd64 and linux/armv7 (compatible with raspberry pi).


Configuration of the docker containers can be done by setting environment variable, such that a config file is not required anymore.
All configuration parameters of the config.yml can be reused with the prefix `HUEKIT_` in order to be set via the environment.
An example configuration can be found in the docker-compose.yml.
The following environment variables are available:


| Name | Description |
|------|-------------|
| `HUEKIT_LOG_LEVEL` | Set the verbosity of the service. |
| `HUEKIT_LOG_FORMAT` | Decide, if you want `json` or `text` logs |
| `HUEKIT_BRIDGE_ADDRESS` | IP address of the hue bridge  |
| `HUEKIT_HOMEKIT_PIN` | Pin, that must be entered in homekit for pairing with huekit |
| `HUEKIT_HOMEKIT_PORT` | Port that huekit will listen on for homekit  |


## 🤝 Contributing

If you are missing features or find some annoying bugs please feel free to submit an issue or a bugfix within a pull request :)


## 📝 License

© 2020 Daniel Jankowski


This project is licensed under the MIT license.


Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:


The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.


THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
