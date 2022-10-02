# backup
Tool for creation backups.

## usage
* create `.backup.yaml` (by default config must store at /etc/backup/config.yaml) with next content:
```yaml
#drivers section, describes possible storages of backup`s archives
drivers:
    s3: # s3 like storages configurations
        default: # name of configuration 
          profile: <aws_shared_profile_name>
          bucket: <bucket_name>
          url: <endpoint_url>
          path: <path_on_bucket>
    fs: # just store archives in a directory
        default:
            output_path: new/dir
```
* build backup
```shell
go build -o backup cmd/backup/main.go
```
* run with args: dir, driver, archive type
```shell
# path started from ~ or . will be resolved to absolute
# s3 - driver name, in example this is eq s3.default
# you can describe many driver configurations 
# and use it in args in following format: <driver_type>.<driver_name>
# archive type is optional, one of zip,tar,tgz, first is default
./backup ~/my/dir s3 tgz --config=my_cfg.yaml
```

## schedule
For scheduling used [cron](https://github.com/robfig/cron) library.
You must add to config root jobs section in following format:
```yaml
jobs:
    - name: "Documents" # job to upload a directory
      target: "~/Documents" #target path
      driver: s3.default #driver name
      schedule: "@every 1m" #see format of intervals in lib above
      archive: tgz #archive type
      timeout: 5m #timeout
```
After it, run one of commands:
```bash
./backup jobs #prints out jobs section
./backup schedule start # starts a schedule process
```