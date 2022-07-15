# backup
Tool for creation backups.

## usage
* create `.backup.yaml` at user`s home dir with next content
```yaml
#drivers section, describes possible storages of backup`s archives
drivers:
    s3: # s3 like storages configurations
        default: # name of configuration 
            access_key: <access_key_id>
            secret_key: <secret_key>
            bucket: <bucket_name>
            url: <endpoint_url>
            region: <region>
            path: <path_on_bucket>
    dir: # just store archives in a directory
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