const entities = [
    {
        key: 'productTypes', title: 'Типы товаров', route: '/product-types', id: 'type_id',
        columns: [
            ['type_id', 'ID'], ['name', 'Наименование'], ['for_adult', '18+']
        ],
        fields: [
            { name: 'name', label: 'Наименование', type: 'text', required: true },
            { name: 'for_adult', label: 'Товар 18+', type: 'checkbox' }
        ],
        filters: [
            { name: 'search', label: 'Поиск по названию', type: 'text' },
            { name: 'for_adult', label: '18+', type: 'select', options: [
                { value: '', label: 'Все' }, { value: 'true', label: 'Да' }, { value: 'false', label: 'Нет' }
            ] }
        ]
    },
    {
        key: 'products', title: 'Товары', route: '/products', id: 'product_id',
        columns: [
            ['product_id', 'ID'], ['type_name', 'Тип'], ['name', 'Товар'], ['manufacturer', 'Производитель'], ['weight', 'Вес'], ['volume', 'Объем']
        ],
        fields: [
            { name: 'type_id', label: 'Тип товара', type: 'select', ref: 'productTypes', required: true },
            { name: 'name', label: 'Наименование', type: 'text', required: true },
            { name: 'manufacturer', label: 'Производитель', type: 'text' },
            { name: 'weight', label: 'Вес', type: 'number', step: '0.001' },
            { name: 'volume', label: 'Объем', type: 'number', step: '0.001' }
        ],
        filters: [
            { name: 'search', label: 'Поиск по товару/производителю', type: 'text' },
            { name: 'type_id', label: 'Тип товара', type: 'select', ref: 'productTypes', empty: 'Все типы' }
        ]
    },
    {
        key: 'prices', title: 'Цены', route: '/prices', id: 'price_id',
        columns: [
            ['price_id', 'ID'], ['product_name', 'Товар'], ['date_start', 'Дата начала'], ['date_end', 'Дата окончания'], ['discount', 'Скидка, %'], ['full_price', 'Полная цена'], ['total_price', 'Итоговая цена']
        ],
        fields: [
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', required: true },
            { name: 'date_start', label: 'Дата начала', type: 'date', required: true },
            { name: 'date_end', label: 'Дата окончания', type: 'date' },
            { name: 'discount', label: 'Скидка, %', type: 'number', step: '0.01' },
            { name: 'full_price', label: 'Полная цена', type: 'number', step: '0.01', required: true },
            { name: 'total_price', label: 'Итоговая цена', type: 'number', step: '0.01', required: true, readonly: true }
        ],
        filters: [
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', empty: 'Все товары' },
            { name: 'date_from', label: 'Дата с', type: 'date' },
            { name: 'date_to', label: 'Дата по', type: 'date' }
        ]
    },
    {
        key: 'suppliers', title: 'Поставщики', route: '/suppliers', id: 'supplier_id',
        columns: [
            ['supplier_id', 'ID'], ['name', 'Поставщик'], ['inn', 'ИНН'], ['kpp', 'КПП'], ['ogrn', 'ОГРН'], ['phone', 'Телефон'], ['email', 'Email']
        ],
        fields: [
            { name: 'name', label: 'Наименование', type: 'text', required: true },
            { name: 'inn', label: 'ИНН', type: 'text' },
            { name: 'kpp', label: 'КПП', type: 'text' },
            { name: 'ogrn', label: 'ОГРН', type: 'text' },
            { name: 'phone', label: 'Телефон', type: 'text' },
            { name: 'email', label: 'Email', type: 'email' }
        ],
        filters: [
            { name: 'search', label: 'Поиск по названию, ИНН или email', type: 'text' }
        ]
    },
    {
        key: 'supplies', title: 'Закупки', route: '/product-supplies', id: 'supply_id',
        columns: [
            ['supply_id', 'ID'], ['delivery_date', 'Дата поставки'], ['product_name', 'Товар'], ['supplier_name', 'Поставщик'], ['price', 'Цена закупки'], ['quantity', 'Количество']
        ],
        fields: [
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', required: true },
            { name: 'supplier_id', label: 'Поставщик', type: 'select', ref: 'suppliers', required: true },
            { name: 'price', label: 'Цена закупки', type: 'number', step: '0.01' },
            { name: 'quantity', label: 'Количество', type: 'number', step: '0.001' },
            { name: 'delivery_date', label: 'Дата поставки', type: 'date' }
        ],
        filters: [
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', empty: 'Все товары' },
            { name: 'supplier_id', label: 'Поставщик', type: 'select', ref: 'suppliers', empty: 'Все поставщики' },
            { name: 'date_from', label: 'Дата с', type: 'date' },
            { name: 'date_to', label: 'Дата по', type: 'date' }
        ]
    },
    {
        key: 'stock', title: 'Остатки', route: '/stock', id: 'stock_id',
        columns: [
            ['stock_id', 'ID'], ['product_name', 'Товар'], ['quantity', 'Количество']
        ],
        fields: [
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', required: true },
            { name: 'quantity', label: 'Количество', type: 'number', step: '0.001' }
        ],
        filters: [
            { name: 'search', label: 'Поиск по товару', type: 'text' },
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', empty: 'Все товары' }
        ]
    },
    {
        key: 'sales', title: 'Продажи', route: '/sales', id: 'sale_id',
        columns: [
            ['sale_id', 'ID'], ['datetime', 'Дата и время'], ['discount', 'Скидка, %'], ['full_cost', 'Полная стоимость'], ['total_cost', 'Итоговая стоимость']
        ],
        fields: [
            { name: 'datetime', label: 'Дата и время', type: 'datetime-local', required: true },
            { name: 'discount', label: 'Скидка, %', type: 'number', step: '0.01' },
            { name: 'full_cost', label: 'Полная стоимость', type: 'number', step: '0.01', required: true },
            { name: 'total_cost', label: 'Итоговая стоимость', type: 'number', step: '0.01', required: true }
        ],
        filters: [
            { name: 'date_from', label: 'Дата с', type: 'date' },
            { name: 'date_to', label: 'Дата по', type: 'date' }
        ]
    },
    {
        key: 'productSales', title: 'Товары в продажах', route: '/product-sales', id: 'product_sales_id', hiddenInMenu: true,
        columns: [
            ['product_sales_id', 'ID'], ['sale_datetime', 'Дата продажи'], ['product_name', 'Товар'], ['sale_price', 'Цена продажи'], ['quantity', 'Количество']
        ],
        fields: [
            { name: 'sale_id', label: 'Продажа', type: 'select', ref: 'sales', required: true },
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', required: true },
            { name: 'sale_price', label: 'Цена продажи', type: 'number', step: '0.01' },
            { name: 'quantity', label: 'Количество', type: 'number', step: '0.001' }
        ],
        filters: [
            { name: 'sale_id', label: 'Продажа', type: 'select', ref: 'sales', empty: 'Все продажи' },
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', empty: 'Все товары' }
        ]
    },
    {
        key: 'cancellations', title: 'Списания', route: '/cancellations', id: 'cancellation_id',
        columns: [
            ['cancellation_id', 'ID'], ['datetime', 'Дата и время'], ['product_name', 'Товар'], ['quantity', 'Количество']
        ],
        fields: [
            { name: 'datetime', label: 'Дата и время', type: 'datetime-local', required: true },
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', required: true },
            { name: 'quantity', label: 'Количество', type: 'number', step: '0.001' }
        ],
        filters: [
            { name: 'product_id', label: 'Товар', type: 'select', ref: 'products', empty: 'Все товары' },
            { name: 'date_from', label: 'Дата с', type: 'date' },
            { name: 'date_to', label: 'Дата по', type: 'date' }
        ]
    }
];

const reports = [
    {
        key: 'stock', title: 'Товары магазина в наличии', route: '/reports/stock-in-hand', period: false,
        columns: [
            ['product_name', 'Товар'], ['type_name', 'Тип'], ['manufacturer', 'Производитель'], ['quantity', 'Количество'], ['current_price', 'Текущая цена'], ['stock_value', 'Стоимость остатков']
        ]
    },
    {
        key: 'supplies', title: 'Поставки за период', route: '/reports/supplies', period: true,
        columns: [
            ['delivery_date', 'Дата'], ['product_name', 'Товар'], ['supplier_name', 'Поставщик'], ['price', 'Цена'], ['quantity', 'Количество'], ['total_supply_cost', 'Сумма']
        ]
    },
    {
        key: 'cancellations', title: 'Списания за период', route: '/reports/cancellations', period: true,
        columns: [
            ['datetime', 'Дата и время'], ['product_name', 'Товар'], ['quantity', 'Количество'], ['estimated_price', 'Оценочная цена'], ['estimated_loss', 'Оценочные потери']
        ]
    },
    {
        key: 'profit', title: 'Информация о расходах и прибыли за период', route: '/reports/profit', period: true,
        columns: [
            ['indicator', 'Показатель'], ['amount', 'Сумма']
        ]
    }
];

const state = {
    currentEntity: entities[0],
    currentReport: reports[0],
    references: {},
    rows: [],
    currentPrices: {},
    productsById: {}
};

const modal = new bootstrap.Modal(document.getElementById('entityModal'));
const saleItemsModal = new bootstrap.Modal(document.getElementById('saleItemsModal'));

document.addEventListener('DOMContentLoaded', async () => {
    bindBaseEvents();
    renderEntityList();
    renderReportList();
    await loadReferences();
    await selectEntity(entities[0].key);
    await selectReport(reports[0].key, false);
});

function bindBaseEvents() {
    document.getElementById('tablesModeBtn').addEventListener('click', () => setMode('tables'));
    document.getElementById('reportsModeBtn').addEventListener('click', () => setMode('reports'));
    document.getElementById('addBtn').addEventListener('click', () => openForm());
    document.getElementById('refreshBtn').addEventListener('click', () => loadEntityRows());
    document.getElementById('runReportBtn').addEventListener('click', () => loadReportRows());
    document.getElementById('entityForm').addEventListener('submit', submitForm);
}

function setMode(mode) {
    const tables = mode === 'tables';
    document.getElementById('tablesSection').classList.toggle('d-none', !tables);
    document.getElementById('reportsSection').classList.toggle('d-none', tables);
    document.getElementById('tablesModeBtn').classList.toggle('active', tables);
    document.getElementById('reportsModeBtn').classList.toggle('active', !tables);
}

function renderEntityList() {
    const list = document.getElementById('entityList');
    list.innerHTML = '';
    entities.filter(entity => !entity.hiddenInMenu).forEach(entity => {
        const button = document.createElement('button');
        button.className = 'list-group-item list-group-item-action';
        button.textContent = entity.title;
        button.dataset.key = entity.key;
        button.addEventListener('click', () => selectEntity(entity.key));
        list.appendChild(button);
    });
}

function renderReportList() {
    const list = document.getElementById('reportList');
    list.innerHTML = '';
    reports.forEach(report => {
        const button = document.createElement('button');
        button.className = 'list-group-item list-group-item-action';
        button.textContent = report.title;
        button.dataset.key = report.key;
        button.addEventListener('click', () => selectReport(report.key, true));
        list.appendChild(button);
    });
}

async function selectEntity(key) {
    const entity = entities.find(item => item.key === key);
    if (!entity) return;
    state.currentEntity = entity;
    document.querySelectorAll('#entityList .list-group-item').forEach(item => item.classList.toggle('active', item.dataset.key === key));
    document.getElementById('entityTitle').textContent = entity.title;
    renderFilters();
    await loadReferences();
    await loadEntityRows();
}

async function selectReport(key, run = true) {
    const report = reports.find(item => item.key === key);
    if (!report) return;
    state.currentReport = report;
    document.querySelectorAll('#reportList .list-group-item').forEach(item => item.classList.toggle('active', item.dataset.key === key));
    document.getElementById('reportTitle').textContent = report.title;
    renderReportFilters();
    if (run) await loadReportRows();
}

function renderFilters() {
    const form = document.getElementById('filterForm');
    form.innerHTML = '';
    state.currentEntity.filters.forEach(filter => {
        form.appendChild(renderInput(filter, null, true));
    });

    const col = document.createElement('div');
    col.className = 'col-12 col-md-auto';
    col.innerHTML = '<button type="submit" class="btn btn-outline-primary">Найти</button>';
    form.appendChild(col);
    form.onsubmit = event => {
        event.preventDefault();
        loadEntityRows();
    };
}

function renderReportFilters() {
    const form = document.getElementById('reportFilterForm');
    form.innerHTML = '';
    if (!state.currentReport.period) {
        form.innerHTML = '<div class="col-12 text-muted">Для этого отчета период не требуется.</div>';
        return;
    }
    form.appendChild(renderInput({ name: 'date_from', label: 'Дата с', type: 'date' }, null, true));
    form.appendChild(renderInput({ name: 'date_to', label: 'Дата по', type: 'date' }, null, true));
}

async function loadReferences() {
    const refRoutes = [
        ['productTypes', '/product-types', item => `${item.name}${item.for_adult ? ' (18+)' : ''}`],
        ['products', '/products', item => item.name],
        ['suppliers', '/suppliers', item => item.name],
        ['sales', '/sales', item => `${item.datetime} / ${formatMoney(item.total_cost)} руб.`]
    ];

    state.productsById = {};

    for (const [key, route, labelFn] of refRoutes) {
        try {
            const rows = await apiRequest(route);
            if (key === 'products') {
                rows.forEach(row => {
                    state.productsById[row.product_id] = row;
                });
            }
            state.references[key] = rows.map(row => ({
                value: row[referenceIdByKey(key)],
                label: labelFn(row)
            }));
        } catch (error) {
            state.references[key] = [];
        }
    }

    try {
        const prices = await apiRequest('/prices');
        state.currentPrices = buildCurrentPriceMap(prices);
    } catch (error) {
        state.currentPrices = {};
    }
}

function buildCurrentPriceMap(prices) {
    const result = {};
    const today = new Date().toISOString().slice(0, 10);

    prices.forEach(price => {
        const productId = price.product_id;
        if (!productId) return;

        const dateStart = String(price.date_start || '');
        const dateEnd = String(price.date_end || '');
        const isActual = dateStart <= today && (!dateEnd || dateEnd >= today);
        const totalPrice = toNumber(price.total_price);
        const fullPrice = toNumber(price.full_price);
        const value = totalPrice > 0 ? totalPrice : fullPrice;

        if (!Number.isFinite(value) || value <= 0) return;

        if (!result[productId]) {
            result[productId] = {
                value,
                dateStart,
                isActual
            };
            return;
        }

        const current = result[productId];
        if (isActual && !current.isActual) {
            result[productId] = { value, dateStart, isActual };
            return;
        }

        if (isActual === current.isActual && dateStart > current.dateStart) {
            result[productId] = { value, dateStart, isActual };
        }
    });

    return result;
}

function referenceIdByKey(key) {
    return {
        productTypes: 'type_id',
        products: 'product_id',
        suppliers: 'supplier_id',
        sales: 'sale_id'
    }[key];
}

async function loadEntityRows() {
    try {
        const params = getFormValues(document.getElementById('filterForm'));
        const rows = await apiRequest(`${state.currentEntity.route}${buildQuery(params)}`);
        state.rows = rows;
        renderTable(document.getElementById('dataTable'), state.currentEntity.columns, rows, true);
        clearAlert();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function loadReportRows() {
    try {
        const params = getFormValues(document.getElementById('reportFilterForm'));
        const rows = await apiRequest(`${state.currentReport.route}${buildQuery(params)}`);
        renderTable(document.getElementById('reportTable'), state.currentReport.columns, rows, false);
        clearAlert();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

function renderTable(table, columns, rows, withActions) {
    const visibleColumns = columns.filter(([key]) => !isGuidColumn(key));

    if (!rows.length) {
        table.innerHTML = `<tbody><tr><td class="empty-state">Нет данных для отображения</td></tr></tbody>`;
        return;
    }

    const thead = document.createElement('thead');
    const headRow = document.createElement('tr');
    visibleColumns.forEach(([, label]) => {
        const th = document.createElement('th');
        th.textContent = label;
        headRow.appendChild(th);
    });
    if (withActions) {
        const th = document.createElement('th');
        th.textContent = 'Действия';
        th.className = 'text-end';
        headRow.appendChild(th);
    }
    thead.appendChild(headRow);

    const tbody = document.createElement('tbody');
    rows.forEach(row => {
        const tr = document.createElement('tr');
        visibleColumns.forEach(([key]) => {
            const td = document.createElement('td');
            td.textContent = formatValue(row[key], key);
            tr.appendChild(td);
        });
        if (withActions) {
            const td = document.createElement('td');
            td.className = 'text-end text-nowrap';

            if (state.currentEntity.key === 'sales') {
                const productsBtn = document.createElement('button');
                productsBtn.className = 'btn btn-sm btn-outline-success me-2';
                productsBtn.textContent = 'Товары';
                productsBtn.addEventListener('click', () => openSaleItems(row));
                td.appendChild(productsBtn);
            }

            const editBtn = document.createElement('button');
            editBtn.className = 'btn btn-sm btn-outline-primary me-2';
            editBtn.textContent = 'Изменить';
            editBtn.addEventListener('click', () => openForm(row));
            td.appendChild(editBtn);

            const deleteBtn = document.createElement('button');
            deleteBtn.className = 'btn btn-sm btn-outline-danger';
            deleteBtn.textContent = 'Удалить';
            deleteBtn.addEventListener('click', () => deleteRow(row));
            td.appendChild(deleteBtn);

            tr.appendChild(td);
        }
        tbody.appendChild(tr);
    });

    table.innerHTML = '';
    table.appendChild(thead);
    table.appendChild(tbody);
}

async function openSaleItems(sale) {
    const table = document.getElementById('saleItemsTable');
    const title = document.getElementById('saleItemsTitle');
    const saleId = sale.sale_id;

    title.textContent = `Товары продажи от ${formatValue(sale.datetime, 'datetime')}`;
    table.innerHTML = '<tbody><tr><td class="empty-state">Загрузка данных...</td></tr></tbody>';
    saleItemsModal.show();

    try {
        const rows = await apiRequest(`/product-sales?sale_id=${encodeURIComponent(saleId)}`);
        renderSaleItemsTable(rows);
    } catch (error) {
        table.innerHTML = `<tbody><tr><td class="empty-state text-danger">${escapeHtml(error.message)}</td></tr></tbody>`;
    }
}

function renderSaleItemsTable(rows) {
    const table = document.getElementById('saleItemsTable');

    if (!rows.length) {
        table.innerHTML = '<tbody><tr><td class="empty-state">В этой продаже пока нет товаров.</td></tr></tbody>';
        return;
    }

    const columns = [
        ['product_name', 'Товар'],
        ['sale_price', 'Цена продажи'],
        ['quantity', 'Количество'],
        ['line_total', 'Сумма']
    ];

    const thead = document.createElement('thead');
    const headRow = document.createElement('tr');
    columns.forEach(([, label]) => {
        const th = document.createElement('th');
        th.textContent = label;
        headRow.appendChild(th);
    });
    thead.appendChild(headRow);

    const tbody = document.createElement('tbody');
    rows.forEach(row => {
        const tr = document.createElement('tr');
        const lineTotal = Number(row.sale_price || 0) * Number(row.quantity || 0);
        const rowWithTotal = {
            ...row,
            line_total: Number.isFinite(lineTotal) ? lineTotal.toFixed(2) : '—'
        };

        columns.forEach(([key]) => {
            const td = document.createElement('td');
            td.textContent = formatValue(rowWithTotal[key], key);
            tr.appendChild(td);
        });
        tbody.appendChild(tr);
    });

    table.innerHTML = '';
    table.appendChild(thead);
    table.appendChild(tbody);
}

function isGuidColumn(key) {
    return key.endsWith('_id');
}

const moneyColumnKeys = new Set([
    'current_price',
    'stock_value',
    'price',
    'total_supply_cost',
    'estimated_price',
    'estimated_loss',
    'amount',
    'full_price',
    'total_price',
    'sale_price',
    'line_total',
    'full_cost',
    'total_cost'
]);

const quantityColumnKeys = new Set([
    'quantity',
    'weight',
    'volume'
]);

function formatValue(value, key = '') {
    if (value === true) return 'Да';
    if (value === false) return 'Нет';
    if (value === null || value === undefined || value === '') return '—';

    if (moneyColumnKeys.has(key)) {
        return formatMoney(value);
    }

    if (quantityColumnKeys.has(key)) {
        return formatQuantity(value);
    }

    return value;
}

function formatQuantity(value) {
    const number = toNumber(value);
    if (!Number.isFinite(number)) return value;
    return number.toLocaleString('ru-RU', {
        minimumFractionDigits: 0,
        maximumFractionDigits: 3
    });
}

function openForm(row = null) {
    const entity = state.currentEntity;

    if (entity.key === 'sales') {
        openSaleForm(row);
        return;
    }

    document.getElementById('entityForm').dataset.mode = 'generic';
    document.getElementById('modalTitle').textContent = row ? `Редактирование: ${entity.title}` : `Добавление: ${entity.title}`;
    document.getElementById('editingId').value = row ? row[entity.id] : '';

    const fields = document.getElementById('formFields');
    fields.innerHTML = '';
    entity.fields.forEach(field => fields.appendChild(renderInput(field, row, false)));

    if (entity.key === 'prices') {
        bindPriceFormAutoCalculation();
    }

    modal.show();
}


function bindPriceFormAutoCalculation() {
    const form = document.getElementById('entityForm');
    const fullPriceInput = form.querySelector('[name="full_price"]');
    const discountInput = form.querySelector('[name="discount"]');
    const totalPriceInput = form.querySelector('[name="total_price"]');

    if (!fullPriceInput || !discountInput || !totalPriceInput) return;

    fullPriceInput.min = '0';
    discountInput.min = '0';
    discountInput.max = '100';
    totalPriceInput.readOnly = true;
    totalPriceInput.classList.add('bg-light');

    const recalculatePrice = () => {
        if (String(fullPriceInput.value).trim() === '') {
            totalPriceInput.value = '';
            return;
        }

        const fullPrice = Math.max(toNumber(fullPriceInput.value), 0);
        const discount = Math.min(Math.max(toNumber(discountInput.value), 0), 100);
        const totalPrice = fullPrice * (1 - discount / 100);
        totalPriceInput.value = totalPrice.toFixed(2);
    };

    fullPriceInput.addEventListener('input', recalculatePrice);
    discountInput.addEventListener('input', recalculatePrice);
    recalculatePrice();
}

async function openSaleForm(row = null) {
    const isEdit = Boolean(row);
    const saleId = isEdit ? row.sale_id : '';

    document.getElementById('entityForm').dataset.mode = 'sale-save';
    document.getElementById('modalTitle').textContent = isEdit ? 'Редактирование продажи' : 'Добавление продажи';
    document.getElementById('editingId').value = saleId;

    const fields = document.getElementById('formFields');
    fields.innerHTML = `
        <div class="col-12 col-md-6">
            <label class="form-label" for="saleDatetime">Дата и время</label>
            <input class="form-control" id="saleDatetime" name="datetime" type="datetime-local" required>
        </div>
        <div class="col-12 col-md-6">
            <label class="form-label" for="saleDiscount">Скидка, %</label>
            <input class="form-control" id="saleDiscount" name="discount" type="number" min="0" max="100" step="0.01" value="0">
        </div>

        <div class="col-12">
            <div class="d-flex flex-wrap justify-content-between align-items-center gap-2 mt-2 mb-2">
                <h6 class="mb-0">Товары в продаже</h6>
                <button class="btn btn-outline-primary btn-sm" type="button" id="addSaleItemBtn">Добавить товар</button>
            </div>
            <div id="saleItemsContainer"></div>
            <div class="form-text">Цена подтягивается из таблицы «Цены». При необходимости ее можно изменить вручную.</div>
        </div>

        <div class="col-12 col-md-6">
            <label class="form-label" for="saleFullCost">Полная стоимость</label>
            <input class="form-control" id="saleFullCost" name="full_cost" type="number" step="0.01" readonly required>
        </div>
        <div class="col-12 col-md-6">
            <label class="form-label" for="saleTotalCost">Итоговая стоимость</label>
            <input class="form-control" id="saleTotalCost" name="total_cost" type="number" step="0.01" readonly required>
        </div>
    `;

    document.getElementById('saleDatetime').value = isEdit ? normalizeInputValue(row.datetime, 'datetime-local') : toLocalDateTimeInputValue(new Date());
    document.getElementById('saleDiscount').value = isEdit ? String(row.discount || 0) : '0';
    document.getElementById('saleDiscount').addEventListener('input', recalculateSaleTotals);
    document.getElementById('addSaleItemBtn').addEventListener('click', () => addSaleItemRow());

    if (isEdit) {
        fields.insertAdjacentHTML('afterbegin', '<div class="col-12" id="saleItemsLoading"><div class="alert alert-light border mb-0">Загрузка товаров продажи...</div></div>');
        modal.show();

        try {
            const saleItems = await apiRequest(`/product-sales?sale_id=${encodeURIComponent(saleId)}`);
            document.getElementById('saleItemsLoading')?.remove();

            if (saleItems.length) {
                saleItems.forEach(item => addSaleItemRow(item));
            } else {
                addSaleItemRow();
            }
            recalculateSaleTotals();
        } catch (error) {
            document.getElementById('saleItemsLoading')?.remove();
            showAlert(`Не удалось загрузить товары продажи: ${error.message}`, 'danger');
            addSaleItemRow();
            recalculateSaleTotals();
        }
        return;
    }

    addSaleItemRow();
    recalculateSaleTotals();
    modal.show();
}

function addSaleItemRow(item = {}) {
    const container = document.getElementById('saleItemsContainer');
    if (!container) return;

    const row = document.createElement('div');
    row.className = 'sale-item-row row g-2 align-items-end mb-2 p-2 border rounded-3 bg-light';

    const productOptions = [{ value: '', label: 'Выберите товар' }].concat(state.references.products || []);
    row.dataset.productSalesId = item.product_sales_id || '';
    row.innerHTML = `
        <input type="hidden" class="sale-product-sales-id" value="${escapeHtml(item.product_sales_id || '')}">
        <div class="col-12 col-lg-5">
            <label class="form-label">Товар</label>
            <select class="form-select sale-product" required>
                ${productOptions.map(option => `
                    <option value="${escapeHtml(option.value)}" ${String(option.value) === String(item.product_id || '') ? 'selected' : ''}>${escapeHtml(option.label)}</option>
                `).join('')}
            </select>
        </div>
        <div class="col-6 col-lg-2">
            <label class="form-label">Количество</label>
            <input class="form-control sale-quantity" type="number" min="0.001" step="0.001" value="${escapeHtml(item.quantity ?? 1)}" required>
        </div>
        <div class="col-6 col-lg-2">
            <label class="form-label">Цена</label>
            <input class="form-control sale-price" type="number" min="0" step="0.01" value="${escapeHtml(item.sale_price ?? '')}" required>
        </div>
        <div class="col-6 col-lg-2">
            <label class="form-label">Сумма</label>
            <input class="form-control sale-line-total" type="text" value="0.00" readonly>
        </div>
        <div class="col-6 col-lg-1 d-grid">
            <button class="btn btn-outline-danger" type="button" title="Удалить товар">×</button>
        </div>
    `;

    const productSelect = row.querySelector('.sale-product');
    const quantityInput = row.querySelector('.sale-quantity');
    const priceInput = row.querySelector('.sale-price');
    const removeBtn = row.querySelector('.btn-outline-danger');

    productSelect.addEventListener('change', () => {
        const productId = productSelect.value;
        const currentPrice = getCurrentProductPrice(productId);
        if (currentPrice > 0) {
            priceInput.value = currentPrice.toFixed(2);
        }
        recalculateSaleTotals();
    });
    quantityInput.addEventListener('input', recalculateSaleTotals);
    priceInput.addEventListener('input', recalculateSaleTotals);
    removeBtn.addEventListener('click', () => {
        row.remove();
        if (!container.querySelector('.sale-item-row')) {
            addSaleItemRow();
        }
        recalculateSaleTotals();
    });

    container.appendChild(row);

    if (productSelect.value && !priceInput.value) {
        const currentPrice = getCurrentProductPrice(productSelect.value);
        if (currentPrice > 0) {
            priceInput.value = currentPrice.toFixed(2);
        }
    }

    recalculateSaleTotals();
}

function getCurrentProductPrice(productId) {
    return Number(state.currentPrices?.[productId]?.value || 0);
}

function recalculateSaleTotals() {
    const rows = Array.from(document.querySelectorAll('#saleItemsContainer .sale-item-row'));
    let fullCost = 0;

    rows.forEach(row => {
        const price = toNumber(row.querySelector('.sale-price')?.value);
        const quantity = toNumber(row.querySelector('.sale-quantity')?.value);
        const lineTotal = price * quantity;
        row.querySelector('.sale-line-total').value = Number.isFinite(lineTotal) ? lineTotal.toFixed(2) : '0.00';
        fullCost += Number.isFinite(lineTotal) ? lineTotal : 0;
    });

    const discount = Math.min(Math.max(toNumber(document.getElementById('saleDiscount')?.value), 0), 100);
    const totalCost = fullCost * (1 - discount / 100);

    const fullInput = document.getElementById('saleFullCost');
    const totalInput = document.getElementById('saleTotalCost');
    if (fullInput) fullInput.value = fullCost.toFixed(2);
    if (totalInput) totalInput.value = totalCost.toFixed(2);
}

function renderInput(field, row = null, compact = false) {
    const col = document.createElement('div');
    col.className = compact ? 'col-12 col-md-3' : 'col-12 col-md-6';

    const value = row ? (row[field.name] ?? '') : '';
    const id = `${field.name}_${Math.random().toString(16).slice(2)}`;

    if (field.type === 'checkbox') {
        col.className = compact ? 'col-12 col-md-3' : 'col-12';
        col.innerHTML = `
            <div class="form-check mt-4">
                <input class="form-check-input" type="checkbox" id="${id}" name="${field.name}" ${value === true ? 'checked' : ''}>
                <label class="form-check-label" for="${id}">${field.label}</label>
            </div>
        `;
        return col;
    }

    const required = field.required ? 'required' : '';
    const label = `<label class="form-label" for="${id}">${field.label}</label>`;

    if (field.type === 'select') {
        const options = field.options || buildReferenceOptions(field);
        col.innerHTML = `${label}<select class="form-select" id="${id}" name="${field.name}" ${required}>${options.map(option => `
            <option value="${escapeHtml(option.value)}" ${String(option.value) === String(value) ? 'selected' : ''}>${escapeHtml(option.label)}</option>
        `).join('')}</select>`;
        return col;
    }

    const step = field.step ? `step="${field.step}"` : '';
    const min = field.min !== undefined ? `min="${field.min}"` : '';
    const max = field.max !== undefined ? `max="${field.max}"` : '';
    const readonly = field.readonly ? 'readonly' : '';
    const inputValue = normalizeInputValue(value, field.type);
    col.innerHTML = `${label}<input class="form-control" id="${id}" name="${field.name}" type="${field.type}" value="${escapeHtml(inputValue)}" ${step} ${min} ${max} ${readonly} ${required}>`;
    return col;
}

function buildReferenceOptions(field) {
    const result = [];
    if (field.empty !== undefined) {
        result.push({ value: '', label: field.empty });
    }
    return result.concat(state.references[field.ref] || []);
}

function normalizeInputValue(value, type) {
    if (!value) return '';
    if (type === 'datetime-local' && String(value).includes(' ')) {
        return String(value).replace(' ', 'T').slice(0, 16);
    }
    return String(value);
}

async function submitForm(event) {
    event.preventDefault();

    if (event.target.dataset.mode === 'sale-save') {
        await submitSaleWithItems();
        return;
    }

    const entity = state.currentEntity;
    const id = document.getElementById('editingId').value;
    const body = getFormValues(event.target, entity.fields);

    try {
        if (id) {
            await apiRequest(`${entity.route}/${id}`, { method: 'PUT', body: JSON.stringify(body) });
            showAlert('Запись обновлена.', 'success');
        } else {
            await apiRequest(entity.route, { method: 'POST', body: JSON.stringify(body) });
            showAlert('Запись добавлена.', 'success');
        }
        modal.hide();
        await loadReferences();
        await loadEntityRows();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function submitSaleWithItems() {
    recalculateSaleTotals();

    const saleId = document.getElementById('editingId').value;
    const isEdit = Boolean(saleId);
    const datetime = document.getElementById('saleDatetime').value;
    const discount = toNumber(document.getElementById('saleDiscount').value);
    const fullCost = toNumber(document.getElementById('saleFullCost').value);
    const totalCost = toNumber(document.getElementById('saleTotalCost').value);
    const items = collectSaleItems();

    if (!datetime) {
        showAlert('Укажите дату и время продажи.', 'danger');
        return;
    }

    if (!items.length) {
        showAlert('Добавьте хотя бы один товар в продажу.', 'danger');
        return;
    }

    const invalidItem = items.find(item => !item.product_id || item.quantity <= 0 || item.sale_price < 0);
    if (invalidItem) {
        showAlert('Проверьте товары: товар должен быть выбран, количество должно быть больше нуля, цена не должна быть отрицательной.', 'danger');
        return;
    }

    try {
        const saleBody = {
            datetime: datetime ? `${datetime}:00Z` : null,
            discount,
            full_cost: fullCost,
            total_cost: totalCost
        };

        const finalSaleId = isEdit
            ? saleId
            : (await apiRequest('/sales', {
                method: 'POST',
                body: JSON.stringify(saleBody)
            })).sale_id;

        if (isEdit) {
            await apiRequest(`/sales/${encodeURIComponent(saleId)}`, {
                method: 'PUT',
                body: JSON.stringify(saleBody)
            });
        }

        const existingRows = isEdit
            ? await apiRequest(`/product-sales?sale_id=${encodeURIComponent(saleId)}`)
            : [];
        const existingIds = new Set(existingRows.map(item => item.product_sales_id));
        const usedIds = new Set();

        for (const item of items) {
            const body = {
                sale_id: finalSaleId,
                product_id: item.product_id,
                sale_price: item.sale_price,
                quantity: item.quantity
            };

            if (item.product_sales_id) {
                usedIds.add(item.product_sales_id);
                await apiRequest(`/product-sales/${encodeURIComponent(item.product_sales_id)}`, {
                    method: 'PUT',
                    body: JSON.stringify(body)
                });
            } else {
                await apiRequest('/product-sales', {
                    method: 'POST',
                    body: JSON.stringify(body)
                });
            }
        }

        for (const oldId of existingIds) {
            if (!usedIds.has(oldId)) {
                await apiRequest(`/product-sales/${encodeURIComponent(oldId)}`, { method: 'DELETE' });
            }
        }

        showAlert(isEdit ? 'Продажа и список товаров обновлены.' : 'Продажа и товары продажи добавлены.', 'success');
        modal.hide();
        await loadReferences();
        await loadEntityRows();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

function collectSaleItems() {
    const rows = Array.from(document.querySelectorAll('#saleItemsContainer .sale-item-row'));
    return rows.map(row => ({
        product_sales_id: row.querySelector('.sale-product-sales-id')?.value || '',
        product_id: row.querySelector('.sale-product')?.value || '',
        quantity: toNumber(row.querySelector('.sale-quantity')?.value),
        sale_price: toNumber(row.querySelector('.sale-price')?.value)
    })).filter(item => item.product_id || item.quantity || item.sale_price || item.product_sales_id);
}

function getFormValues(form, fields = null) {
    const result = {};
    const formData = new FormData(form);
    const sourceFields = fields || Array.from(form.querySelectorAll('[name]')).map(input => ({
        name: input.name,
        type: input.type
    }));

    sourceFields.forEach(field => {
        const input = form.querySelector(`[name="${field.name}"]`);
        if (!input) return;

        if (field.type === 'checkbox') {
            result[field.name] = input.checked;
            return;
        }

        const raw = formData.get(field.name);
        if (raw === null || String(raw).trim() === '') {
            result[field.name] = null;
            return;
        }

        if (field.type === 'number') {
            result[field.name] = Number(raw);
        } else if (field.type === 'date') {
            result[field.name] = `${raw}T00:00:00Z`;
        }  else if (field.type === 'datetime-local') {
            result[field.name] = raw ? `${raw}:00Z` : null;
        } else {
            result[field.name] = String(raw).trim();
        }
    });
    return result;
}

function toNumber(value) {
    if (value === null || value === undefined || value === '') return 0;
    const normalized = String(value).replace(',', '.');
    const number = Number(normalized);
    return Number.isFinite(number) ? number : 0;
}

function formatMoney(value) {
    const number = toNumber(value);
    return number.toFixed(2);
}

function toLocalDateTimeInputValue(date) {
    const pad = value => String(value).padStart(2, '0');
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
}

async function deleteRow(row) {
    const entity = state.currentEntity;
    const id = row[entity.id];
    const ok = confirm('Удалить выбранную запись? Если на нее есть ссылки из других таблиц, база данных не позволит удалить ее.');
    if (!ok) return;

    try {
        await apiRequest(`${entity.route}/${id}`, { method: 'DELETE' });
        showAlert('Запись удалена.', 'success');
        await loadReferences();
        await loadEntityRows();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

function showAlert(message, type) {
    document.getElementById('alertBox').innerHTML = `
        <div class="alert alert-${type} alert-dismissible fade show" role="alert">
            ${escapeHtml(message)}
            <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Закрыть"></button>
        </div>
    `;
}

function clearAlert() {
    document.getElementById('alertBox').innerHTML = '';
}

function escapeHtml(value) {
    return String(value ?? '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#039;');
}
