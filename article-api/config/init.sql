
CREATE TABLE IF NOT EXISTS articles (
  id SERIAL not null unique,
  title VARCHAR(500) not null,
  date VARCHAR not null,
  body TEXT not null,
  tags TEXT[],
  primary key(id)
);

CREATE TABLE IF NOT EXISTS tags (
    tag VARCHAR PRIMARY KEY,
    count INTEGER NOT NULL,
    articles INTEGER[],
    related_tags TEXT[]
);

insert into articles(title, body, date, tags)
values
    ('latest science shows that potato chips are better for you than sugar',
     'some text, potentially containing simple markup about how potato chips are great ...', 
     '2023-04-05',
     '{"health", "fitness", "science"}'),
    ('Another Article', 
    'Yet another blog post about something exciting',
    '2023-04-06',
     '{"lifestyle", "fitness", "yoga"}'),
    ('One more article about health and well-being',
     'some text, potentially containing simple markup about how live a healthy lifestyle ...', 
     '2023-04-07',
     '{"health", "medical", "science"}');

insert into tags(tag, count, articles, related_tags)
values
    ('health', 
    17, 
    '{1,7}',
    '{"fitness", "science"}');