<?php

require_once 'bootstrap.php';

use Thrift\Transport\TSocketPool;
use Thrift\Transport\TBufferedTransport;
use Thrift\Protocol\TBinaryProtocol;
use Thrift\Exception\TTransportException;
use Thrift\Exception\TProtocolException;
use fun\rpc\FunServantClient;
use fun\rpc\Context;
use fun\rpc\TCacheMissed;
use fun\rpc\TMongoMissed;
use fun\rpc\TMemcacheData;

try {
    $sock = new TSocketPool(array('localhost'), array(9001));
    $sock->setDebug(1);
    $sock->setSendTimeout(400000);
    $sock->setRecvTimeout(400000);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 4096, 4096);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('rid' => "123", 'reason' => 'call.init.567', 'host' => 'server1', 'ip' => '12.3.2.1'));

    // redis
    $r = $client->rd_call($ctx, 'SET', 'default', 'test_php', array('2,3,4,',));
    var_dump($r);
    $r = $client->rd_call($ctx, 'GET', 'default', 'test_php', array());
    var_dump($r);

    for ($i = 0; $i < 500; $i++) {
        $lockKey = "foo";
        var_dump($client->gm_lock($ctx, 'just a test', $lockKey));
        $client->gm_unlock($ctx, 'just a test', $lockKey);
    }

    $t1 = microtime(TRUE);
    // game get unique name with len 3
    //for ($i = 0; $i < 2; $i ++) {
    //for ($i = 0; $i < 658; $i ++) {
    for ($i = 0; $i < 50000000; $i ++) {
        //$name = $client->ping($ctx);
        $name = $client->gm_name3($ctx);
        echo "$i $name\n";
        //usleep(10000);
        //sleep(1);
    }

    $ok = $client->zk_create($ctx, "/maintain/global", "");
    var_dump($ok);
    $nodes = $client->zk_children($ctx, "/maintain");
    print_r($nodes);
    $ok = $client->zk_del($ctx, "/maintain/global");
    var_dump($ok);

    $transport->close();
} catch (Exception $ex) {
    print 'Something went wrong: ' . $ex->getMessage() . "\n";
}

echo microtime(TRUE) - $t1, "\n";
