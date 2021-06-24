var CLIENT = {
	hash: "",
	lang: "",
	ro_hash: "",
	sort: "",
	sort_type: "",
	hide_attention: false,
	version: "",
	nodes: [],
	nodesSummary: [],
	AllNodes: [],
	wallets: [],
	prices: [],
	statFetched1: false,
	statFetched2: false,
};

var Cached = {Pages:{},JS:{},};
var conf = {defaultLang: "en_US"};


function setLanguage(content, locale = false, cb = false) {
		if (locale == false) {
			var lang = localStorage.getItem("lang")
			if (typeof lang != "string" || lang == false) {
				CLIENT.lang = conf.defaultLang
			} else {
				CLIENT.lang = lang
			}
		} else {
			localStorage.setItem('lang', locale);	
			CLIENT.lang = locale
		}
		var elem = "LANG_" + CLIENT.lang
		if (typeof Cached.JS[elem] !== 'undefined') {
			delete LANG
			if (cb === false) {
				evalCompilePrint(content, Cached.JS[elem], "content")
				return
			}
			evalCompilePrint(content, Cached.JS[elem], "content", function() {cb(CLIENT.lang)})
		} else {
			$.get("/stat/lang/" + CLIENT.lang + '.js?' + Math.random(0,100000), function(data) {
				delete LANG
				if (cb === false) {
					evalCompilePrint(content, Cached.JS[elem], "content")
					return
				}
				evalCompilePrint(content, data, "content", function() {cb(CLIENT.lang)})
			});
		}
		return
}

function setLanguageView(view, locale) {
	setLanguage(view, locale, function() {
		let attention = `<div>${LANG.attention.text}<br><a href="/login/auth/${CLIENT.hash}/">/login/auth/${CLIENT.hash}/</a><br><input type="button" value="${LANG.attention.btn_text}" class="attention_yes" onclick="hide_attention()"></div>`
		$('#attention').html(attention)
		settingParameters()
		getfullstack()
		calcNodesInfo()
	})
	return
}

function startView(hash_out) {
	$('#content').remove()
	CLIENT.hash = hash_out
	$('.wrapper').prepend('<div class="content in" id="content"></div>')
	$('.wrapper').prepend(view_switchLang)
	setLanguage(view, "", function() {
			let attention = `<div>${LANG.attention.text}<br><a href="/login/auth/${CLIENT.hash}/">/login/auth/${CLIENT.hash}/</a><br><input type="button" value="${LANG.attention.btn_text}" class="attention_yes" onclick="hide_attention()"></div>`
			$('#attention').html(attention)
			settingParameters()
			autoupdater()
			$('#lang_' + CLIENT.lang).addClass('active')
	})
}

function evalCompilePrint(source, data, dest_id, cb = false) {
	$.globalEval(data)
	var t = Handlebars.compile(source)
	Cached.JS["LANG_" + CLIENT.lang] = data
	xdata = {LANG};
	$('#' + dest_id).html(t(xdata))
	if (cb !== false) {
		cb()
	}
}

function walletsInfoUpdate() {
	var n = CLIENT.wallets.length
	if (n < 1) {
		$('#wallets_nf').show()
		return
	}
	let wallets = ``
	for (i = 0; i < n; i++) {
		var nv = i+1
		if (CLIENT.wallets[i].balance < 0) {
			wallets += `<div class="wallet waiting" id="wallet-${CLIENT.wallets[i].id}"><p style="font-weight: bold">${LANG.wallet_tracker.walletname_label} ${nv}</p><p><a href="https://explorer.nkn.org/detail/address/${CLIENT.wallets[i].nkn_wallet}" rel="noreferrer" target="_blank" title="Explorer">${CLIENT.wallets[i].nkn_wallet}</a></p><p>${LANG.wallet_tracker.balance_label}: ${LANG.wallet_tracker.wait_for_update}</p></div>`
		} else {
			var usd_val = (CLIENT.wallets[i].balance * CLIENT.prices.usd).toFixed(2)
			wallets += `<div class="wallet" id="wallet-${CLIENT.wallets[i].id}"><p style="font-weight: bold">${LANG.wallet_tracker.walletname_label} ${nv}</p><p><a href="https://explorer.nkn.org/detail/address/${CLIENT.wallets[i].nkn_wallet}" rel="noreferrer" target="_blank" title="Explorer">${CLIENT.wallets[i].nkn_wallet}</a></p><p>${LANG.wallet_tracker.balance_label}: ${CLIENT.wallets[i].balance} NKN [ ${usd_val}$ ]</p></div>`
		}
	}
	$('#user_wallets').html(wallets)
}

function checkVersion(data) {
	if (CLIENT.version == "") {
		CLIENT.version = data.Version
		$('#site_version').html("Version: " + data.Version)
		return
	}
	if (CLIENT.version != data.Version) {
		let message = `<div id="update_needed"><h1>${LANG.outdated.title}</h1><span>${LANG.outdated.text_refresh}</span><br><span>${LANG.outdated.text_curver}: <strong>${CLIENT.version}</strong><br>${LANG.outdated.text_lastver}: <strong>${data.Version}</strong></span></div>`
		$('body').html(message)
	}
	return
}

function getfullstack() {
	api_query("getfullstack", true, "post", false,
		function(data) {
			if (data.Code != 0) {
				if (data.Code in LANG.answers) {
					genErr(LANG.answers[data.Code])
				} else {
					genErr(data.Value)
					return
				}
			}
			if (data.DAEMON.Code == 0) {
				checkVersion(data.DAEMON.Value)
			}
			//Check for prices
			if (data.PRICES.Code == 0) {
				CLIENT.prices = data.PRICES.Value
			} else {
				CLIENT.prices["usd"] = 0
			}

			//Check for wallets
			if (data.WALLETS.Code == 3) {
				CLIENT.wallets = []
				addWalletLabels()
				walletsInfoUpdate()
			}
			if (data.WALLETS.Code == 0) {
				CLIENT.wallets = data.WALLETS.Value
				addWalletLabels()
				walletsInfoUpdate()
			}

			//Check for NKN Network
			netstatusWorker(data.NETSTATUS)

			//Check for user nodes
			if (data.NODES.Code == 0) {
				$('#nodes_nf').hide()
				CLIENT.nodes = data.NODES.Value
			} else if (data.NODES.Code == 3) {
				$('#nodes_nf').show()
				CLIENT.nodes = [];
			}
			sortedBy()
			parseNodes()
			calcNodesInfo()
		}, function() {
			genErr(LANG.answers.ajaxMainErr)
		})
	return
}


function saveSettings() {
	var data = {wallets: []};

	//pull wallets
	var failed = false
	n = $('#settingsWalletsList').find("[id*=setwal-]").each(function(el, item) {
		var ival = $.trim($(item).val())
		if (ival.length != 36 && ival.length != 0) {
			settingsErr(LANG.answers.settings.incorrect)
			failed = true
		}
		if (ival.length == 36) {
			data.wallets[el] = ival
		}
	})
	if (failed == true) {
		return
	}
	api_query("savemysettings", true, "post", data, function(data) {
		if (data.Err == true) {
			settingsErr(LANG.answers.settings[data.Code])
			return
		}
		if (data.Code == 0) {
			getmywallets()
			$('#settingsErr').hide()
			$('#settingsModal').hide(300)
		}
		return
	}, function() {
		settingsErr(LANG.answers.ajaxMainErr)
	})
}

function getmywallets() {
	api_query("getmywallets", true, "post", false,
	function(data) {
		if (data.Err === true) {
			settingsErr(LANG.answers.getmywallets[data.Code])	
			return
		}
		if (data.Code == 3) {
			CLIENT.wallets = []
			addWalletLabels()
			walletsInfoUpdate()
		}
		if (data.Code == 0) {
			CLIENT.wallets = data.Value
			addWalletLabels()
			walletsInfoUpdate()
		}
	},
	function(){
		settingsErr(LANG.answers.ajaxMainErr)
	})
}

function addWalletLabels(add_field = false) {
	var n = 1;
	var w = ""
	var modal = $('#settingsModal')
	var wloading = $('#settings_wallets_loading')
	if ($(modal).is(":visible") && !($(wloading).is(":visible")) && add_field == false) {
		return
	}
	if (add_field != false) {
		n = $('#settingsWalletsList').find("[id*=setwal-]").length

		if (n >= 3) {
			return
		}
		n++
		let wallabel = `<div style="margin: 20px 0 0 0;"><p>${LANG.modal.settings.label_wal} ${n}:</p><input id="setwal-${n}" type="text" class="inputtext" value="" placeholder="${LANG.modal.settings.label_wal_placeholder}"></div>`
		$('#settingsWalletsList').append(wallabel)
		if (n >= 3) {
			$('#set_addmorewallets').hide()
			return
		} else {
			$('#set_addmorewallets').show()
		}
		return
	}
	var l = CLIENT.wallets.length
	if (l > 0) {
			for (i = 0; i < l; i++) {
				if (n == 3) {
					$('#set_addmorewallets').hide()
				}
				if (!($("#wal" + n).length)) {
					let wallabel = `<div style="margin: 20px 0 0 0;"><p>${LANG.modal.settings.label_wal} ${n}:</p><input id="setwal-${n}" type="text" class="inputtext" value="${CLIENT.wallets[i].nkn_wallet}" placeholder="${LANG.modal.settings.label_wal_placeholder}"></div>`
					w += wallabel
					n++
				} else {
					continue
				}
			}

	} else {
		if (!($("#wal" + n).length)) {
			let wallabel = `<div style="margin: 20px 0 0 0;"><p>${LANG.modal.settings.label_wal} ${n}:</p><input id="setwal-${n}" type="text" class="inputtext" value="" placeholder="${LANG.modal.settings.label_wal_placeholder}"></div>`
			w = wallabel
		} else {
			return
		}
	}
	$('#settingsWalletsList').html(w)
	return

}


function showModal(divid) {
	$('#' + divid + 'Modal').show(300)
	return
}

function reversetogglechkbox() {
	if ($("#control-all").prop("checked")) {
		$('#nodes_table').find('input:checkbox').prop('checked', '')
	} else {
		$('#nodes_table').find('input:checkbox').prop('checked', 'checked')
	}
	return
}

function togglechkbox() {
	if ($("#control-all").prop("checked")) {
		$('#nodes_table').find('input:checkbox').prop('checked', 'checked')
	} else {
		$('#nodes_table').find('input:checkbox').prop('checked', '')
	}
	return
}

function switchTab(data) {
	var s = "addNodesSingle"
	var m = "addNodesMultiple"
	var set_wal = "settingsWallets"
	var set_not = "settingsNotifications"
	switch(data) {
		case "single":
			$('#switch-single').addClass('active')
			$('#switch-multiple').removeClass('active')
			$('#'+m).hide(300, function() {
				$('#'+s).show(300)
			})
		break
		case "multiple":
			$('#switch-multiple').addClass('active')
			$('#switch-single').removeClass('active')
			$('#'+s).hide(300, function() {
				$('#'+m).show(300)
			})
		break
		case "wallets":
			$('#switch-wallets').addClass('active')
			$('#switch-notifications').removeClass('active')
			$('#'+set_not).hide(300, function() {
				$('#'+set_wal).show(300)
			})
		break
		case "notifications":
			$('#switch-notifications').addClass('active')
			$('#switch-wallets').removeClass('active')
			$('#'+set_wal).hide(300, function() {
				$('#'+set_not).show(300)
			})
		break
	}
}

function addNodeErr(data) {
	$('#addNodesErr').text(data)
	$('#addNodesErr').show()
	return
}

function settingsErr(data) {
	$('#settingsErr').text(data)
	$('#settingsErr').show()
	return
}


function addNode() {
	$('#addNodeButton').prop('disabled', 'disabled')
	if($('#addNodesSingle').is(":visible")) {
		var nodeIP = $('#nodeIP').val()
		var nodeName = $('#nodeName').val()
		var data = {multiple: 0, IP: nodeIP, Name: nodeName};

		api_query("addnodes", true, "post", data, function(data) {
			$('#addNodeButton').prop('disabled', '')
			if (data.Code == 0) {
				$('#nodeIPList').val("")
				$('#nodeNameList').val("")
				getmynodesstat();
				$('#addNodesErr').hide()
				$('#addNodesModal').hide(300)
				return
			}
			if (data.Code in LANG.answers.addnodes) {
				addNodeErr(LANG.answers.addnodes[data.Code])
			} else {
				addNodeErr(data.Value)
			}
			return
		},
		function(data) {
			$('#addNodeButton').prop('disabled', '')
			addNodeErr(LANG.answers.ajaxMainErr)
		})

	} else if ($('#addNodesMultiple').is(":visible")) {
		var nodeIP = $('#nodeIPList').val()
		var nodeName = $('#nodeNameList').val()
		data = {multiple: 1, IP: nodeIP, Name: nodeName};
		api_query("addnodes", true, "post", data, function(data){
			$('#addNodeButton').prop('disabled', '')
			if (data.Code == 0) {
				getmynodesstat();
				$('#nodeIPList').val("")
				$('#nodeNameList').val("")
				$('#addNodesModal').hide(300)
				return
			} 
			if (data.Code in LANG.answers.addnodes) {
				addNodeErr(LANG.answers.addnodes[data.Code])
			} else {
				addNodeErr(data.Value)
			}
			return
		},
		function(data) {
			$('#addNodeButton').prop('disabled', '')
			addNodeErr(LANG.answers.ajaxMainErr)
		})
	} else {
		$('#addNodeButton').prop('disabled', '')
		addNodeErr(LANG.answers.ajaxMainErr)
	}
}

function settingParameters() {
	localStorage.setItem('hash', CLIENT.hash);	
	var hide_out = localStorage.getItem('hide_attention');
	if (hide_out == null || hide_out == false) {
		$('#attention').show(100)
	} else {
		CLIENT.hide_attention = true
	}
	var s = localStorage.getItem('sort');
	var s_type = localStorage.getItem('sort_type');
	if (s == null || CLIENT.sort_type == null) {
		CLIENT.sort = "t_name"
		CLIENT.sort_type = "ASC"
	}
	if (s != null && s_type != null) {
		CLIENT.sort = s
		CLIENT.sort_type = s_type
	}
	if (CLIENT.sort != "") {
		sort_nodes(CLIENT.sort, CLIENT.sort_type)
	}
	return
}
function EnterById() {
	$("#b-enter2").prop('disabled', true)
	var id = $("#hashId").val()
	CLIENT.hash = id
	api_query("auth", true, "post", false,
		function(data) {
			if (data.Err === true) {
				genErr(LANG.answers.enter[data.Code])	
				$("#b-enter2").prop('disabled', false)
				return
			}
			//do enter
			if (data.Code == 0 && data.Value.length == 64) {
				startView(data.Value)
			}
		},
		function(){
			genErr(LANG.answers.ajaxMainErr)
			$('#b-enter2').prop('disabled', false);
	})
	return
}
function genId(el) {
	var oldVal = $(el).val()
	$(el).val(LANG.index.gen_btn_change).end().prop('disabled', true)
	api_query("genId", false, "post", false,
		function(data, el) {
			if (data.Err === true) {
				genErr(LANG.answers.genId[data.Code])
				$(el).prop('disabled', false).end().val(oldVal)
				return
			}
			//do gen and enter
			if (data.Code == 0 && data.Value.length == 64) {
				startView(data.Value)
			}
		},
		function(){
			genErr(LANG.answers.ajaxMainErr)
	})
}
function rmnodes() {
	if ($("#nodes_table input:checkbox:checked").length > 0) {
		var data = {};
		var tmp_arr = [];
    	$("#nodes_table input:checkbox:checked").each(function(i, el) {
    		var tmp = $(el).val()
    		if ($.isNumeric(tmp)) {
    			tmp_arr[i] = tmp
    		}
    	})
    	if (tmp_arr.length > 0) {
    		var ids = Object.assign({}, tmp_arr);
    	}
	}
	datax = {ids:ids};
	api_query("rmnodes", true, "post", datax,
		function(data) {
			if (data.Err === true) {
				genErr(LANG.answers.rmnodes[data.Code])
				return
			} else if (data.Code == 0) {
				$('#control-all').prop('checked', '')
				getmynodesstat()
				calcNodesInfo()
			}
		},
		function(){
			genErr(LANG.answers.ajaxMainErr)
	})
	
}
function toggleEnter() {
	$("#b-enter").hide(150, function() {
		var enterInput = '<input type="text" id="hashId" maxlength="66" value="" placeholder="' + LANG.index.enter_input_id + '" class="t-enter"><br><input type="button" value="' + LANG.index.enter_btn_change + '" id="b-enter2" class="b-enter" onclick="EnterById()">'
		$("#s-enter").html(enterInput)
	})
}

function genErr(errtext) {
	$("#error").hide(100, function() { $(this).show(150, function() { $(this).text(errtext)})})
}

function api_query(method, hash_needed = false, qtype = "get", datax = false, cb_done = false, cb_fail = false) {
	var id = ""
	var url = "/api/" + method + "/"

	if (hash_needed && datax == false) {
		datax = {hash: CLIENT.hash};
	} else if (hash_needed && datax != false) {
		datax["hash"] = CLIENT.hash
	}
	$.ajax({
		url: url,
		dataType: "json",
		type: qtype,
		data: datax
		}).done(function(data) {
		if (cb_done !== false) {
			cb_done(data)
		}
		}).fail(function(data) {
		if (cb_fail !== false) {
	 		cb_fail(data)
		}
		})
 	return
}
function sort_nodes(id, type = null) {
	if (type == null) {
		if (CLIENT.sort != id) {
			CLIENT.sort_type = "ASC"
		}
		if (CLIENT.sort_type == "ASC" && CLIENT.sort == id) {
			CLIENT.sort_type = "DESC"
		} else {
			CLIENT.sort_type = "ASC"
		}
	}
	CLIENT.sort = id
	$('#tr_top').children('div').each(function(i,elem){
		$(elem).children('span').empty()
	})
	if (CLIENT.sort_type == "ASC") {
		$('#'+id).html($('#'+id).text() + "<span>&#9660;</span>")
	} else {
		$('#'+id).html($('#'+id).text() + "<span>&#9650;</span>")
	}
	localStorage.setItem('sort', CLIENT.sort);
	localStorage.setItem('sort_type', CLIENT.sort_type);
	sortedBy();
	parseNodes();
	return
}
function sortedBy() {
	keys = [];
	var m = new Map();
	if (CLIENT.sort == "t_ip") {
		for (const [key, value] of Object.entries(CLIENT.nodes)) {
			keys.push(value.ip)
			m.set(value.ip, key)
		}
		if (CLIENT.sort_type == "ASC") {
			keys.sort((a, b) => {
			const num1 = Number(a.split(".").map((num) => (`000${num}`).slice(-3) ).join(""));
			const num2 = Number(b.split(".").map((num) => (`000${num}`).slice(-3) ).join(""));
			return num2-num1;
			});
		} else {
			keys.sort((a, b) => {
			const num1 = Number(a.split(".").map((num) => (`000${num}`).slice(-3) ).join(""));
			const num2 = Number(b.split(".").map((num) => (`000${num}`).slice(-3) ).join(""));
			return num1-num2;
			});
		}
		var sorting = [];
		keys.forEach(function(k) {
			var id = m.get(k)
			sorting.push(CLIENT.nodes[id])
		})
		CLIENT.nodes = sorting
	} else if (CLIENT.sort == "t_name") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return (b.name > a.name) ? 1 : ((a.name > b.name) ? -1 : 0);} )
		} else {
			CLIENT.nodes.sort(function(a,b) {return (a.name > b.name) ? 1 : ((b.name > a.name) ? -1 : 0);})
		}
	} else if (CLIENT.sort == "t_status") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return (b.SyncState > a.SyncState) ? 1 : ((a.SyncState > b.SyncState) ? -1 : 0);} )
		} else {
			CLIENT.nodes.sort(function(a,b) {return (a.SyncState > b.SyncState) ? 1 : ((b.SyncState > a.SyncState) ? -1 : 0);})
		}
	} else if (CLIENT.sort == "t_height") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return b.Height - a.Height;})
		} else {
			CLIENT.nodes.sort(function(a,b) {return a.Height - b.Height;})
		}
	} else if (CLIENT.sort == "t_uptime") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return b.Uptime - a.Uptime;})
		} else {
			CLIENT.nodes.sort(function(a,b) {return a.Uptime - b.Uptime;})
		} 
	} else if (CLIENT.sort == "t_proposal") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return b.ProposalSubmitted - a.ProposalSubmitted;})
		} else {
			CLIENT.nodes.sort(function(a,b) {return a.ProposalSubmitted - b.ProposalSubmitted;})
		}
	} else if (CLIENT.sort == "t_relay") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return b.RelaysPerHour - a.RelaysPerHour;})
		} else {
			CLIENT.nodes.sort(function(a,b) {return a.RelaysPerHour - b.RelaysPerHour;})
		}
	} else if (CLIENT.sort == "t_relay10") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return b.RelaysPerHour10 - a.RelaysPerHour10;})
		} else {
			CLIENT.nodes.sort(function(a,b) {return a.RelaysPerHour10 - b.RelaysPerHour10;})
		} 
	} else if (CLIENT.sort == "t_relay60") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return b.RelaysPerHour60 - a.RelaysPerHour60;})
		} else {
			CLIENT.nodes.sort(function(a,b) {return a.RelaysPerHour60 - b.RelaysPerHour60;})
		} 
	} else if (CLIENT.sort == "t_version") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return (b.Version > a.Version) ? 1 : ((a.Version > b.Version) ? -1 : 0);} )
		} else {
			CLIENT.nodes.sort(function(a,b) {return (a.Version > b.Version) ? 1 : ((b.Version > a.Version) ? -1 : 0);})
		}
	} else if (CLIENT.sort == "t_latestup") {
		if (CLIENT.sort_type == "ASC") {
			CLIENT.nodes.sort(function(a,b) {return b.Currtimestamp - a.Currtimestamp;})
		} else {
			CLIENT.nodes.sort(function(a,b) {return a.Currtimestamp - b.Currtimestamp;})
		}
	}
	return
}
function clearNodesTable() {
	$('#nodes_table').children('div').each(function(i, item){
		if ($(item).attr('id').match(/Node\-([0-9]+)/g, '')) {
			$(item).remove()
		}
	})
}
function parseNodes() {
	var sumUptime = 0
	var sumRelaysPerHour = 0
	var sumRelaysPerHour10 = 0
	var sumRelaysPerHour60 = 0
	var averageRelays = 0
	var averageUptime = 0
	var sumOffline = 0
	var sumProposal = 0
	var sumNodes = 0
	var sumActiveNodes = 0
	var RelaysViewK = 0
	clearNodesTable();
	$.each(CLIENT.nodes, function(key, item) {
		if (item.SyncState != "PERSIST_FINISHED") {
			var cl = "warning"
		} else if (item.SyncState == "_OFFLINE_") {
			var cl = "error"	
		} else {
			var cl = "mining"
			sumActiveNodes++
		}
		sumNodes++
		var d = new Date(item.latest_update + "+0300")
		var latest_update = d.toLocaleDateString() + " / " + d.toLocaleTimeString()
		if (typeof item.Err === "undefined") {
			sumUptime += item.Uptime
			sumProposal += item.ProposalSubmitted
			sumRelaysPerHour += item.RelaysPerHour
			RelaysViewK = (item.RelaysPerHour/1000).toFixed(2)
			if (item.RelaysPerHour10 > 0) {
				sumRelaysPerHour10 += item.RelaysPerHour10
				Relays10ViewK = (item.RelaysPerHour10/1000).toFixed(2) + "k"
			} else {
				Relays10ViewK = "N/A";
			}
			if (item.RelaysPerHour60 > 0) {
				sumRelaysPerHour60 += item.RelaysPerHour60
				Relays60ViewK = (item.RelaysPerHour60/1000).toFixed(2) + "k"
			} else {
				Relays60ViewK = "N/A";
			}
			if (item.Version != "") {
				VersionView = item.Version
			} else {
				VersionView = "N/A";
			}
			if (item.Uptime >= 3600*24) {
				UptimeView = (item.Uptime/(3600*24)).toFixed(2) + " d"
			} else if (item.Uptime >= 3600) {
				UptimeView = (item.Uptime/3600).toFixed(2) + " h"
			} else {
				UptimeView = item.Uptime + " s"
			}
			let node = `<div class="tr ${cl}" id="Node-${item.node_id}"><div class="td"><input type="checkbox" id="controlNode-${item.node_id}" name="controlNode-${item.node_id}" value="${item.node_id}"></div><div class="td">${item.name}</div><div class="td">${item.ip}</div><div class="td">${item.SyncState}</div><div class="td">${item.ProposalSubmitted}</div><div class="td">${item.Height}</div><div class="td">${UptimeView}</div><div class="td">${RelaysViewK}k</div><div class="td">${Relays10ViewK}</div><div class="td">${Relays60ViewK}</div><div class="td">${VersionView}</div><div class="td">${latest_update}</div></div>`;
			$('#nodes_table').append(node)
		} else {
			if (item.Err == 1) {
				var status = "OFFLINE"
				sumOffline++
				cl = "error"
			}
			if (item.Err == 2) {
				var status = LANG.nodes_tables.b_node_waiting
				cl = "waiting"
			}
			if (item.Err == 1 || item.Err == 2) {
				CLIENT.nodes[key]['Height'] = -1
				CLIENT.nodes[key]['ProposalSubmitted'] = -1
				CLIENT.nodes[key]['Uptime'] = -1
				let node = `<div class="tr ${cl}" id="Node-${key}"><div class="td"><input type="checkbox" id="controlNode-${item.node_id}" name="controlNode-${item.node_id}" value="${item.node_id}"></div><div class="td">${item.name}</div><div class="td">${item.ip}</div><div class="td">${status}</div><div class="td">N/A</div><div class="td">N/A</div><div class="td">N/A</div><div class="td">N/A</div><div class="td">N/A</div><div class="td">N/A</div><div class="td">N/A</div><div class="td">N/A</div><div class="td">${latest_update}</div></div>`;
				$('#nodes_table').append(node)
			}
		}
	})
	if (sumActiveNodes > 0) {
		if (sumRelaysPerHour > 0 && (sumNodes-sumOffline) > 0) {
			averageRelays = sumRelaysPerHour/(sumNodes-sumOffline)
			averageUptime = sumUptime/(sumNodes-sumOffline)
		}
	}
	CLIENT.nodesSummary['averageRelays'] = averageRelays
	CLIENT.nodesSummary['averageUptime'] = averageUptime
	CLIENT.nodesSummary['Proposal'] = sumProposal
	CLIENT.nodesSummary['RelaysPerHour'] = sumRelaysPerHour
	CLIENT.nodesSummary['RelaysPerHour10'] = sumRelaysPerHour10
	CLIENT.nodesSummary['RelaysPerHour60'] = sumRelaysPerHour60
	CLIENT.nodesSummary['Nodes'] = sumNodes
	CLIENT.nodesSummary['ActiveNodes'] = sumActiveNodes
	if (sumNodes > 0) {
		CLIENT.statFetched1 = true
	}
	return
}
function getmynodesstat() {
	api_query("getmynodesstat", true, "post", false, function(data) {
		if (data.Code == 0) {
			$('#nodes_nf').hide()
			CLIENT.nodes = data.Value
		} else if (data.Code == 3) {
			$('#nodes_nf').show()
			CLIENT.nodes = [];
		}
		if (data.Err === true) {
			genErr(LANG.answers.getmynodesstat[data.Code])
			return
		}
		sortedBy()
		parseNodes()
		return
	}, function() {
		genErr(LANG.answers.ajaxMainErr)
	})

	return
}
function calcNodesInfo() {
	if (CLIENT.statFetched1 == false || CLIENT.statFetched2 == false) {
		return
	}
	var r = CLIENT.nodesSummary['RelaysPerHour'] == 0 ? "N/A": (CLIENT.nodesSummary['RelaysPerHour']/1000).toFixed(2) + "k"
	var r10 = CLIENT.nodesSummary['RelaysPerHour10'] == 0 ? "N/A": (CLIENT.nodesSummary['RelaysPerHour10']/1000).toFixed(2) + "k"
	var r60 = CLIENT.nodesSummary['RelaysPerHour60'] == 0 ? "N/A": (CLIENT.nodesSummary['RelaysPerHour60']/1000).toFixed(2) + "k"
	var au = CLIENT.nodesSummary['averageUptime']
	if (au >= 3600*24) {
		UptimeView = (au/(3600*24)).toFixed(2) + " d"
	} else if (au == 0) {
		UptimeView = "-"
	} else if (au >= 3600) {
		UptimeView = (au/3600).toFixed(2) + " h"
	} else {
		UptimeView = au.toFixed(2) + " s"
	}
	let sumstat = `<div class="tr" id="sum_tr"><div class="td"></div><div class="td">${LANG.nodes_tables.b_sum_label}</div><div class="td">*</div><div class="td">*</div><div class="td">${CLIENT.nodesSummary['Proposal']}</div><div class="td">*</div><div class="td">${UptimeView}</div><div class="td">${r}</div><div class="td">${r10}</div><div class="td">${r60}</div><div class="td">*</div><div class="td">*</div></div>`
	$('#sum_tr').remove();
	$('#tr_top').after(sumstat)

	var blocks_per_day = CLIENT.nodesSummary['average_blocksPerDay'] == 0 ? 3850 : CLIENT.nodesSummary['average_blocksPerDay']
	var month = blocks_per_day * 30;
	var control_percentage = CLIENT.nodesSummary['RelaysPerHour']/CLIENT.AllNodes['relays_per_hour']*100;
	var waitProposalMonth = month/100*control_percentage;
	var inactiveNodes = CLIENT.nodesSummary['Nodes'] - CLIENT.nodesSummary['ActiveNodes'];
	if (control_percentage != 0) {
		var waitOneProposal = 1440/(blocks_per_day/100*control_percentage)/60;
		if (waitOneProposal > 24) {
			waitOneProposal = waitOneProposal/24;
			var wop_sense = " " + LANG.sense_days
		} else if (waitOneProposal < 1) {
			waitOneProposal = waitOneProposal*60;
			var wop_sense =  " " + LANG.sense_minutes
		} else {
			var wop_sense =  " " + LANG.sense_hours
		}
	} else {
		var waitOneProposal = 0
		var wop_sense = ""
	}

	var inactives
	var sm
	if (CLIENT.nodesSummary['ActiveNodes'] > 1) {
		sm = LANG.nodes_tables.aremining_label
	} else {
		sm = LANG.nodes_tables.ismining_label
	}
	if (CLIENT.nodesSummary['ActiveNodes'] < CLIENT.nodesSummary['Nodes']) {
		inactives = CLIENT.nodesSummary['Nodes'] - CLIENT.nodesSummary['ActiveNodes'];
		if (inactives > 1) {
			genErr(inactives + " " + LANG.warn_not_mining_multiple)
		} else {
			genErr(inactives + " " + LANG.warn_not_mining)
		}
	} else {
		$('#error').hide(300)
	}

	var averageRelays = Math.floor(CLIENT.nodesSummary['averageRelays']).toLocaleString();
	$('#sum-NodesCount').text(CLIENT.nodesSummary['ActiveNodes'] + " " + LANG.sense_of + " " + CLIENT.nodesSummary['Nodes'] + " " + sm)
	$('#sum-NetworkControl').text(control_percentage.toFixed(5) + "%")
	$('#sum-AllRelays').text((CLIENT.nodesSummary['RelaysPerHour']).toLocaleString() + " " + LANG.sense_relayh)
	$('#sum-AllRelays10').text((CLIENT.nodesSummary['RelaysPerHour10']).toLocaleString() + " " + LANG.sense_relayh)
	$('#sum-AverageRelays').text(averageRelays + " " + LANG.sense_relayh)
	$('#sum-waitProposalMonth').text("≈" + waitProposalMonth.toFixed(2) + " / " + LANG.sense_month)

	var wait_per_month = waitProposalMonth * 11.09
	var wait_per_month_usd = (wait_per_month * CLIENT.prices.usd).toFixed(2)
	$('#sum-waitNKNMonth').text("≈" + (wait_per_month).toFixed(2) + " NKN / " + LANG.sense_month + " [ " + wait_per_month_usd  + "$ ]")
	$('#sum-waitOneProposal').text(LANG.sense_every + " ≈" + waitOneProposal.toFixed(2) + wop_sense)
	$('#jNST_client').show()
	CLIENT.statFetched2 = false	
	return
}


function netstatusWorker(data) {
	$.each(data.Value, function(key, item) {
			var sense = ""
			CLIENT.AllNodes[key] = item
			switch(key) {
				case 'average_uptime':
					if (item <= 3600) {
						item = item
						sense = LANG.sense_seconds
					} else if(item/3600 <= 24) {
						item = item/3600
						sense = LANG.sense_hours
					} else if (item/3600/24 <= 365) {
						item = item/3600/24
						sense = LANG.sense_days
					} else if (item/3600/24 > 365) {
						sense = LANG.sense_years
					}
					item = item.toFixed(2)
				break;
				case 'relays_per_hour':
					sense = LANG.sense_relayh
				break;
				case 'average_relays':
					sense = LANG.sense_relayh
				break;
				case 'latest_update':
					d = new Date(item + "+0300")
					item = d.toLocaleDateString() + " / " + d.toLocaleTimeString()
					sense = LANG.sense_timezone
				break;
			}
			if (typeof item === "number" && key != "average_uptime" && key != "last_height" && key != "average_blockTime" && key != "average_blocksPerDay") {
				item = item.toLocaleString()
			}
			if (key != "last_height" && key != "last_timestamp" && key != "average_blockTime" && key != "average_blocksPerDay") {
  				$('#ns-' + key).text(item)
  				$('#ns-' + key + '-sense').text(sense)
  			} else {
  				CLIENT.nodesSummary[key] = item
  			}
		});
	$('#jNST_loading').hide(300, function() {$('#jNST').slideDown(1500)})
	CLIENT.statFetched2 = true
	return
}

function getnetstatus() {
	api_query("getnetstatus", true, "post", false, function(data) {
		if (data.Code !== 0)  {
			return false
		}
		netstatusWorker(data)
	})
	return
}
function hide_attention() {
	CLIENT.hide_attention = true
	localStorage.setItem('hide_attention', true)
	$('#attention').hide(100)
	return
}
function logout() {
	localStorage.clear()
	document.location = '/'
}
function autoupdater() {
	//NEW API FUNC - getfullstack
	getfullstack()
	var tid = setInterval(getfullstack, 60000);
	var tid4 = setInterval(calcNodesInfo, 3000);
}
