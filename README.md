# mina-archive-proxy

HTTP API service to provide access to Mina archive PostgreSQL database.

## Usage

To see a full list of available options, run:

```bash
mina-archive-proxy --help
```

Output:

```bash
Usage of ./mina-archive-proxy:
  -cors-enabled
    	Enable CORS on the server
  -db string
    	Database connection string
  -debug
    	Enable debug mode
  -ledger-enabled
    	Enable staking ledger dump endpoint (default true)
  -mina-bin string
    	Full path to Coda binary (default "coda")
  -version
    	Show version
```

## API Reference

| Method | Path             | Description
|--------|------------------|---------------------------------------------------
| GET    | /                | Blockchain stats
| GET    | /status          | Archive db status
| GET    | /chain           | Get current canonical chain
| GET    | /blocks          | Get all blocks
| GET    | /blocks/:hash    | Get block details by hash
| GET    | /block_producers | Get all block producer keys and number of blocks
| GET    | /public_keys     | Get all public keys
| GET    | /public_keys/:id | Get public keys details by value
| GET    | /staking_ledger  | Get staking ledger dump

### Swagger

To see Swagger documentation locally execute the following command:

```bash
make swagger
```

## License

Apache License v2.0