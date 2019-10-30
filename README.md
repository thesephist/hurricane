# Hurricane

Zero-configuration, read-only JSON API in front of an Airtable base.

## Usage

You can compile Hurricane into a single executable binary with

```sh
go build cmd/hurricane.go -o ./hurricane
```

Run Hurricane with `./hurricane`. Hurricane takes no command line arguments, but reads a couple of configuration strings from the environment:

```
HURRICANE_BASE_ID   Airtable Base ID
HURRICANE_API_KEY   Airtable API KEY

(optional)
HURRICANE_PORT      The network port on which Hurricane should listen
```

## API

For a given base and each table, Hurricane exposes a simple API:

### `/` (root path)

List all records in the table / default view as a JSON array.

### `/recXXXXXX`

Retrieve a specific record as a JSON object.

### `/view/XXXXXX`

Retrieve all records in a view as a JSON array.

## Caching

Airtable's API is rate-limited per API key, so Hurricane by default will cache queries to the Airtable API for 60 seconds. You can tweak the `Cache` instance in the main file to change this.
