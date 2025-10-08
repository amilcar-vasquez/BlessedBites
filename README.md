# BlessedBites

**BlessedBites** is a web-based application designed to streamline the operations of a small to mid-sized food restaurant. It provides functionality to manage menu items, take orders, and more.

---

## 📌 Project Overview

This project includes:
- A backend database to store and manage data related to menu, orders, and users.
- A server-side rendered frontend using templates.
- Admin pages to manage menu items and categories.

---

## 📂 Documentation

- Database schema diagram: `documentation/BlessedBites.png`
- SQL dump: `documentation/blessed_bites_dump.sql`

---

## 🔐 Database Credentials (for testing/development)

> ⚠️ **Important:** These are for **testing purposes** only.

```
Database: blessed_bites
User: blessed_bites
Password: Matthew3.5:6
Host: localhost
Port: 4000
```

---

## 🛠️ Local development (live reload)

Use the `make dev` helper to run the app locally with live reload via `air`.

1. Ensure you have an `.env` file in the project root with your `JOURNAL_DB_DSN` set to your Supabase/Postgres connection string (or set `DB_DSN` directly in your environment). Example `.env`:

```env
# Example (Supabase or remote Postgres)
JOURNAL_DB_DSN=postgresql://<user>:<password>@<host>:<port>/<db>
```

2. Run live-reload dev server (it will use the Supabase DSN from `.env`):

```bash
make dev
```

What `make dev` does:
- Loads variables from `.env` into the environment and sets `DB_DSN` from `JOURNAL_DB_DSN` if present.
- Sets `APP_ENV=development` which toggles development behavior in the server (notably CSRF cookie secure flag).
- Runs `air` using the bundled `./bin/air` if present, otherwise the system `air`.

Notes:
- The app will now connect to whatever DSN you provide in `JOURNAL_DB_DSN` (we removed the local postgres service from `docker-compose.yml`).
- Make sure your Supabase instance allows connections from your environment and that credentials are correct.

---

If you want a short `dev` script or a `docker-compose.override.yml` created, tell me and I will add it.
