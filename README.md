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


## 🔧 Usage/Installation

- Download the released package for your operating system or follow the build instructions
- Retrieve the ip address of the bridge in the hue app (Settings > Hue Bridges > i near the [Bridge Name])
- Place the config.yml, that is contained in the release package or from [./configs/config.yml](./configs/config.yml), near the retrieved binary.
- Open the config.yml and insert the ip address in the `''` behind the `bridge_address` key
- Run `./huekit`
- Check, if it says, that you need to press the link button. If so, press the button to authenticate huekit at your hue bridge
- Open the Home app on your apple device and click on the top right `+` icon to add an accessory.
- Tap on "I Don't Have A Code or Cannot Scan"
- Select `HueKit` and insert the code from the `config.yml`, that is configured after the `homekit_pin` key


**Hint** In order to reset the huekit, remove the `huekit_data` directory near the binary.


## 🏗 Build

In order to build the binary, just run `make build`. The binary will be placed in the `./bin` directory.


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
