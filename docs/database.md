To make things persistent between sessions, `exatorrent` uses database. But because `exatorrent` uses awesome `sqlite3` which comes embedded within `exatorrent`, you barely notice it.


Note that once you start using one Database, you must stick to it. You cannot jump between Database Implementations and if you do, you loose data.

### Data Stored in Database

1. Users and their data
2. Trackers
3. State of Torrent
4. Piece Completion State of Torrent
5. State of Files of Torrent
6. Lock State of Torrent
