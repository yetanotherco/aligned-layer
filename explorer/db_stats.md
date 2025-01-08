This document contains a list of commands to run on the explorer to view some relevant DB stats:

Info gathered from:

https://hexdocs.pm/ecto_psql_extras/readme.html

1) Attach to the running node:
```
make explorer_attach
```

2) execute any of the following commands:

## Cache hit
This command provides information on the efficiency of the buffer cache, for both index reads (index hit rate) as well as table reads (table hit rate). A low buffer cache hit ratio can be a sign that the Postgres instance is too small for the workload.
```
EctoPSQLExtras.cache_hit(Explorer.Repo)
```
Similar:
```
EctoPSQLExtras.index_cache_hit(Explorer.Repo)
EctoPSQLExtras.table_cache_hit(Explorer.Repo)
```

## Index Usage
This command provides information on the efficiency of indexes, represented as what percentage of total scans were index scans. A low percentage can indicate under indexing, or wrong data being indexed.
```
EctoPSQLExtras.index_usage(Explorer.Repo)
```

Similar:
This command displays indexes that have < 50 scans recorded against them, and are greater than 5 pages in size, ordered by size relative to the number of index scans.
```
EctoPSQLExtras.unused_indexes(Explorer.Repo, args: [min_scans: 20])
```


```
EctoPSQLExtras.duplicate_indexes(Explorer.Repo)
```
```
EctoPSQLExtras.null_indexes(Explorer.Repo, args: [min_relation_size_mb: 10])
```

## Outliers
This command displays statements, obtained from pg_stat_statements, ordered by the amount of time to execute in aggregate. This includes the statement itself, the total execution time for that statement, the proportion of total execution time for all statements that statement has taken up, the number of times that statement has been called, and the amount of time that statement spent on synchronous I/O (reading/writing from the file system).

Typically, an efficient query will have an appropriate ratio of calls to total execution time, with as little time spent on I/O as possible. Queries that have a high total execution time but low call count should be investigated to improve their performance. Queries that have a high proportion of execution time being spent on synchronous I/O should also be investigated.
```
EctoPSQLExtras.outliers(Explorer.Repo, args: [limit: 20])
```

Similar:
```
EctoPSQLExtras.calls(Explorer.Repo, args: [limit: 20])
```

## Sequential Scans
This command displays the number of sequential scans recorded against all tables, descending by count of sequential scans. Tables that have very high numbers of sequential scans may be under-indexed, and it may be worth investigating queries that read from these tables.
```
EctoPSQLExtras.seq_scans(Explorer.Repo)
```

## Current connections
```
EctoPSQLExtras.connections(Explorer.Repo)
```

---

Here are some more general queries

## DB settings
This method displays values for selected PostgreSQL settings. You can compare them with settings recommended by [PGTune](https://pgtune.leopard.in.ua/#/) and tweak values to improve performance.
```
EctoPSQLExtras.db_settings(Explorer.Repo)
```

## Table Size
This command displays the size of each table and materialized view in the database, in MB. 
```
EctoPSQLExtras.table_size(Explorer.Repo)
```

Similar:
```
EctoPSQLExtras.table_indexes_size(Explorer.Repo)
```
```
EctoPSQLExtras.total_table_size(Explorer.Repo)
```

---

Here are some queries which im unsure if they add value:

## Locks
This command displays queries that have taken out an exclusive lock on a relation. Exclusive locks typically prevent other operations on that relation from taking place, and can be a cause of "hung" queries that are waiting for a lock to be granted.
```
EctoPSQLExtras.locks(Explorer.Repo)
```

Similar:

This command displays all the current locks, regardless of their type.
```
EctoPSQLExtras.all_locks(Explorer.Repo)
```

## Bloat
This command displays an estimation of table "bloat" – space allocated to a relation that is full of dead tuples, that has yet to be reclaimed. Tables that have a high bloat ratio, typically 10 or greater, should be investigated to see if vacuuming is aggressive enough, and can be a sign of high table churn.
```
EctoPSQLExtras.bloat(Explorer.Repo)
```


## Long running queries
This command displays _currently running_ queries, that have been running for longer than 5 minutes.
```
EctoPSQLExtras.long_running_queries(Explorer.Repo, args: [threshold: "200 milliseconds"])
```

## Mandelbrot
```
EctoPSQLExtras.mandelbrot(Explorer.Repo)
```


EctoPSQLExtras.query(:cache_hit, Explorer.Repo)


EctoPSQLExtras.diagnose(Explorer.Repo)