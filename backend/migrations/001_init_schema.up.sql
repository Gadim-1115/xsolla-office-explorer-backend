-- Enable pgcrypto for UUID generation if needed
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. ZONES TABLE
-- id strictly matches Figma SVG layer IDs (e.g., zone_gameroom_1, zone_kitchen)
CREATE TABLE IF NOT EXISTS zones (
                                     id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    is_locked_default BOOLEAN NOT NULL DEFAULT FALSE,
    difficulty VARCHAR(20) NULL,
    media_url TEXT NOT NULL,
    description TEXT NOT NULL,
    CONSTRAINT chk_zone_difficulty CHECK (
    (is_locked_default = TRUE AND difficulty IN ('EASY', 'MEDIUM', 'MEDIUM_HARD', 'HARD'))
    OR
(is_locked_default = FALSE AND difficulty IS NULL)
    )
    );

-- 2. QUESTS TABLE
-- 1:1 relationship with locked zones.
-- code_blocks stores [{"id": "b1", "code": "..."}]
-- correct_order stores ["b1", "b2", "b3"]
CREATE TABLE IF NOT EXISTS quests (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id VARCHAR(50) NOT NULL UNIQUE REFERENCES zones(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    problem_statement TEXT NOT NULL,
    code_blocks JSONB NOT NULL,
    correct_order JSONB NOT NULL
    );

-- 3. HOTSPOTS TABLE
-- Relative coordinates (0.00 to 100.00) ensure responsive scaling across map and room modals.
CREATE TABLE IF NOT EXISTS hotspots (
                                        id VARCHAR(50) PRIMARY KEY,
    zone_id VARCHAR(50) NULL REFERENCES zones(id) ON DELETE CASCADE,
    category VARCHAR(30) NOT NULL,
    x_percent NUMERIC(5,2) NOT NULL,
    y_percent NUMERIC(5,2) NOT NULL,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    CONSTRAINT chk_hotspot_category CHECK (
                                              category IN ('COFFEE_MACHINE', 'FRIDGE', 'DISPENSER', 'PROJECTOR', 'VIEWPOINT', 'EASTER_EGG')
    ),
    CONSTRAINT chk_hotspot_coordinates CHECK (
                                                 x_percent >= 0.00 AND x_percent <= 100.00 AND
                                                 y_percent >= 0.00 AND y_percent <= 100.00
                                             )
    );

CREATE INDEX IF NOT EXISTS idx_hotspots_zone_id ON hotspots(zone_id);

-- 4. USER PROGRESS TABLE
-- Stores idempotent sets of unlocked zone_id and claimed hotspot_id strings.
CREATE TABLE IF NOT EXISTS user_progress (
                                             user_id UUID PRIMARY KEY,
                                             unlocked_zones JSONB NOT NULL DEFAULT '[]'::jsonb,
                                             found_eggs JSONB NOT NULL DEFAULT '[]'::jsonb,
                                             updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_user_progress_unlocked_zones ON user_progress USING GIN (unlocked_zones);
CREATE INDEX IF NOT EXISTS idx_user_progress_found_eggs ON user_progress USING GIN (found_eggs);