CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS cascade;
-- Create Team and User Management related relations
CREATE TABLE IF NOT EXISTS cascade.team (
    team_id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    team_name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS cascade.users (
    user_id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    email VARCHAR(100) NOT NULL UNIQUE,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin', 'manager')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS cascade.user_team (
    user_id UUID NOT NULL,
    team_id UUID NOT NULL,
    PRIMARY KEY (user_id, team_id),
    FOREIGN KEY (user_id) REFERENCES cascade.users(user_id),
    FOREIGN KEY (team_id) REFERENCES cascade.team(team_id)
);
-- Create Chaos Engineering related relations
CREATE TABLE IF NOT EXISTS cascade.scenarios (
    scenario_id UUID NOT NULL DEFAULT (uuid_generate_v4()),
    version INT NOT NULL DEFAULT 1,
    description TEXT,
    team_id UUID NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (scenario_id, version),
    FOREIGN KEY (team_id) REFERENCES cascade.team(team_id)
);
CREATE TABLE IF NOT EXISTS cascade.sessions (
    session_id SERIAL PRIMARY KEY,
    scenario_id UUID NOT NULL,
    version INT NOT NULL,
    user_id UUID NOT NULL,
    start_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP,
    status VARCHAR(20) NOT NULL CHECK (
        status IN ('queued', 'running', 'completed', 'failed')
    ),
    FOREIGN KEY (scenario_id, version) REFERENCES cascade.scenarios(scenario_id, version),
    FOREIGN KEY (user_id) REFERENCES cascade.users(user_id)
);
-- Add indexes
CREATE INDEX idx_scenario_team_id ON cascade.scenarios(team_id);
CREATE INDEX idx_session_scenario_id ON cascade.sessions(scenario_id);
CREATE INDEX idx_session_user_id ON cascade.sessions(user_id);
CREATE INDEX idx_user_team_user_id ON cascade.user_team(user_id);
CREATE INDEX idx_user_team_team_id ON cascade.user_team(team_id);
