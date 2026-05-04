# Key-Value Storage — pytania sprawdzające

## Log (append-only)

1. Dlaczego bazy danych używają logu typu append-only zamiast nadpisywać istniejące rekordy? Podaj trzy powody.
2. Dlaczego zapis sekwencyjny jest szybszy od losowego — i czy ta różnica istnieje tylko na dyskach magnetycznych?
3. Jak niezmienność (immutability) plików segmentów upraszcza odzyskiwanie po awarii?

## Indeksy

4. Czym jest indeks i dlaczego potrzebujesz czasem kilku indeksów dla tych samych danych?
5. Jaki jest główny koszt utrzymywania indeksu? Dlaczego?

## Hash Index / Bitcask

6. Opisz najprostszą strategię indeksowania w stylu Bitcask — co zawiera hash table, a co plik z danymi?
7. Do jakiego typu obciążenia (workload) Bitcask jest szczególnie dobrze dopasowany?
8. Wymień dwa główne ograniczenia indeksu hashowego.
9. Dlaczego zapytania zakresowe (range queries, np. `user:000`–`user:999`) są nieefektywne w hash indexie?

## Segmenty i kompakcja

10. Po co dzielić log na segmenty zamiast trzymać jeden ogromny plik?
11. Co dokładnie robi proces kompakcji i jaki jest jego efekt dla liczby segmentów oraz wydajności odczytu?
12. W jakiej kolejności sprawdzasz segmenty przy lookupie klucza i dlaczego właśnie tak?

## Format pliku

13. Dlaczego CSV nie jest optymalnym formatem dla logu? Co stosuje się zamiast tego i dlaczego nie wymaga escape'owania?

## Szczegóły implementacyjne

14. Czym jest tombstone i dlaczego nie wystarczy „po prostu usunąć” rekord z pliku?
15. Co Bitcask robi, żeby przyspieszyć recovery po restarcie? Co by się stało bez tej optymalizacji?
16. W jaki sposób Bitcask wykrywa rekordy zapisane częściowo (np. po crashu w trakcie zapisu)?
17. Dlaczego typowa implementacja używa jednego wątku piszącego, ale odczyty mogą być współbieżne?

## Część zaawansowana (OTHERS)

18. Skoro w segmencie typu dziennik kolejność par klucz-wartość zapisanych w czasie ma znaczenie (późniejsze nadpisuje wcześniejsze), to dlaczego mówi się jednocześnie, że „kolejność par klucz-wartość w pliku nie ma znaczenia”? Jak to pogodzić?
19. Jeśli wymusimy, by pary w segmencie były posortowane wg klucza (SSTable), to jakie problemy powstają z sekwencyjnym zapisem? Jak intuicyjnie można to obejść? (otwarte — pomyśl)
20. Jakie zalety daje sortowanie kluczy w segmencie w porównaniu z niesortowanym logiem Bitcaska? (pomyśl o ograniczeniach z pkt 8–9)
