# mina-archive-proxy

Archive data API for Mina

## Usage

```bash
Usage of ./mina-archive-proxy:
  -coda-bin string
    	Full path to Coda binary
  -ledger-enabled
    	Enable staking ledger dump endpoint (default true)
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