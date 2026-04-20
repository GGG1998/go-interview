# Dictionary — Key Concepts

| Term | Definition |
|------|-----------|
| **Log** | An append-only sequence of records. New entries are added at the end; existing entries are never modified or deleted in place. |
| **Append-only** | A write strategy where data is only ever added, never overwritten. Enables simpler crash recovery and faster sequential writes. |
| **Index** | Separate metadata structure that maps keys to their location in the data file. Speeds up reads at the cost of slower writes. |
| **Hash table / Hash index** | A data structure mapping each key to a byte offset in the data file. Stored in memory for fast lookups. |
| **Byte offset** | The position (in bytes from the start) of a value inside a file. The hash index stores this offset for each key. |
| **Segment** | A fixed-size chunk of the log. When a segment reaches its size limit, it is closed and a new one is opened. |
| **Compaction** | The process of merging segments and removing duplicate keys, keeping only the most recent value for each key. Reclaims disk space. |
| **Tombstone** | A special deletion marker appended to the log when a key is deleted. Compaction uses it to discard all previous values for that key. |
| **Snapshot** | A point-in-time copy of an in-memory hash table saved to disk. Used by Bitcask to speed up recovery after a crash. |
| **Checksum** | A value computed from data used to detect corruption. Bitcask stores checksums alongside log entries to detect partial writes caused by crashes. |
| **Sequential write** | Writing data one block after another in order. Much faster than random writes, especially on HDDs. |
| **Random write** | Writing data to arbitrary positions on disk. Slower due to seek time on HDDs and write amplification on SSDs. |
| **Bitcask** | A storage engine (used in Riak) that implements a hash index over an append-only log. Optimized for high write throughput with frequent updates. |
| **Concurrency control** | Mechanisms to ensure safe access to data by multiple threads or processes simultaneously. |
| **Crash recovery** | The process of restoring a consistent database state after an unexpected failure or restart. |
