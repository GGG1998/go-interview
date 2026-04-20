# Dictionary — Key Concepts

| Term | Definition |
|------|-----------|
| **SSTable** (Sorted String Table) | A segment file where key-value pairs are sorted by key and each key appears at most once per file. |
| **LSM Tree** (Log-Structured Merge-Tree) | A storage engine architecture combining a memtable, SSTable segments, and compaction. Optimized for high write throughput. |
| **Memtable** | An in-memory balanced tree (e.g. red-black tree) that buffers writes before they are flushed to disk as an SSTable. |
| **Sparse index** | An index that stores only one entry per block/range of keys (not every key). Works because sorted data allows scanning forward from a known position. |
| **Compaction** | The background process of merging SSTable segments, discarding overwritten and deleted values, and keeping only the most recent value per key. |
| **Size-tiered compaction** | A compaction strategy where newer, smaller SSTables are merged into older, larger ones over time. |
| **Leveled compaction** | A compaction strategy where keys are split across levels; older data moves to lower levels, enabling incremental compaction with lower memory usage. |
| **Bloom filter** | A memory-efficient probabilistic data structure that can definitively say a key does *not* exist in a dataset, avoiding unnecessary disk reads for missing keys. |
| **B-tree** | A balanced tree data structure that stores sorted key-value pairs in fixed-size pages on disk. Standard in relational databases (PostgreSQL, MySQL, etc.). |
| **Page** | A fixed-size unit of storage in a B-tree (typically 4 KB). Reads and writes operate one page at a time. |
| **Branching factor** | The number of child page references on a single B-tree page. A higher branching factor means a shorter tree and fewer disk reads per lookup. |
| **Leaf page** | The bottom-level pages in a B-tree that contain actual values (or references to them). |
| **Page split** | When a B-tree page is full and a new key must be inserted, the page is split into two half-full pages and the parent is updated. |
| **Write-Ahead Log (WAL)** | An append-only file where every B-tree modification is recorded before being applied. Used to restore a consistent state after a crash. |
| **Merge-sort** | The algorithm used when merging SSTable segments — read one entry at a time from each input file, always pick the smallest key, write it to the output. |
| **Red-black tree / AVL tree** | Self-balancing binary search trees used to implement the memtable. Guarantee O(log n) inserts and lookups while maintaining sorted order. |
