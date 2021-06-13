<script>
	let index_src = `<div style="width: 180px;"><img src="/stat/images/nkn-logo.png" alt="NKN" style="max-width: 170px;"></div><div style="text-align: center; width: 250px;"><div id="error"></div><div id="setLang" style="margin-top: 15px"><span style="margin: 5px">{{LANG.index.label_select_lang}}:</span><select onchange="setLanguage(index_src, $(this).val(), function(lang) {$('#opt_lang_' + lang).prop('selected', 'selected')})">{{#each LANG.lang_values}}<option id="opt_lang_{{@key}}" value="{{@key}}">{{this}}</option>{{/each}}</select></div><input type="button" class="b-gen" value="{{LANG.index.gen_btn_val}}" onclick="genId(this)"><br><span id="s-enter"><input type="button" id="b-enter" class="b-enter" value="{{LANG.index.enter_btn_val}}" onclick="toggleEnter()"></span></div>
	{{{include "pages/modal/donate.tpl"}}}|m
`

	let view_switchLang = `<div class="switchLang" id="switchLang"><div id="site_version"></div>
<span onclick="setLanguageView(view, 'en_US')" id="lang_en_US">EN</span><span onclick="setLanguageView(view, 'ru_RU')" id="lang_ru_RU">RU</span></div><div class="logo">
<div><img src="/stat/images/nkn-logo.png" alt=""></div><div>. OVH</div></div><div style="clear: both;"></div>`

let view = `{{{include "pages/nav.tpl"}}}|m
<div style="clear: both;"></div>
<div class="attention_id" id="attention">
</div>

<div style="clear: both;"></div>
<h2 style="text-align: center; margin: 7px;">{{LANG.nodes_tables.NST.title}}</h2>
<div class="NetworkStatus" id="NS">
<div id="jNST_loading">
	<img src="/stat/images/nkn_loading.png" alt="loading">
</div>
<div id="jNST">
<div class="NetworkStatusTable" id="NST">
<div class="tr">
<div class="td" style="text-align: right; width: 50%">{{LANG.nodes_tables.NST.sum_relays}}:</div>
<div class="td"><span id="ns-relays_per_hour"></span> <span id="ns-relays_per_hour-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{LANG.nodes_tables.NST.sum_nodes}}:</div>
<div class="td"><span id="ns-nodes_count"></span> <span id="ns-nodes_count-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{LANG.nodes_tables.NST.sum_persist}}:</div>
<div class="td"><span id="ns-persist_nodes_count"></span> <span id="ns-persist_nodes_count-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{LANG.nodes_tables.NST.sum_average_relays}}:</div>
<div class="td"><span id="ns-average_relays"></span> <span id="ns-average_relays-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{LANG.nodes_tables.NST.sum_average_uptime}}:</div>
<div class="td"><span id="ns-average_uptime"></span> <span id="ns-average_uptime-sense"></span></div>
</div>
<div class="tr">
<div class="td" style="text-align: right;">{{LANG.nodes_tables.NST.last_update}}:</div>
<div class="td"><span id="ns-latest_update"></span> <span id="ns-latest_update-sense"></span></div>
</div>
</div>
</div>
</div>

<div id="jWST">
<h2 style="text-align: center; margin: 7px;">{{LANG.wallet_tracker.title}}</h2>
<div id="user_wallets">
<div class="wallets_nf">{{LANG.wallet_tracker.wallets_nf}}</div>
</div>
</div>

<div id="jNST_client">
<h2 style="text-align: center; margin: 7px;">{{LANG.nodes_tables.NST_client.title}}</h2>
<div class="nodes_summary_table" id="nodes_summary_table">

<div class="tr">
<div class="td">{{LANG.nodes_tables.NST_client.sum_nodes}}:</div>
<div class="td"><span id="sum-NodesCount">N/A</span></div>
</div>
<div class="tr">
<div class="td" style="">{{LANG.nodes_tables.NST_client.network_control}}:</div>
<div class="td" style=""><span id="sum-NetworkControl">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{LANG.nodes_tables.NST_client.sum_relays}}</div>
<div class="td"><span id="sum-AllRelays">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{LANG.nodes_tables.NST_client.sum_average_relays}}:</div>
<div class="td"><span id="sum-AverageRelays">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{LANG.nodes_tables.NST_client.waitproposal}}:</div>
<div class="td"><span id="sum-waitProposalMonth">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{LANG.nodes_tables.NST_client.waitproposal_nkn}}:</div>
<div class="td"><span id="sum-waitNKNMonth">N/A</span></div>
</div>
<div class="tr">
<div class="td">{{LANG.nodes_tables.NST_client.average_rewards_interval}} :</div>
<div class="td"><span id="sum-waitOneProposal">N/A</span></div>
</div>

</div>
</div>

<div id="error"></div>

<h2 style="text-align: center; margin: 7px;">{{LANG.nodes_tables.title}}</h2>
<div class="js_body">

<div class="nodes_table" id="nodes_table">

<div class="tr" id="tr_top">
<div class="td" id="check" style="width: 10px !important; min-width: 10px !important;"><input type="checkbox" name="control-all" id="control-all" onchange="togglechkbox()"></div>
<div class="td" id="t_name" onclick="sort_nodes(this.id)">{{LANG.nodes_tables.col_name}}</div>
<div class="td" id="t_ip" onclick="sort_nodes(this.id)">{{LANG.nodes_tables.col_ip}}</div>
<div class="td" id="t_status" onclick="sort_nodes(this.id)">{{LANG.nodes_tables.col_status}}</div>
<div class="td" id="t_proposal" onclick="sort_nodes(this.id)">{{LANG.nodes_tables.col_proposal}}</div>
<div class="td" id="t_height" onclick="sort_nodes(this.id)">{{LANG.nodes_tables.col_height}}</div>
<div class="td" id="t_uptime" onclick="sort_nodes(this.id)">{{LANG.nodes_tables.col_uptime}}</div>
<div class="td" id="t_relay" onclick="sort_nodes(this.id)">R/H</div>
<div class="td" id="t_relay10" onclick="sort_nodes(this.id)">R/H 10m</div>
<div class="td" id="t_relay60" onclick="sort_nodes(this.id)">R/H 60m</div>
<div class="td" id="t_relay360" onclick="sort_nodes(this.id)">R/H 360m</div>
<div class="td" id="t_relay1440" onclick="sort_nodes(this.id)">R/H 1440m</div>
<div class="td" id="t_latestup" onclick="sort_nodes(this.id)">{{LANG.nodes_tables.col_update}}</div>
</div>
</div>
<div class="nodes_nf" id="nodes_nf">{{LANG.nodes_tables.no_nodes_label}}</div>
<div style="margin: 5px;"><input type="button" class="rmbutton" value="{{LANG.nodes_tables.remove_btn_label}}" onclick="rmnodes()"></div>
</div>
<div class="clear: both"></div>


{{{include "pages/modal/addnodes.tpl"}}}|m
{{{include "pages/modal/reference.tpl"}}}|m
{{{include "pages/modal/roadmap.tpl"}}}|m
{{{include "pages/modal/changelog.tpl"}}}|m
{{{include "pages/modal/settings.tpl"}}}|m
{{{include "pages/modal/donate.tpl"}}}|m
`

$(document).mouseup(function (e){
		var m = $('.modal');
		if (m.is(":visible")) {
			var modal = $(".modal-content");
			if (!modal.is(e.target) && modal.has(e.target).length === 0){
				m.hide(300)
			}
		}
})

$(document).ready(function() {
	Handlebars.registerHelper('if_eq', function(arg1, arg2, options) {
    	return (arg1 == arg2) ? options.fn(this) : options.inverse(this);
	});


	var hash_out = localStorage.getItem('hash');
	if (hash_out != null && typeof hash_out === "string") {
		startView(hash_out)
		$('#switchLang').children('span').click(function(e) {
		if (e.target !== this) {
			return
		}
		$('#switchLang').children('span').removeClass('active')
		$(this).addClass('active')
	})

	$('#check').on('click', function(e) {
		if (e.target !== this) {
			return
		}
		reversetogglechkbox()
	});
	} else {
		$('.wrapper').prepend('<div class="content index" id="content"></div>')
		setLanguage(index_src, "", function(lang) {$('#opt_lang_' + lang).prop('selected', 'selected');	$('#donatelink').text(LANG.donate.title)})
	}
	$('#donatelink').click(function() {
		showModal('donate')
	})
})
</script>
