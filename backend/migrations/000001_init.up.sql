--Типы товаров
create table if not exists product_types(
	type_id uuid primary key,
	name varchar(100),
	for_adult bool
);

--Товары
create table if not exists products(
	product_id uuid primary key,
	type_id uuid not null references product_types(type_id),
	name varchar(100) not null,
	manufacturer varchar(100),
	weight decimal(12,3),
	volume decimal(12,3)
);

--Цены
create table if not exists prices(
	price_id uuid primary key,
	product_id uuid not null references products(product_id),
	
	date_start date not null,
	date_end date,

	discount decimal(5,2),
	full_price decimal(10,2) not null,
	total_price  decimal(10,2) not null
);

--Поставщики
create table if not exists suppliers(
	supplier_id uuid primary key,
	name varchar(150) not null,
	inn varchar(12),
	kpp varchar(9),
	ogrn varchar(13),
	phone varchar(11),
	email varchar(50)
);

--Закупки
create table if not exists product_supplies(
	supply_id uuid primary key,
	product_id uuid not null references products(product_id),
	supplier_id uuid not null references suppliers(supplier_id),
	
	price decimal(10,2),
	quantity decimal(12,3),
	
	delivery_date date
);

--Остатки
create table if not exists stock(
    stock_id uuid primary key,
    product_id uuid references products(product_id),

    quantity numeric(12,3)
);

--Продажи
create table if not exists sales(
	sale_id uuid primary key,
	datetime TIMESTAMP not null,
	
	discount decimal(5,2),
	full_cost decimal(10,2) not null,
	total_cost  decimal(10,2) not null
);

--Связь товаров и продаж
create table if not exists product_sales(
	product_sales_id uuid primary key,
	sale_id uuid not null references sales(sale_id),
	product_id uuid not null references products(product_id),
	
	sale_price numeric(10,2),
	quantity decimal(12,3)
);

--Списания
create table if not exists cancellation(
	cancellation_id uuid primary key,
	datetime TIMESTAMP not null,
	product_id uuid not null references products(product_id),
	quantity decimal(12,3)
);