#pushd 
cd $PWD/sql/schema/
goose postgres "postgres://postgres:spitfire@localhost:5432/chirpy" down
#popd
