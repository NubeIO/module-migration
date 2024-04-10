# README

## How to build

```bash
go build -o module-migration
```

## How to migrate

### Pre-requisite

1. Download `module-migration` binary file
2. wget https://github.com/NubeIO/module-migration/releases/download/v0.0.1/module-migration-0.0.1-79bf27f4.armv7.zip`
3. unzip module-migration-0.0.1-79bf27f4.armv7.zip`
4. Hit below commands from superuser (`sudo su`) or add `sudo` command at prefix

### How to migrate ROS

#### Steps

1. `rm -r /data/rubix-os/data/plugins/*`
2. Upgrade ROS into latest
3. Upgrade necessary all modules
4. SSH in into the device
5. Hit command: `./module-migration migrate-ros`
6. Restart ROS (`systemctl restart nubeio-rubix-os.service` or restart from RCE)

### How to migrate rubix-edge-wires

#### Prerequisite

1. Wires needs to be running in the system

#### Steps

1. Hit command: `./module-migration migrate-wires`
2. Upgrade rubix-edge-wires into latest
3. Copy `/data/backup/migration/rubix-edge-wires/backup.json` file's content
4. Paste the content in the Wires Sheet
