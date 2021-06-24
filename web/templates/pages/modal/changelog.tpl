<div id="changelogModal" class="modal">
<div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">{{LANG.modal.changelog.title}}</h3> <a href="#close" title="{{LANG.modal.close}}" class="close" onclick="$('#changelogModal').hide(200)">X</a></div><div class="modal-body">
{{#if_eq LANG.current_lang "ru_RU"}}
	<h4>21 Мая - 24 Июня 2021 (Версия 1.0.10)</h4>
	<p>- Множество важных изменений во внутренней логике работы.</p>
	<p>- Открыт исходный код проекта: <a href="https://github.com/bytestostring/nkn.ovh" target="_blank">github.com/bytestostring/nkn.ovh</a></p>
	<h4>20 Мая 2021</h4>
	<p>- Добавлено автоматическое удаление неактивных узлов в статусе "Offline". Если узел находится более семи дней в статусе "Offline", он будет удален автоматически. Если узел ни разу не ответил корректно с момента добавления, он будет удален через 6 часов после добавления.</p>
	<p>- [Backend] Множество изменений во внутренней логике опроса узлов.</p>
	<p>- [Frontend] Исправлены ошибки при сортировке по столбцу в таблице узлов	.</p>
	<p>- [Frontend] Добавлено блокирующее уведомление, которое становится активным, если версия страницы не соответствует версии сайта.</p>
	<p>- Время обхода пользовательских узлов уменьшено с 10 до 2 минут.</p>
	<h4>14 Мая 2021</h4>
	<p>- [Улучшение] Добавлен мониторинг кошельков.</p>
	<p>- Незначительные изменения, исправлено несколько багов.</p>
	<h4>7 Мая 2021</h4>
	<p>- Исправлен баг из-за которого статус узлов, которые не в сети, отображался некорректно в таблице узлов.</p>
	<h4>6 Мая 2021</h4>
	<p>- [Backend] Исправлен баг с гороутинами, Golang daemon работал некорректно после предыдущего обновления.</p>
	<p>- [Backend] Незначительные изменения.</p>
	<p>- [Backend] Добавлен повторный опрос узлов, которые не ответили в заданное время, для уменьшения количества ошибок связанных с несвоевременными ответами.</p>
	<p>- Изменена переменная blocks_per_day, теперь она принимает динамическое значение среднего количества блоков в день с начала запуска основной сети блокчейна и влияет на некоторые расчеты.</p>
	<p>- Добавлена интерпретация даты и времени в часовом поясе пользователя.</p>
	<p>- Исправлен баг с расчетом среднего количества релеев для узлов пользователя.</p>
	<h4>4 Мая 2021</h4>
	<p>- Добавлена поддержка интерфейса на английском, добавлено быстрое переключение языков</p>
	<h4>2 Мая 2021</h4>
	<p>- Добавлен лимит: максимальное количество узлов на ID NKNC: 5000</p>
	<p>- Исправлен расчет 1/6/24-часовых релеев узлов</p>
	<p>- Исправлена переменная в расчете доходности узлов (variable blocks_per_day reduced from 4320 to 3700)</p>
	<p>- Добавлен throttling по IP адресу на генерацию ID.</p>
	<p>- Незначительные исправления в NKNC Daemon (Golang).</p>
{{/if_eq}}
{{#if_eq LANG.current_lang "en_US"}}
	<h4>21 May - 24 June 2021 (Version 1.0.10)</h4>
	<p>- Many important changes in internal logic of the site.</p>
	<p>- The project is open source: <a href="https://github.com/bytestostring/nkn.ovh" target="_blank">github.com/bytestostring/nkn.ovh</a></p>
	<h4>20 May 2021</h4>
	<p>- Added automatic removal of inactive nodes. If a node has been Offline for more than seven days, it will be removed. If a node hasn't responded correctly to the server within 6 hours since the node was added, it will also be removed.</p>
	<p>- [Backend] Many changes in the internal logic of user nodes' polling.</p>
	<p>- [Frontend] Fixed bugs that were causing issues with sorting in the nodes table.</p>
	<p>- [Frontend] Added a blocking notification that appears when the page's version on the user side does not match the latest site version.</p>
	<p>- The polling time of user nodes was reduced from 10 to 2 minutes</p>

	<h4>14 May 2021</h4>
	<p>- [Improvement] Added a wallet tracker.</p>
	<p>- Minor changes, fixed some bugs.</p>
	<h4>7 May 2021</h4>
	<p>- Fixed a bug causing wrong status for offline nodes in the user nodes table.</p>
	<h4>6 May 2021</h4>
	<p>- [Backend] Fixed a bug related to goroutines: golang daemon was working incorrectly after the previous update.</p>
	<p>- [Backend] Minor changes.</p>
	<p>- [Backend] Added repetitive polling for nodes that have not responded on schedule, it reduces the number of incorrect nodes status after the first poll.</p>
	<p>- Changed the variable "blocks_per_day": now it has a dynamic value of the average amount of blocks per day since the Mainnet launch.</p>
	<p>- Added representation of dates and times in user's local timezone.</p>
	<p>- Fixed a bug in the calculation of the average amount of relays for user nodes.</p>
	<h4>4 May 2021</h4>
	<p>- Added English language support and fast switch between languages.</p>
	<h4>2 May 2021</h4>
	<p>- New limit: Maximum amount of nodes per NKNC ID is now 5000</p>
	<p>- Fixed the calculation for 1/6/24-hour data of user nodes.</p>
	<p>- Fixed a variable: the variable "blocks_per_day" was reduced from 4320 to 3700.</p>
	<p>- Added a feature: throttling/delay for ID generation by IP.</p>
	<p>- Minor changes in NKNC Daemon (Golang).</p>
{{/if_eq}}
</div>
</div>
</div>
</div>
