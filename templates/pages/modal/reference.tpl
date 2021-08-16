<div id="referenceModal" class="modal">
<div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">{{.LANG.Modal.reference.title}}</h3> <a href="#close" title="{{.LANG.Modal.control.close}}" class="close" onclick="closeModal('reference')">X</a></div><div class="modal-body">
{{if eq .LANG.CurrentLang "en_US"}}
	<p> - The site updates automatically every minute.</p>
	<p> - NKN Network statistics is updated once per 50 minutes</p>
	<p> - The user nodes statistics is updated once per 2 minutes</p>
	<p> - User nodes relative performance and rewards are calculated based on the ratio of the user's active nodes performance to the total network one.</p>
	<p><strong> - R/H</strong> – means relays per hour for node uptime period</p>
	<p><strong> - R/H Nm</strong> – means relays per hour for the latest N minutes of node uptime.</p>
	<p><strong>-  N/A</strong> – means no information is available at this moment..</p>
{{end}}
{{if eq .LANG.CurrentLang "ru_RU"}}
	<p> - Сайт имеет автообновление полученной статистики с интервалом в 1 минуту.</p>
	<p> - Статистика по общей сети обновляется 1 раз в 50 минут</p>
	<p> - Статистика по узлам пользователей обновляется каждые 2 минуты</p>
	<p> - Относительная производительность пользовательских узлов и количества наград рассчитывается,  исходя из суммы релеев активных узлов.</p>
	<p><strong> - R/H</strong> – Количество релеев за час с учетом полного срока аптайма ноды</p>
	<p><strong> - R/H Nm</strong> – Количество релеев за час. Число N означает, что количество релеев указано за последние N минут работы ноды.</p>
	<p><strong>-  N/A</strong> – Информация недоступна на текущий момент времени.</p>
{{end}}
</div>
</div>
</div>
</div>
