export COMMENTO_ORIGIN=http://localhost:18080
export COMMENTO_PORT=18080
export COMMENTO_POSTGRES=postgres://postgres:postgres@localhost:5432/commento?sslmode=disable
export COMMENTO_CDN_PREFIX=$COMMENTO_ORIGIN
export COMMENTO_GITHUB_KEY=27641858a67a126b3a22
export COMMENTO_GITHUB_SECRET=2b2c4d464d8a76154644463229bbd3ef2bd3d0ed

mkdir -p build/devel
(cd api && make devel ) && cp api/build/devel/commento build/devel/commento
build/devel/commento

# create user test with LOGIN password 'test'
# ALTER ROLE postgres SET password_encryption = 'scram-sha-256';
# ALTER ROLE postgres WITH PASSWORD 'postgres';
