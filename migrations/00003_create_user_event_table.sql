-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_events (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    event_type INT NOT NULL,
    user_id INT NOT NULL,
    entity_id INT NOT NULL,
    value jsonb
);

CREATE INDEX user_events_entity_id_idx ON user_events(entity_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_events;

DROP INDEX user_events_entity_id_idx;
-- +goose StatementEnd
