go build -o lemin .
mv -f lemin ./lemin_test
cd lemin_test
source lemin_key.sh
shopt -s expand_aliases
audit
lemin