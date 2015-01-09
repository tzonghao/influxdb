# The Influx Query Language Specification

## Introduction

This is a reference for the Influx Query Language ("InfluxQL").

InfluxQL is a SQL-like query language for interacting with InfluxDB.  It was lovingly crafted to feel familiar to those coming from other
SQL or SQL-like environments while providing features specific to storing
and analyzing time series data.

## Notation

This specification uses the same notation used by Google's Go programming language, which can be found at http://golang.org.  The syntax is specified in Extended Backus-Naur Form ("EBNF"):

```
Production  = production_name "=" [ Expression ] "." .
Expression  = Alternative { "|" Alternative } .
Alternative = Term { Term } .
Term        = production_name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .
```

Notation operators in order of increasing precedence:

```
|   alternation
()  grouping
[]  option (0 or 1 times)
{}  repetition (0 to n times)
```

## Characters & Digits

```
newline             = /* the Unicode code point U+000A */ .
unicode_char        = /* an arbitrary Unicode code point except newline */ .
ascii_letter        = "A" .. "Z" | "a" .. "z" .
decimal_digit       = "0" .. "9" .
```

## Database name

Database names are more limited than other identifiers because they appear in URLs.

```
db_name             = ascii_letter { ascii_letter | decimal_digit | "_" | "-" } .
```

## Identifiers

```
identifier          = unquoted_identifier | quoted_identifier .
unquoted_identifier = ascii_letter { ascii_letter | decimal_digit | "_" | "." } .
quoted_identifier   = `"` unicode_char { unicode_char } `"` .
```

## Keywords

```
ALL        ALTER    AS          ASC          BEGIN
BY         CREATE   CONTINUOUS  DATABASE     DEFAULT
DELETE     DESC     DROP        DURATION     END
EXISTS     EXPLAIN  FIELD       FROM         GRANT
GROUP      IF       INNER       INSERT       INTO
KEYS       LIMIT    LIST        MEASUREMENT  MEASUREMENTS
ON         ORDER    PASSWORD    POLICY       PRIVILEGES
QUERIES    QUERY    READ        REPLICATION  RETENTION
EVOKE      SELECT   SERIES      TAG          TO
USER       VALUES   WHERE       WITH         WRITE
```

## Literals

### Numbers

```
int_lit             = decimal_lit .
decimal_lit         = ( "1" .. "9" ) { decimal_digit } .
float_lit           = decimals "." decimals .
decimals            = decimal_digit { decimal_digit } .
```

### Strings

```
string_lit          = '"' { unicode_char } '"' .
```

### Durations

```
duration_lit        = decimals duration_unit .
duration_unit       = "u" | "µ" | "s" | "h" | "d" | "w" | "ms" .
```

## Queries

A query is composed of one or more statements separated by a semicolon.

```
query               = statement { ; statement } .

statement           = alter_retention_policy_stmt |
					            create_continuous_query_stmt |
                      create_database_stmt |
                      create_retention_policy_stmt |
                      create_user_stmt |
                      delete_stmt |
                      drop_continuous_query_stmt |
                      drop_database_stmt |
                      drop_series_stmt |
                      drop_user_stmt |
                      grant_stmt |
                      list_continuous_queries_stmt |
                      list_databases_stmt |
                      list_field_key_stmt |
                      list_field_value_stmt |
                      list_measurements_stmt |
                      list_series_stmt |
                      list_tag_key_stmt |
                      list_tag_value_stmt |
                      revoke_stmt |
                      select_stmt .
```

## Statements

### ALTER RETENTION POLICY

```
alter_retention_policy_stmt  = "ALTER RETENTION POLICY" policy_name "ON"
                               db_name retention_policy_option
                               [ retention_policy_option ]
                               [ retention_policy_option ] .

policy_name                  = identifier .

retention_policy_option      = retention_policy_duration |
                               retention_policy_replication |
                               "DEFAULT" .

retention_policy_duration    = "DURATION" duration_lit .
retention_policy_replication = "REPLICATION" int_lit
```

#### Examples:

```sql
-- Set default retention policy for mydb to 1h.cpu.
ALTER RETENTION POLICY "1h.cpu" ON mydb DEFAULT;

-- Change duration and replication factor.
ALTER RETENTION POLICY policy1 ON somedb DURATION 1h REPLICATION 4
```

### CREATE CONTINUOUS QUERY

```
create_continuous_query_stmt = "CREATE CONTINUOUS QUERY" query_name "ON" db_name
                               "BEGIN" select_stmt "END" .

query_name                   = identifier .
```

#### Examples:

```sql
CREATE CONTINUOUS QUERY 10m_event_count
ON db_name
BEGIN
  SELECT count(value)
  INTO 10m.events
  FROM events
  GROUP BY time(10m)
END;

-- this selects from the output of one continuous query and outputs to another series
CREATE CONTINUOUS QUERY 1h_event_count
ON db_name
BEGIN
  SELECT sum(count) as count
  INTO 1h.events
  FROM events
  GROUP BY time(1h)
END;
```

### CREATE DATABASE

```
create_database_stmt = "CREATE DATABASE" db_name
```

#### Example:

```sql
CREATE DATABASE foo
```

### CREATE RETENTION POLICY

```
create_retention_policy_stmt = "CREATE RETENTION POLICY" policy_name "ON"
                               db_name retention_policy_duration
                               retention_policy_replication
                               [ "DEFAULT" ] .
```

#### Examples

```sql
-- Create a retention policy.
CREATE RETENTION POLICY "10m.events" ON somedb DURATION 10m REPLICATION 2;

-- Create a retention policy and set it as the default.
CREATE RETENTION POLICY "10m.events" ON somedb DURATION 10m REPLICATION 2 DEFAULT;
```

### CREATE USER

```
create_user_stmt = "CREATE USER" user_name "WITH PASSWORD" password .
```

#### Example:

```sql
CREATE USER jdoe WITH PASSWORD "1337password";
```

### DELETE

```
delete_stmt  = "DELETE" from_clause where_clause .
```

#### Example:

```sql
DELETE FROM 
```

### GRANT

```
grant_stmt = "GRANT" privilege [ on_clause ] to_clause
```

#### Examples:

```sql
-- grant cluster admin privileges
GRANT ALL TO jdoe;

-- grant read access to a database
GRANT READ ON mydb TO jdoe;
```

## Clauses (sadly, we haven't implemented `SANTA` yet)

```
from_clause  = "FROM" measurements .

where_clause = "WHERE" 

on_clause    = db_name .

to_clause    = user_name .
```

## Other

```
measurements =

user_name        = identifier .

password         = identifier .

privilege        = "ALL" [ "PRIVILEGES" ] | "READ" | "WRITE" .
```
