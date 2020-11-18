# Get the current height of the blockchain and post json to stdout
#
#
# Example:
#      > bash get_height.sh
#      { "height":  54 }
#
#

FILE=parsing/height/height
if ! [ -f "$FILE" ]; then
    # shellcheck disable=SC2164
    cd parsing/height
    go build
    cd ../..
fi

./parsing/height/height
