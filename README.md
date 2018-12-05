# tf

A command line tool for calculate the relevance between TERMs and FILEs.

It was designed for generate some decent searching result from filesystem without external index nor particular tokenizing algorithm.

The program's name 'tf' mean the original algorithm Term-Frequency. But the program not fully respect the TF. It added some tweaks to enhance the relevance accuracy.



## Usage

```bash
tf <term> -f <filename>
tf <term1> <term2> <term3> -f <filename1> -f <filename2>
echo <filename> | tf <term>
find . -iname '*.txt' | tf <term1> <term2> | sort -n
```

Output example:

```bash
  0.14094721 file1.txt
 13.01023817 file2.txt
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
