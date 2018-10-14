# tagdb
Filesystem tags implementation.

tagdb allows you to "attach" keywords to arbitrary files and then ask for files
with certain keyword attached. Pretty simple, huh?

Why you would need keywords on files? First let's answer another question: Why
Twitter does have #hashtags thing? To help you find similar tweets. tagdb works
the same way and helps you find similar or somehow related files.

To help you understand why tagdb can be useful let's discuss few usage examples.
First one is IT-related e-books collection (I actually have over 3 GiB of
them). There are books about security,
programming, algorithms, system management, cryptography, machine learning,
networking, etc... You got the idea. How would you organize this collection
into hierarchy of directories? Consider that many books cover more than one topic.
My answer: You don't have to. Place all books into one directory and use tagdb. 
[Grokking Algorithms] book (btw pretty good book) is about algorithms? Great,
label it with "algorithms" keyword:
```
tagctl tag Grokking_Algorithms.png -t algorithms
```
Now you can easily get list of all algorithms-related books with following command:
```
tagctl query algorithms
```

Another great example is a photo gallery. Nobody likes `IMG_XXXXX.png` file
names assigned by cameras, right?
Sort photos by location/year/etc using tagdb!
```
tagctl tag IMG_00001.png IMG_00002.png ... -t location
tagctl query year_2017
tagctl query location
```

See? tagdb can be very useful! Now go install it!

### Installation

```
go get github.com/foxcpp/tagdb/tagctl
```
And... that's all. Now you have CLI utility `tagctl` installed in
`$GOPATH/bin`.

Binaries for Linux/Windows/etc can be found in
[Releases](https://github.com/foxcpp/tagdb/releases) section.

### Usage

See `tagctl help` for CLI usage hints.

Few examples:
```
$ tagctl tag Documents/agreement.odt.asc -t work
$ tagctl untag Pictures/avatar.xcf -t todo
$ tagctl tags
work
todo
$ tagctl query work
/home/user/Documents/agreement.odt.asc
```

### License

MIT.

### Similar Projects

* DBFS by SourceForge
* [TMSU](https://github.com/oniony/TMSU) by oniony

[Grokking Algorithms]: https://www.amazon.com/Grokking-Algorithms-illustrated-programmers-curious/dp/1617292230
