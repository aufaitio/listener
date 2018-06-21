# listener

Micro service that manages registered repositories and their dependencies via service webhooks.

## Config

The listener config should be in Yaml format and named app.yaml.

```yaml
# Default config values set by application. Outlined to illustrate config structure.
db:
    host: localhost
    port: 27017
    name: aufait
    username: null
    password: null

errorFile: ./config/errors
port: 8080
```

## Development

### CLI

#### Usage

`./server [--configPath=<path>]`

#### Options

```
-h --help			 Show this message
--version			 Show version info
--configPath=<path>  Path to app.yaml config file [default: config]
```

#### Examples
```bash
# Pass --config if you need to override config options.
mongod --dbpath <dbPath>
./server
```
