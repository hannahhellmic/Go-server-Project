Для начала
```
go build -o client
```

1. Создание аккаунта
```
-port 1323 -cmd create -name /enter name/ -amount /enter amount/
```
2. Получение данных об аккауте
```
-port 1323 -cmd get -name /enter name/
```
3. Удаление аккаунта
```
-port 1323 -cmd delete -name /enter name/
```
4. Изменение баланса
```
-port 1323 -cmd patch -name /enter name/ -sum_change /enter balance change/
```
5. Изменение имени
```
-port 1323 -cmd change -name /enter name/ -new_name /enter new name
```
6. Перевод средств
```
-port 1323 -cmd transfer -name_from /enter account name to transfer from/ -name_to /enter account name to transfer to/ -amount /enter amount/
```
7. Вывести список аккаунтов
```
-port 1323 -cmd get_all
```
8. Список операций
```
-port 1323 -cmd transactions -name /enter name/
```
