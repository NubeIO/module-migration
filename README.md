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
    1. Hit command: `sudo ./module-migration generate-csv --external-token <external_token(need to escape $ by \$)> --ip <ip-or-ros-domain> --port <int> --schema <http|https>`
    2. Modify `./migration.csv`

2. Migrate
    1. Hit command: `sudo ./module-migration migrate --ssh-username <ssh-username> --ssh-password <ssh-password>`

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
4. Enable Modules in bulk from RCE
5. Upgrade RubixEdgeWires in bulk from RCE
6. Migrate older data into newer one
    - Make `ROS Migration State = false`
    - Make `Wires Migration State = false`
    - Make `Plugin Deletion State = true`
    - Run Migration
