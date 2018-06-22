
set -e
# generate a report of test coverage. Produces a nice html detail of this.
# see: progress/cover.html for output.

LOG_FILE="progress/test_log"
COVER_FILE="progress/cover.out"
COVER_HTML="progress/cover.html"
CURRENT_DIR=$(pwd  | sed 's/\//\\\//g')

echo >> $LOG_FILE
echo `date '+%Y-%m-%d %H:%M:%S'` >> $LOG_FILE

# log all test results in dir and generate a coverage profile.
go test ./... -coverprofile $COVER_FILE >> $LOG_FILE

# replace absolute paths so cover tool builds html. dont understand why
# it dosent work without this. may be my GO_PATH but it creates the
# correct absoulute path for the files then when the cover tool tries to
# produce the output it tries to write the abspath + GO_ROOT and GO_PATH.
sed -i '' "s/_$CURRENT_DIR/dito/g" $COVER_FILE
sed -i '' "s/_$CURRENT_DIR/dito/g" $LOG_FILE

# generate interactive coverage page.
go tool cover -html=$COVER_FILE -o $COVER_HTML

