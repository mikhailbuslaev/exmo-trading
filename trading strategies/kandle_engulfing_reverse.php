<?php
$path = $_SERVER['DOCUMENT_ROOT'];
$path .= "/functions.php";
include_once($path);
function solver($kandles) {
    $kandles_count=count($kandles);
    for($i=0;$i<$kandles_count;++$i) {
        $kandles[$i][body_length]=$kandles[$i][c]-$kandles[$i][o];
        if ($kandles[$i][body_length]<0) {
            $kandles[$i][color]=-1;
            $kandles[$i][body_length]=(-1)*$kandles[$i][body_length];
        } elseif ($kandles[$i][body_length]>0) {
            $kandles[$i][color]=1;
        } else {
            $kandles[$i][color]=0;
        }
    }
    for ($j=25; $j < $kandles_count; ++$j) { 
        $trend[$j]=0;
        for ($i=$j-25; $i < $j; ++$i) { 
            $trend[$j]=$trend[$j]+$kandles[$i][color];
        }
    }
    for ($i=25;$i<$kandles_count;++$i) { 
        if ($kandles[$i][body_length]>2*$kandles[$i-1][body_length] &&
         $kandles[$i][color]!=$kandles[$i-1][color] && $kandles[$i-1][color]!=0) {
            if ($kandles[$i][color]==-1 && $trend[$i]>0) {
                sellOrder($kandles[$i][c], $kandles[$i][t], $kandles[$i][body_length]);
            } elseif ($kandles[$i][color]==1 && $trend[$i]<0) {
                buyOrder($kandles[$i][c], $kandles[$i][t], $kandles[$i][body_length]);
            }
        }
    }
}

?>