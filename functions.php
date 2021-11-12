<?php
class Data {
    public $kandles;
    public function getKandles($opentime, $closetime, $resolution) {
        $data_url = 'https://api.exmo.com/v1.1/candles_history?symbol=BTC_USD&resolution='. $resolution.'&from=' . $opentime . '&to='. $closetime;
        $output = json_decode(file_get_contents($data_url), TRUE);
        $output = $output[candles];
        return $output;
    }
}
function buyOrder($price, $today, $body_legth) {
    echo '<br>';
    echo 'buy at '.$price.' '.date("d-m-Y H:i:s", $today/1000).' '.$body_legth;
};
function sellOrder($price, $today, $body_legth) {
    echo '<br>';
    echo 'sell at '.$price.' '.date("d-m-Y H:i:s", $today/1000).' '.$body_legth;
};
?>