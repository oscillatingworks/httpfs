httpfs
======

httpfs is a HTTP server that translate requests into basic UNIX files and directories
operations, using native Go system calls.

HTTP->Unix Mappings
-------------------

| HTTP Request                   | UNIX operation                   |
| ------------------------------ | -------------------------------- |
|  GET /$FILE_PATH               | `cat $HOME/$FILE_PATH`           |
|  GET /$DIR_PATH                | `ls $HOME/DIR_PATH`              |
|  POST /$SOME_PATH type=file    | `touch $HOME/SOME_PATH`          |
|  POST /$SOME_PATH type=dir     | `mkdir -p $HOME/SOME_PATH`       |
|  PUT /$FILE_PATH content=text  | `echo "text" >> $HOME/FILE_PATH` |
|  DELETE /$FILE_PATH            | `rm $HOME/$FILE_PATH`            |
|  DELETE /$DIR_PATH             | `rm -rf $HOME/DIR_PATH`          |

Install
-------

```
go get -u github.com/oscillatingworks/httpfs
httpfs
```

LICENSE
-------

MIT
