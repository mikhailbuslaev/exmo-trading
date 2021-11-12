<?php
//--------------------------PDO POSTGRESQL CONNECTION-----------------//
function postgre_connect($your_host, $your_port, $dbname) 
{
try {
	$pdo = new PDO('pgsql:host='.$your_host.';port='.$your_port.';dbname='.$dbname.';user=postgres;password=postgres');
	$pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
}
catch(PDOException $e) {
	echo $e->getMessage();
}
return $pdo;
}
//------------------GETTING DATA FUNCTIONS-------------------------//
function get_open_price($your_data)
{
	$open_price = array();
for ($i=0; $i < count($your_data); $i++) { 
	$open_price[$i] = $your_data[$i][o];
	}
	return $open_price;
}

function get_close_price($your_data)
{
	$close_price = array();
for ($i=0; $i < count($your_data); $i++) { 
	$close_price[$i] = $your_data[$i][c];
	}
	return $close_price;
}

function get_volume($your_data)
{
	$volume = array();
for ($i=0; $i < count($your_data); $i++) { 
	$volume[$i] = $your_data[$i][v];
	}
	return $volume;
}

function get_highest_price($your_data)
{
	$highest_price = array();
for ($i=0; $i < count($your_data); $i++) { 
	$highest_price[$i] = $your_data[$i][h];
	}
	return $highest_price;
}

function get_lowest_price($your_data)
{
	$lowest_price = array();
for ($i=0; $i < count($your_data); $i++) { 
	$lowest_price[$i] = $your_data[$i][l];
	}
	return $lowest_price;
}

function get_kandle_time($your_data)
{
	$kandle_time = array();
for ($i=0; $i < count($your_data); $i++) { 
	$kandle_time[$i] = $your_data[$i][t];
	}
	return $kandle_time;
}
//-------------------------MEDIUM AVERAGE FUNCTION-----------------------------//
function get_MA($your_points, $your_num_of_candles) {
for ($j=0; $j < (count($your_points)-$your_num_of_candles); $j++) { 
	$sum = 0;

	for ($i=$j; $i < ($your_num_of_candles+$j); $i++) { 
		$sum = $sum + $your_points[$i];
	}
	$average[$j+$your_num_of_candles] = $sum/$your_num_of_candles;
}
return $average;
}
//-------------------------GETTING INFORMATION ABOUT TREND FUNCTION-----------------------------//
function get_trend($your_points, $smoothness_of_trend) {
$ma500 = get_MA($your_points, $smoothness_of_trend);
	for ($i=$smoothness_of_trend+2; $i < count($your_points); $i++) { 
	if ($ma500[$i]<$ma500[$i-1] & $ma500[$i-1]<$ma500[$i-2]) {
	$trend[$i] = "bear";
	} else {
		if ($ma500[$i]>$ma500[$i-1] & $ma500[$i-1]>$ma500[$i-2]) {
		$trend[$i] = "bull";
		} else {
			$trend[$i] = "flat";
		}
	}

}
return $trend;
}
//$smoothness_of_trend ~ 1000-500-250 are normal values


//----------------------------------------------PROFIT CALCULATE FUNCTION-------------------------------------------//
function profit_calculate($order_book) {
$profit=0;
for ($i=1; $i < count($order_book['long_history'])+1; $i++) { 
	$profit = $profit+($order_book['long_history'][$i]['close_price']-$order_book['long_history'][$i]['open_price']);
}
for ($i=1; $i < count($order_book['short_history'])+1; $i++) { 
	$profit = $profit+($order_book['short_history'][$i]['open_price']-$order_book['short_history'][$i]['close_price']);
}
return $profit;
}

//------------------------------------------------TREND ANALYSIS FUNCTION-----------------------------------------//
function  trend_analysis($your_points, $kandle_time, $trend_smoothness) {

try {
	$pdo = new PDO('pgsql:host=mbesl;port=5432;dbname=trading_bot_memory;user=postgres;password=postgres');
	$pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
}
catch(PDOException $e) {
	echo $e->getMessage();
}



$order_data_long_insert = $pdo->prepare("insert into order_data(open_date, open_price, type) values(?, ?, 'long')");


$order_data_short_insert = $pdo->prepare("insert into order_data(open_date, open_price, type) values(?, ?, 'short')");



$order_data_long_update = $pdo->prepare("update order_data set close_date = ?, close_price = ? where type = 'long' and close_price is null" );


$order_data_short_update = $pdo->prepare("update order_data set close_date = ?, close_price = ? where type = 'short' and close_price is null" );




$long_order_status_request = $pdo->prepare("update order_status set status = ? where name = 'long'");
$long_order_status_request->bindParam(1, $long_order);

$short_order_status_request = $pdo->prepare("update order_status set status = ? where name = 'long'");
$short_order_status_request->bindParam(1, $short_order);


$get_long_order_status = "select status from order_status where name = 'long'";
$get_short_order_status = "select status from order_status where name = 'short'";

$long_order = $pdo->query($get_long_order_status);
$short_order = $pdo->query($get_short_order_status);

$long_order = $long_order->fetchAll(PDO::FETCH_COLUMN, 0);
$short_order = $short_order->fetchAll(PDO::FETCH_COLUMN, 0);

$long_order=$long_order[0];
$short_order=$short_order[0];

$trend = get_MA($your_points, $trend_smoothness);

$m_request="select value from order_book where name='m'";
$n_request="select value from order_book where name='n'";

$m_update = $pdo->prepare("update order_book set value = ? where name = 'm'");
$m_update->bindParam(1, $m);

$n_update = $pdo->prepare("update order_book set value = ? where name = 'n'");
$n_update->bindParam(1, $n);

$m = $pdo->query($m_request);
$m = $m->fetchAll(PDO::FETCH_COLUMN, 0);
$m = $m[0];
$n = $pdo->query($n_request);
$n = $n->fetchAll(PDO::FETCH_COLUMN, 0);
$n = $n[0];
for ($i=count($your_points)-1000; $i < count($your_points); $i++) { 
	if ($trend[$i] > $trend[$i-1] & ($trend[$i-1] < $trend[$i-2] || $trend[$i-2] < $trend[$i-3]) & $long_order == 'false' & $short_order == 'false') {
		$order_data_long_insert->bindParam(1, date('d-m-Y H:i:s', $kandle_time[$i]/1000));
		$order_data_long_insert->bindParam(2, $your_points[$i]);
		$order_data_long_insert->execute();
		$long_order='true';
		$long_order_status_request->execute();
		$m = $m+1;
		$m_update->execute();

	} elseif ($trend[$i] < $trend[$i-1] & ($trend[$i-1] > $trend[$i-2] || $trend[$i-2] > $trend[$i-3]) & $long_order == 'true' & $short_order == 'false') {
		$order_data_long_update->bindParam(1, date('d-m-Y H:i:s', $kandle_time[$i]/1000));
		$order_data_long_update->bindParam(2, $your_points[$i]);
		$order_data_long_update->execute();
		$long_order='false';
		$long_order_status_request->execute();

	} elseif ($trend[$i] < $trend[$i-1] & ($trend[$i-1] > $trend[$i-2] || $trend[$i-2] > $trend[$i-3]) & $long_order == 'false' & $short_order == 'false') {
		$order_data_short_insert->bindParam(1, date('d-m-Y H:i:s', $kandle_time[$i]/1000));
		$order_data_short_insert->bindParam(2, $your_points[$i]);
		$order_data_short_insert->execute();
		$short_order='true';
		$short_order_status_request->execute();
		$n = $n+1;
		$n_update->execute();

	} elseif ($trend[$i] > $trend[$i-1] & ($trend[$i-1] < $trend[$i-2] || $trend[$i-2] < $trend[$i-3]) & $long_order == 'false' & $short_order == 'true') {
		$order_data_short_update->bindParam(1, date('d-m-Y H:i:s', $kandle_time[$i]/1000));
		$order_data_short_update->bindParam(2, $your_points[$i]);
		$order_data_short_update->execute();
		$short_order = 'false';
		$short_order_status_request->execute();

	}
}

}
?>
