1. Dlaczego bazy danych używają logu typu append-only zamiast nadpisywać istniejące rekordy? Podaj trzy powody.

Wydajność zapisu (Szybkość): Tak jak wspomniałeś – zapis dopisujący na końcu jest operacją sekwencyjną. Dzięki temu unikamy kosztownego "szukania" miejsca na dysku (random seek), co sprawia, że baza danych może przyjmować dane niemal z prędkością fizycznych możliwości nośnika.

Spójność i Recovery (Łatwiejsze odzyskiwanie): Jeśli system ulegnie awarii w trakcie zapisu, w modelu append-only po prostu ucinasz uszkodzony koniec pliku. Gdybyś nadpisywał dane (in-place update) i nastąpiłaby awaria, stary rekord mógłby zostać częściowo nadpisany nowym, co doprowadziłoby do utraty obu wersji danych i korupcji pliku.

Współbieżność (Prosta synchronizacja): Ponieważ pliki są niezmienne (tylko dopisujemy), nie musisz martwić się o to, że jeden proces czyta rekord, który w tej samej milisekundzie inny proces modyfikuje w środku pliku. To drastycznie upraszcza implementację blokad (locks).


1a. "Chyba nie do końca wiem co znaczy zapis losowy"

Zapis sekwencyjny to pisanie zdanie po zdaniu, od góry do dołu strony. Twoja ręka płynnie przesuwa się w jednym kierunku.

Zapis losowy (random access) to sytuacja, w której musiałbyś dopisać jedno słowo na 5. stronie, potem jedno zdanie na 120. stronie, a potem poprawić literkę na 40. stronie. Twoja ręka ciągle skacze po całej książce, szukając konkretnego miejsca.

2. 

HDD (Dyski magnetyczne): Tutaj różnica jest gigantyczna. Dysk ma fizyczną głowicę i wirujący talerz. Przy zapisie losowym głowica musi fizycznie przesunąć się nad odpowiednią ścieżkę i poczekać, aż talerz obróci się w odpowiednie miejsce (tzw. seek time). To trwa wieki (w skali komputerowej). Przy sekwencyjnym głowica stoi w miejscu, a dane "przepływają" pod nią z ogromną prędkością.

SSD (Dyski półprzewodnikowe): Tu nie ma ruchomych części, więc mogłoby się wydawać, że to nie ma znaczenia. A jednak ma. Choć SSD radzą sobie z dostępem losowym znacznie lepiej niż HDD, to zapis sekwencyjny wciąż jest szybszy. Wynika to z architektury chipów pamięci Flash – łatwiej i szybciej jest zapisać duży, ciągły blok danych niż wiele małych kawałków rozproszonych po różnych komórkach pamięci.

3.

Jasna granica: Plik, który już został zamknięty (zakończony), jest „święty”. Masz pewność, że dane w nim są kompletne i poprawne.

Tylko ostatni segment jest „niepewny”: Podczas recovery po awarii musisz sprawdzić tylko sam koniec ostatniego, aktywnego segmentu.

Łatwe odrzucanie śmieci: Jeśli na końcu logu znajdziesz niepełny rekord (bo prąd padł w połowie zapisu), po prostu go ucinasz. Reszta danych w tym segmencie i wszystkie poprzednie segmenty są nienaruszone.

4.

Indeks to niezależna struktura metadanych. Zawiera informacje o dostępie interesujących nas danych. Na przykład mamy 1000000000 rekordów, to indeks będzie zawierać offsety częściowe offsety nakierowujące np przez dziel i zwyciężaj. 

5.

Koszt indeksu to budowanie indeksu za każdym razem jak dodamy wartość w miejscu, który będzie wymagał przepisania i zdefiniowania jeszcze raz offsetów/wskaźników na nowo.

Dodatkowo ten zapis to juz nie jest jeden zapis logu, tylko indexu i logu. Świetnie to czyta, ale spowalnia zapis.

Drugą kwestią jest pamięć. Czy to ram w zalezności od architektury lub dysk

6.

W RAM-ie trzymają hash table key value. Key faktycznie nasze np: id, a value offset w logu.
Sam plik zawiera już zawartośc danych

7. Najlepszy do appendowania, oraz modyfikowania tych samych kluczy
Czytanie jest związane z Losowością przydzielanego miejsca w pamięci - 
to jest wada, chcąc odczytać sekwencyjnie userów od 100 do 1000,
będziemy skakać po różnych miejscach pamięci

8. RAM jest ograniczony, w przypadku miliardowych zapisów, możemy ją wyczerpać.
Zapytanie zakresowe to katastrofa

9. W odpowiedzi 7 jest zapis

10. 

Dla optymalizacji/dla zarządzania zasobami:

- wyszukiwanie jest skompresowane
- Stare logi są mergowane i aktualizowane

11.

to proces "wyrzucania śmieci". Ponieważ w logu mamy wiele wersji tego samego klucza (bo go aktualizowaliśmy), kompakcja:

Przegląda stare segmenty.
Zostawia tylko najnowszą wartość dla każdego klucza.
Usuwa wszystkie stare, nieaktualne wpisy.

Efekt: Z kilku dużych segmentów powstaje jeden mniejszy. Wydajność odczytu rośnie, bo system ma mniej plików do przeszukania.

12.

Od góry do dołu. Bo sprawdzamy od najnowszych wartości, może segment będzie go zawierać

14. 

Chodzi o to, ze system jest append-only, zanim byśmy odśmiecili, a system by padł
mógłby ożywić dane

---

Chce zaznaczyć, że LOG to nie jest jeden duży plik tylko kontener, w którym są trzymane informacje o segmentach
przez ich indeksy np w pamięci(Bitcask).

Dlaczego?
- Przy dużym jednym logu, mergowanie duplikatów stałoby się koszmarem, bo przepisywalibyśmy
jeden duży plik: Ogromne I/O, długi czas, ryzyko przerwania procesu
- Jak odtworzyć stan po crashu, jeśli plik zostanie uszkodzony? Mamy brak izolacji awarii
