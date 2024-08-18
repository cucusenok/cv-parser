1. `cp .env.example ./app/.env`
1. `docker compose build`
1. `docker compose up -d postgres`
1. `docker exec -i code-parser psql -U postgres -d parser < only_words.sql`
