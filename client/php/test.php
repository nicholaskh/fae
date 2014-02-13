<?php

$GLOBALS['THRIFT_ROOT'] = '/opt/app/thrift/lib/php';
$GLOBALS['SERVANT_ROOT'] = '../../servant/gen-php/fun/rpc';
require_once $GLOBALS['THRIFT_ROOT'].'/Thrift.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Base/TBase.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Exception/TException.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Exception/TTransportException.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Exception/TProtocolException.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Exception/TApplicationException.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Protocol/TBinaryProtocol.php';
require_once $GLOBALS['THRIFT_ROOT'].'/StringFunc/TStringFunc.php';
require_once $GLOBALS['THRIFT_ROOT'].'/StringFunc/Core.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Factory/TStringFuncFactory.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Type/TType.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Type/TMessageType.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TSocket.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TSocketPool.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TBufferedTransport.php';
require_once $GLOBALS['SERVANT_ROOT'].'/FunServant.php';
require_once $GLOBALS['SERVANT_ROOT'].'/Types.php';

use Thrift\Transport\TSocketPool;
use Thrift\Transport\TBufferedTransport;
use Thrift\Protocol\TBinaryProtocol;
use Thrift\Exception\TTransportException;
use Thrift\Exception\TProtocolException;
use fun\rpc\FunServantClient;
use fun\rpc\Context;
use fun\rpc\TCacheMissed;
use fun\rpc\TMongoMissed;

try {
    $sock = new TSocketPool(array('localhost'), array(9001));
    $sock->setDebug(0);
    $sock->setSendTimeout(1000);
    $sock->setRecvTimeout(2500);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 1024, 1024);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('caller' => "from php test.php"));

    // ping
    $return = $client->ping($ctx);
    echo "[Client] ping received: ", $return, "\n";

    // mc
    echo '[Client] mc_set received: ', $client->mc_set($ctx, 'hello-php', 'world 世界', 120), "\n";
    echo '[Client] mc_get received: ', $client->mc_get($ctx, 'hello-php'), "\n";
    try {
        echo '[Client] mc_get hello-non-exist received: ', $client->mc_get($ctx, 'hello-non-exist'), "\n";
    } catch (TCacheMissed $ex) {
        echo $ex->getMessage(), "\n";
    }

    // dlog
    echo '[Client] dlog received: ', $client->dlog($ctx, 'error', 'ae', 
        json_encode(array('hello'=>'world'))), "\n";

    // lc
    echo '[Client] lc_set received: ', $client->lc_set($ctx, 'hello-php-lc', 'world 世界'), "\n";
    echo '[Client] lc_get received: ', $client->lc_get($ctx, 'hello-php-lc'), "\n";
    echo '[Client] lc_del received: ', $client->lc_del($ctx, 'hello-php-lc'), "\n";

    // mg.insert
    $doc = array(
        "name" => "funky.php",
        "gendar" => "M",
        "abtype" => array(
            "payment" => "a",
            "tutorial" => "b",
        )
    );
    echo '[Client] mg_insert received: ', $client->mg_insert($ctx, 'db1', 'usertest', 0, 
        bson_encode($doc)), "\n";

    // mg.findOne
    try {
        $idmap = $client->mg_find_one($ctx, 'default', 'idmap', 0,
            bson_encode(array('snsid' => '100003391571259')), bson_encode(''));
        echo "[Client] mg_find_one received: \n";
        print_r(bson_decode($idmap));
    } catch (TMongoMissed $ex) {
        echo $ex->getMessage(), "\n";
    }

    // mg.count
    echo "[Client] mg_count received:", $client->mg_count($ctx, 'default', 'idmap', 0,
        bson_encode(array('uid' => array('$gte' => 1)))), "\n";
    echo "[Client] mg_count received:", $client->mg_count($ctx, 'default', 'idmap', 0,
        bson_encode(array('uid' => array('$gte' => 100000)))), "\n";

    // mg.findAll
    echo "[Client] mg_find_all received: \n";
    try {
        $docs = $client->mg_find_all($ctx, 'default', 'idmap', 0,
            bson_encode(array('uid' => array('$gte' => 1))), bson_encode(array()),
            0, 0, array());
        $r = array();
        foreach ($docs as $doc) {
            $r[] = bson_decode($doc);
        }
        print_r($r);
    } catch (TProtocolException $ex) {
        print_r($ex);
    }

    $transport->close();
} catch (TException $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}
