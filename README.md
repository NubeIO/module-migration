# README

### How to build

```bash
go build -o module-migration
```

### How to migrate

#### Pre-requisite

1. Download `module-migration` binary file

#### How to migrate ROS

1. Upgrade ROS into latest
2. Upgrade necessary all modules
3. SSH in into the device
4. Hit command: `./module-migration migrate-ros`
5. Restart ROS

### How to migrate rubix-edge-wires

1. Hit command: `./module-migration migrate-wires`
2. Upgrade rubix-edge-wires into latest
3. Copy `/data/backup/migration/rubix-edge-wires/backup.json` file's content
4. Paste the content in the Wires Sheet
