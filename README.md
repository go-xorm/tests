[中文](https://github.com/go-xorm/tests/blob/master/README_CN.md)

Xorm is a simple and powerful ORM for Go.

[![Build Status](https://drone.io/github.com/go-xorm/tests/status.png)](https://drone.io/github.com/go-xorm/tests/latest)  [![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/go-xorm/xorm) [![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/lunny/xorm/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

This is the tests project for xorm.

Please add your test codes here if you want to pull request to xorm.

# How to run
* test sqlite

    ./sqlite3.sh

* test mysql

    create empty database xorm_test, xorm_test1, xorm_test3 on your mysql and make an account root, and let passwd is empty on localhost, and then

    ./mysql.sh

* test postgres

    create empty database xorm_test, xorm_test1, xorm_test3 on your postgres and make an account root, and let passwd is empty on localhost, and then

    ./postgres.sh
