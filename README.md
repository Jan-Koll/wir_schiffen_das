# Installation

```bash
cp .env.example .env
docker compose up
```

After installation, create the db_schema or import existing backup

**Import database schema:**
```bash
cat db/dump_empty.sql | docker exec -i wir_schiffen_das-db-1 psql -U db_admin -d construction
```
**Create a new DB Backup:**
```bash
docker exec -t wir_schiffen_das-db-1 pg_dumpall -c -U db_admin > dump_`date +%d-%m-%Y"_"%H_%M_%S`.sql
```
