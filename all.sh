#!/bin/sh
exec curl --cacert cacert.pem -E gw1-both.pem "https://eq4.crest.dn42:8080/all"
