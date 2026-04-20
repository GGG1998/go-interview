# LSM Trees (Log-Structured Merge-Tree)

## SSTable — Sorted String Table

An SSTable is a segment file where key-value pairs are **sorted by key**, and each key appears exactly once per merged segment file (compaction guarantees this).

### Merging segments

Merging is efficient even when files are larger than available memory. It uses a **merge-sort approach**:
1. Read all input files simultaneously.
2. Pick the smallest key across all files.
3. Write it to the output file.
4. Repeat until all files are consumed.

The result is a new merged segment file, also sorted by key.

**Duplicate keys across segments:** Each segment covers writes from a specific time period, so a newer segment's value always wins. Keep the value from the newest segment, discard older ones.

> **Digression:** Yes, segments can be gigabytes in size. If the key `Gabriel Domanowski` exists in 10 files, the merge process reads the first key from each file, picks the smallest, writes it out, advances that file's pointer, and repeats. It never loads whole files — just one entry at a time from each. The segment's sparse index tells you where to start reading in each file.

---

## Sparse Index

With sorted data, you no longer need an index entry for every key in memory. A **sparse index** is enough — one entry per few kilobytes. To find a key:
1. Look up the nearest known key in the sparse index.
2. Scan forward from that offset until you find the key (or confirm it's absent).

Records can also be **grouped into compressed blocks**. Each sparse index entry points to the start of a compressed block. This reduces both disk space and I/O operations.

---

## Memtable — In-Memory Write Buffer

Maintaining a sorted structure on disk is hard. Instead, writes go to an in-memory balanced tree (e.g. a red-black tree or AVL tree), called a **memtable**.

**Write path:**
1. Add the key-value pair to the memtable.
2. When the memtable exceeds a threshold (typically a few MB), flush it to disk as a new SSTable segment.
3. While flushing, new writes go to a fresh memtable instance.

**Read path:**
1. Check the memtable first.
2. Then check the newest SSTable segment on disk.
3. Continue through older segments until found.

**Problem:** If the database crashes, the memtable (not yet written to disk) is lost. Solution: keep a separate append-only log on disk solely for crash recovery — discard it once the memtable is successfully flushed.

**Background:** Periodically run a compaction process to merge segments and remove overwritten or deleted values.

---

## LSM Tree

The combination of memtable + SSTable segments + compaction is known as an **LSM tree (Log-Structured Merge-Tree)**.

Storage engines using this approach are sometimes called **LSM storage engines** (e.g. LevelDB, RocksDB, Cassandra, HBase).

### Advantages
- **Very high write throughput** — all disk writes are sequential.
- **Efficient range queries** — data is sorted, so scanning a key range is straightforward.

---

## Performance Optimizations

### Bloom filter
Looking up a key that doesn't exist requires checking the memtable and every segment on disk — expensive. A **Bloom filter** is a memory-efficient probabilistic data structure that can definitively say a key does *not* exist, avoiding unnecessary disk reads.

### Compaction strategies

| Strategy | Description |
|----------|-------------|
| **Size-tiered** | Newer, smaller SSTables are progressively merged into older, larger ones. |
| **Leveled** | Keys are split into smaller SSTables across separate "levels". Older data moves to lower levels. Allows incremental compaction with lower memory usage. |
