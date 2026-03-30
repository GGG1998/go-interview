1. Jakie są przyczyny większej latency
- Wysokie użycie CPU (High CPU Utilization): Gdy system jest przeciążony, scheduler Go (planista) ma trudności z przypisaniem wątków OS do gorutyn, co powoduje kolejkowanie zadań.
- Zbyt wiele gorutyn (Kolejki lokalne/globalne): Tworzenie milionów gorutyn jednocześnie zwiększa konkurencję o zasoby. Jeśli kolejki lokalne (Local Run Queues) są pełne, gorutyny muszą czekać w wolniejszej, globalnej kolejce.
- Blokujące operacje (Blocking Operations): Częste blokowanie gorutyn przez operacje I/O, sieciowe lub długie obliczenia bez udostępniania CPU innym (brak preemption).
- Garbage Collection (GC): Przy ogromnej liczbie gorutyn (zwłaszcza z dużymi stosami), Garbage Collector może działać dłużej, wstrzymując na chwilę wykonywanie zadań (tzw. "stop-the-world" lub duże zużycie CPU przez GC).
- Nadmierne przełączanie kontekstu (Context Switching): Mimo że gorutyny są tanie, przełączanie między nimi (zwłaszcza przy ich ogromnej liczbie) generuje narzut czasowy. 

2. Skutki większej latencji
- Wzrost P99/P95 (Tail Latency): Średnia latencja może wyglądać dobrze, ale użytkownicy odczują "zacięcia" w najwolniejszych 1% lub 5% przypadków (tzw. tail latency), co jest krytyczne w mikrousługach.
- Spadek przepustowości (Throughput): Aplikacja przetwarza mniej żądań na sekundę, ponieważ rdzenie procesora są zajęte zarządzaniem kolejkami zamiast faktyczną pracą.
- Niestabilne działanie: System staje się nieprzewidywalny, a czas odpowiedzi na żądania rośnie w sposób niekontrolowany. 

3. Jak to naprawić?
- Optymalizacja kodu: Zmniejszenie liczby blokad (mutexów) i lepsze zarządzanie kanałami.
- Kontrola liczby gorutyn: Użycie puli robotników (worker pool) zamiast tworzenia nowej gorutyny dla każdego najmniejszego zadania.
- Debugowanie: Użycie narzędzi takich jak go tool pprof do analizy blokad (block profile) i trace do sprawdzenia działania schedulera.