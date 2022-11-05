To make things persistent between sessions, `exatorrent` uses database. But because `exatorrent` uses awesome `sqlite3` which comes embedded within `exatorrent`, you barely notice it.
Instead of `sqlite3`, you can also use `postgresql` if you want to. Know that both `sqlite` and `postgresql` are provided as choices of Database Implementations. When in doubt, use `sqlite3`(which is used by `exatorrent` as default
,i.e, don't worry about configuring Database at all.


Note that once you start using one Database, you must stick to it. You cannot jump between Database Implementations and if you do, you loose data.

### Data Stored in Database

1. Users and their data
2. Trackers
3. State of Torrent
4. Piece Completion State of Torrent
5. State of Files of Torrent
6. Lock State of Torrent


## Postgresql

Normally you will be fine  using `sqlite3` which `exatorrent` uses by default, which you don't need to setup and configure. But, Postgresql is also provided as choice. If you want to try out Postgresql as Database (Note that you can't switch back to sqlite later), follow instructions below :

* Create Postgresql Database. Remember it's credentials.
* Remember format for Connection URL: `postgres://username:password@localhost:5432/database_name`
* You need to pass connection URL to `exatorrent`. You can either set `DATABASE_URL` environment variable as connection URL or you can write connection URL to file at `<exatorrent-directory>/config/psqlconfig.txt`. If you want sample connection URL written at `psqlconfig.txt`, pass `-psql` flag to `exatorrent`
