-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           TEXT NOT NULL,
    total_seats     INT NOT NULL,
    available_seats INT NOT NULL DEFAULT 0,
    event_at        TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bookings(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID REFERENCES events(id),
    status      TEXT NOT NULL CHECK( status in ('pending', 'confirmed', 'cancelled')) DEFAULT 'pending',
    telegram_id INT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_bookings_event_id ON bookings(event_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS bookings;
-- +goose StatementEnd
