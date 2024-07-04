# README

It will perform a bulk migration according to the host records and status listed in the CSV.

## How to build

```bash
go build -o module-migration
GITHUB_TOKEN=<GITHUB_TOKEN> bash build.bash
```

## How to build (for armv7)

```bash
GITHUB_TOKEN=<GITHUB_TOKEN> bash build.bash
```

## How to migrate

### Pre-requisite

1. We need to have `sqlite3` in local
2. Download `module-migration` binary file
3. `wget https://github.com/NubeIO/module-migration/releases/download/v1.0.0/module-migration-1.0.0-728cca4e.amd64.zip`
4. `unzip module-migration-1.0.0-728cca4e.amd64.zip`

### How to

#### Steps

1. Generate host list.
    1. Hit command: `./module-migration generate-csv --external-token <external_token(need to escape $ by \$)> --ip <ip-or-ros-domain> --port <int> --schema <http|https>`
    2. Modify `./migration.csv`

2. Migrate
    1. Hit command: `./module-migration migrate --ssh-username <ssh-username> --ssh-password <ssh-password> --ssh-port <22 by default>`

### Workflow in migration

1. Delete all plugins from selected devices (so on ROS upgrade it doesn't have to upgrade plugins -- to make it faster)
    - Generate CSV from module-migration
    - Filter records on that CSV file for the migration
    - Make `ROS Migration State = true`
    - Make `Wires Migration State = true`
    - Make `Plugin Deletion State = false`
    - Run Migration
2. Upgrade ROS in bulk from RCE
3. Install Modules in bulk from RCE
4. Upgrade RubixEdgeWires in bulk from RCE
5. Migrate older data into newer one
    - Make `ROS Migration State = false`
    - Make `Wires Migration State = false`
    - Make `Plugin Deletion State = true`
    - Run Migration
6. Enable Modules in bulk from RCE
