CREATE SEQUENCE IF NOT EXISTS news_articles_id_seq;

-- Table Definition
CREATE TABLE "expenses" (
    id SERIAL PRIMARY KEY,
	title TEXT,
	amount FLOAT,
	note TEXT,
	tags TEXT[]
    PRIMARY KEY ("id")
);

INSERT INTO "news_articles" ("id", "title", "content", "author") VALUES (1, 'test-title', 'test-content', 'test-author');