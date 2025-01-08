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
EctoPSQLExtras.index_cache_hit(YourApp.Repo)
EctoPSQLExtras.table_cache_hit(YourApp.Repo)
```

## Index Usage
This command provides information on the efficiency of indexes, represented as what percentage of total scans were index scans. A low percentage can indicate under indexing, or wrong data being indexed.
```
EctoPSQLExtras.index_usage(Explorer.Repo)
```

## Outliers
This command displays statements, obtained from pg_stat_statements, ordered by the amount of time to execute in aggregate. This includes the statement itself, the total execution time for that statement, the proportion of total execution time for all statements that statement has taken up, the number of times that statement has been called, and the amount of time that statement spent on synchronous I/O (reading/writing from the file system).

Typically, an efficient query will have an appropriate ratio of calls to total execution time, with as little time spent on I/O as possible. Queries that have a high total execution time but low call count should be investigated to improve their performance. Queries that have a high proportion of execution time being spent on synchronous I/O should also be investigated.
```
EctoPSQLExtras.outliers(YourApp.Repo, args: [limit: 20])
```






---

Here are some more general queries

## DB settings
This method displays values for selected PostgreSQL settings. You can compare them with settings recommended by [PGTune](https://pgtune.leopard.in.ua/#/) and tweak values to improve performance.
```
EctoPSQLExtras.db_settings(YourApp.Repo)
```

---

Here are some queries which im unsure if they add value:

## Locks
This command displays queries that have taken out an exclusive lock on a relation. Exclusive locks typically prevent other operations on that relation from taking place, and can be a cause of "hung" queries that are waiting for a lock to be granted.
```
EctoPSQLExtras.locks(YourApp.Repo)
```

Similar:

This command displays all the current locks, regardless of their type.
```
EctoPSQLExtras.all_locks(YourApp.Repo)
```





EctoPSQLExtras.query(:cache_hit, Explorer.Repo)

EctoPSQLExtras.long_running_queries(Explorer.Repo, args: [threshold: "200 milliseconds"])

EctoPSQLExtras.diagnose(Explorer.Repo)