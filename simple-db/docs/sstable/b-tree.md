# B-Trees

## Overview

B-trees, like SSTables, store key-value pairs **sorted by key**, enabling efficient lookups and range queries. But the underlying design is fundamentally different from log-structured indexes.

---

## Pages (Blocks)

B-trees divide the database into **fixed-size pages** (traditionally 4 KB, sometimes larger). Reads and writes operate one page at a time — this maps directly to how disk hardware works.

Each page can reference other pages (like pointers, but on disk rather than in memory). These references form a **tree of pages**.

---

## Tree Structure

- **Root page** — the entry point for every lookup. Contains a set of keys and references to child pages.
- **Internal pages** — each covers a range of keys; keys between references define the boundaries of sub-ranges.
- **Leaf pages** — contain the actual values (or references to where values are stored).

**Branching factor** — the number of child page references on a single page. A high branching factor keeps the tree short.

### Why this matters for databases

The main bottleneck in databases is disk reads. A "tall and thin" tree (like a binary tree) might require 30–40 disk reads to find one record. A B-tree is **short and wide** — because each page holds many keys, the tree depth stays low (typically 3–4 levels), which is enough to index 256 TB of data.

---

## Operations

**Lookup:** Start at the root, follow references through internal pages, arrive at the correct leaf page.

**Update:** Find the leaf page containing the key, modify the value in place, write the page back to disk. All references to the page remain valid.

**Insert:** Find the page whose range covers the new key, add the key to it. If the page is full, **split** it into two half-full pages and update the parent page to reflect the new key ranges.

---

## Reliability

Overwriting a page is a hardware-level operation:
- **HDD:** Move the read head, wait for the disk to rotate to the correct sector, overwrite.
- **SSD:** More complex — requires erasing and rewriting relatively large blocks.

If a crash happens mid-write (e.g. during a page split that updates multiple pages), the index can be left in a corrupted state.

**Solution: Write-Ahead Log (WAL)** — an append-only file where every modification is recorded before being applied to the tree. On recovery, the WAL is replayed to restore a consistent state.

**Concurrency:** Multiple threads accessing the B-tree simultaneously require careful locking. LSM trees handle this more simply because segments are immutable.
