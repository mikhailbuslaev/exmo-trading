
<?php
header('refresh: 900');
//--------------------GET KANDLES FROM EXPO API--------------------------------//

include 'exmo.php';//NEED FOR TRADING
include 'functions.php';//NEED FOR GETTING POINTS FROM DATA
include 'trading strategies/kandle_engulfing_reverse.php';
$data = new Data();

$data->kandles=$data->getKandles(1635714000, time(), 15);
solver($data->kandles);
echo '<pre>';
?>

