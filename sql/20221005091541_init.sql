-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id         INTEGER     PRIMARY KEY ON CONFLICT IGNORE
                           NOT NULL,
    first_seen INTEGER     NOT NULL,
    lang       VARCHAR (2) DEFAULT en
);
CREATE TABLE data (
    upd            INTEGER PRIMARY KEY ON CONFLICT FAIL
                           NOT NULL,
    btcd           REAL    NOT NULL,
    ethd           REAL    NOT NULL,
    btcd_ch        REAL    NOT NULL,
    ethd_ch        REAL    NOT NULL,
    stablecoinscap REAL    NOT NULL,
    totalcap       REAL    NOT NULL,
    totalcap_ch    REAL    NOT NULL,
    liquid_usd     REAL    NOT NULL,
    liquid_usd_ch  REAL    NOT NULL,
    liquid_num     INTEGER NOT NULL,
    oi             INTEGER NOT NULL,
    oi_ch          REAL    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE data;
-- +goose StatementEnd
