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
    $sock->setSendTimeout(4000);
    $sock->setRecvTimeout(4000);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 1024, 1024);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('rid' => "123", 'reason' => 'call.init.121', 'host' => 'server1', 'ip' => '12.3.2.1'));

    // mysql select
    echo "\nDEMO SELECT\n";
    echo "===============================\n";
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'SELECT uid from UserInfo where uid>?', array(1), '');
    echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
    print_r($rows);

    // mysql update
    echo "\nDEMO UPDATE\n";
    echo "===============================\n";
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'UPDATE UserInfo set power=power+1 where uid=?', array(1), 'UserInfo:1');
    echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
    print_r($rows);

    // mysql merge blob column
    echo "\nDEMO MERGE\n";
    echo "===============================\n";
    $merged = $client->my_merge($ctx, 'AllianceShard', 'Rally', 1, 'alliance_id=51 and uid=50', 
        'Rally:' . json_encode(array(
            'alliance_id' => 51,
            'uid' => 50,
        )),
        'slots_info', 
        json_encode(
            array(
                'info' => array( 
                    "88" => time(),
                )
            )));
    print_r($merged);
    print_r(json_decode($merged->newVal, TRUE));

    // mysql transation
    echo "\nDEMO transtaion\n";
    echo "===============================\n";
    $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'BEGIN', NULL, '');
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'UPDATE UserInfo set power=power+1 where uid=?', array(1), 'UserInfo:1');
    $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'COMMIT', NULL, '');
    //$client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'ROLLBACK', NULL);
    echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
    print_r($rows);

    $transport->close();
} catch (Exception $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}

