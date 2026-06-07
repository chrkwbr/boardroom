env "local" {
  src = "file://db/schema.sql"
  dev = "docker://postgres/17/dev?search_path=public"
  url = "postgres://boardroom:boardroom@localhost:5433/boardroom?sslmode=disable"

  migration {
    dir = "file://migrations"
  }
}

