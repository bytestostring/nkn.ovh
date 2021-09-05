<div id="referenceModal" class="modal">
<div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">{{.LANG.Modal.reference.title}}</h3> <a href="#close" title="{{.LANG.Modal.control.close}}" class="close" onclick="closeModal('reference')">X</a></div><div class="modal-body">
{{if or (eq .LANG.CurrentLang "en_US") (eq .LANG.CurrentLang "zn_CN") }}
	<h3>Basic information</h3>
	<p> - The site updates automatically every minute.</p>
	<p> - NKN Network statistics is updated once per 20 minutes.</p>
	<p> - The user nodes statistics is updated once per 2 minutes.</p>
	<p> - User nodes relative performance and rewards are calculated based on the ratio of the user's active nodes performance to the total network one.</p>
	<p><strong> - R/H</strong> – means relays per hour for node uptime period.</p>
	<p><strong> - R/H Nm</strong> – means relays per hour for the latest N minutes of node uptime.</p>
	<p><strong>-  N/A</strong> – means no information is available at this moment.</p>
	<p><br></p>
	<h3>Node status:</h3>
		<p><strong> - Out of NKN Network</strong> means node is disconnected from the main NKN network and operates offline. It can't get rewards from mining and its statistics doesn't affect the summary statistics on the nodes table.</p>
		<p><strong> - PERSIST_FINISHED</strong> means node is synced and mining now.</p>
		<p><strong> - WAIT_FOR_SYNCING</strong> means node is waiting for synchronization.</p>
		<p><strong> - SYNC_STARTED</strong> means node has started synchronizing.</p>
		<p><strong> - SYNC_FINISHED</strong> means node has finished synchronizing.</p>
		<p><strong> - OFFLINE</strong> means node is offline.</p>
{{end}}
{{if eq .LANG.CurrentLang "ru_RU"}}
	<h3>Основная информация</h3>
	<p> - Сайт имеет автообновление полученной статистики с интервалом в 1 минуту.</p>
	<p> - Статистика по общей сети обновляется 1 раз в 20 минут.</p>
	<p> - Статистика по узлам пользователей обновляется каждые 2 минуты.</p>
	<p> - Относительная производительность пользовательских узлов и количества наград рассчитывается,  исходя из суммы релеев активных узлов.</p>
	<p><strong> - R/H</strong> – Количество релеев за час с учетом полного срока аптайма ноды.</p>
	<p><strong> - R/H Nm</strong> – Количество релеев за час. Число N означает, что количество релеев указано за последние N минут работы ноды.</p>
	<p><strong>-  N/A</strong> – Информация недоступна на текущий момент времени.</p>
	<p><br></p>
	<h3>Статусы узла:</h3>
		<p><strong> - Out of NKN Network</strong> - узел отключен от основной сети NKN и осуществляет работу автономно.
Узел с таким статусом не может получать вознаграждения за майнинг. Статистика узла не суммируется к общей статистике в таблице узлов.</p>
		<p><strong> - PERSIST_FINISHED</strong> - узел работает и осуществляет майнинг.</p>
		<p><strong> - WAIT_FOR_SYNCING</strong> - узел ожидает начала процесса синхронизации.</p>
		<p><strong> - SYNC_STARTED</strong> - узел находится в процессе синхронизации.</p>
		<p><strong> - SYNC_FINISHED</strong> - узел закончил процесс синхронизации.</p>
		<p><strong> - OFFLINE</strong> - узел не работает (отключен от сети Интернет).</p>
{{end}}
</div>
</div>
</div>
</div>
