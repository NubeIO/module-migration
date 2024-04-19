# README

## How to build

```bash
go build -o module-migration
```

## How to migrate

### Pre-requisite

1. Download `module-migration` binary file
2. wget https://github.com/NubeIO/module-migration/releases/download/v0.0.2/module-migration-0.0.3-3b1c3956.armv7.zip`
3. unzip module-migration-0.0.3-3b1c3956.armv7.zip`



### How to

#### Steps

1. Generate host list.
   1. Hit command: `sudo ./module-migration generate-csv`. Generates host list in `/data/migration/migration.csv`
   2. Modify `/data/migration/migration.csv`.

2. Migrate 
   1. Hit command: `sudo ./module-migration migrate`

