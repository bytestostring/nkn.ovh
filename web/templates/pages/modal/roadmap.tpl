<div id="roadModal" class="modal">
<div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">{{LANG.modal.roadmap.title}}</h3> <a href="#close" title="{{LANG.modal.close}}" class="close" onclick="$('#roadModal').hide(200)">X</a></div><div class="modal-body">
{{#if_eq LANG.current_lang "ru_RU"}}
	<h4>Май - Июнь 2021</h4>
	<p>0. Добавить throttling по IP-адресу во избежание атак (завершено).<p>
	<p>1. Добавить возможность смены языка. English/Russian (завершено)</p>
	<p>2. Добавить более информативные уведомления о узлах в статусе Offline, добавить звуковое оповещение.</p>
	<p>3. Добавить возможность предоставления read-only доступа по токену (по ссылке).</p>
	<p>4. Добавить возможность устанавливать пароль на hash-link.</p>
	<p>5. Улучшить демон, сократить интервал опроса для узлов, которые не ответили в заданное время, для более корректного отображения информации. (завершено)</p>
	<p>6. Добавить постраничную навигацию по таблице узлов</p>
{{/if_eq}}
{{#if_eq LANG.current_lang "en_US"}}
	<h4>May - June 2021</h4>
	<p>0. Add delay and limits on the amount of created hashlinks by IP. (completed)</p>
	<p>1. Add English language support with a fast switch between languages. (completed)</p>
	<p>2. Add a feature: improved notification system (with sound notifications and support of messengers)</p>
	<p>3. Add a feature: Read-only access by tokens which can be provided by the user.</p>
	<p>4. Add a feature: custom password for hashlink (ID) at the user's discretion (optional).</p>
	<p>5. Improve golang daemon (backend): reduce the polling interval for nodes that have not responded on schedule, and add repetitive polling for offline nodes (completed)</p>
	<p>6. Add a feature: page navigation for the user nodes table.</p>
{{/if_eq}}</div>
</div>
</div>
</div>
