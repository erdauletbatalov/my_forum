-- user TABLE --
CREATE TABLE IF NOT EXISTS "user" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "username" TEXT NOT NULL UNIQUE,
  "email" TEXT NOT NULL UNIQUE,
  "password" TEXT NOT NULL
);

-- post TABLE --
CREATE TABLE IF NOT EXISTS "post" (
  "id" INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
  "user_id" INTEGER NOT NULL,
  "title" TEXT NOT NULL UNIQUE,
  "content" TEXT NOT NULL,
  CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES user(id)
);

-- comment TABLE --
CREATE TABLE IF NOT EXISTS "comment" (
  "id" INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
  "user_id" INTEGER NOT NULL,
  "post_id" INTEGER NOT NULL,
  "content" TEXT NOT NULL,
  CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES user(id),
  CONSTRAINT fk_post
    FOREIGN KEY (post_id)
    REFERENCES post(id)
);

-- comment TABLE --
CREATE TABLE IF NOT EXISTS "vote" (
  "id" INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
  "post_id" INTEGER NOT NULL,
  "comment_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "vote_obj" INTEGER NOT NULL,
  "vote_type" INTEGER NOT NULL,
  CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES user(id),
  CONSTRAINT fk_post
    FOREIGN KEY (post_id)
    REFERENCES post(id),
  CONSTRAINT fk_comment
    FOREIGN KEY (comment_id)
    REFERENCES comment(id)
);