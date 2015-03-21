# neosearch-import

# Build

```
go get -u github.com/NeowayLabs/neosearch
go build -v -tags leveldb
```

# usage

```
./neosearch-import
[General options]
     --file, -f: Read NeoSearch JSON database from file. (Required)
   --create, -c: Create new index database
     --name, -n: Name of index database
 --data-dir, -d: Data directory
     --help, -h: Display this help
```

Indexing the sample file:
```
mkdir /tmp/data
./neosearch-import -f samples/operating_systems.json  -c -d /tmp/data -n operating-systems
```


