<div id="donateModal" class="modal">
<div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">{{.LANG.Modal.donate.title}}</h3> <a href="#close" title="{{.LANG.Modal.control.close}}" class="close" onclick="closeModal('donate')">X</a></div><div class="modal-body">
{{if or (eq .LANG.CurrentLang "en_US") (eq .LANG.CurrentLang "zn_CN") }}
	<p>If you wanna help this project:</p>
{{end}}
{{if eq .LANG.CurrentLang "ru_RU"}}
	<p>Если Вы хотите помочь проекту:</p>
{{end}}
	<p>NKN Mainnet address: <a href="https://explorer.nkn.org/detail/address/NKNZKKF9u1MUQWnK272YoFiMTn5tjZh7uRQE/1" rel="noreferrer" target="_blank">NKNZKKF9u1MUQWnK272YoFiMTn5tjZh7uRQE</a></p>
</div>
</div>
</div>
</div>
