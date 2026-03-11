-- Kirona — init.sql
-- Exécuté une seule fois à la création du container PostgreSQL.
-- Crée une base par service pour isoler les schémas.

CREATE DATABASE kirona_auth;
CREATE DATABASE kirona_user;
CREATE DATABASE kirona_catalog;
CREATE DATABASE kirona_people;
CREATE DATABASE kirona_social;
CREATE DATABASE kirona_watchparty;
CREATE DATABASE kirona_notification;

-- Droits
GRANT ALL PRIVILEGES ON DATABASE kirona_auth       TO kirona;
GRANT ALL PRIVILEGES ON DATABASE kirona_user       TO kirona;
GRANT ALL PRIVILEGES ON DATABASE kirona_catalog    TO kirona;
GRANT ALL PRIVILEGES ON DATABASE kirona_people     TO kirona;
GRANT ALL PRIVILEGES ON DATABASE kirona_social     TO kirona;
GRANT ALL PRIVILEGES ON DATABASE kirona_watchparty TO kirona;
GRANT ALL PRIVILEGES ON DATABASE kirona_notification TO kirona;
