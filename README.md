<h1 align="center">
    <img src="https://github.com/stelmanjones/wrc/blob/main/docs/wrc-logo.svg" alt="Logo" width="356" height="125">
  </a>
</h1>

<div align="center">
A small and easy to use library for WRC 23 telemetry data.
</div>

<div align="center">
<br />

[![version](https://img.shields.io/github/v/tag/stelmanjones/wrc?style=flat-square&label=version
)](LICENSE)

[![PRs welcome](https://img.shields.io/badge/PRs-welcome-ff69b4.svg?style=flat-square)](https://github.com/stelmanjones/wrc/issues?q=is%3Aissue+is%3Aopen)
[![activity](https://img.shields.io/github/last-commit/stelmanjones/wrc?style=flat-square&logo=github
)](https://github.com/stelmanjones/wrc/commits)

</div>

<details open="open">
<summary>Table of Contents</summary>

- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

</details>

---

## Getting Started


### Installation

```sh
# Run this to install
go get github.com/stelmanjones/wrc
```


### Usage

```go

var (
  	conn, err = net.ListenPacket("udp4", ":6969")
  Client = wrc.New()
)

func main() {

    go Client.Run()

    for {
      p,err := Client.Latest()
      if err != nil {
       // Error 
      }
      log.Info("Packet: ", "RPM",p.VehicleEngineRpmCurrent)

      spd,err := Client.AverageSpeedKmph()
      log.Infof("Average speed: %f",spd)

    }
}

```
</br>

## Issues

See the [open issues](https://github.com/stelmanjones/wrc/issues) for a list of proposed features (and known issues).

## Contributing

First off, thanks for taking the time to contribute! Contributions are what makes the open-source community such an amazing place to learn, inspire, and create. Any contributions you make will benefit everybody else and are **greatly appreciated**.

Please try to create bug reports that are:

- _Reproducible._ Include steps to reproduce the problem.
- _Specific._ Include as much detail as possible: which version, what environment, etc.
- _Unique._ Do not duplicate existing opened issues.
- _Scoped to a Single Bug._ One bug per report.

## License

This project is licensed under the **MIT license**. Feel free to edit and distribute this repo as you like.
