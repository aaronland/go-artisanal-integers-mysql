# go-artisanal-integers-mysql

No, really.

## Caveats

### DSN

Although [the
documentation](https://github.com/go-sql-driver/mysql#dsn-data-source-name) for
the `go-sql-driver/mysql` package suggests that you can pass a DSN string
_without_ a password (for example `root:/@integers`) in practice that does not
appear to be case.

Specifically any attempt to query the database will fail with a `"Error 1046: No
database selected\n"` error.

## Schema

```
CREATE TABLE `integers` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `stub` char(1) NOT NULL DEFAULT '',
  PRIMARY KEY  (`id`),
  UNIQUE KEY `stub` (`stub`)
) ENGINE=MyISAM;
```

## Performance

Running `intd` backed by MySQL on a vanilla Vagrant machine (running Ubuntu 14.04) on a laptop against 500 concurrent users, using siege:

```
$> siege -c 500 http://localhost:8080
** SIEGE 3.0.5
** Preparing 500 concurrent users for battle.
The server is now under siege...^C
Lifting the server siege...      done.

Transactions:			58285 hits
Availability:			100.00 %
Elapsed time:			70.71 secs
Data transferred:		0.32 MB
Response time:			0.02 secs
Transaction rate:		824.28 trans/sec
Throughput:			0.00 MB/sec
Concurrency:			14.98
Successful transactions:	58217
Failed transactions:		0
Longest transaction:		1.70
Shortest transaction:		0.00
```

## See also

* https://github.com/aaronland/go-artisanal-integers
