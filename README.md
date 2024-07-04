# README

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

1. Download `module-migration` binary file
2. `wget https://github.com/NubeIO/module-migration/releases/download/v0.0.3/module-migration-0.0.3-3b1c3956.armv7.zip`
3. `unzip module-migration-0.0.3-3b1c3956.armv7.zip`

### How to migrate ROS

#### Steps

1. `rm -r /data/rubix-os/data/plugins/*`
2. Upgrade ROS into latest
3. Upgrade necessary all modules
4. SSH in into the device
5. Hit command: `sudo ./module-migration migrate-ros`
6. Restart ROS (`systemctl restart nubeio-rubix-os.service` or restart from RCE)

### How to migrate rubix-edge-wires

#### Prerequisite

1. Wires needs to be running in the system

#### Steps

1. Hit command: `sudo ./module-migration migrate-wires`
2. Upgrade rubix-edge-wires into latest
3. Copy `/data/backup/migration/rubix-edge-wires/backup.json` file's content
4. Paste the content in the Wires Sheet
