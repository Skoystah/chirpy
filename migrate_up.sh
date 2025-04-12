#pushd 
cd $PWD/sql/schema/
if [[ "$(uname -s)" == "Darwin" ]]; then
    goose postgres "postgres://geert:@localhost:5432/chirpy" up
else
    goose postgres "postgres://postgres:spitfire@localhost:5432/chirpy" up
fi
#popd
