#!/bin/bash

date
find . -type f -name "*.go" > cscope.files
find . -type f -name "*.proto" >> cscope.files

echo "excludes *.pb.go ..."
mv cscope.files cscope.files.bak
sed -e '/\.pb\.go/d;/vendor\//d' cscope.files.bak > cscope.files
rm -f cscope.files.bak



if [ $# -eq 1 ]; then
    # remove test related files
    echo "excludes some test files..."
    mv cscope.files cscope.files.bak
    sed -e '/mock/d;/_test/d;/UnitTest/d;/\/test\//d' cscope.files.bak > cscope.files
    #sed -e '/flashlibs/d;/memcached/d;/\/test\//d' cscope.files.bak > cscope.files
    rm -f cscope.files.bak
fi
cscope -b

#find  . -name \*.py -print | xargs etags
#rsync -avr --exclude=.git/* --exclude=cscope.* --exclude=TAGS  mc-io  mc-control  mc99:/home/shzyao/

date

