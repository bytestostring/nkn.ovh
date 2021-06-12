<div id="settingsModal" class="modal">
<div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">{{LANG.modal.settings.title}}</h3> <a href="#close" title="{{LANG.modal.close}}" class="close" onclick="$('#settingsModal').hide(200)">X</a></div><div class="modal-body">
<div style="text-align: center"><span class="switch active" id="switch-wallets" onclick="switchTab('wallets')">{{LANG.modal.settings.label_wallets}}</span><span class="switch" id="switch-notifications" onclick="switchTab('notifications')">{{LANG.modal.settings.label_notifications}}</span></div>
<div id="settingsErr"></div>


<div id="settingsWallets" style="margin: 20px; margin-left: 0px !important">
 - {{LANG.modal.settings.wallets_info}}
<div id="settingsWalletsList">
<div id="settings_wallets_loading"><h3>{{LANG.modal.settings.wallets_loading}}</h3></div>
</div>
<div class="add_more_wallet_btn" onclick="addWalletLabels(1)" id="set_addmorewallets"><img src="/stat/images/addmore.webp" style="max-width: 20px; vertical-align: bottom;" alt="Add more">&nbsp;{{LANG.modal.settings.addmore_btn}}</div>
</div>

<div id="settingsNotifications" style="display: none">
<h3 style="text-align: center; margin: 20px">{{LANG.modal.settings.tmp_devel}}</h3>
</div>
<div style="text-align: center; margin: 10px">
<input type="button" value="{{LANG.modal.settings.save_btn}}" onclick="saveSettings()" id="saveSettings_btn" class="sendbutton">
</div>
</div>
</div>
</div>
</div>
