Тестовое задание в компанию MEDODS.<br/>

Описание:<br/>
1)Структура проекта:<br/>
В проекте находятся следующие папки: cmd, config, elements, handlers, pkg.<br/>
cmd - содержит исполняемый файл main.go.<br/>
config - содержит config.yml с данными для работы программы.<br/>
elements - содержит элемент сервер.<br/>
handlers - содержит файл с путями по которым будут идти запросы :/get_token для получения токенов и /refresh_token для их обновления.<br/>
в pkg находятся следующие папки repository,service,tests,tokens.<br/>
В repisitory хранятся скрипты отвечающие за взаимодействие с БД.<br/>
В service вызываются промежуточные функции и функции взаимодействия с БД.<br/>
в tests лежат тесты.<br/>
в tokens функции по работе с токенами.<br/>
Логика приложения следующая: запрос -> handler -> service -> repository.<br/>

2)Структура БД:<br/>
В БД 2 таблицы users и refresh_tokens.<br/>
В users хранится id, guid, email.<br/>
В refresh_tokens храниться user_id, refresh_token,available,create_at.<br/>
Связаны таблицы по id -> user_id.<br/>

Краткое описание работы программы:<br/>
Токены связаны по uuid который генерируется в процессе создания токенов.<br/>
Во время refresh token помимо проверки на валидность токена идёт сравнение token_guid из access_token и refresh_token (для того чтобы понимать что токены, были выпущены вместе).<br/>
Также во время refresh token мы смотрим в БД доступен ли данный refresh_token (На случай если надо будет отозвать токены).<br/>
Если пользователь отправил разные ip, то в ответе json в поле email вернётся почта пользователя (На которую надо будет отправить warning)<br/>

!!Как работать с программой<br/>
В БД находитсья пользователь с guid: b758e00b-9c42-4f2b-a84d-0af47b937d17. Поэтому для успешного запроса должно использоваться это guid<br/>
При запросе надо передаваться ip пользователя<br/>

Формат запроса на /get_token<br/>
{<br/>
    "guid":"b758e00b-9c42-4f2b-a84d-0af47b937d17",<br/>
    "ip":"1.2.3"<br/>
}<br/>

Формат запроса на /refresh_token
{
    "access_token":"",
    "refresh_token":"",
    "ip":"1.2.4"
}

Docker:<br/>
Чтобы использовать docker нужно прописать docker-compose up находясь в папке medods. !!!Так как docker не видит postgre на localhost необходимо в config.yml который лежит в папке config<br/>
Поменять db host на postgres (надо раскомментить строку, и закоментить с localhost). По умолчание в config подключение к ДБ идет по localhost.<br/>

Tests:<br/>
Тесты лежат в паке tests. в ней два файла в all_test.go запускаются тесты, в test_table.go лежать сами тесты.<br/>
