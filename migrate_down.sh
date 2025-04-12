#pushd 
cd $PWD/sql/schema/
if [[ "$(uname -s)" == "Darwin" ]]; then
    goose postgres "postgres://geert:@localhost:5432/chirpy" down
else
    goose postgres "postgres://postgres:spitfire@localhost:5432/chirpy" down
fi
#popd
