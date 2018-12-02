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

The input files be considered as using `utf8` encoding. No matter what the true encoding it is.



## Install

```bash
git clone https://gitlab.com/visig/tf $(go env GOPATH)/src/gitlab.com/visig/tf
cd $(go env GOPATH)/src/gitlab.com/visig/tf
go get ./...
./maintain.sh install
```



## License

MIT
