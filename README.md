# ipgeolocation CLI Documentation

The `ipgeolocation` CLI lets you interact with the [ipgeolocation.io](https://ipgeolocation.io) API directly from your terminal. It supports looking up IP information (single or bulk), configuring API keys, and formatting output in multiple styles.

---

## Table of Contents

- [Installation](#installation)
- [Global Flags](#global-flags)
- [Commands](#commands)
    - [1. `config`](#1-config)
    - [2. `ipgeo`](#2-ipgeo)
    - [3. `bulk-ipgeo`](#3-bulk-ipgeo)
    - [4. `ip-security`](#4-ip-security)
    - [5. `bulk-ip-security`](#5-bulk-ip-security)
    - [6. `asn`](#6-asn)
    - [7. `abuse`](#7-abuse)
    - [8. `timezone`](#8-timezone)
    - [9. `time-conversion`](#9-time-conversion)
    - [10. `astronomy`](#10-astronomy)
    - [11. `astronomy-timeseries`](#11-astronomy-timeseries)
    - [12. `parse-user-agent`](#12-parse-user-agent)
    - [13. `parse-bulk-user-agents`](#13-parse-bulk-user-agents)
- [License](#license)




## Installation

### 1. Go install

To install `ipgeolocation` using `go install`, run:
```bash
go install github.com/IPGeolocation/cli@latest
```

Make sure `$GOBIN` or `$GOPATH/bin` is in your `PATH`, then run:

```bash
ipgeolocation --help
```
---

### 2. Download and build from source
```bash
git clone https://github.com/IPGeolocation/cli.git
cd cli
go build -o ipgeolocation .
./ipgeolocation --help
```

## Global Flags
These flags are available for all commands:

| Flag         | Description |
|--------------|-------------|
| `-h, --help` | Show help for the command. |

You can also check the version for `ipgeolocation` using the `--version` flag:

```bash
ipgeolocation --version
```
---

## Commands

### 1. `config`
Configure your API key for `ipgeolocation.io`.

#### Usage
```bash
ipgeolocation config --apikey=<your-key>
```

#### Flags
| Flag        | Type   | Description |
|-------------|--------|-------------|
| `--apikey`  | string | Your API key from [ipgeolocation.io](https://ipgeolocation.io). |

---

### 2. `ipgeo`
Lookup geolocation information for a **single IP address or domain** using the `ipgeolocation.io` API.

#### Usage
```bash
ipgeolocation ipgeo [flags]
```

#### Flags
| Flag         | Type         | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--ip`       | string       | `""`      | IPv4, IPv6, or domain name (e.g. `8.8.8.8`, `google.com`). |
| `--include`  | string[]  | `[]`      | Include extra fields (e.g. `security,timezone,currency`). |
| `--excludes` | string[]  | `[]`      | Exclude fields from output. |
| `--fields`   | string[]  | `[]`      | Return only specific fields (e.g. `ip,organization`). |
| `--lang`     | string       | `""`      | Response language. |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Examples
Get info about your current IP:
```bash
ipgeolocation ipgeo
```

Lookup a specific IP:
```bash
ipgeolocation ipgeo --ip 8.8.8.8
```

Lookup a domain with extra fields:
```bash
ipgeolocation ipgeo --ip google.com --include security,timezone,currency
```

Get only specific fields in YAML:
```bash
ipgeolocation ipgeo --ip 1.1.1.1 --fields location --output yaml
```

Exclude unnecessary fields:
```bash
ipgeolocation ipgeo --ip 8.8.8.8 --excludes currency,time_zone
```

---

### 3. `bulk-ip-geo`
Lookup geolocation information for **multiple IPs** in one request.

#### Usage
```bash
ipgeolocation bulk-ip-geo [flags]
```

#### Flags
| Flag             | Type          | Default   | Description |
|------------------|--------------|-----------|-------------|
| `--ips`          | string[]  | `[]`      | Comma-separated list of IPs. Example: `--ips 8.8.8.8,1.1.1.1` |
| `--file`         | string       | `""`      | Path to a text file containing IPs (one per line). |
| `--include`      | string[]  | `[]`      | Include extra fields (e.g. `location,time_zone`). |
| `--exclude`      | string[]  | `[]`      | Exclude fields (e.g. `currency`). |
| `--fields`       | string[]  | `[]`      | Return only specific fields (e.g. `location`). |
| `--lang`         | string       | `""`      | Response language (if supported). |
| `--output`       | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |
| `--output-file`  | string       | `""`      | Save output to JSON file. Example: `--output-file results.json` |


For further information, please visit [IP Geolocation API Documentation](https://ipgeolocation.io/ip-location-api.html#documentation-overview).
#### Examples
Lookup 3 IP addresses:
```bash
ipgeolocation bulk-ip-geo --ips 8.8.8.8,1.1.1.1,192.30.253.112
```

Lookup from a file:
```bash
ipgeolocation bulk-ip-geo --file=ips.txt
```

Include location and timezone:
```bash
ipgeolocation bulk-ip-geo --file=ips.txt --include=location,time_zone
```

Output as YAML:
```bash
ipgeolocation bulk-ip-geo --ips=8.8.8.8 --output=yaml
```

Save results to JSON file:
```bash
ipgeolocation bulk-ip-geo --ips=8.8.8.8,1.1.1.1 --output-file=output.json
```

---

#### Output Formats

- **pretty** (default): Human-readable formatted JSON.  
- **raw**: Raw API response.  
- **table**: Tabular display of common fields.  
- **yaml**: YAML-formatted output.  
- **json file**: If `--output-file` is provided, results are saved to a `.json` file.  

---

### 4. `ip-security`
Lookup IP security information using the `ipgeolocation.io` API.

#### Usage
```bash
ipgeolocation ip-security [flags]
```

#### Flags
| Flag         | Type          | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--ip`       | string       | `""`      | IPv4 or IPv6 address. |
| `--include`  | string[]  | `[]`      | Include extra fields in output. |
| `--exclude`  | string[]  | `[]`      | Exclude fields from output. |
| `--fields`   | string[]  | `[]`      | Return only specific fields (e.g. `ip,organization`). |
| `--lang`     | string       | `""`      | Response language. |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Examples
Get info about your current IP:
```bash
ipgeolocation ip-security
```

Lookup a specific IP:
```bash
ipgeolocation ip-security --ip 8.8.8.8
```

Lookup security information for a specific IP with including additional fields:
```bash
ipgeolocation ip-security --ip 8.8.8.8 --include=location,time_zone
```

Output as YAML:
```bash
ipgeolocation ip-security --ip 8.8.8.8 --output=yaml
```

---
For further information, please visit [IP Security API Documentation](https://ipgeolocation.io/ip-security-api.html#documentation-overview).

---

### 5. `bulk-ip-security`

Lookup IP security information for **multiple IPs** in one request. 

#### Usage
```bash
ipgeolocation bulk-ip-security [flags]
```

#### Flags
| Flag             | Type          | Default   | Description |
|------------------|--------------|-----------|-------------|
| `--ips`          | string[]  | `[]`      | Comma-separated list of IPs. Example: `--ips 8.8.8.8,1.1.1.1` |
| `--file`         | string       | `""`      | Path to a text file containing IPs (one per line). |
| `--include`      | string[]  | `[]`      | Include extra fields (e.g. `location,time_zone`). |
| `--exclude`      | string[]  | `[]`      | Exclude fields (e.g. `currency`). |
| `--fields`       | string[]  | `[]`      | Return only specific fields (e.g. `location`). |
| `--lang`         | string       | `""`      | Response language (if supported). |
| `--output`       | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |
| `--output-file`  | string       | `""`      | Save output to JSON file. Example: `--output-file results.json` |

#### Examples
Lookup 3 IP addresses:
```bash
ipgeolocation bulk-ip-security --ips 8.8.8.8,1.1.1.1,192.30.253.112
```

Lookup from a file:
```bash
ipgeolocation bulk-ip-security --file=ips.txt
```

Include location and timezone:
```bash
ipgeolocation bulk-ip-security --file=ips.txt --include=location,time_zone
```

Output as YAML:
```bash
ipgeolocation bulk-ip-security --ips=8.8.8.8 --output=yaml
```

Save results to JSON file:
```bash
ipgeolocation bulk-ip-security --ips=8.8.8.8,1.1.1.1 --output-file=output.json
```

---
For further information, please visit [Bulk IP Security API Documentation](https://ipgeolocation.io/ip-security-api.html#documentation-overview).

---

### 6. `asn`
Lookup ASN (Autonomous System Number) information using the `ipgeolocation.io` API.

#### Usage
```bash
ipgeolocation asn [flags]
```

#### Flags
| Flag         | Type          | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--ip`       | string       | `""`      | IPv4 or IPv6 address. |
| `--include`  | string[]  | `[]`      | Include extra fields in output. |
| `--exclude`  | string[]  | `[]`      | Exclude fields from output. |
| `--fields`   | string[]  | `[]`      | Return only specific fields (e.g. `ip,organization`). |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Examples
Get ASN info about your current IP:
```bash
ipgeolocation asn
```

Lookup a specific ASN:
```bash
ipgeolocation asn --asn 12345
```

Lookup a specific IP with including additional fields:
```bash
ipgeolocation asn --ip 8.8.8.8 --include routes,peers,upstreams,downstreams,whois_response
```

Exclude unnecessary fields:
```bash
ipgeolocation asn --ip 8.8.8.8 --exclude asn.date_allocated,asn.allocation_status
```

Get only specific fields in table format:
```bash
ipgeolocation asn --ip 8.8.8.8 --fields asn.organization,asn.country,asn.downstreams --output table
```

Output as YAML:
```bash
ipgeolocation asn --ip 8.8.8.8 --output=yaml
```

---
For further information, please visit [ASN API Documentation](https://ipgeolocation.io/asn-api.html#documentation-overview).

---

### 7. `abuse`
Lookup abuse information using the `ipgeolocation.io` API.

#### Usage
```bash
ipgeolocation abuse [flags]
```

#### Flags
| Flag         | Type          | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--ip`       | string       | `""`      | IPv4 or IPv6 address. |
| `--exclude`  | string[]  | `[]`      | Exclude fields from output. |
| `--fields`   | string[]  | `[]`      | Return only specific fields (e.g. `ip,organization`). |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Examples
Get abuse info about your current IP:
```bash
ipgeolocation abuse
```

Lookup a specific IP:
```bash
ipgeolocation abuse --ip 8.8.8.8
```

Exclude unnecessary fields:
```bash
ipgeolocation abuse --ip 8.8.8.8 --exclude abuse.handle,abuse.emails
```

Get only specific fields in table format:
```bash
ipgeolocation abuse --ip 8.8.8.8 --fields abuse.role,abuse.emails --output table
```

Output as YAML:
```bash
ipgeolocation abuse --ip 8.8.8.8 --output=yaml
```

---
For further information, please visit [Abuse API Documentation](https://ipgeolocation.io/ip-abuse-contact-api.html#documentation-overview).

---

### 8. `timezone`
Lookup timezone information using the `ipgeolocation.io` API.

#### Usage
```bash
ipgeolocation timezone [flags]
```

#### Flags
| Flag         | Type          | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--ip`       | string       | `""`      | IPv4 or IPv6 address. |
| `--tz`  | string | `""`      | Timezone. |
| `--location`  | string  | `""`      | Location name (e.g. New York) |
| `--latitude`   | float64  | `0`      | Latitude (e.g. 37.7749). |
| `--longitude`  | float64  | `0`      | Longitude (e.g. -122.4194). |
| `--lang`     | string       | `""`      | Response language (if supported). |
| `--iata`     | string       | `""`      | IATA code (e.g. DXB). |
| `--icao`     | string       | `""`      | ICAO code (e.g. KATL). |
| `--lo`       | string       | `""`      | LO code (e.g. DEBER). |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Examples
Get timezone info about your current IP:
```bash
ipgeolocation timezone
```

Lookup a specific IP:
```bash
ipgeolocation timezone --ip 8.8.8.8
```

Lookup a specific timezone:
```bash
ipgeolocation timezone --tz America/New_York
```

Lookup a specific location:
```bash
ipgeolocation timezone --location New York
```

Lookup a specific latitude and longitude:
```bash
ipgeolocation timezone --latitude 37.7749 --longitude -122.4194
```

Lookup a specific IATA code:
```bash
ipgeolocation timezone --iata DXB
```

Lookup a specific ICAO code:
```bash
ipgeolocation timezone --icao KATL
```

Lookup a specific LO code:
```bash
ipgeolocation timezone --lo DEBER
```

Output as YAML:
```bash
ipgeolocation timezone --ip 8.8.8.8 --output=yaml
```
For further information, please visit [Timezone API Documentation](https://ipgeolocation.io/timezone-api.html#documentation-overview).

---


### 9. `time-conversion`
Convert between timezones using the `ipgeolocation.io` API.

#### Usage
```bash
ipgeolocation time-conversion [flags]
```

#### Flags
| Flag         | Type          | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--tz_from`     | string       | `""`      | Timezone to convert from. |
| `--tz_to`       | string       | `""`      | Timezone to convert to. |
| `--location_from`     | string       | `""`      | Location to convert from. |
| `--location_to`       | string       | `""`      | Location to convert to. |
| `--lat_from`     | float64       | `0`      | Latitude to convert from. |
| `--long_from`    | float64       | `0`      | Longitude to convert from. |
| `--lat_to`       | float64       | `0`      | Latitude to convert to. |
| `--long_to`      | float64       | `0`      | Longitude to convert to. |
| `--iata_from`    | string       | `""`      | IATA code to convert from. |
| `--iata_to`      | string       | `""`      | IATA code to convert to. |
| `--icao_from`    | string       | `""`      | ICAO code to convert from. |
| `--icao_to`      | string       | `""`      | ICAO code to convert to. |
| `--lo_from`      | string       | `""`      | LO code to convert from. |
| `--lo_to`        | string       | `""`      | LO code to convert to. |
| `--time`     | string       | `""`      | Time to convert. |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Examples
Convert from one timezone to another:
```bash
ipgeolocation time-conversion --tz_from America/New_York --tz_to Europe/London
```

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
For further information, please visit [Time Conversion API Documentation](https://ipgeolocation.io/timezone-api.html#convert-time-bw-time-zones).

---

### 10. `astronomy`
The `astronomy` command uses the ipgeolocation.io Astronomy API
to fetch astronomy-related data such as:

- Sunrise and sunset
- Solar noon
- Moonrise and moonset
- Moon phase
- Day length
- Timezone-based or coordinate-based location support

You can specify the location using IP, city name, coordinates, or timezone.

API Reference: https://ipgeolocation.io/astronomy-api.html


#### Usage
```bash
ipgeolocation astronomy [flags]
```

#### Flags
| Flag         | Type          | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--ip`       | string       | `""`      | IPv4 or IPv6 address. |
| `--location`     | string       | `""`      | City name. |
| `--latitude` | float64      | `0`       | Latitude. |
| `--longitude`| float64      | `0`       | Longitude. |
| `--lang`     | string       | `""`      | Response language (if supported). |
| `--tz` | string       | `""`      | Timezone. |
| `--elevation` | float64      | `0`       | Elevation. |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Examples
Get astronomy info about your current IP:
```bash
ipgeolocation astronomy
```

Get astronomy info about a specific location:
```bash
ipgeolocation astronomy --location "New York"
```

Get astronomy info about a specific timezone:
```bash
ipgeolocation astronomy --tz "America/New_York"
```

Get astronomy info about a specific latitude and longitude:
```bash
ipgeolocation astronomy --latitude 37.7749 --longitude -122.4194
```

Output as raw:
```bash
ipgeolocation astronomy --ip 8.8.8.8 --output=raw
```
For further information, please visit [Astronomy API Documentation](https://ipgeolocation.io/astronomy-api.html#documentation-overview).

---


### 11. `astronomy-timeseries`
The `astronomy-timeseries` command uses the ipgeolocation.io Astronomy API
to fetch astronomy-related data such as:

- Sunrise and sunset
- Solar noon
- Moonrise and moonset
- Moon phase
- Day length
- Timezone-based or coordinate-based location support

#### Usage
```bash
ipgeolocation astronomy-timeseries [flags]
```

#### Flags
| Flag         | Type          | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--ip`       | string       | `""`      | IPv4 or IPv6 address. |
| `--location`     | string       | `""`      | City name. |
| `--latitude` | float64      | `0`       | Latitude. |
| `--longitude`| float64      | `0`       | Longitude. |
| `--lang`     | string       | `""`      | Response language (if supported). |
| `--start-date` | string       | `""`      | Start date (e.g. 2023-01-01) Only YYYY-MM-DD. |
| `--end-date` | string       | `""`      | End date (e.g. 2023-12-31) Only YYYY-MM-DD |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Note: 
- The `start-date` and `end-date` flags are required.
- The `start-date` and `end-date` should be in YYYY-MM-DD format.
- The `start-date` and `end-date` should be less than 90 days.



#### Examples
Get astronomy info about your current IP:
```bash
ipgeolocation astronomy-timeseries --start-date 2023-01-01 --end-date 2023-01-02
```

Get astronomy info about a specific location:
```bash
ipgeolocation astronomy-timeseries --location "New York" --start-date 2023-01-01 --end-date 2023-01-02
```

Get astronomy info about a specific latitude and longitude:
```bash
ipgeolocation astronomy-timeseries --latitude 37.7749 --longitude -122.4194 --start-date 2023-01-01 --end-date 2023-01-02
```

Output as raw:
```bash
ipgeolocation astronomy-timeseries --ip 8.8.8.8 --start-date 2023-01-01 --end-date 2023-01-02 --output=raw
```
For further information, please visit [Astronomy Timeseries API Documentation](https://ipgeolocation.io/astronomy-api.html#time-series-lookup).

---

### 12. `parse-user-agent`
The `parse-user-agent` command uses the ipgeolocation.io User Agent Parser API
to parse user agent strings and extract relevant information such as:

- Device type
- Operating system
- Browser
- Browser version
- Device vendor
- Device model

#### Usage
```bash
ipgeolocation parse-user-agent [flags]
```

#### Flags
| Flag         | Type          | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--user-agent`       | string       | `""`      | User agent string. |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Examples
Parse a user agent string:
```bash
ipgeolocation parse-user-agent --user-agent "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
```

Output as raw:
```bash
ipgeolocation parse-user-agent --user-agent "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36" --output=raw
```
For further information, please visit [User Agent Parser API Documentation](https://ipgeolocation.io/user-agent-api.html#documentation-overview).

---

### 13. `parse-bulk-user-agents`
The `parse-bulk-user-agents` command uses the ipgeolocation.io Bulk User-Agent Parser API
to parse multiple user agent strings and extract relevant information such as:

- Device type
- Operating system
- Browser
- Browser version
- Device vendor
- Device model

#### Usage
```bash
ipgeolocation parse-bulk-user-agents [flags]
```

#### Flags
| Flag         | Type          | Default   | Description |
|--------------|--------------|-----------|-------------|
| `--user-agents`       | string[]       | `[]`      | User agent strings. |
| `--output`   | string       | `pretty`  | Output format: `pretty`, `raw`, `table`, `yaml`. |

#### Note: 
- The `user-agents` flag is required.
- The `user-agents` flag should be an array of user agent strings.

#### Examples
Parse multiple user agent strings:
```bash
ipgeolocation parse-bulk-user-agents --user-agents "Mozilla/5.0 (Windows NT 10.0; Win64; x64)","curl/7.64.1"
```

Output as raw:
```bash
ipgeolocation parse-bulk-user-agents --user-agents "Mozilla/5.0 (Windows NT 10.0; Win64; x64)","curl/7.64.1" --output=raw
```
For further information, please visit [Bulk User Agent Parser API Documentation](https://ipgeolocation.io/user-agent-api.html#parse-bulk-ua-strings).

---


## License
This project is released under the [MIT License](LICENSE).
