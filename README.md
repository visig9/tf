# tf

A command line tool for calculate term-frequency of files.

It was designed for get some decent search result in filesystem directly without prepared external index nor particular tokenizing algorithm.

A higher score meaning a higher relevance.



## Usage

    tf -t <term> <filename>...
    tf -t <term1> -t <term2> <filename>...
    echo <filename> | tf -t <term>
    find . -type f -iname '*.txt' | tf -t <term> | sort

The input files be considered as using utf8 encoding. No matter what the real encoding be used.



## Install

    go get gitlab.com/visig/tf



## License

MIT
