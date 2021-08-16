<div id="settingsModal" class="modal">
<div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">{{.LANG.Modal.settings.title}}</h3> <a href="#close" title="{{.LANG.Modal.control.close}}" class="close" onclick="closeModal('settings')">X</a></div><div class="modal-body">
<div style="text-align: center"><span class="switch active" id="switch-settings-wallets" onclick="switchTab('wallets')">{{.LANG.Modal.settings.label_wallets}}</span><span class="switch" id="switch-settings-notifications" onclick="switchTab('notifications')">{{.LANG.Modal.settings.label_notifications}}</span></div>
<div id="settingsErr"></div>


<div id="settingsWallets" style="margin: 20px; margin-left: 0px !important">
 - {{.LANG.Modal.settings.wallets_info}}
<div id="settingsWalletsList">
<div id="settings_wallets_loading"><h3>{{.LANG.Modal.settings.wallets_loading}}</h3></div>
</div>
<div class="add_more_wallet_btn" onclick="addWalletLabels(1)" id="set_addmorewallets"><img src="/static/images/addmore.webp" style="max-width: 20px; vertical-align: bottom;" alt="Add more">&nbsp;{{.LANG.Modal.settings.addmore_btn}}</div>
</div>

<div id="settingsNotifications" style="display: none">
<h3 style="text-align: center; margin: 20px">{{.LANG.Modal.settings.tmp_devel}}</h3>
</div>
<div style="text-align: center; margin: 10px">
<input type="button" value="{{.LANG.Modal.settings.save_btn}}" onclick="saveSettings()" id="saveSettings_btn" class="sendbutton">
</div>
</div>
</div>
</div>
</div>
