-- +goose Up
-- +migrate Up

CREATE TABLE AUDIT_LOG (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    UserId BIGINT,
    TableName NVARCHAR(255),
    Operation NVARCHAR(255),
    RecordId NVARCHAR(255),
    OldValues TEXT,
    NewValues TEXT,
    IpAddress NVARCHAR(255),
    ActionTime DATETIME
);

CREATE TABLE BATCH_PROCESSING (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    ProcessType NVARCHAR(255),
    Status NVARCHAR(255),
    TotalRecords INT,
    ProcessedRecords INT,
    MatchedRecords INT,
    StartTime DATETIME,
    EndTime DATETIME,
    ErrorMessage NVARCHAR(255),
    InitiatedBy INT,
    FilePath NVARCHAR(255)
);

CREATE TABLE MASTER_LOCAL_BLACKLIST (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    Nama NVARCHAR(255),
    Alias1 NVARCHAR(255),
    Alias2 NVARCHAR(255),
    Alias3 NVARCHAR(255),
    Alias4 NVARCHAR(255),
    Type NVARCHAR(255),
    TempatLahir NVARCHAR(255),
    TanggalLahir DATE,
    Ktp NVARCHAR(255),
    Npwp NVARCHAR(255),
    NoPaspor NVARCHAR(255),
    CreatedAt DATETIME,
    UpdatedAt DATETIME,
    IsActive BIT
);

CREATE TABLE MASTER_NASABAH (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    CifNumber NVARCHAR(255) UNIQUE,
    NamaNasabah NVARCHAR(255),
    TempatLahir NVARCHAR(255),
    TanggalLahir DATE,
    Ktp NVARCHAR(255),
    Npwp NVARCHAR(255),
    NoPaspor NVARCHAR(255),
    StatusNasabah NVARCHAR(255),
    CreatedAt DATETIME,
    UpdatedAt DATETIME
);

CREATE TABLE MASTER_TERORIS (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    Nama NVARCHAR(255),
    Alias1 NVARCHAR(255),
    Alias2 NVARCHAR(255),
    Alias3 NVARCHAR(255),
    Alias4 NVARCHAR(255),
    Type NVARCHAR(255),
    TempatLahir NVARCHAR(255),
    TanggalLahir DATE,
    Ktp NVARCHAR(255),
    Npwp NVARCHAR(255),
    NoPaspor NVARCHAR(255),
    CreatedAt DATETIME,
    UpdatedAt DATETIME,
    IsActive BIT
);

CREATE TABLE MASTER_WMD (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    Nama NVARCHAR(255),
    Alias1 NVARCHAR(255),
    Alias2 NVARCHAR(255),
    Alias3 NVARCHAR(255),
    Alias4 NVARCHAR(255),
    Alias5 NVARCHAR(255),
    Alias6 NVARCHAR(255),
    Alias7 NVARCHAR(255),
    Alias8 NVARCHAR(255),
    Alias9 NVARCHAR(255),
    Alias10 NVARCHAR(255),
    Type NVARCHAR(255),
    TempatLahir NVARCHAR(255),
    TanggalLahir DATE,
    Ktp NVARCHAR(255),
    Npwp NVARCHAR(255),
    NoPaspor NVARCHAR(255),
    CreatedAt DATETIME,
    UpdatedAt DATETIME,
    IsActive BIT
);

CREATE TABLE MASTER_MATCHING (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    Name NVARCHAR(255),
    WatchlistSource NVARCHAR(255),
    IsIndividual BIT,
    Description NVARCHAR(255),
    IsActive BIT
);

INSERT INTO MASTER_MATCHING (Name, WatchlistSource, IsIndividual, Description, IsActive)
VALUES
    ('Matching Data Teroris (Individu)', 'MASTER_TERORIS', 1, 'Pemadanan Data Watchlist Daftar Teroris dan Orgainsasi Teroris (DTTOT) Individu',1),
    ('Matching Data Teroris (Korporasi)', 'MASTER_TERORIS', 0, 'Pemadanan Data Watchlist Daftar Teroris dan Orgainsasi Teroris (DTTOT) Korporasi',1),
    ('Matching Data WMD (Individu)', 'MASTER_WMD', 1, 'Pemadanan Data Watchlist Senjata Pemusnah Massal (WMD) Individu',1),
    ('Matching Data WMD (Korporasi)', 'MASTER_WMD', 0, 'Pemadanan Data Watchlist Senjata Pemusnah Massal (WMD) Korporasi',1),
    ('Matching Data Local Blacklist (Individu)', 'MASTER_LOCAL_BLACKLIST', 1, 'Pemadanan Data Watchlist Local Blacklist Individu',1),
    ('Matching Data Local Blacklist (Korporasi)', 'MASTER_LOCAL_BLACKLIST', 0, 'Pemadanan Data Watchlist Local Blacklist Korporasi',1);

CREATE TABLE MASTER_MATCHING_CONFIG (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    MatchingId BIGINT,
    FieldName NVARCHAR(255) UNIQUE,
    FieldWeight DECIMAL(18,0),
    MatchingAlgorithm NVARCHAR(255),
    IsActive BIT,
    CreatedBy INT,
    CreatedAt DATETIME,
    UpdatedAt DATETIME
);

CREATE TABLE MASTER_MATCHING_CONFIG_CHANGE_LOG (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    MatchingConfigId BIGINT,
    FieldName NVARCHAR(255),
    OldValue NVARCHAR(255),
    NewValue NVARCHAR(255),
    CreatedBy INT,
    CreatedAt DATETIME,
    ApprovedBy INT,
    ApprovedAt DATETIME,
    ChangedBy INT,
    ChangeDate DATETIME
);

CREATE TABLE MATCHING_DETAILS (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    MatchingResultId BIGINT,
    FieldName NVARCHAR(255),
    CustomerValue NVARCHAR(255),
    WatchlistValue NVARCHAR(255),
    FieldScore DECIMAL(18,0),
    FieldWeight DECIMAL(18,0),
    AlgorithmUsed NVARCHAR(255)
);

CREATE TABLE MATCHING_RESULTS (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    BatchId BIGINT,
    CifNumber NVARCHAR(255),
    CustomerName NVARCHAR(255),
    WatchlistId BIGINT,
    WatchlistSource NVARCHAR(255),
    SimilarityScore DECIMAL(18,0),
    Status NVARCHAR(255),
    ProcessDate DATE,
    ProcessTime DATETIME,
    CreatedAt DATETIME
);

CREATE TABLE PROCESSING_LOG (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    BatchId BIGINT,
    LogLevel NVARCHAR(255),
    LogMessage NVARCHAR(255),
    ErrorDetails NVARCHAR(255),
    LogTime DATETIME
);

CREATE TABLE SYSTEM_CONFIG (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    ConfigKey NVARCHAR(255) UNIQUE,
    ConfigValue NVARCHAR(255),
    ConfigType NVARCHAR(255),
    Description NVARCHAR(255),
    CreatedBy INT,
    CreatedAt DATETIME,
    UpdatedAt DATETIME
);

CREATE TABLE USER_SESSIONS (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    UserId BIGINT,
    SessionToken NVARCHAR(255) UNIQUE,
    ExpiresAt DATETIME,
    IpAddress NVARCHAR(255),
    UserAgent NVARCHAR(255),
    CreatedAt DATETIME
);

CREATE TABLE USERS (
    Id BIGINT PRIMARY KEY IDENTITY(1,1),
    Username NVARCHAR(255) UNIQUE,
    Email NVARCHAR(255) UNIQUE,
    PasswordHash NVARCHAR(255),
    FullName NVARCHAR(255),
    Role NVARCHAR(255),
    IsActive BIT,
    LastLogin DATETIME,
    CreatedBy INT,
    CreatedAt DATETIME,
    UpdatedAt DATETIME
);
