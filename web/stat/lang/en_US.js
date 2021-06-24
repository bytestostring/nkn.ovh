var LANG = {
	lang_values: {"en_US": "English", "ru_RU": "Russian"},
	current_lang: "en_US",
	sense_timezone: "(your local timezone)",
	sense_of: "of",
	sense_days: "days",
	sense_hours: "hours",
	sense_seconds: "seconds",
	sense_minutes: "minutes",
	sense_month: "month",
	sense_years: "years",
	sense_every: "every",
	sense_relayh: "relays / h",
	warn_not_mining: "Node is not mining now",
	warn_not_mining_multiple: "Nodes are not mining now",
	nav: {
		addnodes: "Add nodes",
		reference: "Reference",
		road: "Roadmap",
		changelog: "Changelog",
		settings: "Settings",
		logout: "Logout",
	},
	donate: {
		title: "Donate"
	},
	attention: {
		text: "Note! Save the link below for access to your account, add this address into your bookmarks!",
		btn_text: "Got it! Don't remind me again"
	},
	index: {
		gen_btn_val: "Generate a new ID",
		gen_btn_change: "Please wait",
		enter_btn_val: "Log in with ID",
		enter_btn_change: "Log in",
		enter_input_id: "Paste your ID",
		label_select_lang: "Language"
	},
	wallet_tracker: {
		title: "Wallet tracker",
		balance_label: "Balance",
		walletname_label: "Wallet",
		wait_for_update: "Waiting for update",
		wallets_nf: "No wallets added"
	},
	nodes_tables: {
		NST: {
			title: "NKN Network",
			sum_relays: "Summary performance",
			sum_nodes: "Total Nodes",
			sum_persist: "PERSIST_FINISHED nodes",
			sum_average_relays: "Node's average performance",
			sum_average_uptime: "Node's average uptime",
			last_update: "The latest update"
		},
		NST_client: {
			title: "Your nodes performance",
			sum_nodes: "Nodes",
			network_control: "Your network control",
			sum_relays: "Summary performance",
			sum_average_relays: "Average performance",
			waitproposal: "Estimated rewards (blocks)",
			waitproposal_nkn: "Estimated rewards in NKN",
			average_rewards_interval: "Average interval of rewards"
		},
		title: "Your nodes statistics",
		no_nodes_label: "You have no nodes added",
		col_name: "Name",
		col_ip: "IP address",
		col_status: "Status",
		col_proposal: "Proposed",
		col_height: "Block height",
		col_uptime: "Uptime",
		col_version: "Version",
		col_update: "Latest update",
		b_sum_label: "Total",
		b_node_waiting: "Waiting for update",
		relayh_label: "relay / h",
		ismining_label: "is mining",
		aremining_label: "are mining",
		remove_btn_label: "Remove checked"
	},
	modal: {
		close: "Close",
		addNodes: {
			title: "Add nodes",
			label_single: "Single",
			label_multiple: "Multiple",
			label_nodeIP: "IP address",
			label_nodeIP_placeholder: "IPv4 address only",
			label_nodeName: "Node name",
			label_nodeName_placeholder: "Optional, 32 characters max",
			label_recommend: "We recommend you use short names (less than 10 characters)",
			label_multi_nodeIP: "IP list",
			label_multi_prompt: "IP addresses must be separated by either commas, spaces or line breaks. Only one of the ways can be used simultaneously.",
			label_multi_nodeName: "Nodes prefix (name)",
			add_btn: "Add node(s)"
		},
		reference: {
			title: "Reference"
		},
		roadmap: {
			title: "Roadmap"
		},
		donate: {
			title: "Donate",
		},
		changelog: {
			title: "Changelog"
		},
		settings: {
			title: "Settings",
			label_wallets: "Wallets",
			label_notifications: "Notifications",
			label_wal_placeholder: "NKN Address",
			label_wal: "Wallet",
			wallets_info: "You can add up to 3 wallets to your account.",
			wallets_loading: "List of your wallets is loading from our server! Please wait!",
			tmp_devel: "The section is being developed. Please try again later",
			addmore_btn: "Add more",
			save_btn: "Save settings"
		}
	},
	outdated: {
		title: "The page is outdated",
		text_refresh: "Please, press F5 or CTRL+F5 to refresh the page.",
		text_curver: "The displayed page's version",
		text_lastver: "The latest version"
	},
	answers: {
		unknown: "Unknown error",
		ajaxMainErr: "Your AJAX request failed",
		1: "A database error occurred",
		2: "A database error occurred",
		4: "An error occured",
		enter: {
			1: "A database error occurred",
			2: "A database error occurred",
			254: "Incorrect ID length",
			255: "ID not found"
		},
		genId: {
			1: "A database error occurred",
			2: "A database error occurred",
			3: "You have already created 3 ID. Please wait 30 minutes"
		},
		rmnodes: {
			1: "A database error occurred",
			2: "A database error occurred",
			3: "Incorrect POST query"
		},
		addnodes: {
			1: "A database error occurred",
			2: "A database error occurred",
			3: "Incorrect POST query",
			4: "Incorrect IP in your list",
			5: "Incorrect delimiter (less than two IP found)",
			6: "Node name is too long (more that 32 symbols)",
			7: "You cannot add over 5000 nodes to your account"
		},
		getmynodesstat: {
			1: "A database error occurred",
			2: "A database error occurred"
		},
		getmywallets: {
			1: "A database error occurred",
			2: "A database error occurred"
		},
		settings: {
			1: "A database error occurred",
			2: "A database error occurred",
			4: "Wrong the array length",
			5: "Wrong the address length",
			6: "Wrong an NKN address",
			incorrect: "Rejected! Incorrect wallet address"
		}
	}
};
