# Broker Project

# To deploy this broker you will need to create a PostgrSQL database with name "broker" and a table "messages" with fileds: identifier integer primary key, body text, creation timestamp without timezone, expiration integer, subject text.

It has a simple strcuture, it can be used with multiple pods which sync using redis, the subscription does not support multiple pod deployment, but the publish and fetch methods work fine that way. Some metrics are collected with prometheus which can be accessed through services in kubernetes, the code has been instrumented for Jaeger sampling as well.
