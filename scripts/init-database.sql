-- init-databases.sql
-- This script runs ONCE when PostgreSQL container first starts.
-- It creates the three necessary databases and enables shared extensions/functions in each.

-- =====================================================
-- 1. CREATE DATABASE
-- =====================================================
CREATE DATABASE telex_waitlist;

-- =====================================================
-- 2. APPLY EXTENSIONS AND FUNCTIONS TO EACH DB
-- =====================================================

-- Connect to user_service_db
\c user_service_db;

-- Enable UUID and text search extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create the shared update_at function in this DB
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
