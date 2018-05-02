# [WIP] tagdb

Filesystem tags implementation.

Sometimes it's hard to create good hiearchy for your files. Sadly, all modern
filesystems require this. ...

### Installation

```
go get github.com/foxcpp/tagdb/tagctl
```
And... that's all. 

### Usage

```
$ tagctl addtag Documents/agreement.odt.asc -t work
$ tagctl remtag Pictures/avatar.xcf -t todo
$ tagctl tags Documents/agreement.odt.asc
work important
$ tagctl query 'work & important'
/home/user/Documents/agreement.odt.asc
```

### License

MIT.
