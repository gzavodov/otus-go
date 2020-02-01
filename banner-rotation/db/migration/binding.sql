CREATE TABLE banner_binding
(
    id bigserial PRIMARY KEY,
    banner_id bigint NOT NULL,
    slot_id bigint NOT NULL
);

CREATE UNIQUE INDEX ix_banner_binding_banner_slot ON banner_binding (banner_id, slot_id);