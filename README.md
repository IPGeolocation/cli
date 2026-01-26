# IPGeolocation.io CLI Tool

## Overview

The official `ipgeolocation` **Command Line Interface (CLI)** provides a fast, lightweight, and script-friendly way to access **[IPGeolocation.io](https://ipgeolocation.io)**'s set of APIs directly from your terminal. It enables developers, DevOps engineers, network operators, and security teams to retrieve IP intelligence without writing any code or handling raw HTTP requests.

Using a single executable, you can query **IP location, security intelligence, ASN details, timezone, astronomy, abuse contact information, and user agent parsing** for your own IP or for any IPv4, IPv6 address, or domain name. The CLI outputs structured, machine-readable data that integrates seamlessly into shell scripts, cron jobs, CI/CD pipelines, monitoring systems, and log-processing workflows.

- [IP Location API](https://ipgeolocation.io/ip-location-api.html): Get all-in-one unified solution for **location** (city, locality, state, country, etc.), **currency**, **network** (AS number, ASN name, organization, asn type, date of allocation, company/ISP name, company type, company domain), **timezone** , **useragent** string parsing, **security** (threat score, is_tor, is_bot, proxy_provider, cloud_provider), and **abuse contact** (route/CIDR network, country, address, email, phone numbers) information.
- [IP Security API](https://ipgeolocation.io/ip-security-api.html): Get security, network, location, hostname, timezone and useragent parsing.
- [ASN API](https://ipgeolocation.io/asn-api.html): Get ASN details along with peers, upstreams, downstreams, routes, and raw WHOIS.
- [Abuse Contact API](https://ipgeolocation.io/ip-abuse-contact-api.html): Get abuse emails, phone numbers, kind, organization, route/CIDR network and country.
- [Astronomy API](https://ipgeolocation.io/astronomy-api.html): Get sunrise, sunset, moonrise, moonset, moon phases with precise twilight period times in combination with location information.
- [Timezone API](https://ipgeolocation.io/timezone-api.html): Get timezone name, multiple time formats, daylight saving status and its details along with location information.
- [Timezone Convert API](https://ipgeolocation.io/timezone-api.html): Convert time between timezone names, geo coordinates, location addresses, IATA codes, ICAO codes, or UN/LOCODE.
- [User Agent API](https://ipgeolocation.io/user-agent-api.html): Get browser, Operating System, and device info from single or multiple Useragent string parsing.

The `ipgeolocation` CLI is built for speed and automation, making it easy to integrate IP intelligence into shell scripts, bash pipelines, cron jobs, and CI/CD workflows. It enriches server logs, firewall events, and SIEM data, supports on the fly lookups during incident response, and enables automated fraud detection, geo based routing, and access control without language specific dependencies. Designed for system administrators, DevOps engineers, and developers, the CLI delivers consistent and reliable insights directly from the terminal.

Based on:
- API version: 2.0.0

**Official Release:**
- Available on [**pkg.go.dev**](https://pkg.go.dev/github.com/IPGeolocation/cli) and [**GitHub Releases**](https://github.com/IPGeolocation/cli/releases)
- Source Code: [**GitHub Repository**](https://github.com/IPGeolocation/cli)

## Table of Contents
1. [Requirements](#requirements)
2. [Installation](#installation)
    - [Go install](#go-install)
    - [Download and build from source](#download-and-build-from-source)
    - [Using Download Prebuilt Binaries](#using-download-prebuilt-binaries)
    - [Troubleshooting](#troubleshooting)
3. [API Plan Tiers and Documentation](#api-plan-tiers-and-documentation)
4. [Fields and commands Availability](#fields-and-commands-availability)
5. [Authentication Setup](#authentication-setup)
   - [How to Get Your API Key](#how-to-get-your-api-key)
   - [ApiKeyAuth](#apikeyauth)
6. [Global Flags](#global-flags)
7. [Commands](#commands)
    - [`config` Command](#config-command)
      - [`config` Usage](#config-usage)
      - [Flags for `config`](#flags-for-config)
    - [`ipgeo` Command](#ipgeo-command)
      - [`ipgeo` Usage](#ipgeo-usage)
      - [Flags for `ipgeo`](#flags-for-ipgeo)
      - [Developer Plan Examples](#developer-plan-examples)
      - [Standard Plan Examples](#standard-plan-examples)
      - [Advanced Plan Examples](#advanced-plan-examples)
    - [`bulk-ip-geo` Command](#bulk-ip-geo-command)
      - [`bulk-ip-geo` Usage](#bulk-ip-geo-usage)
      - [Flags for `bulk-ip-geo`](#flags-for-bulk-ip-geo)
    - [Output Formats](#output-formats)
    - [`ip-security` Command](#ip-security-command)
      - [`ip-security` Usage](#ip-security-usage)
      - [Flags for `ip-security`](#flags-for-ip-security)
      - [Get Default Fields in `ip-security`](#get-default-fields-in-ip-security)
      - [Include Multiple Optional Fields](#include-multiple-optional-fields)
      - [Request with Field Filtering `ip-security`](#request-with-field-filtering-ip-security)
    - [`bulk-ip-security` Command](#bulk-ip-security-command)
      - [`bulk-ip-security` Usage](#bulk-ip-security-usage)
      - [Flags for `bulk-ip-security`](#flags-for-bulk-ip-security)
      - [`bulk-ip-security` Examples](#bulk-ip-security-examples)
    - [`asn` Command](#asn-command)
      - [`asn` Usage](#asn-usage)
      - [Flags for `asn`](#flags-for-asn)
      - [Get ASN Information of your IP](#get-asn-information-of-your-ip)
      - [Get ASN Information by ASN Number](#get-asn-information-by-asn-number)
      - [Combine All objects using Include](#combine-all-objects-using-include)
    - [`abuse` Command](#abuse-command)
      - [`abuse` Usage](#abuse-usage)
      - [Flags for `abuse`](#flags-for-abuse)
      - [Get abuse info about your current IP](#get-abuse-info-about-your-current-ip)
      - [Lookup a specific IP](#lookup-a-specific-ip)
      - [Exclude unnecessary fields](#exclude-unnecessary-fields)
      - [Lookup Abuse Contact with Specific Fields](#lookup-abuse-contact-with-specific-fields)
    - [`timezone` Command](#timezone-command)
      - [`timezone` Usage](#timezone-usage)
      - [Flags for `timezone`](#flags-for-timezone)
      - [Get timezone info about your current IP](#get-timezone-info-about-your-current-ip)
      - [Get Timezone for a specific IP Address](#get-timezone-for-a-specific-ip-address)
      - [Get Timezone by Timezone Name](#get-timezone-by-timezone-name)
      - [Get Timezone from Any Address](#get-timezone-from-any-address)
      - [Get Timezone from Location Coordinates](#get-timezone-from-location-coordinates)
      - [Get Timezone and Airport Details from IATA Code](#get-timezone-and-airport-details-from-iata-code)
      - [Get Timezone and City Details from UN/LOCODE](#get-timezone-and-city-details-from-unlocode)
    - [`time-conversion` Command](#time-conversion-command)
      - [`time-conversion` Usage](#time-conversion-usage)
      - [Flags for `time-conversion`](#flags-for-time-conversion)
      - [Convert Current Time from One Timezone to Another](#convert-current-time-from-one-timezone-to-another)
    - [`astronomy` Command](#astronomy-command)
      - [`astronomy` Usage](#astronomy-usage)
      - [Flags for `astronomy`](#flags-for-astronomy)
      - [Lookup Astronomy API by Coordinates](#lookup-astronomy-api-by-coordinates)
      - [Lookup Astronomy API by IP Address](#lookup-astronomy-api-by-ip-address)
      - [Lookup Astronomy API by Location String](#lookup-astronomy-api-by-location-string)
    - [`astronomy-timeseries` Command](#astronomy-timeseries-command)
      - [`astronomy-timeseries` Usage](#astronomy-timeseries-usage)
      - [Flags for `astronomy-timeseries`](#flags-for-astronomy-timeseries)
      - [Get astronomy timeseries info for an IP address](#get-astronomy-timeseries-info-for-an-ip-address)
      - [Get astronomy timeseries about a specific lat and lon](#get-astronomy-timeseries-about-a-specific-lat-and-lon)
      - [Get astronomy timeseries about a specific location](#get-astronomy-timeseries-about-a-specific-location)
    - [`parse-user-agent` Command](#parse-user-agent-command)
      - [`parse-user-agent` Usage](#parse-user-agent-usage)
      - [Flags for `parse-user-agent`](#flags-for-parse-user-agent)
      - [Parse a user agent string](#parse-a-user-agent-string)
    - [`parse-bulk-user-agents` Command](#parse-bulk-user-agents-command)
      - [`parse-bulk-user-agents` Usage](#parse-bulk-user-agents-usage)
      - [Flags for `parse-bulk-user-agents`](#flags-for-parse-bulk-user-agents)
      - [Parse multiple user agent strings](#parse-multiple-user-agent-strings)
- [License](#license)

## Requirements

- Go 1.18+
- API Key from [IPGeolocation.io](https://ipgeolocation.io)

## Installation
### Go install

To install `ipgeolocation` using `go install`, run:
```bash
go install github.com/IPGeolocation/cli/cmd/ipgeolocation@latest
```

Make sure `$GOBIN` or `$GOPATH/bin` is in your `PATH`, then run:

```bash
ipgeolocation --help
```

### Download and build from source
```bash
git clone https://github.com/IPGeolocation/cli.git
cd cli
go build -o ipgeolocation ./cmd/ipgeolocation
./ipgeolocation --help
```
Here you will need to add `$GOBIN` or `$GOPATH/bin` is in your `PATH` as well.

### Using Download Prebuilt Binaries

#### Overview
These are prebuilt binaries for the IPGeolocation CLI tool, version **1.0.3**. Users can download these files directly from [GitHub Releases](https://github.com/IPGeolocation/cli/releases) or from the table below without needing to build from source.

The CLI provides geolocation information, timezone, user-agent parsing, bulk IP lookups, and more.

#### Prebuilt Binaries

| Platform | Architecture | File Name / Downoad Link                                                                                                                             |
|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| Linux    | amd64        | [**ipgeolocation-1.0.3-linux-amd64.tar.gz**](https://github.com/IPGeolocation/cli/releases/download/v1.0.3/ipgeolocation-1.0.3-linux-amd64.tar.gz)   |
| Linux    | arm64        | [**ipgeolocation-1.0.3-linux-arm64.tar.gz**](https://github.com/IPGeolocation/cli/releases/download/v1.0.3/ipgeolocation-1.0.3-linux-arm64.tar.gz)   |
| macOS    | amd64        | [**ipgeolocation-1.0.3-darwin-amd64.tar.gz**](https://github.com/IPGeolocation/cli/releases/download/v1.0.3/ipgeolocation-1.0.3-darwin-amd64.tar.gz) |
| macOS    | arm64        | [**ipgeolocation-1.0.3-darwin-arm64.tar.gz**](https://github.com/IPGeolocation/cli/releases/download/v1.0.3/ipgeolocation-1.0.3-darwin-arm64.tar.gz) |
| Windows  | amd64        | [**ipgeolocation-1.0.3-windows-amd64.zip**](https://github.com/IPGeolocation/cli/releases/download/v1.0.3/ipgeolocation-1.0.3-windows-amd64.zip)     |

#### Installation Instructions

**Linux:**
1. Download the `.tar.gz` file for your architecture.
2. Extract it to a folder in your PATH, e.g., `/usr/local/bin`:
   ```bash
   tar -xzf ipgeolocation-1.0.3-linux-amd64.tar.gz -C /usr/local/bin
   ```
3. Rename the binary for simplicity:
   ```bash
   mv /usr/local/bin/ipgeolocation-1.0.3-linux-amd64 /usr/local/bin/ipgeolocation
   ```
4. Make the binary executable:
   ```bash
   chmod +x /usr/local/bin/ipgeolocation
   ```
5. Verify installation:
   ```bash
   ipgeolocation --help
   ```

**macOS:**
1. Download the `.tar.gz` file for your architecture (amd64 or arm64).
2. Extract to a folder in your PATH, e.g., `/usr/local/bin`:
   ```bash
   tar -xzf ipgeolocation-1.0.3-darwin-amd64.tar.gz -C /usr/local/bin
   ```
3. Rename the binary:
   ```bash
   mv /usr/local/bin/ipgeolocation-1.0.3-darwin-amd64 /usr/local/bin/ipgeolocation
   ```
4. Make executable:
   ```bash
   chmod +x /usr/local/bin/ipgeolocation
   ```
5. Verify installation:
   ```bash
   ipgeolocation --help
   ```

**Windows:**
1. Download the `.zip` file.
2. Extract the `ipgeolocation-1.0.3-windows-amd64.exe` to a folder included in your system `PATH`.
3. Rename the binary to `ipgeolocation.exe` for convenience.
4. Open Command Prompt and verify:
   ```cmd
   ipgeolocation --help
   ```

> [!NOTE]
> - Ensure execution permissions on Linux/macOS.
> - Recommended folder for binaries: `/usr/local/bin` or any folder in your PATH.
> - Prebuilt binaries include version **1.0.3** in their filename. Rename them after extraction for easier usage.
> - For updates, check GitHub Releases.

### Troubleshooting
- **Command not found:** Ensure the binary is in a folder included in your `PATH`.
- **Execution permission error:** Run `chmod +x <binary>` on Linux/macOS.
- **Wrong architecture:** Download the binary matching your OS and CPU architecture.
- **Go install issues:** Use `GOPROXY=direct` if Go module proxy caching creates issues.

## API Plan Tiers and Documentation

The documentation below corresponds to the four available API tier plans:

- **Developer Plan** (Free): [Full Documentation](https://ipgeolocation.io/ip-location-api.html#Free)
- **Standard Plan**: [Full Documentation](https://ipgeolocation.io/ip-location-api.html#Standard)
- **Advance Plan**: [Full Documentation](https://ipgeolocation.io/ip-location-api.html#Advance)
- **Security Plan**: [Full Documentation](https://ipgeolocation.io/ip-security-api.html#documentation-overview)

For a detailed comparison of what each plan offers, visit the [Pricing Page](https://ipgeolocation.io/pricing.html).

## Fields and commands Availability
IP Geolocation offers four plans from billing point of view: **Free, Standard, Security, Advance**. The availability of each command, over all plans are presented below.

| Sub Command              | Details                                                                                                                      | Free | Standard | Security | Advance |
|--------------------------|------------------------------------------------------------------------------------------------------------------------------|:----:|:--------:|:--------:|:-------:|
| `config`                 | Set up and check API Key.                                                                                                    |  ✔   |    ✔     |    ✔     |    ✔    |
| `ipgeo`                  | Get geolocation data for a single IP address, along with network, currency, abuse, timezone, security, asn, useragent, etc.  |  ✔   |    ✔     |    ✖     |    ✔    |
| `bulk-ip-geo`            | Get geolocation data for multiple IP addresses in a single API request, along with same data as in `ipgeo`.                  |  ✖   |    ✔     |    ✖     |    ✔    |
| `ip-security`            | Get security information (VPN, TOR, proxy, etc.) for a single IP, along with network, timezone, location, and currency, etc. |  ✖   |    ✖     |    ✔     |    ✖    |
| `bulk-ip-security`       | Get security threat intelligence for multiple IP addresses, along with same data as in `ip-security`.                        |  ✖   |    ✖     |    ✔     |    ✖    |
| `asn`                    | Get details of any AS number or IP address associated ASN.                                                                   |  ✖   |    ✖     |    ✖     |    ✔    |
| `abuse`                  | Get abuse reporting contact information for a given IP address.                                                              |  ✖   |    ✖     |    ✖     |    ✔    |
| `astronomy`              | Get sunrise, sunset, moonrise, moonset, and related data for a location.                                                     |  ✔   |    ✔     |    ✔     |    ✔    |
| `astronomy-timeseries`   | Get astronomy information for given date range at once.                                                                      |  ✔   |    ✔     |    ✔     |    ✔    |
| `timezone`               | Get timezone details using IP address, city, coordinates, IATA, ICAO, UNLOCODE, or timezone ID.                              |  ✔   |    ✔     |    ✔     |    ✔    |
| `time-conversion`        | Convert time between two specified timezones using city, coordinates, IATA, ICAO, UNLOCODE, or timezone ID.                  |  ✔   |    ✔     |    ✔     |    ✔    |
| `parse-user-agent`       | Parse single User Agent string.                                                                                              |  ✔   |    ✔     |    ✔     |    ✔    |
| `parse-bulk-user-agents` | Parse multiple User Agent String at once.                                                                                    |  ✖   |    ✔     |    ✔     |    ✔    |

> [!TIP]
> The availability of fields in every API endpoint across all API plans is provided in the **_Reference Table_** within each respective API Documentation. e.g., for IPGeolocationApi, please visit [https://ipgeolocation.io/ip-location-api.html#reference-to-ipgeolocation-api-response](https://ipgeolocation.io/ip-location-api.html#reference-to-ipgeolocation-api-response).

## Authentication Setup
To authenticate API requests, you need to get an API key from [ipgeolocation.io](https://ipgeolocation.io/).

### How to Get Your API Key

1. **Sign up** here: [https://app.ipgeolocation.io/signup](https://app.ipgeolocation.io/signup)
2. **(optional)** Verify your email, if you signed up using email.
3. **Log in** to your account: [https://app.ipgeolocation.io/login](https://app.ipgeolocation.io/login)
4. After logging in, navigate to your **Dashboard** to find your API key: [https://app.ipgeolocation.io/dashboard](https://app.ipgeolocation.io/dashboard)

### ApiKeyAuth
Once you've obtained the api key, configure your API client as follows:

The client must configure the authentication and authorization parameters in accordance with the API server security policy.
```bash
ipgeolocation config --apikey=<your-key>
```

## Global Flags
These flags are available for all commands:

| Flag         | Description                |
|--------------|----------------------------|
| `-h, --help` | Show help for the command. |

> [!TIP]
> You can also check the version for `ipgeolocation` using the `--version` flag:

```bash
ipgeolocation --version
```

## Commands

### `config` Command
Configure your API key for **[ipgeolocation.io](https://ipgeolocation.io/)**.

#### `config` Usage
```bash
ipgeolocation config --apikey=<your-key>
```

#### Flags for `config`
| Flag       | Type   | Description                                                     |
|------------|--------|-----------------------------------------------------------------|
| `--apikey` | string | Your API key from [ipgeolocation.io](https://ipgeolocation.io). |


### `ipgeo` Command
Lookup geolocation information for a **single IP address or domain** from the `ipgeolocation.io` API.

#### `ipgeo` Usage
```bash
ipgeolocation ipgeo [flags]
```

#### Flags for `ipgeo`
| Flag         | Type     | Default  | Description                                                                     |
|--------------|----------|----------|---------------------------------------------------------------------------------|
| `--ip`       | string   | `""`     | IPv4, IPv6, or domain name (e.g. `8.8.8.8`, `google.com`).                      |
| `--include`  | string[] | `[]`     | Include extra fields (e.g. `hostname,dma,security,abuse,time_zone,user_agent`). |
| `--fields`   | string[] | `[]`     | Return only specific fields (e.g. `location,network.asn.organization`).         |
| `--excludes` | string[] | `[]`     | Exclude fields from output.                                                     |
| `--lang`     | string   | `""`     | Response language.                                                              |
| `--output`   | string   | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`.                                |

> [!NOTE]
> Available language options can be found here: [https://ipgeolocation.io/ip-location-api.html#response-in-multiple-languages](https://ipgeolocation.io/ip-location-api.html#response-in-multiple-languages)

#### Developer Plan Examples
Get info about your current IP:
```bash
ipgeolocation ipgeo
```

##### Get Default Fields in developer plan
Lookup a specific IP:
```bash
ipgeolocation ipgeo --ip 8.8.8.8
```
Sample output:
```json
{
  "country_metadata": {
    "calling_code": "+1",
    "languages": [
      "en-US",
      "es-US",
      "haw",
      "fr"
    ],
    "tld": ".us"
  },
  "currency": {
    "code": "USD",
    "name": "US Dollar",
    "symbol": "$"
  },
  "ip": "8.8.8.8",
  "location": {
    "city": "Mountain View",
    "continent_code": "NA",
    "continent_name": "North America",
    "country_capital": "Washington, D.C.",
    "country_code2": "US",
    "country_code3": "USA",
    "country_emoji": "🇺🇸",
    "country_flag": "https://ipgeolocation.io/static/flags/us_64.png",
    "country_name": "United States",
    "country_name_official": "United States of America",
    "district": "Santa Clara",
    "geoname_id": "6301403",
    "is_eu": false,
    "latitude": "37.42240",
    "longitude": "-122.08421",
    "state_code": "US-CA",
    "state_prov": "California",
    "zipcode": "94043-1351"
  }
}
```

Filtering response with the use of `--fields` and `--excludes`:
```bash
ipgeolocation ipgeo --ip 8.8.8.8 --fields=location --excludes location.continent_code,location.continent_name
```
Sample output:
```json
{
  "ip": "8.8.8.8",
  "location": {
    "city": "Mountain View",
    "country_capital": "Washington, D.C.",
    "country_code2": "US",
    "country_code3": "USA",
    "country_emoji": "🇺🇸",
    "country_flag": "https://ipgeolocation.io/static/flags/us_64.png",
    "country_name": "United States",
    "country_name_official": "United States of America",
    "district": "Santa Clara",
    "geoname_id": "6301403",
    "is_eu": false,
    "latitude": "37.42240",
    "longitude": "-122.08421",
    "state_code": "US-CA",
    "state_prov": "California",
    "zipcode": "94043-1351"
  }
}
```

#### Standard Plan Examples
##### Get Default Fields in standard plan
Lookup domain name with default fields
```bash
ipgeolocation ipgeo --ip google.com
```
Sample output:
```json
{
  "country_metadata": {
    "calling_code": "+49",
    "languages": [
      "de"
    ],
    "tld": ".de"
  },
  "currency": {
    "code": "EUR",
    "name": "Euro",
    "symbol": "€"
  },
  "domain": "google.com",
  "ip": "142.250.185.174",
  "location": {
    "city": "Frankfurt",
    "continent_code": "EU",
    "continent_name": "Europe",
    "country_capital": "Berlin",
    "country_code2": "DE",
    "country_code3": "DEU",
    "country_emoji": "🇩🇪",
    "country_flag": "https://ipgeolocation.io/static/flags/de_64.png",
    "country_name": "Germany",
    "country_name_official": "Federal Republic of Germany",
    "district": "Frankfurt",
    "geoname_id": "6463469",
    "is_eu": true,
    "latitude": "50.11208",
    "longitude": "8.68341",
    "state_code": "DE-HE",
    "state_prov": "Hesse",
    "zipcode": "60311"
  },
  "network": {
    "asn": {
      "as_number": "AS15169",
      "country": "US",
      "organization": "Google LLC"
    },
    "company": {
      "name": "Google LLC"
    }
  }
}
```
##### Retrieving Geolocation Data in Multiple Languages
Here is an example to get the geolocation data for IP address '2001:4230:4890::1' in French `fr` language:
```bash
ipgeolocation ipgeo --ip 2001:4230:4890::1 --lang fr
```

Sample output:
```json
{
  "country_metadata": {
    "calling_code": "+230",
    "languages": [
      "en-MU",
      "bho",
      "fr"
    ],
    "tld": ".mu"
  },
  "currency": {
    "code": "MUR",
    "name": "Mauritius Rupee",
    "symbol": "₨"
  },
  "ip": "2001:4230:4890:0:0:0:0:1",
  "location": {
    "city": "Quatre Bornes",
    "continent_code": "AF",
    "continent_name": "Afrique",
    "country_capital": "Port Louis",
    "country_code2": "MU",
    "country_code3": "MUS",
    "country_emoji": "🇲🇺",
    "country_flag": "https://ipgeolocation.io/static/flags/mu_64.png",
    "country_name": "Maurice",
    "country_name_official": "",
    "district": "Quatre Bornes",
    "geoname_id": "1106777",
    "is_eu": false,
    "latitude": "-20.24304",
    "longitude": "57.49631",
    "state_code": "MU-PW",
    "state_prov": "Wilhems des plaines",
    "zipcode": "72201"
  },
  "network": {
    "asn": {
      "as_number": "AS52095",
      "country": "CZ",
      "organization": "Netassist International s.r.o."
    },
    "company": {
      "name": "African Network Information Center AfriNIC Ltd"
    }
  }
}
```

##### Include HostName, Timezone and User-Agent
```bash
ipgeolocation ipgeo --ip 219.100.37.207 --include hostname,time_zone
```
Sample output:
```json
{
  "country_metadata": {
    "calling_code": "+81",
    "languages": [
      "ja"
    ],
    "tld": ".jp"
  },
  "currency": {
    "code": "JPY",
    "name": "Yen",
    "symbol": "¥"
  },
  "hostname": "public-vpn-13-15.vpngate.v4.open.ad.jp",
  "ip": "219.100.37.207",
  "location": {
    "city": "Tokyo",
    "continent_code": "AS",
    "continent_name": "Asia",
    "country_capital": "Tokyo",
    "country_code2": "JP",
    "country_code3": "JPN",
    "country_emoji": "🇯🇵",
    "country_flag": "https://ipgeolocation.io/static/flags/jp_64.png",
    "country_name": "Japan",
    "country_name_official": "Japan",
    "district": "Koto",
    "geoname_id": "12149683",
    "is_eu": false,
    "latitude": "35.68467",
    "longitude": "139.80881",
    "state_code": "JP-13",
    "state_prov": "Tokyo Metropolis",
    "zipcode": "135-0022"
  },
  "network": {
    "asn": {
      "as_number": "AS36599",
      "country": "US",
      "organization": "SoftEther Telecommunication Research Institute, LLC"
    },
    "company": {
      "name": "SoftEther Corporation"
    }
  },
  "time_zone": {
    "current_time": "2026-01-13 17:30:42.949+0900",
    "current_time_unix": 1768293042.949,
    "current_tz_abbreviation": "JST",
    "current_tz_full_name": "Japan Standard Time",
    "dst_end": {},
    "dst_exists": false,
    "dst_savings": 0,
    "dst_start": {},
    "dst_tz_abbreviation": "",
    "dst_tz_full_name": "",
    "is_dst": false,
    "name": "Asia/Tokyo",
    "offset": 9,
    "offset_with_dst": 9,
    "standard_tz_abbreviation": "JST",
    "standard_tz_full_name": "Japan Standard Time"
  }
}
```
> [!NOTE]
>
> The IP Geolocation API supports hostname lookup for all paid subscriptions. However, this is not included by default. To enable hostname resolution, use the `--include` parameter with one of the following options:
>
> - `hostname`: Performs a quick lookup using the internal hostname database. If no match is found, the IP is returned as-is. This is fast but may produce incomplete results.
> - `liveHostname`: Queries live sources for accurate hostname resolution. This may increase response time.
> - `hostnameFallbackLive`: Attempts the internal database first, and falls back to live sources if no result is found. This option provides a balance of speed and reliability.

#### Advanced Plan Examples
##### Include DMA, Abuse and Security
```bash
ipgeolocation ipgeo --ip 8.8.8.8 --include dma,abuse,security
```

Sample output:
```json
{
  "ip": "8.8.8.8",
  "location": {
    "continent_code": "NA",
    "continent_name": "North America",
    "country_code2": "US",
    "country_code3": "USA",
    "country_name": "United States",
    "country_name_official": "United States of America",
    "country_capital": "Washington, D.C.",
    "state_prov": "California",
    "state_code": "US-CA",
    "district": "Santa Clara",
    "city": "Mountain View",
    "zipcode": "94043-1351",
    "latitude": "37.42240",
    "longitude": "-122.08421",
    "is_eu": false,
    "geoname_id": "6301403",
    "accuracy_radius": "25.388",
    "locality": "Mountain View",
    "country_emoji": "🇺🇸",
    "country_flag": "https://ipgeolocation.io/static/flags/us_64.png",
    "dma_code": "807"
  },
  "country_metadata": {
    "calling_code": "+1",
    "tld": ".us",
    "languages": [
      "en-US",
      "es-US",
      "haw",
      "fr"
    ]
  },
  "network": {
    "asn": {
      "as_number": "AS15169",
      "organization": "Google LLC",
      "country": "US",
      "asn_name": "GOOGLE",
      "type": "BUSINESS",
      "domain": "google.com",
      "date_allocated": "",
      "allocation_status": "",
      "num_of_ipv4_routes": "1026",
      "num_of_ipv6_routes": "106",
      "rir": "ARIN"
    },
    "connection_type": "",
    "company": {
      "name": "Google LLC",
      "type": "business",
      "domain": "google.com"
    }
  },
  "currency": {
    "code": "USD",
    "name": "US Dollar",
    "symbol": "$"
  },
  "security": {
    "threat_score": 0,
    "is_tor": false,
    "is_proxy": false,
    "proxy_type": "",
    "proxy_provider": "",
    "is_anonymous": false,
    "is_known_attacker": false,
    "is_spam": false,
    "is_bot": false,
    "is_cloud_provider": false,
    "cloud_provider": ""
  },
  "abuse": {
    "route": "8.8.8.0/24",
    "country": "",
    "handle": "ABUSE5250-ARIN",
    "name": "Abuse",
    "organization": "Abuse",
    "role": "abuse",
    "kind": "group",
    "address": "1600 Amphitheatre Parkway\nMountain View\nCA\n94043\nUnited States",
    "emails": [
      "network-abuse@google.com"
    ],
    "phone_numbers": [
      "+1-650-253-0000"
    ]
  }
}
```

> [!NOTE]
> All features available in the Free plan are also included in the Standard and Advanced plans. Similarly, all features of the Standard plan are available in the Advanced plan.

Get only specific fields in YAML:
```bash
ipgeolocation ipgeo --ip 1.1.1.1 --fields location --output yaml
```
Sample response:
```yaml
ip: 1.1.1.1
location:
    accuracy_radius: "8.778"
    city: Hong Kong
    confidence: high
    continent_code: AS
    continent_name: Asia
    country_capital: Hong Kong
    country_code2: HK
    country_code3: HKG
    country_emoji: "\U0001F1ED\U0001F1F0"
    country_flag: https://ipgeolocation.io/static/flags/hk_64.png
    country_name: Hong Kong
    country_name_official: Hong Kong Special Administrative Region of the People's Republic of China
    district: Wan Chai District
    geoname_id: "10106797"
    is_eu: false
    latitude: "22.27683"
    locality: Hong Kong
    longitude: "114.17612"
    state_code: ""
    state_prov: Hong Kong SAR
    zipcode: ""

```
Similarly, `raw`, `table` and `pretty` formats are also available. Or one can parse the simple response with `| jq`.

### `bulk-ip-geo` Command
Lookup geolocation information for **multiple IP addresses** in one request.

#### `bulk-ip-geo` Usage
```bash
ipgeolocation bulk-ip-geo [flags]
```

#### Flags for `bulk-ip-geo`

| Flag            | Type     | Default  | Description                                                   |
|-----------------|----------|----------|---------------------------------------------------------------|
| `--ips`         | string[] | `[]`     | Comma-separated list of IPs. Example: `--ips 8.8.8.8,1.1.1.1` |
| `--file`        | string   | `""`     | Path to a text file containing IPs (one per line).            |
| `--include`     | string[] | `[]`     | Include extra fields (e.g. `location,time_zone`).             |
| `--excludes`    | string[] | `[]`     | Exclude fields (e.g. `currency`).                             |
| `--fields`      | string[] | `[]`     | Return only specific fields (e.g. `location`).                |
| `--lang`        | string   | `""`     | Response language (if supported).                             |
| `--output`      | string   | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`.              |
| `--output-file` | string   | `""`     | Save output to JSON file. Example: `--output-file results`    |


For further information, please visit [IP Geolocation API Documentation](https://ipgeolocation.io/ip-location-api.html#documentation-overview).

Lookup 3 IP addresses:
```bash
ipgeolocation bulk-ip-geo --ips 8.8.8.8,1.1.1.1,192.30.253.112
```

Lookup from a file:
```bash
ipgeolocation bulk-ip-geo --file=ips.txt --output-file results
```

Results will be written to the `results.json` file.

> [!NOTE]
> All the `include`, `exclude`, `fields` parameters can be used just like `ipgeo` command.

Output as YAML:
```bash
ipgeolocation bulk-ip-geo --ips=8.8.8.8 --output=yaml
```

Save results to JSON file:
```bash
ipgeolocation bulk-ip-geo --ips=8.8.8.8,1.1.1.1 --output-file=output.json
```

#### Output Formats
- **pretty** (default): Human-readable formatted JSON.  
- **raw**: Raw API response.  
- **table**: Tabular display of common fields.  
- **yaml**: YAML-formatted output.  
- **json file**: If `--output-file` is provided, results are saved to a `.json` file.  

### `ip-security` Command
Lookup IP security information using the `ipgeolocation.io` API.

#### `ip-security` Usage
```bash
ipgeolocation ip-security [flags]
```

#### Flags for `ip-security`
| Flag         | Type     | Default  | Description                                                                                                        |
|--------------|----------|----------|--------------------------------------------------------------------------------------------------------------------|
| `--ip`       | string   | `""`     | IPv4 or IPv6 address.                                                                                              |
| `--include`  | string[] | `[]`     | Include extra fields in output. (e.g., `location,network,currency,time_zone,user_agent,country_metadata,hostname`) |
| `--excludes` | string[] | `[]`     | Exclude fields from output.                                                                                        |
| `--fields`   | string[] | `[]`     | Return only specific fields (e.g. `network.asn.organization`).                                                     |
| `--lang`     | string   | `""`     | Response language.                                                                                                 |
| `--output`   | string   | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`.                                                                   |

#### Get Default Fields in `ip-security`
Get info about your current IP:
```bash
ipgeolocation ip-security
```

Lookup a specific IP:
```bash
ipgeolocation ip-security --ip 8.8.8.8
```
Sample output:
```json
{
  "ip": "2.56.188.34",
  "security": {
    "threat_score": 80,
    "is_tor": false,
    "is_proxy": true,
    "proxy_type": "VPN",
    "proxy_provider": "Nord VPN",
    "is_anonymous": true,
    "is_known_attacker": true,
    "is_spam": false,
    "is_bot": false,
    "is_cloud_provider": true,
    "cloud_provider": "Packethub S.A."
  }
}
```

#### Include Multiple Optional Fields

Lookup security information for a specific IP with including additional fields:
```bash
ipgeolocation ip-security --ip 2.56.188.34  --include location,network,currency,time_zone,user_agent,country_metadata,hostname
```

Output as YAML:
```bash
ipgeolocation ip-security --ip 8.8.8.8 --output=yaml
```

> [!NOTE]
> You can get all the available fields in standard plan in combination with security data, when subscribed to security plan.

For full API specifications, refer to the [official IP Security API documentation](https://ipgeolocation.io/ip-security-api.html#documentation-overview).

#### Request with Field Filtering `ip-security`
```bash
ipgeolocation ip-security --ip 195.154.221.54 --fields security.is_tor,security.is_proxy,security.is_bot,security.is_spam
```
Sample output:
```json
{
  "ip": "195.154.221.54",
  "security": {
    "is_bot": false,
    "is_proxy": true,
    "is_spam": false,
    "is_tor": false
  }
}
```

### `bulk-ip-security` Command

Lookup IP security information for **multiple IP addresses** in one request. 

#### `bulk-ip-security` Usage
```bash
ipgeolocation bulk-ip-security [flags]
```

#### Flags for `bulk-ip-security`
| Flag            | Type     | Default  | Description                                                                                             |
|-----------------|----------|----------|---------------------------------------------------------------------------------------------------------|
| `--ips`         | string[] | `[]`     | Comma-separated list of IPs. Example: `--ips 8.8.8.8,1.1.1.1`                                           |
| `--file`        | string   | `""`     | Path to a text file containing IPs (one per line).                                                      |
| `--include`     | string[] | `[]`     | Include extra fields (e.g. `location,network,currency,time_zone,user_agent,country_metadata,hostname`). |
| `--excludes`    | string[] | `[]`     | Exclude fields (e.g. `currency`).                                                                       |
| `--fields`      | string[] | `[]`     | Return only specific fields (e.g. `location`).                                                          |
| `--lang`        | string   | `""`     | Response language (if supported).                                                                       |
| `--output`      | string   | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`.                                                        |
| `--output-file` | string   | `""`     | Save output to JSON file. Example: `--output-file results`                                              |

#### `bulk-ip-security` Examples
Lookup 3 IP addresses:
```bash
ipgeolocation bulk-ip-security --ips 8.8.8.8,1.1.1.1,192.30.253.112
```

Lookup from a file:
```bash
ipgeolocation bulk-ip-security --file=ips.txt --output-file=output
```

Include location and timezone:
```bash
ipgeolocation bulk-ip-security --file=ips.txt --include=location,network,time_zone
```

Output as YAML:
```bash
ipgeolocation bulk-ip-security --ips=8.8.8.8 --output=yaml
```

Save results to JSON file:
```bash
ipgeolocation bulk-ip-security --ips=8.8.8.8,1.1.1.1 --output-file=output
```

For further information, please visit [Bulk IP Security API Documentation](https://ipgeolocation.io/ip-security-api.html#bulk-ip-security-lookup-api).

### `asn` Command
Lookup ASN (Autonomous System Number) information using the `ipgeolocation.io` API.

#### `asn` Usage
```bash
ipgeolocation asn [flags]
```

#### Flags for `asn`
| Flag         | Type     | Default  | Description                                                                             |
|--------------|----------|----------|-----------------------------------------------------------------------------------------|
| `--ip`       | string   | `""`     | IPv4 or IPv6 address.                                                                   |
| `--include`  | string[] | `[]`     | Include extra fields in output.(e.g., `peers, downstreams, upstreams, whois_response`)  |
| `--excludes` | string[] | `[]`     | Exclude fields from output.                                                             |
| `--fields`   | string[] | `[]`     | Return only specific fields (e.g. `ip,organization`).                                   |
| `--output`   | string   | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`.                                        |

> [!NOTE]
> ASN API is only available in the Advanced Plan

#### Get ASN Information of your IP
Get ASN information about your current IP address:
```bash
ipgeolocation asn
```

#### Get ASN Information by ASN Number
Lookup a specific ASN. The AS number can be prefixed with 'AS' or 'as':
```bash
ipgeolocation asn --asn 15169
```
Sample output:
```json
{
  "asn": {
    "as_number": "AS15169",
    "organization": "Google LLC",
    "country": "US",
    "asn_name": "GOOGLE",
    "type": "BUSINESS",
    "domain": "google.com",
    "date_allocated": "",
    "allocation_status": "",
    "num_of_ipv4_routes": "1026",
    "num_of_ipv6_routes": "106",
    "rir": "ARIN"
  }
}
```

#### Combine All objects using Include
Lookup a specific IP with including additional fields:
```bash
ipgeolocation asn --asn 12 --include routes,peers,upstreams,downstreams,whois_response
```
Sample output:
```json
{
  "asn": {
    "as_number": "AS12",
    "organization": "New York University",
    "country": "US",
    "asn_name": "NYU-DOMAIN",
    "type": "EDUCATION",
    "domain": "nyu.edu",
    "date_allocated": "",
    "allocation_status": "",
    "num_of_ipv4_routes": "11",
    "num_of_ipv6_routes": "1",
    "rir": "ARIN",
    "routes": [
      "192.76.177.0/24",
      "...",
      "216.165.103.0/24"
    ],
    "upstreams": [
      {
        "as_number": "AS3269",
        "description": "Telecom Italia S.p.A.",
        "country": "IT"
      },
      "...",
      {
        "as_number": "AS137",
        "description": "Consortium GARR",
        "country": "IT"
      }
    ],
    "downstreams": [
      {
        "as_number": "AS394666",
        "description": "NYU Langone Health",
        "country": "US"
      },
      {
        "as_number": "AS54965",
        "description": "Polytechnic Institute of NYU",
        "country": "US"
      }
    ],
    "peers": [
      {
        "as_number": "AS3269",
        "description": "Telecom Italia S.p.A.",
        "country": "IT"
      },
      "...",
      {
        "as_number": "AS54965",
        "description": "Polytechnic Institute of NYU",
        "country": "US"
      }
    ],
    "whois_response": "<raw-whois-response>"
  }
}
```

Exclude unnecessary fields:
```bash
ipgeolocation asn --ip 8.8.8.8 --exclude asn.date_allocated,asn.allocation_status
```
Sample output:
```json
{
  "asn": {
    "as_number": "AS15169",
    "asn_name": "GOOGLE",
    "country": "US",
    "domain": "google.com",
    "num_of_ipv4_routes": "1024",
    "num_of_ipv6_routes": "135",
    "organization": "Google LLC",
    "rir": "ARIN",
    "type": "BUSINESS"
  },
  "ip": "8.8.8.8"
}
```

Get only specific fields in table format:
```bash
ipgeolocation asn --ip 8.8.8.8 --fields asn.organization,asn.country,asn.downstreams --output table
```

Output as YAML:
```bash
ipgeolocation asn --ip 8.8.8.8 --output=yaml
```

For further information, please visit [ASN API Documentation](https://ipgeolocation.io/asn-api.html#documentation-overview).

### `abuse` Command
Lookup IP abuse contact information using the `ipgeolocation.io` API.

#### `abuse` Usage
```bash
ipgeolocation abuse [flags]
```

#### Flags for `abuse`
| Flag         | Type     | Default  | Description                                           |
|--------------|----------|----------|-------------------------------------------------------|
| `--ip`       | string   | `""`     | IPv4 or IPv6 address.                                 |
| `--excludes` | string[] | `[]`     | Exclude fields from output.                           |
| `--fields`   | string[] | `[]`     | Return only specific fields (e.g. `ip,organization`). |
| `--output`   | string   | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`.      |

> [!NOTE]
> Abuse Contact API is only available in the Advanced Plan

#### Get abuse info about your current IP
```bash
ipgeolocation abuse
```

#### Lookup a specific IP
```bash
ipgeolocation abuse --ip 8.8.8.8
```

#### Exclude unnecessary fields
```bash
ipgeolocation abuse --ip 8.8.8.8 --exclude abuse.handle,abuse.emails
```

#### Lookup Abuse Contact with Specific Fields
```bash
ipgeolocation abuse --ip 8.8.8.8 --fields abuse.role,abuse.emails --output table
```
Sample output:
```text
Ip                  : 8.8.8.8
Abuse:
  Emails:
    [0]:
      network-abuse@google.com
  Role                : abuse
```

Get output as YAML:
```bash
ipgeolocation abuse --ip 8.8.8.8 --output=yaml
```
Sample output:
```yaml
abuse:
    address: |-
        1600 Amphitheatre Parkway
        Mountain View
        CA
        94043
        United States
    country: ""
    emails:
        - network-abuse@google.com
    handle: ABUSE5250-ARIN
    kind: group
    name: Abuse
    organization: Abuse
    phone_numbers:
        - +1-650-253-0000
    role: abuse
    route: 8.8.8.0/24
ip: 8.8.8.8
```

For further information, please visit [Abuse API Documentation](https://ipgeolocation.io/ip-abuse-contact-api.html#documentation-overview).

### `timezone` Command
Lookup timezone information using the `ipgeolocation.io` API using different query types like IP address, latitude/longitude, and timezone ID.

For full API specifications, refer to the [Timezone API documentation](https://ipgeolocation.io/timezone-api.html#documentation-overview).

#### `timezone` Usage
```bash
ipgeolocation timezone [flags]
```

#### Flags for `timezone`
| Flag          | Type    | Default  | Description                                      |
|---------------|---------|----------|--------------------------------------------------|
| `--ip`        | string  | `""`     | IPv4 or IPv6 address.                            |
| `--tz`        | string  | `""`     | IANA Timezone identifier.                        |
| `--location`  | string  | `""`     | Location name (e.g. New York)                    |
| `--latitude`  | float64 | `0`      | Latitude (e.g. 37.7749).                         |
| `--longitude` | float64 | `0`      | Longitude (e.g. -122.4194).                      |
| `--lang`      | string  | `""`     | Response language (if supported).                |
| `--iata`      | string  | `""`     | IATA code (e.g. DXB).                            |
| `--icao`      | string  | `""`     | ICAO code (e.g. KATL).                           |
| `--lo`        | string  | `""`     | LO code (e.g. DEBER).                            |
| `--output`    | string  | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Get timezone info about your current IP
```bash
ipgeolocation timezone
```

#### Get Timezone for a specific IP Address
```bash
ipgeolocation timezone --ip 8.8.8.8
```
Sample output:
```json
{
  "ip": "8.8.8.8",
  "location": {
    "city": "Mountain View",
    "continent_code": "NA",
    "continent_name": "North America",
    "country_code2": "US",
    "country_code3": "USA",
    "country_name": "United States",
    "country_name_official": "United States of America",
    "district": "Santa Clara",
    "is_eu": false,
    "latitude": "37.42240",
    "longitude": "-122.08421",
    "state_code": "US-CA",
    "state_prov": "California",
    "zipcode": "94043-1351"
  },
  "time_zone": {
    "current_tz_abbreviation": "PST",
    "current_tz_full_name": "Pacific Standard Time",
    "date": "2026-01-13",
    "date_time": "2026-01-13 04:10:37",
    "date_time_txt": "Tuesday, January 13, 2026 04:10:37",
    "date_time_unix": 1768306237.124,
    "date_time_wti": "Tue, 13 Jan 2026 04:10:37 -0800",
    "date_time_ymd": "2026-01-13T04:10:37-0800",
    "dst_end": {
      "date_time_after": "2026-11-01 TIME 01",
      "date_time_before": "2026-11-01 TIME 02",
      "duration": "-1H",
      "gap": false,
      "overlap": true,
      "utc_time": "2026-11-01 TIME 09"
    },
    "dst_exists": true,
    "dst_savings": 0,
    "dst_start": {
      "date_time_after": "2026-03-08 TIME 03",
      "date_time_before": "2026-03-08 TIME 02",
      "duration": "+1H",
      "gap": true,
      "overlap": false,
      "utc_time": "2026-03-08 TIME 10"
    },
    "dst_tz_abbreviation": "PDT",
    "dst_tz_full_name": "Pacific Daylight Time",
    "is_dst": false,
    "month": 1,
    "name": "America/Los_Angeles",
    "offset": -8,
    "offset_with_dst": -8,
    "standard_tz_abbreviation": "PST",
    "standard_tz_full_name": "Pacific Standard Time",
    "time_12": "04:10:37 AM",
    "time_24": "04:10:37",
    "week": 3,
    "year": 2026,
    "year_abbr": "26"
  }
}
```

#### Get Timezone by Timezone Name
```bash
ipgeolocation timezone --tz America/New_York
```
Sample output:
```json
{
  "time_zone": {
    "current_tz_abbreviation": "EST",
    "current_tz_full_name": "Eastern Standard Time",
    "date": "2026-01-13",
    "date_time": "2026-01-13 07:12:48",
    "date_time_txt": "Tuesday, January 13, 2026 07:12:48",
    "date_time_unix": 1768306368.774,
    "date_time_wti": "Tue, 13 Jan 2026 07:12:48 -0500",
    "date_time_ymd": "2026-01-13T07:12:48-0500",
    "dst_end": {
      "date_time_after": "2026-11-01 TIME 01",
      "date_time_before": "2026-11-01 TIME 02",
      "duration": "-1H",
      "gap": false,
      "overlap": true,
      "utc_time": "2026-11-01 TIME 06"
    },
    "dst_exists": true,
    "dst_savings": 0,
    "dst_start": {
      "date_time_after": "2026-03-08 TIME 03",
      "date_time_before": "2026-03-08 TIME 02",
      "duration": "+1H",
      "gap": true,
      "overlap": false,
      "utc_time": "2026-03-08 TIME 07"
    },
    "dst_tz_abbreviation": "EDT",
    "dst_tz_full_name": "Eastern Daylight Time",
    "is_dst": false,
    "month": 1,
    "name": "America/New_York",
    "offset": -5,
    "offset_with_dst": -5,
    "standard_tz_abbreviation": "EST",
    "standard_tz_full_name": "Eastern Standard Time",
    "time_12": "07:12:48 AM",
    "time_24": "07:12:48",
    "week": 3,
    "year": 2026,
    "year_abbr": "26"
  }
}
```

#### Get Timezone from Any Address
```bash
ipgeolocation timezone --location "New York, USA"
```

Sample output:
```json
{
  "location": {
    "city": "Syracuse",
    "country_name": "United States",
    "latitude": "43.06923",
    "locality": "Destiny USA Mall",
    "location_string": "New York, USA",
    "longitude": "-76.17249",
    "state_prov": "New York"
  },
  "time_zone": {
    "current_tz_abbreviation": "EST",
    "current_tz_full_name": "Eastern Standard Time",
    "date": "2026-01-13",
    "date_time": "2026-01-13 07:14:38",
    "date_time_txt": "Tuesday, January 13, 2026 07:14:38",
    "date_time_unix": 1768306478.948,
    "date_time_wti": "Tue, 13 Jan 2026 07:14:38 -0500",
    "date_time_ymd": "2026-01-13T07:14:38-0500",
    "dst_end": {
      "date_time_after": "2026-11-01 TIME 01",
      "date_time_before": "2026-11-01 TIME 02",
      "duration": "-1H",
      "gap": false,
      "overlap": true,
      "utc_time": "2026-11-01 TIME 06"
    },
    "dst_exists": true,
    "dst_savings": 0,
    "dst_start": {
      "date_time_after": "2026-03-08 TIME 03",
      "date_time_before": "2026-03-08 TIME 02",
      "duration": "+1H",
      "gap": true,
      "overlap": false,
      "utc_time": "2026-03-08 TIME 07"
    },
    "dst_tz_abbreviation": "EDT",
    "dst_tz_full_name": "Eastern Daylight Time",
    "is_dst": false,
    "month": 1,
    "name": "America/New_York",
    "offset": -5,
    "offset_with_dst": -5,
    "standard_tz_abbreviation": "EST",
    "standard_tz_full_name": "Eastern Standard Time",
    "time_12": "07:14:38 AM",
    "time_24": "07:14:38",
    "week": 3,
    "year": 2026,
    "year_abbr": "26"
  }
}
```

#### Get Timezone from Location Coordinates
```bash
ipgeolocation timezone --latitude 37.7749 --longitude -122.4194
```
Sample output:
```json
{
  "time_zone": {
    "current_tz_abbreviation": "PST",
    "current_tz_full_name": "Pacific Standard Time",
    "date": "2026-01-13",
    "date_time": "2026-01-13 04:15:42",
    "date_time_txt": "Tuesday, January 13, 2026 04:15:42",
    "date_time_unix": 1768306542.847,
    "date_time_wti": "Tue, 13 Jan 2026 04:15:42 -0800",
    "date_time_ymd": "2026-01-13T04:15:42-0800",
    "dst_end": {
      "date_time_after": "2026-11-01 TIME 01",
      "date_time_before": "2026-11-01 TIME 02",
      "duration": "-1H",
      "gap": false,
      "overlap": true,
      "utc_time": "2026-11-01 TIME 09"
    },
    "dst_exists": true,
    "dst_savings": 0,
    "dst_start": {
      "date_time_after": "2026-03-08 TIME 03",
      "date_time_before": "2026-03-08 TIME 02",
      "duration": "+1H",
      "gap": true,
      "overlap": false,
      "utc_time": "2026-03-08 TIME 10"
    },
    "dst_tz_abbreviation": "PDT",
    "dst_tz_full_name": "Pacific Daylight Time",
    "is_dst": false,
    "month": 1,
    "name": "America/Los_Angeles",
    "offset": -8,
    "offset_with_dst": -8,
    "standard_tz_abbreviation": "PST",
    "standard_tz_full_name": "Pacific Standard Time",
    "time_12": "04:15:42 AM",
    "time_24": "04:15:42",
    "week": 3,
    "year": 2026,
    "year_abbr": "26"
  }
}
```

#### Get Timezone and Airport Details from IATA Code
```bash
ipgeolocation timezone --iata DXB
```
Sample output:
```json
{
  "airport_details": {
    "city": "Dubai",
    "continent_code": "AS",
    "country_code": "AE",
    "elevation_ft": 62,
    "faa_code": "",
    "iata_code": "DXB",
    "icao_code": "OMDB",
    "latitude": "25.25280",
    "longitude": "55.36440",
    "name": "Dubai International Airport",
    "state_code": "AE-DU",
    "type": "large_airport"
  },
  "time_zone": {
    "current_tz_abbreviation": "GST",
    "current_tz_full_name": "Gulf Standard Time",
    "date": "2026-01-13",
    "date_time": "2026-01-13 16:16:46",
    "date_time_txt": "Tuesday, January 13, 2026 16:16:46",
    "date_time_unix": 1768306606.662,
    "date_time_wti": "Tue, 13 Jan 2026 16:16:46 +0400",
    "date_time_ymd": "2026-01-13T16:16:46+0400",
    "dst_end": {},
    "dst_exists": false,
    "dst_savings": 0,
    "dst_start": {},
    "dst_tz_abbreviation": "",
    "dst_tz_full_name": "",
    "is_dst": false,
    "month": 1,
    "name": "Asia/Dubai",
    "offset": 4,
    "offset_with_dst": 4,
    "standard_tz_abbreviation": "GST",
    "standard_tz_full_name": "Gulf Standard Time",
    "time_12": "04:16:46 PM",
    "time_24": "16:16:46",
    "week": 3,
    "year": 2026,
    "year_abbr": "26"
  }
}
```
> [!NOTE]
> Similarly, you can fetch Airport Details and Timezone from using any ICAO code as well.

#### Get Timezone and City Details from UN/LOCODE
```bash
ipgeolocation timezone --lo DEBER
```
Sample ouput:
```json
{
  "lo_code_details": {
    "city": "Berlin",
    "country_code": "DE",
    "country_name": "",
    "latitude": "52.51667",
    "lo_code": "DEBER",
    "location_type": "Port, Rail Terminal, Road Terminal, Airport, Postal Exchange",
    "longitude": "13.38333",
    "state_code": "BE"
  },
  "time_zone": {
    "current_tz_abbreviation": "CET",
    "current_tz_full_name": "Central European Standard Time",
    "date": "2026-01-13",
    "date_time": "2026-01-13 13:18:15",
    "date_time_txt": "Tuesday, January 13, 2026 13:18:15",
    "date_time_unix": 1768306695.348,
    "date_time_wti": "Tue, 13 Jan 2026 13:18:15 +0100",
    "date_time_ymd": "2026-01-13T13:18:15+0100",
    "dst_end": {
      "date_time_after": "2026-10-25 TIME 02",
      "date_time_before": "2026-10-25 TIME 03",
      "duration": "-1H",
      "gap": false,
      "overlap": true,
      "utc_time": "2026-10-25 TIME 01"
    },
    "dst_exists": true,
    "dst_savings": 0,
    "dst_start": {
      "date_time_after": "2026-03-29 TIME 03",
      "date_time_before": "2026-03-29 TIME 02",
      "duration": "+1H",
      "gap": true,
      "overlap": false,
      "utc_time": "2026-03-29 TIME 01"
    },
    "dst_tz_abbreviation": "CEST",
    "dst_tz_full_name": "Central European Summer Time",
    "is_dst": false,
    "month": 1,
    "name": "Europe/Berlin",
    "offset": 1,
    "offset_with_dst": 1,
    "standard_tz_abbreviation": "CET",
    "standard_tz_full_name": "Central European Standard Time",
    "time_12": "01:18:15 PM",
    "time_24": "13:18:15",
    "week": 3,
    "year": 2026,
    "year_abbr": "26"
  }
}
```

Output as YAML:
```bash
ipgeolocation timezone --ip 8.8.8.8 --output=yaml
```

### `time-conversion` Command
Convert between timezones using the `ipgeolocation.io` API to convert a specific time from one timezone to another using timezone identifiers and optional date/time inputs.

For more details, refer to official documentation: [Timezone Converter API](https://ipgeolocation.io/timezone-api.html#convert-time-between-time-zones).

#### `time-conversion` Usage
```bash
ipgeolocation time-conversion [flags]
```

#### Flags for `time-conversion`
| Flag              | Type    | Default  | Description                                      |
|-------------------|---------|----------|--------------------------------------------------|
| `--tz_from`       | string  | `""`     | Timezone to convert from.                        |
| `--tz_to`         | string  | `""`     | Timezone to convert to.                          |
| `--location_from` | string  | `""`     | Location to convert from.                        |
| `--location_to`   | string  | `""`     | Location to convert to.                          |
| `--lat_from`      | float64 | `0`      | Latitude to convert from.                        |
| `--long_from`     | float64 | `0`      | Longitude to convert from.                       |
| `--lat_to`        | float64 | `0`      | Latitude to convert to.                          |
| `--long_to`       | float64 | `0`      | Longitude to convert to.                         |
| `--iata_from`     | string  | `""`     | IATA code to convert from.                       |
| `--iata_to`       | string  | `""`     | IATA code to convert to.                         |
| `--icao_from`     | string  | `""`     | ICAO code to convert from.                       |
| `--icao_to`       | string  | `""`     | ICAO code to convert to.                         |
| `--lo_from`       | string  | `""`     | LO code to convert from.                         |
| `--lo_to`         | string  | `""`     | LO code to convert to.                           |
| `--time`          | string  | `""`     | Time to convert.                                 |
| `--output`        | string  | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Convert Current Time from One Timezone to Another
```bash
ipgeolocation time-conversion --tz_from America/New_York --tz_to Europe/London
```
Sample output:
```json
{
  "original_time": "2026-01-13 07:22:43",
  "converted_time": "2026-01-13 12:22:43",
  "diff_hour": 5,
  "diff_min": 300
}
```
> [!NOTE]
> You can convert time from any timezone to another timezone using location coordinates (Latitude and Longitude), location addresses, IATA codes, ICAO codes and UN/LUCODE.

Convert from one location to another:
```bash
ipgeolocation time-conversion --location_from New York --location_to London
```

Convert from one latitude and longitude to another:
```bash
ipgeolocation time-conversion --lat_from 37.7749 --long_from -122.4194 --lat_to 51.509865 --long_to -0.118092
```

Convert from one IATA code to another:
```bash
ipgeolocation time-conversion --iata_from DXB --iata_to JFK
```

Convert from one ICAO code to another:
```bash
ipgeolocation time-conversion --icao_from KATL --icao_to EWR
```

Convert from one LO code to another:
```bash
ipgeolocation time-conversion --lo_from DEBER --lo_to JFK
```

Output as YAML:
```bash
ipgeolocation time-conversion --tz_from America/New_York --tz_to Europe/London --output=yaml
```

### `astronomy` Command
The `astronomy` command uses the ipgeolocation.io Astronomy API to fetch astronomy-related data such as:

- Sunrise and sunset
- Solar noon
- Moonrise and moonset
- Moon phase
- Day length
- Timezone-based or coordinate-based location support

You can specify the location using IP, city name, coordinates, or timezone.

For further information, please visit [Astronomy API Documentation](https://ipgeolocation.io/astronomy-api.html#documentation-overview).

#### `astronomy` Usage
```bash
ipgeolocation astronomy [flags]
```

#### Flags for `astronomy`
| Flag          | Type    | Default  | Description                                      |
|---------------|---------|----------|--------------------------------------------------|
| `--ip`        | string  | `""`     | IPv4 or IPv6 address.                            |
| `--location`  | string  | `""`     | City name.                                       |
| `--latitude`  | float64 | `0`      | Latitude.                                        |
| `--longitude` | float64 | `0`      | Longitude.                                       |
| `--lang`      | string  | `""`     | Response language (if supported).                |
| `--tz`        | string  | `""`     | Timezone.                                        |
| `--elevation` | float64 | `0`      | Elevation.                                       |
| `--output`    | string  | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Lookup Astronomy API by Coordinates
Get astronomy info about a specific latitude and longitude:
```bash
ipgeolocation astronomy --latitude 37.7749 --longitude -122.4194
```
Sample output:
```json
{
  "astronomy": {
    "current_time": "05:07:18.340",
    "date": "2026-01-13",
    "day_length": "09:50",
    "evening": {
      "astronomical_twilight_begin": "18:14",
      "astronomical_twilight_end": "18:45",
      "blue_hour_begin": "17:30",
      "blue_hour_end": "17:52",
      "civil_twilight_begin": "17:13",
      "civil_twilight_end": "17:41",
      "golden_hour_begin": "16:33",
      "golden_hour_end": "17:30",
      "nautical_twilight_begin": "17:41",
      "nautical_twilight_end": "18:14"
    },
    "mid_night": "00:18",
    "moon_altitude": 16.07194223979877,
    "moon_angle": 301.4632703453228,
    "moon_azimuth": 138.83883667429063,
    "moon_distance": 405270.18931629026,
    "moon_illumination_percentage": "-23.90",
    "moon_parallactic_angle": -34.64165014870525,
    "moon_phase": "WANING_CRESCENT",
    "moon_status": "-",
    "moonrise": "03:17",
    "moonset": "12:53",
    "morning": {
      "astronomical_twilight_begin": "05:51",
      "astronomical_twilight_end": "06:23",
      "blue_hour_begin": "06:44",
      "blue_hour_end": "07:06",
      "civil_twilight_begin": "06:55",
      "civil_twilight_end": "07:23",
      "golden_hour_begin": "07:06",
      "golden_hour_end": "08:04",
      "nautical_twilight_begin": "06:23",
      "nautical_twilight_end": "06:55"
    },
    "night_begin": "18:45",
    "night_end": "05:51",
    "solar_noon": "12:18",
    "sun_altitude": -26.63562944381545,
    "sun_azimuth": 97.35878744233418,
    "sun_distance": 147120665.60394213,
    "sun_status": "-",
    "sunrise": "07:23",
    "sunset": "17:13"
  },
  "location": {
    "city": "San Francisco",
    "country_name": "United States",
    "elevation": "16",
    "latitude": "37.77490",
    "locality": "South of Market",
    "longitude": "-122.41940",
    "state_prov": "California"
  }
}
```

#### Lookup Astronomy API by IP Address
Get astronomy info about your current IP:
```bash
ipgeolocation astronomy
```
Get astronomy information for a specific IP address:
```bash
ipgeolocation astronomy --ip 8.8.8.8
```
Sample output:
```json
{
  "astronomy": {
    "current_time": "05:09:04.314",
    "date": "2026-01-13",
    "day_length": "09:50",
    "evening": {
      "astronomical_twilight_begin": "18:13",
      "astronomical_twilight_end": "18:44",
      "blue_hour_begin": "17:30",
      "blue_hour_end": "17:51",
      "civil_twilight_begin": "17:12",
      "civil_twilight_end": "17:41",
      "golden_hour_begin": "16:33",
      "golden_hour_end": "17:30",
      "nautical_twilight_begin": "17:41",
      "nautical_twilight_end": "18:13"
    },
    "mid_night": "00:16",
    "moon_altitude": 16.728478792843955,
    "moon_angle": 301.47655895571256,
    "moon_azimuth": 139.37283119065324,
    "moon_distance": 405270.54014008155,
    "moon_illumination_percentage": "-23.89",
    "moon_parallactic_angle": -34.405427656780056,
    "moon_phase": "WANING_CRESCENT",
    "moon_status": "-",
    "moonrise": "03:15",
    "moonset": "12:53",
    "morning": {
      "astronomical_twilight_begin": "05:50",
      "astronomical_twilight_end": "06:21",
      "blue_hour_begin": "06:42",
      "blue_hour_end": "07:04",
      "civil_twilight_begin": "06:53",
      "civil_twilight_end": "07:21",
      "golden_hour_begin": "07:04",
      "golden_hour_end": "08:01",
      "nautical_twilight_begin": "06:21",
      "nautical_twilight_end": "06:53"
    },
    "night_begin": "18:44",
    "night_end": "05:50",
    "solar_noon": "12:17",
    "sun_altitude": -25.978681889301882,
    "sun_azimuth": 97.96459578152059,
    "sun_distance": 147120665.60394216,
    "sun_status": "-",
    "sunrise": "07:21",
    "sunset": "17:12"
  },
  "ip": "8.8.8.8",
  "location": {
    "city": "Mountain View",
    "continent_code": "NA",
    "continent_name": "North America",
    "country_code2": "US",
    "country_code3": "USA",
    "country_name": "United States",
    "country_name_official": "United States of America",
    "district": "Santa Clara",
    "elevation": "3",
    "is_eu": false,
    "latitude": "37.42240",
    "locality": "Charleston Terrace",
    "longitude": "-122.08421",
    "state_code": "US-CA",
    "state_prov": "California",
    "zipcode": "94043-1351"
  }
}
```

#### Lookup Astronomy API by Location String
```bash
ipgeolocation astronomy --location "New York, USA"
```
Sample output:
```json
{
  "astronomy": {
    "current_time": "08:10:53.180",
    "date": "2026-01-13",
    "day_length": "09:23",
    "evening": {
      "astronomical_twilight_begin": "17:59",
      "astronomical_twilight_end": "18:34",
      "blue_hour_begin": "17:12",
      "blue_hour_end": "17:36",
      "civil_twilight_begin": "16:55",
      "civil_twilight_end": "17:24",
      "golden_hour_begin": "16:08",
      "golden_hour_end": "17:12",
      "nautical_twilight_begin": "17:24",
      "nautical_twilight_end": "17:59"
    },
    "mid_night": "00:13",
    "moon_altitude": 23.12401867639373,
    "moon_angle": 301.4902102354391,
    "moon_azimuth": 183.36677049153406,
    "moon_distance": 405270.8986986421,
    "moon_illumination_percentage": "-23.88",
    "moon_parallactic_angle": 2.6870094296532643,
    "moon_phase": "WANING_CRESCENT",
    "moon_status": "-",
    "moonrise": "03:20",
    "moonset": "12:27",
    "morning": {
      "astronomical_twilight_begin": "05:53",
      "astronomical_twilight_end": "06:27",
      "blue_hour_begin": "06:50",
      "blue_hour_end": "07:14",
      "civil_twilight_begin": "07:02",
      "civil_twilight_end": "07:32",
      "golden_hour_begin": "07:14",
      "golden_hour_end": "08:18",
      "nautical_twilight_begin": "06:27",
      "nautical_twilight_end": "07:02"
    },
    "night_begin": "18:34",
    "night_end": "05:53",
    "solar_noon": "12:13",
    "sun_altitude": 4.828766240317727,
    "sun_azimuth": 125.48815354581939,
    "sun_distance": 147120665.60394213,
    "sun_status": "-",
    "sunrise": "07:32",
    "sunset": "16:55"
  },
  "location": {
    "city": "Syracuse",
    "country_name": "United States",
    "elevation": "128",
    "latitude": "43.06923",
    "locality": "Destiny USA Mall",
    "location_string": "New York, USA",
    "longitude": "-76.17249",
    "state_prov": "New York"
  }
}
```

### `astronomy-timeseries` Command
The `astronomy-timeseries` command uses the ipgeolocation.io Astronomy API to fetch astronomy-related data such as:

- Sunrise and sunset
- Solar noon
- Moonrise and moonset
- Moon phase
- Day length
- Timezone-based or coordinate-based location support

For further information, please visit [Astronomy Timeseries API Documentation](https://ipgeolocation.io/astronomy-api.html#time-series-lookup).

#### `astronomy-timeseries` Usage
```bash
ipgeolocation astronomy-timeseries [flags]
```

#### Flags for `astronomy-timeseries`
| Flag           | Type    | Default  | Description                                      |
|----------------|---------|----------|--------------------------------------------------|
| `--ip`         | string  | `""`     | IPv4 or IPv6 address.                            |
| `--location`   | string  | `""`     | City name.                                       |
| `--latitude`   | float64 | `0`      | Latitude.                                        |
| `--longitude`  | float64 | `0`      | Longitude.                                       |
| `--lang`       | string  | `""`     | Response language (if supported).                |
| `--start-date` | string  | `""`     | Start date (e.g. 2023-01-01) Only YYYY-MM-DD.    |
| `--end-date`   | string  | `""`     | End date (e.g. 2023-12-31) Only YYYY-MM-DD       |
| `--output`     | string  | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`. |

> [!NOTE] 
> - The `start-date` and `end-date` flags are required.
> - The `start-date` and `end-date` should be in YYYY-MM-DD format.
> - The `start-date` and `end-date` should be less than 90 days.

#### Get astronomy timeseries info for an IP address
For your own IP address
```bash
ipgeolocation astronomy-timeseries --start-date 2023-01-01 --end-date 2023-01-02
```

#### Get astronomy timeseries about a specific lat and lon
```bash
ipgeolocation astronomy-timeseries --latitude 37.7749 --longitude -122.4194 --start-date 2023-01-01 --end-date 2023-01-02
```

#### Get astronomy timeseries about a specific location
```bash
ipgeolocation astronomy-timeseries --location "New York, USA" --start-date 2023-01-01 --end-date 2023-01-02
```
Sample output:
```json
{
  "astronomy": [
    {
      "date": "2023-01-01",
      "day_length": "09:08",
      "evening": {
        "astronomical_twilight_begin": "17:48",
        "astronomical_twilight_end": "18:23",
        "blue_hour_begin": "17:00",
        "blue_hour_end": "17:24",
        "civil_twilight_begin": "16:42",
        "civil_twilight_end": "17:12",
        "golden_hour_begin": "15:54",
        "golden_hour_end": "17:00",
        "nautical_twilight_begin": "17:12",
        "nautical_twilight_end": "17:48"
      },
      "mid_night": "00:08",
      "moon_phase": "WAXING_GIBBOUS",
      "moon_status": "-",
      "moonrise": "12:59",
      "moonset": "02:38",
      "morning": {
        "astronomical_twilight_begin": "05:53",
        "astronomical_twilight_end": "06:28",
        "blue_hour_begin": "06:51",
        "blue_hour_end": "07:16",
        "civil_twilight_begin": "07:04",
        "civil_twilight_end": "07:34",
        "golden_hour_begin": "07:16",
        "golden_hour_end": "08:22",
        "nautical_twilight_begin": "06:28",
        "nautical_twilight_end": "07:04"
      },
      "night_begin": "18:23",
      "night_end": "05:53",
      "solar_noon": "12:08",
      "sun_status": "-",
      "sunrise": "07:34",
      "sunset": "16:42"
    },
    {
      "date": "2023-01-02",
      "day_length": "09:09",
      "evening": {
        "astronomical_twilight_begin": "17:49",
        "astronomical_twilight_end": "18:23",
        "blue_hour_begin": "17:01",
        "blue_hour_end": "17:25",
        "civil_twilight_begin": "16:43",
        "civil_twilight_end": "17:13",
        "golden_hour_begin": "15:55",
        "golden_hour_end": "17:01",
        "nautical_twilight_begin": "17:13",
        "nautical_twilight_end": "17:49"
      },
      "mid_night": "00:08",
      "moon_phase": "WAXING_GIBBOUS",
      "moon_status": "-",
      "moonrise": "13:25",
      "moonset": "03:45",
      "morning": {
        "astronomical_twilight_begin": "05:53",
        "astronomical_twilight_end": "06:28",
        "blue_hour_begin": "06:52",
        "blue_hour_end": "07:16",
        "civil_twilight_begin": "07:04",
        "civil_twilight_end": "07:34",
        "golden_hour_begin": "07:16",
        "golden_hour_end": "08:22",
        "nautical_twilight_begin": "06:28",
        "nautical_twilight_end": "07:04"
      },
      "night_begin": "18:23",
      "night_end": "05:53",
      "solar_noon": "12:08",
      "sun_status": "-",
      "sunrise": "07:34",
      "sunset": "16:43"
    }
  ],
  "location": {
    "city": "Syracuse",
    "country_name": "United States",
    "elevation": "128",
    "latitude": "43.06923",
    "locality": "Destiny USA Mall",
    "location_string": "New York, USA",
    "longitude": "-76.17249",
    "state_prov": "New York"
  }
}
```

### `parse-user-agent` Command
The `parse-user-agent` command uses the ipgeolocation.io User Agent Parser API to parse user agent strings and extract relevant information such as:

- Device type
- Operating system
- Browser
- Browser version
- Device vendor
- Device model

#### `parse-user-agent` Usage
```bash
ipgeolocation parse-user-agent [flags]
```

#### Flags for `parse-user-agent`
| Flag           | Type   | Default  | Description                                      |
|----------------|--------|----------|--------------------------------------------------|
| `--user-agent` | string | `""`     | User agent string.                               |
| `--output`     | string | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`. |

For further information, please visit [User Agent Parser API Documentation](https://ipgeolocation.io/user-agent-api.html#documentation-overview).

#### Parse a user agent string
```bash
ipgeolocation parse-user-agent --user-agent "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
```
Sample output:
```json
{
  "device": {
    "brand": "Unknown",
    "cpu": "Intel x86_64",
    "name": "Desktop",
    "type": "Desktop"
  },
  "engine": {
    "name": "Blink",
    "type": "Browser",
    "version": "58.0",
    "version_major": "58"
  },
  "name": "Chrome",
  "operating_system": {
    "build": "??",
    "name": "Windows NT",
    "type": "Desktop",
    "version": "10.0",
    "version_major": "10"
  },
  "type": "Browser",
  "user_agent_string": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
  "version": "58.0.3029.110",
  "version_major": "58"
}
```

### `parse-bulk-user-agents` Command
The `parse-bulk-user-agents` command uses the ipgeolocation.io Bulk User-Agent Parser API to parse multiple user agent strings and extract relevant information such as:

- Device type
- Operating system
- Browser
- Browser version
- Device vendor
- Device model

#### `parse-bulk-user-agents` Usage
```bash
ipgeolocation parse-bulk-user-agents [flags]
```

#### Flags for `parse-bulk-user-agents`
| Flag            | Type     | Default  | Description                                      |
|-----------------|----------|----------|--------------------------------------------------|
| `--user-agents` | string[] | `[]`     | User agent strings.                              |
| `--output`      | string   | `pretty` | Output format: `pretty`, `raw`, `table`, `yaml`. |

For further information, please visit [Bulk User Agent Parser API Documentation](https://ipgeolocation.io/user-agent-api.html#parse-bulk-ua-strings).

> [!NOTE] 
> - The `user-agents` flag is required.
> - The `user-agents` flag should be an array of user agent strings.

#### Parse multiple user agent strings
```bash
ipgeolocation parse-bulk-user-agents --user-agents "Mozilla/5.0 (Windows NT 10.0; Win64; x64)","curl/7.64.1"
```
Sample output:
```json
[
  {
    "device": {
      "brand": "Unknown",
      "cpu": "Intel x86_64",
      "name": "Desktop",
      "type": "Desktop"
    },
    "engine": {
      "name": "Mozilla",
      "type": "Browser",
      "version": "5.0",
      "version_major": "5"
    },
    "name": "Windows NT",
    "operating_system": {
      "build": "??",
      "name": "Windows NT",
      "type": "Desktop",
      "version": "10.0",
      "version_major": "10"
    },
    "type": "Browser",
    "user_agent_string": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "version": "10.0",
    "version_major": "10"
  },
  {
    "device": {
      "brand": "Curl",
      "cpu": "Unknown",
      "name": "Curl",
      "type": "Robot"
    },
    "engine": {
      "name": "curl",
      "type": "Robot",
      "version": "7.64.1",
      "version_major": "7"
    },
    "name": "Curl",
    "operating_system": {
      "build": "??",
      "name": "Cloud",
      "type": "Cloud",
      "version": "??",
      "version_major": "??"
    },
    "type": "Robot",
    "user_agent_string": "curl/7.64.1",
    "version": "7.64.1",
    "version_major": "7"
  }
]
```

---

## License
This project is released under the [MIT License](https://github.com/IPGeolocation/cli/blob/main/LICENSE).
