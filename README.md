# Lindexdatacenter_test
тестовое задание для Lindexdatacenter



Стек технологий
```GoLang, docker, github```

**Задача:**

Есть 2 файла с данными о продуктах (наименование, цена, рейтинг) в 2-х форматах - CSV и JSON.
Необходимо написать программу, которая считывает данные из переданного в параметре файле, и
выводит «самый дорогой продукт» и «с самым высоким рейтингом».
Предусмотреть, что файлы могут быть огромными.
Репозиторий должен содержать Dockerfile для сборки готового приложения в docker среде.
Репозиторий необходимо выложить на github и предоставить ссылку.

**Решение:**

Без учета возможной подачи огромных файлов JSON парсер мог бы выглядеть вот так:
```
func parseJSON(file *os.File) {
	var (
		maxPrice, maxRating         float64
		maxPriceName, maxRatingName string
		fileNotEmpty                bool
	)

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	res := make([]Product, 0, 100)
	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Fatalln(err)
	}

	if len(res) > 0 {
		fileNotEmpty = true
	}
	for _, product := range res {
		if product.Price > maxPrice {
			maxPrice = product.Price
			maxPriceName = product.Name
		}
		if product.Rating > maxRating {
			maxRating = product.Rating
			maxRatingName = product.Name
		}
	}
	printResult(maxPrice, maxRating, maxPriceName, maxRatingName, fileNotEmpty)
}
```
Но поскольку в условии указано, что нужно предусмотреть переполнение памяти я использовал библиотеку jstream.
Jstream - это потоковый синтаксический анализатор, с ограничением размера буфера 4 кб, благодаря которому можно обрабатывать объекты без полного прочтения файла.

Для проверки производительности я использовал бенчмарки из стандартной библиотеки.
Результаты бенчмарка при обработке JSON документа на 100 объектов:
```
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
BenchmarkOld-4              7372            179065 ns/op           40504 B/op        304 allocs/op
BenchmarkNew-4              5588            184312 ns/op           44732 B/op       1392 allocs/op
PASS
ok      command-line-arguments  2.810s
```
И по времени работы, и по количеству затрачиваемой памяти новый парсер не сильно отстает от стандартного анмаршаллера, маленькие файлы будут обрабатываться сравнимо по скорости.
При обработке больших файлов (ниже результаты при обработке файла на ~18000 объектов) jstream начинает обгонять по скорости стандартную библиотеку в 2 раза и по использованию памяти в 1.5:
```
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
BenchmarkOld-4                13         126442080 ns/op        10852397 B/op      52520 allocs/op
BenchmarkNew-4                18          61129862 ns/op         6019573 B/op     253732 allocs/op
PASS
ok      command-line-arguments  6.146s
```

Для CSV файла достаточно построчно считывать данные.
