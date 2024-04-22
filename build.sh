
VERSION=$(grep -o 'Version = "[^"]*' main.go | awk -F'"' '{print $2}') 
PROJECT_PACK_DIR=builds/$VERSION
if [ ! -e $PROJECT_PACK_DIR ]; then
    mkdir -p $PROJECT_PACK_DIR
fi 

# Linux
tar -cjf $PROJECT_PACK_DIR/linux-$VERSION.tar.bz2 --exclude=.git --exclude=builds .
tar -cJf $PROJECT_PACK_DIR/linux-$VERSION.tar.xz --exclude=.git --exclude=builds .
tar -czf $PROJECT_PACK_DIR/linux-$VERSION.tar.gz --exclude=.git --exclude=builds .

# Windows
GOOS=windows
GOARCH=386 go build -o MCSCS.exe
zip -r $PROJECT_PACK_DIR/windows-$VERSION-32bits.zip MCSCS.exe
GOARCH=amd64 go build -o MCSCS.exe
zip -r $PROJECT_PACK_DIR/windows-$VERSION-64bits.zip MCSCS.exe