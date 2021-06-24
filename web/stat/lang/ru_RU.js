var LANG = {
	lang_values: {"en_US": "English", "ru_RU": "Русский"},
	current_lang: "ru_RU",
	sense_timezone: "",
	sense_of: "из",
	sense_days: "дней",
	sense_hours: "часов",
	sense_seconds: "секунд",
	sense_minutes: "минут",
	sense_month: "месяц",
	sense_years: "лет",
	sense_every: "каждые",
	sense_relayh: "relays / час",
	warn_not_mining: "узел не майнит!",
	warn_not_mining_multiple: " узла (-ов) не майнят сейчас!",
	nav: {
		addnodes: "Добавить узлы",
		reference: "Справка",
		road: "Roadmap",
		changelog: "Список изменений",
		settings: "Настройки",
		logout: "Выход"
	},
	donate: {
		title: "Пожертвовать"
	},
	attention: {
		text: "Важно! Сохраните ссылку для доступа, добавьте этот адрес в закладки!",
		btn_text: "Сохранил, не показывать"
	},
	index: {
		gen_btn_val: "Сгенерировать новый ID",
		gen_btn_change: "Пожалуйста, ожидайте",
		enter_btn_val: "Войти с помощью ID",
		enter_btn_change: "Войти",
		enter_input_id: "Вставьте свой ID",
		label_select_lang: "Язык"
	},
	wallet_tracker: {
		title: "Мониторинг кошельков",
		balance_label: "Баланс",
		walletname_label: "Кошелек",
		wait_for_update: "Ожидает обновления",
		wallets_nf: "Нет добавленных кошельков"
	},
	nodes_tables: {
		NST: {
			title: "Сеть NKN",
			sum_relays: "Общая производительность",
			sum_nodes: "Общее количество узлов",
			sum_persist: "Узлы в PERSIST_FINISHED",
			sum_average_relays: "Средняя производительность",
			sum_average_uptime: "Среднее время работы",
			last_update: "Последнее обновление"
		},
		NST_client: {
			title: "Производительность ваших узлов",
			sum_nodes: "Узлы",
			network_control: "Ваш контроль сети",
			sum_relays: "Общая производительность",
			sum_average_relays: "Средняя производительность",
			waitproposal: "Ожидаемые награды (блоки)",
			waitproposal_nkn: "Ожидаемые награды в NKN",
			average_rewards_interval: "Средний интервал нахождения наград"
		},
		title: "Статистка ваших узлов",
		no_nodes_label: "Нет добавленных узлов",
		col_name: "Имя",
		col_ip: "IP адрес",
		col_status: "Статус",
		col_proposal: "Proposed",
		col_height: "Высота блока",
		col_uptime: "Время работы",
		col_version: "Версия",
		col_update: "Последнее обновление",
		b_sum_label: "Суммарно",
		b_node_waiting: "Ожидает обновления",
		relayh_label: "relays / час",
		ismining_label: "майнит",
		aremining_label: "майнят",
		remove_btn_label: "Удалить выделенные"
	},
	modal: {
		close: "Закрыть",
		addNodes: {
			title: "Добавить узлы",
			label_single: "Один",
			label_multiple: "Несколько",
			label_nodeIP: "IP адрес",
			label_nodeIP_placeholder: "Только IPv4 адрес",
			label_nodeName: "Имя узла",
			label_nodeName_placeholder: "По желанию, максимум 32 символа",
			label_recommend: "Мы рекомендуем использовать короткие имена (менее 10 символов)",
			label_multi_nodeIP: "Список IP",
			label_multi_prompt: "IP должны быть разделены либо запятыми, либо пробелами, либо переносами строк, но только одним из указанных способов",
			label_multi_nodeName: "Префикс узлов (имя)",
			add_btn: "Добавить"
		},
		reference: {
			title: "Справка"
		},
		roadmap: {
			title: "Roadmap"
		},
		donate: {
			title: "Пожертвовать",
		},
		changelog: {
			title: "Список изменений"
		},
		settings: {
			title: "Настройки",
			label_wallets: "Кошельки",
			label_notifications: "Оповещения",
			label_wal_placeholder: "NKN адрес",
			label_wal: "Кошелек",
			wallets_info: "Вы можете добавить до 3 кошельков в свой аккаунт.",
			wallets_loading: "Список ваших кошельков загружается. Пожалуйста, подождите!",
			tmp_devel: "Этот раздел в разработке. Пожалуйста, попробуйте позже.",
			addmore_btn: "Добавить еще",
			save_btn: "Сохранить настройки"
		}
	},
	outdated: {
		title: "Страница устарела",
		text_refresh: "Пожалуйста, нажмите F5 или CTRL+F5 для обновления страницы.",
		text_curver: "Отображаемая версия",
		text_lastver: "Последняя версия"
	},
	answers: {
		unknown: "Неизвестная ошибка",
		ajaxMainErr: "Неудачный AJAX запрос",
		1: "Ошибка при работе с базой данных",
		2: "Ошибка при работе с базой данных",
		enter: {
			1: "Ошибка при работе с базой данных",
			2: "Ошибка при работе с базой данных",
			254:"Неверная длина ID",
			255: "ID не найден"
		},
		genId: {
			1: "Ошибка при работе с базой данных",
			2: "Ошибка при работе с базой данных",
			3: "Вы уже создали 3 ID. Пожалуйста, подождите 30 минут"
		},
		rmnodes: {
			1: "Ошибка при работе с базой данных",
			2: "Ошибка при работе с базой данных",
			3: "Некорректный POST-запрос"
		},
		addnodes: {
			1: "Ошибка при работе с базой данных",
			2: "Ошибка при работе с базой данных",
			3: "Некорректный POST-запрос",
			4: "Некорректный IP присутствует в вашем листе",
			5: "Некорректный разделитель (менее 2 IP найдено)",
			6: "Имя ноды слишком длинное (более 32 символов)",
			7: "Вы не можете добавить более 5000 узлов в свой аккаунт"
		},
		getmynodesstat: {
			1: "Ошибка при работе с базой данных",
			2: "Ошибка при работе с базой данных"
		},
		getmywallets: {
			1: "Ошибка при работе с базой данных",
			2: "Ошибка при работе с базой данных"
		},
		settings: {
			1: "Ошибка при работе с базой данных",
			2: "Ошибка при работе с базой данных",
			4: "Неверный размер массива",
			5: "Некорректная длина адреса",
			6: "Неверный NKN адрес",
			incorrect: "Отклонено! Некорректный адрес кошелька"
		}
	}
};
