<script>
	var index_src = `<div style="width: 180px;"><img src="/static/images/nkn-logo.png" alt="NKN" style="max-width: 170px;"></div><div style="text-align: center; width: 250px;"><div id="error"></div><div id="setLang" style="margin-top: 15px"><span style="margin: 5px">
	{{ .LANG.Index.label_select_lang }}:</span>
	<select onchange="setLanguage('index_src', this.value)">
	{{ $lang := .LANG.CurrentLang }}
	{{range $key, $val := .LANG.LangValues }}
		{{if eq $lang $key }}
			<option value="{{$key}}" selected="selected">{{$val}}</option>
		{{else}}
			<option value="{{$key}}">{{$val}}</option>
		{{end}}
	{{end}}
	</select>
	</div>
	<input type="button" class="b-gen" value="{{ .LANG.Index.gen_btn_val }}" onclick="genId(this)"><br><span id="s-enter"><input type="button" id="b-enter" class="b-enter" value="{{ .LANG.Index.enter_btn_val }}" onclick="toggleEnter()"></span></div>
	<div id="toggleEnter" style="display: none"><input type="text" id="hashId" maxlength="66" value="" placeholder="{{ .LANG.Index.enter_input_id }}" class="t-enter"><br><input type="button" value="{{ .LANG.Index.enter_btn_change }}" id="b-enter2" class="b-enter" onclick="enterById()"></div>
	{{{include "pages/modal/donate.tpl"}}}|m
	{{{include "pages/modal/disconnect.tpl"}}}|m
	`

	var attention = `<div>{{ .LANG.Attention.text }}<br><a href="/login/auth/{{ .Hash }}/">/login/auth/{{ .Hash }}/</a><br><input type="button" value="{{ .LANG.Attention.btn_text }}" class="attention_yes" onclick="hideAttention()"></div>`

	var view_switchLang = `<div class="switchLang" id="switchLang"><div id="site_version"></div><span onclick="setLanguage('view_src', 'en_US')" id="lang_en_US">EN</span><span onclick="setLanguage('view_src', 'ru_RU')" id="lang_ru_RU">RU</span><span onclick="setLanguage('view_src', 'zn_CN')" id="lang_zn_CN">CN</span></div>`

	var view_logo = `<div class="logo"><div><img src="/static/images/nkn-logo.png" alt=""></div><div>. OVH</div></div><div style="clear: both;"></div>`

	var version_src = `<div id="update_needed">
<h1>{{ .LANG.Outdated.title }}</h1>
<span>{{ .LANG.Outdated.text_refresh }}</span>
<br>
<span>{{ .LANG.Outdated.text_curver }}: 
<strong>{{ .CurVersion }}</strong>
<br>{{ .LANG.Outdated.text_lastver }}: 
<strong>{{ .LastVersion }}</strong>
</span>
</div>`

	var baseloading = `<div id="baseloading"></div>`

	var view_src = `
{{{include "pages/nav.tpl"}}}|m
<div style="clear: both;"></div>
<div id="completedQuery" style="display: none">Query completed</div>
<div class="attention_id" id="attention">
</div>
<div style="clear: both;"></div>
<h2 style="text-align: center; margin: 7px;">{{ .LANG.NodesTables.NST.title }}</h2>
<div class="NetworkStatus" id="NS">
<div id="jNST_loading">
	<img src="/static/images/nkn_loading.png" alt="loading">
</div>
<div id="jNST">
<div class="NetworkStatusTable" id="NST">
<div class="tr">
<div class="td" style="text-align: right; width: 50%">{{ .LANG.NodesTables.NST.sum_relays }}:</div>
<div class="td"><span id="ns-relays_per_hour"></span> <span id="ns-relays_per_hour-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{ .LANG.NodesTables.NST.sum_nodes }}:</div>
<div class="td"><span id="ns-nodes_count"></span> <span id="ns-nodes_count-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{ .LANG.NodesTables.NST.sum_persist }}:</div>
<div class="td"><span id="ns-persist_nodes_count"></span> <span id="ns-persist_nodes_count-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{ .LANG.NodesTables.NST.sum_average_relays }}:</div>
<div class="td"><span id="ns-average_relays"></span> <span id="ns-average_relays-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{ .LANG.NodesTables.NST.sum_average_uptime }}:</div>
<div class="td"><span id="ns-average_uptime"></span> <span id="ns-average_uptime-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{ .LANG.NodesTables.NST.last_update }}:</div>
<div class="td"><span id="ns-latest_update"></span> <span id="ns-latest_update-sense"></span></div>
</div>
</div>
</div>
</div>

<div id="jWST">
<h2 style="text-align: center; margin: 7px;">{{ .LANG.WalletTracker.title }}</h2>
<div id="user_wallets">
<div class="wallets_nf" id="wallets_nf">{{ .LANG.WalletTracker.wallets_nf }}</div>
</div>
</div>

<div id="jNST_client">
<h2 style="text-align: center; margin: 7px;">{{ .LANG.NodesTables.NST_client.title }}</h2>
<div class="nodes_summary_table" id="nodes_summary_table">

<div class="tr">
<div class="td">{{ .LANG.NodesTables.NST_client.sum_nodes }}:</div>
<div class="td"><span id="sum-NodesCount">N/A</span></div>
</div>
<div class="tr">
<div class="td" style="">{{ .LANG.NodesTables.NST_client.network_control }}:</div>
<div class="td" style=""><span id="sum-NetworkControl">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{ .LANG.NodesTables.NST_client.sum_relays }}</div>
<div class="td"><span id="sum-AllRelays">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{ .LANG.NodesTables.NST_client.sum_average_relays }}:</div>
<div class="td"><span id="sum-AverageRelays">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{ .LANG.NodesTables.NST_client.waitproposal }}:</div>
<div class="td"><span id="sum-waitProposalMonth">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{ .LANG.NodesTables.NST_client.waitproposal_nkn }}:</div>
<div class="td"><span id="sum-waitNKNMonth">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{ .LANG.NodesTables.NST_client.average_rewards_interval }} :</div>
<div class="td"><span id="sum-waitOneProposal">N/A</span></div>
</div>

</div>
</div>

<div id="error"></div>

<h2 style="text-align: center; margin: 7px;">{{ .LANG.NodesTables.Other.title }}</h2>
<div class="js_body">

<div id="table_settings">
<div style="display: inline-block; text-align: left; margin: 5px;">Nodes per page:</div>
<div style="display: inline-block">
<select onchange="setEntriesPerPage(this.value)" id="selEntriesPerPage">
<option value="50">50</option>
<option value="100">100</option>
<option value="250">250</option>
<option value="500">500</option>
<option value="100000">All</option>
</select>
</div>
<div id="nodes_pages"></div>
</div>

<div class="refreshinfo" id="refreshStatus"><span id="refreshText">{{ .LANG.Refresher.ok }}</span> <span id="refreshSeconds">-</span> <span id="refreshSense">{{ .LANG.SenseSeconds }}</span></div>

<div style="clear: both;"></div>

<div class="nodes_table" id="nodes_table">
<div class="tr" id="tr_top">
<div class="td" id="check" style="width: 10px !important; min-width: 10px !important;">
<input type="checkbox" name="control-all" id="control-all" onchange="toggleCheckBox()"></div>
<div class="td" id="t_name" onclick="preSortNodes(this)">{{ .LANG.NodesTables.Other.col_name }}</div>
<div class="td" id="t_ip" onclick="preSortNodes(this)">{{ .LANG.NodesTables.Other.col_ip }}</div>
<div class="td" id="t_status" onclick="preSortNodes(this)">{{ .LANG.NodesTables.Other.col_status }}</div>
<div class="td" id="t_proposal" onclick="preSortNodes(this)">{{ .LANG.NodesTables.Other.col_proposal }}</div>
<div class="td" id="t_height" onclick="preSortNodes(this)">{{ .LANG.NodesTables.Other.col_height }}</div>
<div class="td" id="t_uptime" onclick="preSortNodes(this)">{{ .LANG.NodesTables.Other.col_uptime }}</div>
<div class="td" id="t_relay" onclick="preSortNodes(this)">R/H</div>
<div class="td" id="t_relay10" onclick="preSortNodes(this)">R/H 10m</div>
<div class="td" id="t_relay60" onclick="preSortNodes(this)">R/H 60m</div>
<div class="td" id="t_version" onclick="preSortNodes(this)">{{ .LANG.NodesTables.Other.col_version }}</div>
<div class="td" id="t_latestup" onclick="preSortNodes(this)">{{ .LANG.NodesTables.Other.col_update }}</div>
</div>
</div>
<div class="nodes_nf" id="nodes_nf">{{ .LANG.NodesTables.Other.no_nodes_label }}</div>
<div style="margin: 5px;"><input type="button" class="rmbutton" value="{{ .LANG.NodesTables.Other.remove_btn_label }}" onclick="rmNodes()"></div>
</div>
<div class="clear: both"></div>

{{{include "pages/modal/addnodes.tpl"}}}|m
{{{include "pages/modal/reference.tpl"}}}|m
{{{include "pages/modal/changelog.tpl"}}}|m
{{{include "pages/modal/settings.tpl"}}}|m
{{{include "pages/modal/donate.tpl"}}}|m
{{{include "pages/modal/disconnect.tpl"}}}|m
`
</script>
