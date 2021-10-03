# nats-groups

run NATS server: docker run -d -p4222:4222 -p6222:6222 -p8222:8222 -p6543:6543 --name=nats-groups nats:2.4.0-alpine3.14 nats-server -js --port=4222

go run .\main.go

output:

```shell
{"L":"INFO","T":"2021-10-03T18:09:19.168+0300","M":"failed to get consumer","pid":3904,"stream name":"alerts","consumer name":"consumer2","error":"nats:
 consumer not found"}
{"L":"INFO","T":"2021-10-03T18:09:19.198+0300","M":"consumer added","pid":3904,"consumer":{"stream_name":"alerts","name":"durable2","created":"2021-10-0
3T14:43:22.3656767Z","config":{"durable_name":"durable2","deliver_subject":"group2","deliver_group":"group2","deliver_policy":"last","ack_policy":"expli
cit","ack_wait":5000000000,"max_deliver":-1,"replay_policy":"instant","max_ack_pending":20000},"delivered":{"consumer_seq":20,"stream_seq":20,"last_acti
ve":"2021-10-03T15:08:29.6549585Z"},"ack_floor":{"consumer_seq":20,"stream_seq":20,"last_active":"2021-10-03T15:08:29.6558439Z"},"num_ack_pending":0,"nu
m_redelivered":0,"num_waiting":0,"num_pending":0,"cluster":{"leader":"NCUPKLJSMO7PNSE7ZOAX3XC4PBIBDCLPL2LHZCPRZKKEUUP4M5GI6YMW"}}}
{"L":"INFO","T":"2021-10-03T18:09:19.200+0300","M":"subscription added","pid":3904,"sub":{"Subject":"group2","Queue":"group2"}}
{"L":"INFO","T":"2021-10-03T18:09:19.203+0300","M":"failed to get consumer","pid":3904,"stream name":"alerts","consumer name":"consumer1","error":"nats:
 consumer not found"}
{"L":"INFO","T":"2021-10-03T18:09:19.204+0300","M":"consumer added","pid":3904,"consumer":{"stream_name":"alerts","name":"durable1","created":"2021-10-0
3T14:43:22.3614482Z","config":{"durable_name":"durable1","deliver_subject":"group1","deliver_group":"group1","deliver_policy":"last","ack_policy":"expli
cit","ack_wait":5000000000,"max_deliver":-1,"replay_policy":"instant","max_ack_pending":20000},"delivered":{"consumer_seq":20,"stream_seq":20,"last_acti
ve":"2021-10-03T15:08:29.6549648Z"},"ack_floor":{"consumer_seq":20,"stream_seq":20,"last_active":"2021-10-03T15:08:29.6560021Z"},"num_ack_pending":0,"nu
m_redelivered":0,"num_waiting":0,"num_pending":0,"cluster":{"leader":"NCUPKLJSMO7PNSE7ZOAX3XC4PBIBDCLPL2LHZCPRZKKEUUP4M5GI6YMW"}}}
{"L":"INFO","T":"2021-10-03T18:09:19.205+0300","M":"subscription added","pid":3904,"sub":{"Subject":"group1","Queue":"group1"}}
{"L":"INFO","T":"2021-10-03T18:09:19.206+0300","M":"message published","pid":3904,"subject name":"alerts.subject1","message":{"ID":0,"Timestamp":"2021-1
0-03T18:09:19.2062779+03:00"}}
{"L":"INFO","T":"2021-10-03T18:09:19.207+0300","M":"message published","pid":3904,"subject name":"alerts.subject2","message":{"ID":1,"Timestamp":"2021-1
0-03T18:09:19.207324+03:00"}}
{"L":"INFO","T":"2021-10-03T18:09:19.208+0300","M":"subject alerts.subject2, group group2: msg received","pid":3904,"message":{"ID":0,"Timestamp":"2021-
10-03T18:09:19.2062779+03:00"}}
{"L":"INFO","T":"2021-10-03T18:09:19.208+0300","M":"subject alerts.subject1, group group1: msg received","pid":3904,"message":{"ID":0,"Timestamp":"2021-
10-03T18:09:19.2062779+03:00"}}
{"L":"INFO","T":"2021-10-03T18:09:19.208+0300","M":"subject alerts.subject2, group group2: msg received","pid":3904,"message":{"ID":1,"Timestamp":"2021-
10-03T18:09:19.207324+03:00"}}
{"L":"INFO","T":"2021-10-03T18:09:19.208+0300","M":"subject alerts.subject1, group group1: msg received","pid":3904,"message":{"ID":1,"Timestamp":"2021-
10-03T18:09:19.207324+03:00"}}
{"L":"FATAL","T":"2021-10-03T18:09:24.211+0300","M":"===== TIMED OUT =====","pid":3904}
exit status 1

```