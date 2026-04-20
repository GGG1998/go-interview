# SSTable & LSM Trees — Notes

## From Hash Index to Sorted Segments

The previous key-value approach stored segments in arbitrary order and kept a full in-memory hash table. A simple change unlocks much better performance: **require that key-value pairs in each segment file are sorted by key**. This format is called an **SSTable (Sorted String Table)**.

Each key must appear exactly once per merged segment file — compaction already guarantees this.

---

## Merging Segments (Merge-Sort)

Merging is efficient even when files are larger than available memory. The algorithm mirrors merge-sort:

1. Open all input segment files simultaneously.
2. Read the first key from each file.
3. Write the smallest key to the output file.
4. Advance that file's pointer and repeat.

**Duplicate keys:** Each segment covers writes from a specific time window, so a newer segment always wins. Keep the newest value, discard older ones.

> **How it works at scale:** Even with 100 large segments (potentially gigabytes each), the process reads one entry at a time from each file — it never loads whole files into memory. A sparse index tells the process where to start reading in each file. If `Gabriel Domanowski` exists in 10 files, all 10 are compared simultaneously; the newest value wins and the rest are discarded.

---

## Sparse Index

With sorted data, you no longer need an in-memory entry for every key. A **sparse index** is sufficient — one entry per few kilobytes. To find a key:

1. Find the nearest known key in the sparse index.
2. Scan forward from that offset until the key is found (or confirmed absent).

Records can be **grouped into compressed blocks** before writing to disk. Each sparse index entry points to the start of a compressed block. This reduces both disk usage and I/O operations.

---

## Memtable — Write Buffer

Maintaining sorted order on disk is hard. Instead, writes accumulate in an **in-memory balanced tree** (red-black tree or AVL tree), called a **memtable**.

**Write path:**
1. Insert the key-value pair into the memtable.
2. When the memtable exceeds a threshold (typically a few MB), flush it to disk as a new SSTable segment. New writes go to a fresh memtable while flushing.

**Read path:**
1. Check the memtable.
2. Check the newest SSTable segment on disk.
3. Continue through older segments until found.

**Crash problem:** The memtable lives in memory — a crash loses it. Solution: maintain a separate append-only log on disk for crash recovery only. Discard it once the memtable is successfully flushed to an SSTable.

**Background compaction:** Periodically merge segments and remove overwritten or deleted values.

---

## LSM Tree

The full architecture — memtable + SSTable segments + compaction — is called an **LSM tree (Log-Structured Merge-Tree)**. Storage engines using this design are called **LSM storage engines** (e.g. LevelDB, RocksDB, Cassandra).

**Advantages:**
- Very high write throughput — all disk writes are sequential.
- Efficient range queries — data is sorted.

---

## B-Trees

B-trees also store key-value pairs sorted by key, enabling efficient lookups and range queries — but the design is fundamentally different.

Instead of variable-size segments, B-trees divide the database into **fixed-size pages** (typically 4 KB). Reads and writes operate one page at a time, mirroring how disk hardware works.

### Tree structure
- **Root page** — entry point for every lookup; contains keys and references to child pages.
- **Internal pages** — cover a key range; keys between references define sub-range boundaries.
- **Leaf pages** — contain the actual values (or references to them).

**Branching factor** — the number of child references per page. Higher branching = shorter tree = fewer disk reads.

### Why tree depth matters
The main database bottleneck is disk reads. A binary tree might need 30–40 reads to find one record. A B-tree is **short and wide** — typically only 3–4 levels deep, enough to index 256 TB of data.

### Operations
- **Lookup:** Root → internal pages → leaf page.
- **Update:** Find leaf, modify value, write page back to disk.
- **Insert:** Find the right page, add the key. If the page is full, **split** it into two half-full pages and update the parent.

### Reliability
Overwriting a page is a hardware-level operation. A crash mid-write (e.g. during a split) can corrupt the index.

**Solution: Write-Ahead Log (WAL)** — every modification is recorded in an append-only log before being applied to the tree. Used to recover a consistent state after a crash.
