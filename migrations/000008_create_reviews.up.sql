CREATE TABLE reviews (
    id              SERIAL PRIMARY KEY,
    reviewer_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    property_id     INTEGER REFERENCES properties(id) ON DELETE CASCADE,
    reviewee_id     INTEGER REFERENCES users(id) ON DELETE CASCADE,
    rating          SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT reviews_single_target CHECK (
        (property_id IS NOT NULL AND reviewee_id IS NULL) OR
        (property_id IS NULL     AND reviewee_id IS NOT NULL)
    ),
    UNIQUE(reviewer_id, property_id),
    UNIQUE(reviewer_id, reviewee_id)
);

CREATE INDEX idx_reviews_property_id ON reviews(property_id);
CREATE INDEX idx_reviews_reviewee_id ON reviews(reviewee_id);
CREATE INDEX idx_reviews_reviewer_id ON reviews(reviewer_id);
