-- +migrate Up
CREATE TABLE IF NOT EXISTS certificate_revocation_list (
    id integer primary key,
    name text unique not null,
    signature  BLOB unique not null,
    this_update DATE not null,
    next_update DATE,
    url text,
    raw BLOB
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_certificate_revocation_list_name
    ON certificate_revocation_list(name);

CREATE TABLE IF NOT EXISTS revoked_certificate (
    id integer primary key,
    serialnumber text unique not null,
    revocation_date DATE not null,
    reason text not null,
    revocation_list integer not null,
    foreign key (revocation_list) references certificate_revocation_list(id)
       ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_revoked_certificates_serialnumber
    ON revoked_certificate(serialnumber);

-- +migrate Down
DROP TABLE certificate_revocation_list;

DROP TABLE revoked_certificate;
