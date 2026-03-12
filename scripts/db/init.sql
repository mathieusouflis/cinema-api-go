-- Kirona — init.sql
-- Exécuté une seule fois à la création du container PostgreSQL.
-- Crée une base par service pour isoler les schémas.
-- POSTGRES_USER=kirona crée l'utilisateur automatiquement.

-- kirona_auth est créé automatiquement via POSTGRES_DB dans docker-compose
CREATE DATABASE kirona_user;
CREATE DATABASE kirona_catalog;
CREATE DATABASE kirona_people;
CREATE DATABASE kirona_social;
CREATE DATABASE kirona_watchparty;
CREATE DATABASE kirona_notification;
