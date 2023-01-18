# partition_searcher

<h2>Для чего</h2>
<p>
Данный сервис предназначен для многопоточного поиска записей в партициях PostgreSQL 
(также можно шарды для mySQL, но нужно будет скорректировать коннект к базе в settings/database.go).
</p>

<h2>Для кого</h2>
<p>
Не все языки поддерживают многопоточные операции (такие как Python, PHP, Nodejs и т. д.). 
Но зато все эти языки совместимы с GRPC <a href="https://grpc.io/docs/languages/">ссылка на список языков</a>
Для таких языков и предназначен данный сервис, так как поиск по таблицам осуществляется при помощи горутин,
тоесть многопоточно. Сортировка также выполняется в нескольких потоках.
</p>

<h2>Запуск</h2>
<ul>
<li>В config.yml необходимо пробисать коннект к вашим партициям.</li>
<li>
В файле proto/get_records.proto в структуре Record необходимо прописать структуру записи из которых состоит таблица разделённая на партиции.
Необязательно прописывать все поля, достаточно только те, которые вам понадобядся на клиенте.
</li>
<li>Запустить каманду protoc --go_out=./proto --go-grpc_out=./proto proto/get_records.proto из корня.</li>
<li>Запустить GRPC сервер командой <pre>./server/server</pre></li>
</ul>

<h2>Тестирование</h2>
<ul>
<li>
Запустить окружение
<pre>cd example/</pre>
<pre>docker-compose up -d</pre>
</li>
<li>
Запустить миграции (тут будет создана таблица tasks и разделена на партиции и будут добавлены все необходимые триггеры).
<pre>cd example/migrations/</pre>
<pre>./migrations/</pre>
</li>
<li>
Запускаем клиент. Радуемся результату)
<pre>./example/example</pre>
</li>
</ul>