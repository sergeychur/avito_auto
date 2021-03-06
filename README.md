# avito_auto
Тестовое задание на стажировку в Авито
## Реализованная функциональность
* Возможность создавать короткие ссылки.
* Возможность перейти по сохраненному ранее короткому представлению с редиректом на исходный URL.
* Валидация ссылок (проверка URL на корректность и выполнение запроса на URL для проверки существования такого сайта в принципе).
* Возможность задавать кастомные ссылки.

## Выполненные доп. задания, не являющиеся функциональностями
* Проведено нагрузочное тестирование.
* Написаны unit-тесты на хэндлеры, мидлвары и модуль, связанный с конфигом. Покрытие тестами всего проекта - ~40%.
Не тестировались методы, связанные с БД; функции, вызывывающие библиотечные функции; и функции, которые создают
объекты(сервера и т.п.).
 

## Необходимое для запуска ПО
1. Непосредственно сервиса:
    * Docker;
    * docker-compose.
2. Простого UI:
    * npm.
3. Тестов:
    * go 1.15
    
## Инструкция по запуску
1. git clone https://github.com/sergeychur/avito_auto.git
2. Непосредственно сервиса:
    ```
   sudo docker-compose up --build
   ```
3. Простого UI:
    ```
   cd test_ui && npx http-server
   ``` 
4. Тестов:
    ```
    go test -v ./...
    ```

## Инструкция по использованию
После запуска есть два варианта развития событий.
Первый - пользоваться сервисом исключительно как API, второй - пользоваться тестовым UI.
В первом случае пригодится документация к API. Ее можно получить перейдя по следующему адресу:
```
http://localhost:8091/doc/doc.html
```
Во втором случае необходимо лишь открыть в браузере
```
http://localhost:8080/
```
Если у Вас занят порт 8080, есть вероятность, что ```http-server``` может 
запуститься на другом порте. В таком случае необходимо во-первых поменять порт в адресе,
во-вторых добавить в файл config_deploy.json в массив allowed_hosts элемент вида ```http://localhost:<port>/```,
где ```<port>``` - порт, на котором запустился ```http-server```.
На бесплатный хостинг UI я не разместил, так как у меня закончилось бесплатное время на Digital Ocean:)

## Нагрузочное тестирование
Перед тестированием база данных была заполена тестовыми данными. Было сохранено 10000 коротких представлений урлов.

Само нагрузочное тестирование было проведено с помощью утилиты ```ab```.
Нагрузка шла в 10 потоков. Было выполнено 10000 запросов на редирект по полному адресу, соответствующему короткому представлению.
Тесты проводились на 64-битной операционной системе Ubuntu 20.04 LTS на компьютере с 8 Гб оперативной памяти, 
процессором Intel Core i5-7300Q с 4 ядрами и тактовой частотой 2.5 ГГц.
Результат тестирования представлен ниже.

```
Server Software:        
Server Hostname:        localhost
Server Port:            8091

Document Path:          /api/link/google
Document Length:        45 bytes

Concurrency Level:      10
Time taken for tests:   2.409 seconds
Complete requests:      10000
Failed requests:        0
Non-2xx responses:      10000
Total transferred:      2120000 bytes
HTML transferred:       450000 bytes
Requests per second:    4151.41 [#/sec] (mean)
Time per request:       2.409 [ms] (mean)
Time per request:       0.241 [ms] (mean, across all concurrent requests)
Transfer rate:          859.47 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.7      0      13
Processing:     0    2   6.8      1     211
Waiting:        0    2   6.8      1     210
Total:          0    2   6.8      2     211

Percentage of the requests served within a certain time (ms)
  50%      2
  66%      2
  75%      2
  80%      3
  90%      4
  95%      5
  98%      9
  99%     12
 100%    211 (longest request)
```

Все запросы возвращают не 2хх так как выполняется редирект, и код - 303.

Средний RPS - 4151 запросов в секунду.
Среднее время, потраченное на 1 запрос - 2.409 миллисекунд.
