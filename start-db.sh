docker run -it --rm -v `pwd`/data:/var/lib/postgresql/data -p 5432:5432 \
  -e POSTGRES_DB=commento -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres \
  postgres
