create table signon_context
(
    signon_id_pc INTEGER         not null constraint signon_context_pk primary key,
    state        TEXT            not null,
    comment      TEXT default '' not null
);

create table signons
(
    id                     INTEGER constraint signons_pk primary key,
    id_pc                  INTEGER not null,
    created_at             TEXT    not null,
    energy_type            TEXT,
    company                TEXT,
    firstname              TEXT,
    lastname               TEXT,
    zip                    TEXT,
    city                   TEXT,
    street                 TEXT,
    house_no               TEXT,
    pc_state               TEXT,
    desired_delivery_start TEXT,
    meter_no               TEXT,
    malo                   TEXT,
    melo                   TEXT,
    config_id              TEXT
);
