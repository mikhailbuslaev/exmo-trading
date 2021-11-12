
<?php
header('refresh: 900');
//--------------------GET KANDLES FROM EXPO API--------------------------------//

include 'exmo.php';//NEED FOR TRADING
include 'functions.php';//NEED FOR GETTING POINTS FROM DATA

$oldtime = 1633076513;
$newtime = time();
$data_url = 'https://api.exmo.com/v1.1/candles_history?symbol=BTC_USD&resolution=15&from=' . $oldtime . '&to='. $newtime;
$data = json_decode(file_get_contents($data_url), TRUE);
date_default_timezone_set('UTC+3');
echo date('Y-m-d h:i:s A');
echo '<pre>';
$data = $data[candles];
//---------------------GET POINTS FROM KANDLES DATA-----------------------------//
$open_price = get_open_price($data);
$close_price = get_open_price($data);
$volume = get_volume($data);
$highest_price = get_highest_price($data);
$lowest_price = get_lowest_price($data);
$kandle_time = get_kandle_time($data);

//--------------------TREND STRATEGY------------------------------------------//

trend_analysis($open_price, $kandle_time, 500);

?>

