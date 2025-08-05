-- Migration: 001_init_watchlist_system.up.sql
-- Create initial tables for watchlist screening system

-- Users table (created first due to foreign key dependencies)
CREATE TABLE USERS (
   id BIGINT IDENTITY(1,1) PRIMARY KEY,
   username NVARCHAR(255) NOT NULL UNIQUE,
   email NVARCHAR(255) NOT NULL UNIQUE,
   password_hash NVARCHAR(255) NOT NULL,
   full_name NVARCHAR(255),
   role NVARCHAR(50) NOT NULL CHECK (role IN ('SUPER_ADMIN', 'ADMIN_UKPN', 'ADMIN_USER', 'OPERATOR', 'REVIEWER')),
   is_active BIT NOT NULL DEFAULT 1,
   last_login DATETIME2,
   created_by BIGINT,
   created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
   updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),
   FOREIGN KEY (created_by) REFERENCES USERS(id)
);

-- Create indexes for users
CREATE INDEX IX_USERS_username ON USERS(username);
CREATE INDEX IX_USERS_email ON USERS(email);
CREATE INDEX IX_USERS_role ON USERS(role);

-- User sessions
CREATE TABLE USER_SESSIONS (
   id BIGINT IDENTITY(1,1) PRIMARY KEY,
   user_id BIGINT NOT NULL,
   session_token NVARCHAR(255) NOT NULL UNIQUE,
   expires_at DATETIME2 NOT NULL,
   ip_address NVARCHAR(45),
   user_agent NVARCHAR(500),
   created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
   FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE
);

CREATE INDEX IX_USER_SESSIONS_user_id ON USER_SESSIONS(user_id);
CREATE INDEX IX_USER_SESSIONS_token ON USER_SESSIONS(session_token);
CREATE INDEX IX_USER_SESSIONS_expires ON USER_SESSIONS(expires_at);

-- Master customer table
CREATE TABLE MASTER_NASABAH (
                                id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                cif_number NVARCHAR(50) NOT NULL UNIQUE,
                                nama_nasabah NVARCHAR(255),
                                tempat_lahir NVARCHAR(255),
                                tanggal_lahir DATE,
                                ktp NVARCHAR(50),
                                npwp NVARCHAR(50),
                                no_paspor NVARCHAR(50),
                                status_nasabah NVARCHAR(20) CHECK (status_nasabah IN ('ACTIVE', 'INACTIVE')),
                                created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                                updated_at DATETIME2 NOT NULL DEFAULT GETDATE()
);

CREATE INDEX IX_MASTER_NASABAH_cif ON MASTER_NASABAH(cif_number);
CREATE INDEX IX_MASTER_NASABAH_nama ON MASTER_NASABAH(nama_nasabah);
CREATE INDEX IX_MASTER_NASABAH_ktp ON MASTER_NASABAH(ktp);
CREATE INDEX IX_MASTER_NASABAH_status ON MASTER_NASABAH(status_nasabah);

-- Terrorist watchlist
CREATE TABLE MASTER_TERORIS (
                                id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                nama NVARCHAR(255),
                                alias1 NVARCHAR(255),
                                alias2 NVARCHAR(255),
                                alias3 NVARCHAR(255),
                                alias4 NVARCHAR(255),
                                tempat_lahir NVARCHAR(255),
                                tanggal_lahir DATE,
                                ktp NVARCHAR(50),
                                npwp NVARCHAR(50),
                                no_paspor NVARCHAR(50),
                                created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                                updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                                is_active BIT NOT NULL DEFAULT 1
);

CREATE INDEX IX_MASTER_TERORIS_nama ON MASTER_TERORIS(nama);
CREATE INDEX IX_MASTER_TERORIS_active ON MASTER_TERORIS(is_active);

-- WMD (Weapons of Mass Destruction) watchlist
CREATE TABLE MASTER_WMD (
                            id BIGINT IDENTITY(1,1) PRIMARY KEY,
                            nama NVARCHAR(255),
                            alias1 NVARCHAR(255),
                            alias2 NVARCHAR(255),
                            alias3 NVARCHAR(255),
                            alias4 NVARCHAR(255),
                            alias5 NVARCHAR(255),
                            alias6 NVARCHAR(255),
                            alias7 NVARCHAR(255),
                            alias8 NVARCHAR(255),
                            alias9 NVARCHAR(255),
                            alias10 NVARCHAR(255),
                            tempat_lahir NVARCHAR(255),
                            tanggal_lahir DATE,
                            ktp NVARCHAR(50),
                            npwp NVARCHAR(50),
                            no_paspor NVARCHAR(50),
                            created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                            updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                            is_active BIT NOT NULL DEFAULT 1
);

CREATE INDEX IX_MASTER_WMD_nama ON MASTER_WMD(nama);
CREATE INDEX IX_MASTER_WMD_active ON MASTER_WMD(is_active);

-- Local blacklist
CREATE TABLE MASTER_LOCAL_BLACKLIST (
                                        id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                        nama NVARCHAR(255),
                                        alias1 NVARCHAR(255),
                                        alias2 NVARCHAR(255),
                                        alias3 NVARCHAR(255),
                                        alias4 NVARCHAR(255),
                                        tempat_lahir NVARCHAR(255),
                                        tanggal_lahir DATE,
                                        ktp NVARCHAR(50),
                                        npwp NVARCHAR(50),
                                        no_paspor NVARCHAR(50),
                                        created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                                        updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                                        is_active BIT NOT NULL DEFAULT 1
);

CREATE INDEX IX_MASTER_LOCAL_BLACKLIST_nama ON MASTER_LOCAL_BLACKLIST(nama);
CREATE INDEX IX_MASTER_LOCAL_BLACKLIST_active ON MASTER_LOCAL_BLACKLIST(is_active);

-- Matching configuration
CREATE TABLE MATCHING_CONFIG (
                                 id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                 watchlist_source NVARCHAR(50),
                                 field_name NVARCHAR(100) NOT NULL UNIQUE,
                                 field_weight DECIMAL(5,4),
                                 matching_algorithm NVARCHAR(50) CHECK (matching_algorithm IN ('jaro_winkler', 'levenshtein')),
                                 is_active BIT NOT NULL DEFAULT 1,
                                 created_by BIGINT,
                                 created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                                 updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                                 FOREIGN KEY (created_by) REFERENCES USERS(id)
);

CREATE INDEX IX_MATCHING_CONFIG_source ON MATCHING_CONFIG(watchlist_source);
CREATE INDEX IX_MATCHING_CONFIG_active ON MATCHING_CONFIG(is_active);

-- System configuration
CREATE TABLE SYSTEM_CONFIG (
                               id BIGINT IDENTITY(1,1) PRIMARY KEY,
                               config_key NVARCHAR(100) NOT NULL UNIQUE,
                               config_value NVARCHAR(500),
                               config_type NVARCHAR(20) CHECK (config_type IN ('DECIMAL', 'INTEGER', 'STRING')),
                               description NVARCHAR(500),
                               created_by BIGINT,
                               created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                               updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                               FOREIGN KEY (created_by) REFERENCES USERS(id)
);

CREATE INDEX IX_SYSTEM_CONFIG_key ON SYSTEM_CONFIG(config_key);

-- Batch processing
CREATE TABLE BATCH_PROCESSING (
                                  id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                  process_type NVARCHAR(20) CHECK (process_type IN ('MANUAL', 'SCHEDULED')),
                                  status NVARCHAR(20) CHECK (status IN ('RUNNING', 'COMPLETED', 'FAILED')),
                                  total_records INT,
                                  processed_records INT,
                                  matched_records INT,
                                  start_time DATETIME2,
                                  end_time DATETIME2,
                                  error_message NVARCHAR(MAX),
                                  initiated_by BIGINT,
                                  file_path NVARCHAR(500),
                                  FOREIGN KEY (initiated_by) REFERENCES USERS(id)
);

CREATE INDEX IX_BATCH_PROCESSING_status ON BATCH_PROCESSING(status);
CREATE INDEX IX_BATCH_PROCESSING_type ON BATCH_PROCESSING(process_type);
CREATE INDEX IX_BATCH_PROCESSING_start_time ON BATCH_PROCESSING(start_time);

-- Processing log
CREATE TABLE PROCESSING_LOG (
                                id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                batch_id BIGINT,
                                log_level NVARCHAR(10) CHECK (log_level IN ('INFO', 'WARN', 'ERROR')),
                                log_message NVARCHAR(MAX),
                                error_details NVARCHAR(MAX),
                                log_time DATETIME2 NOT NULL DEFAULT GETDATE(),
                                FOREIGN KEY (batch_id) REFERENCES BATCH_PROCESSING(id) ON DELETE CASCADE
);

CREATE INDEX IX_PROCESSING_LOG_batch_id ON PROCESSING_LOG(batch_id);
CREATE INDEX IX_PROCESSING_LOG_level ON PROCESSING_LOG(log_level);
CREATE INDEX IX_PROCESSING_LOG_time ON PROCESSING_LOG(log_time);

-- Matching results
CREATE TABLE MATCHING_RESULTS (
                                  id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                  batch_id BIGINT,
                                  cif_number NVARCHAR(50),
                                  customer_name NVARCHAR(255),
                                  watchlist_id BIGINT,
                                  watchlist_source NVARCHAR(50) CHECK (watchlist_source IN ('DTTOT', 'WMD', 'LOCAL_BLACKLIST')),
                                  similarity_score DECIMAL(5,4),
                                  status NVARCHAR(20) CHECK (status IN ('PENDING', 'REVIEWED', 'APPROVED', 'REJECTED')),
                                  process_date DATE,
                                  process_time DATETIME2,
                                  created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
                                  FOREIGN KEY (batch_id) REFERENCES BATCH_PROCESSING(id),
                                  FOREIGN KEY (cif_number) REFERENCES MASTER_NASABAH(cif_number)
);

CREATE INDEX IX_MATCHING_RESULTS_batch_id ON MATCHING_RESULTS(batch_id);
CREATE INDEX IX_MATCHING_RESULTS_cif ON MATCHING_RESULTS(cif_number);
CREATE INDEX IX_MATCHING_RESULTS_status ON MATCHING_RESULTS(status);
CREATE INDEX IX_MATCHING_RESULTS_source ON MATCHING_RESULTS(watchlist_source);
CREATE INDEX IX_MATCHING_RESULTS_score ON MATCHING_RESULTS(similarity_score);

-- Matching details
CREATE TABLE MATCHING_DETAILS (
                                  id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                  matching_result_id BIGINT NOT NULL,
                                  field_name NVARCHAR(100),
                                  customer_value NVARCHAR(500),
                                  watchlist_value NVARCHAR(500),
                                  field_score DECIMAL(5,4),
                                  field_weight DECIMAL(5,4),
                                  algorithm_used NVARCHAR(50),
                                  FOREIGN KEY (matching_result_id) REFERENCES MATCHING_RESULTS(id) ON DELETE CASCADE
);

CREATE INDEX IX_MATCHING_DETAILS_result_id ON MATCHING_DETAILS(matching_result_id);
CREATE INDEX IX_MATCHING_DETAILS_field ON MATCHING_DETAILS(field_name);

-- Audit log
CREATE TABLE AUDIT_LOG (
                           id BIGINT IDENTITY(1,1) PRIMARY KEY,
                           user_id BIGINT,
                           table_name NVARCHAR(100),
                           operation NVARCHAR(10) CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE', 'SELECT')),
                           record_id NVARCHAR(50),
                           old_values NVARCHAR(MAX),
                           new_values NVARCHAR(MAX),
                           ip_address NVARCHAR(45),
                           action_time DATETIME2 NOT NULL DEFAULT GETDATE(),
                           FOREIGN KEY (user_id) REFERENCES USERS(id)
);

CREATE INDEX IX_AUDIT_LOG_user_id ON AUDIT_LOG(user_id);
CREATE INDEX IX_AUDIT_LOG_table ON AUDIT_LOG(table_name);
CREATE INDEX IX_AUDIT_LOG_operation ON AUDIT_LOG(operation);
CREATE INDEX IX_AUDIT_LOG_time ON AUDIT_LOG(action_time);