BEGIN;

CREATE TABLE IF NOT EXISTS races (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	date TIMESTAMP NOT NULL,
    owner_id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS races_competitors (
	race_id UUID,
	competitor_id UUID,
	register_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY(race_id, competitor_id)
);

CREATE TABLE IF NOT EXISTS teams (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	admin_id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_members (
	team_id UUID PRIMARY KEY REFERENCES teams (id),
	member_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
	id UUID PRIMARY KEY,
	payload JSONB NOT NULL,
	user_id UUID NOT NULL,
	occurred_at TIMESTAMP NOT NULL
);
COMMIT;
