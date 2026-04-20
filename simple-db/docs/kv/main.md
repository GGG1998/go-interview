# Key-Value Storage — Notes

## Log (Append-Only File)

Many databases use an append-only log as their core storage mechanism. A log is a sequence of records to which you can only add new entries — never modify existing ones. It doesn't have to be human-readable; it can be binary and intended only for other programs to read.

### Why append-only instead of overwriting?

- **Sequential writes are fast.** Appending and merging segments are sequential write operations, which are significantly faster than random writes — especially on spinning magnetic disks. SSDs also benefit from sequential writes.
- **Simpler crash recovery.** When segment files are immutable, you don't need to worry about a crash happening mid-update that leaves a file half-old and half-new.
- **Avoids fragmentation.** Merging old segments prevents data files from fragmenting over time.

---

## Indexes

An index is additional metadata stored separately from the data itself. It acts as a signpost that helps locate the data you're looking for.

- If you want to search the same data in several different ways, you may need multiple indexes covering different parts of the data.
- **Every index slows down writes**, because each write must also update the index.

---

## Hash Index (In-Memory Hash Table)

Key-value stores are similar to dictionaries in programming languages, typically implemented as a hash table.

**Simplest indexing strategy:** Keep an in-memory hash table where each key maps to a byte offset in the data file. That offset points to where the value is stored. When you append a new key-value pair, you update the hash table to reflect the new offset.

This is the approach used by **Bitcask**, which is well suited to workloads where values are updated frequently.

### Limitations of hash indexes

- The hash table must fit in memory — not suitable for very large numbers of keys.
- Range queries are inefficient (e.g., scanning all keys from `user:000` to `user:999` requires checking each key individually).

---

## Segments and Compaction

Since data is only appended, disk space would eventually run out. The solution: split the log into **segments** of a fixed size. Once a segment reaches the limit, close it and start writing to a new one.

Old segments can be **compacted** (merged): duplicate keys are removed, keeping only the most recent value. This keeps the number of segments small, so lookups don't need to check many hash tables.

Each segment has its own in-memory hash table. To look up a key, check the newest segment first, then the next, and so on.

---

## File Format

CSV is not optimal for a log. A binary format is faster and simpler: encode the length of a string in bytes, followed by the raw bytes (no escaping needed).

---

## Implementation Details

### Deleting records
To delete a key, append a special **tombstone** record. During compaction, the tombstone signals that all previous values for that key should be discarded.

### Crash recovery
In-memory hash tables are lost on restart. Bitcask speeds up recovery by storing **snapshots** of each segment's hash table on disk, which can be loaded quickly on startup. Without snapshots, recovery would require reading every segment file from start to finish.

### Partially written records
A crash can happen mid-write. Bitcask files include **checksums** to detect and ignore corrupted log entries.

### Concurrency control
Because data is written sequentially to the log, the typical implementation uses a single writer thread. Reads can be concurrent since segment files are immutable.


--- OTHERS

“ każdy przechowujący dane segment o strukturze dziennika jest sekwencją par klucz-wartość. Te pary pojawiają się w kolejności, w jakiej zostały zapisane, a wartości umieszczone dalej w dzienniku są ważniejsze niż wcześniejsze wartości powiązane z tym samym kluczem. Oprócz tego kolejność par klucz-wartość w pliku nie ma znaczenia.”

“Teraz wprowadźmy prostą zmianę w formacie plików z segmentami i dodajmy wymóg, by sekwencja par klucz-wartość była posortowana według klucza. Ten wymóg na pozór uniemożliwia sekwencyjny zapis.”
