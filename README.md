# tf

A command line tool for calculate term-frequency of files.

It was designed for get some decent search result in filesystem directly without prepared external index nor particular tokenizing algorithm.

A higher score meaning a higher relevance.



## Usage

```bash
tf -t <term> <filename>...
tf -t <term1> -t <term2> <filename>...
echo <filename> | tf -t <term>
find . -type f -iname '*.txt' | tf -t <term> | sort
```

Output example:

```bash
0.140947213 file1.txt
0.010238174 file2.txt
```

The input files be considered as using `utf8` encoding. No matter what the true encoding it is.



## Download

### Pre-build Files

<https://gitlab.com/visig/tf/tags>



### Build from Source

Prepare a golang environment, then:

```bash
git clone https://gitlab.com/visig/tf $(go env GOPATH)/src/gitlab.com/visig/tf
cd $(go env GOPATH)/src/gitlab.com/visig/tf
go get ./...
./maintain.sh install
```


## License

MIT
