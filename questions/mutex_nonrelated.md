# RWMutex Struct

```
type RWMutex struct {
    w           Mutex  // Mutex for writer
    writerSem   uint32 // Semafor for writer (waiting for reader)
    readerSem   uint32 // Semafor for reader (waiting for writer)
    readerCount int32  // Number of active reader
    readerWait  int32  // Number of reader, who must finish because writer has priority
}
```

┌────────────────┬─────────────────────────────────────┬──────────────────────────────────────────────────────────────────────────────────┐
│     State      │    Counter Status (readerCount)     │                                   Description                                    │
├────────────────┼─────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────┤             
│ Idle           │ readerCount == 0                    │ No locks are held. Anyone can enter.                                             │                    
├────────────────┼─────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────┤                                                                                                
│ Active Readers │ readerCount > 0                     │ Multiple readers are active. No writer is waiting yet.                           │
├────────────────┼─────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────┤                                                                                                
│ Writer Waiting │ readerCount < 0 (e.g., -1073741821) │ A writer called Lock(). It "flipped" the count to negative to block new readers. │
├────────────────┼─────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────┤                                                                                                
│ Writer Active  │ readerCount < 0 AND readerWait == 0 │ All previous readers have left. The writer now has exclusive access.             │
└────────────────┴─────────────────────────────────────┴──────────────────────────────────────────────────────────────────────────────────┘

