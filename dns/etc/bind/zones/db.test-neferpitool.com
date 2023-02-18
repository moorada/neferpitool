$TTL    604800
@       IN      SOA     ns1.test-neferpitool.com. root.test-neferpitool.com. (
                  3       ; Serial
             604800     ; Refresh
              86400     ; Retry
            2419200     ; Expire
             604800 )   ; Negative Cache TTL
;
; name servers - NS records
     IN      NS      ns1.test-neferpitool.com.

; name servers - A records
ns1.test-neferpitool.com.          IN      A      172.20.0.2

host1.test-neferpitool.com.        IN      A      172.20.0.3
host2.test-neferpitool.com.        IN      A      172.20.0.4
