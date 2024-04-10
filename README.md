# BIOS CLI

## Over REST

```
go run main.go server
```

## Make a systemctl file

### Over REST

```json
{
  "file": "ctl.yaml",
  "args": {
    "name": "driver-bacnet",
    "desc": "a driver-bacnet service"
  }
}
```

### Over CLI

```
go run main.go build ctl.yaml name=driver-bacnet desc="My new service"
```

## Download a GitHub build

### Over REST

```json
{
    "file": "git.yaml",
    "args": {
        "owner": "NubeIO",
        "repo": "driver-bacnet",
        "tag": "v1.0.0-rc.1",
        "arch": "arvm7",
        "location": "./",
        "token": ""
    }
}
```

### Over CLI

```
go run main.go build git.yaml owner=NubeIO repo=driver-bacnet tag=v1.0.0-rc.1 arch=armv7 location=./ token=<TOKEN>
```
