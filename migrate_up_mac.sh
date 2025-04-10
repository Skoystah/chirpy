#pushd 
cd $PWD/sql/schema/
goose postgres "postgres://geert:@localhost:5432/chirpy" up
#popd
