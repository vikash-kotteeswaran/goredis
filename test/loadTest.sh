for i in {1..1000}
do
    echo 'PING' | nc 0.0.0.0 7379
done
