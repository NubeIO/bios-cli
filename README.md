# ros-bios


over rest

```
go run main.go server
```


## make a systemctl file

```
go run main.go build ctl.yaml name=driver-bacnet desc="My new service"
```

```json
{
  "file": "ctl.yaml",
  "args": {
    "name": "driver-bacnet",
    "desc": "a driver-bacnet service"
  }
}
```

## download a github build

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

over cli

```
go run main.go build git.yaml owner=NubeIO repo=driver-bacnet tag=v1.0.0-rc.1 arch=armv7 location="./" token<TOKEN>
```