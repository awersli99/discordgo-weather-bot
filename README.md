# discordgo-weather-bot
A discord bot that displays the current weather statistics from any location in the world.

## Installation

Install Go here: https://golang.org/doc/install

Use 'go get' to install the discord go package

```bash
go get github.com/bwmarrin/discordgo
```
Then use 'go build' to create an executable.

```bash
go build main.go
```

## Usage

```bash
./{executable name} -t {bot token} -k {api.apixu.com api key}
```

## Bot usage

```bash
?weather {location name}
```

If no locations are being found even though the names are correct, the API key is most likely incorrect.
