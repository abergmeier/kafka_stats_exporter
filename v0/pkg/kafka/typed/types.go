package typed

//Stats is the root schema for librdkafka Statistics
//For details see https://github.com/edenhill/librdkafka/blob/master/STATISTICS.md
type Stats struct {
	Name             string                     `json:"name"               kpromlbl:"name"`      //Handle instance name
	ClientId         string                     `json:"client_id"          kpromlbl:"client_id"` //The configured (or default) client.id
	Type             string                     `json:"type"               kpromlbl:"type"`      //Instance type (producer or consumer)
	Ts               int                        `json:"ts"                 kpromcol:"CounterVec,internal monotonic clock (microseconds)"`
	Time             int                        `json:"time"               kpromcol:"CounterVec,Wall clock time in seconds since the epoch"`
	Age              int                        `json:"age"                kpromcol:"CounterVec,Time since this client instance was created (microseconds)"`
	Replyq           int                        `json:"replyq"             kpromcol:"GaugeVec,Number of ops (callbacks%2C events%2C etc) waiting in queue for application to serve with Poll()"`
	MsgCnt           int                        `json:"msg_cnt"            kpromcol:"GaugeVec,Current number of messages in producer queues"`
	MsgSize          int                        `json:"msg_size"           kpromcol:"GaugeVec,Current total size of messages in producer queues"`
	MsgMax           int                        `json:"msg_max"            kpromcol:"CounterVec,Threshold: maximum number of messages allowed allowed on the producer queues"`
	MsgSizeMax       int                        `json:"msg_size_max"       kpromcol:"CounterVec,Threshold: maximum total size of messages allowed on the producer queues"`
	Tx               int                        `json:"tx"                 kpromcol:"CounterVec,Total number of requests sent to Kafka brokers"`
	TxBytes          int                        `json:"tx_bytes"           kpromcol:"CounterVec,Total number of bytes transmitted to Kafka brokers"`
	Rx               int                        `json:"rx"                 kpromcol:"CounterVec,Total number of responses received from Kafka brokers"`
	RxBytes          int                        `json:"rx_bytes"           kpromcol:"CounterVec,Total number of bytes received from Kafka brokers"`
	Txmsgs           int                        `json:"txmsgs"             kpromcol:"CounterVec,Total number of messages transmitted (produced) to Kafka brokers"`
	TxmsgBytes       int                        `json:"txmsg_bytes"        kpromcol:"CounterVec,Total number of message bytes (including framing%2C such as per-Message framing and MessageSet/batch framing) transmitted to Kafka brokers"`
	Rxmsgs           int                        `json:"rxmsgs"             kpromcol:"CounterVec,Total number of messages consumed%2C not including ignored messages (due to offset%2C etc)%2C from Kafka brokers."`
	RxmsgBytes       int                        `json:"rxmsg_bytes"        kpromcol:"CounterVec,Total number of message bytes (including framing) received from Kafka brokers"`
	SimpleCnt        int                        `json:"simple_cnt"         kpromcol:"GaugeVec,Internal tracking of legacy vs new consumer API state"`
	MetadataCacheCnt int                        `json:"metadata_cache_cnt" kpromcol:"GaugeVec,Number of topics in the metadata cache."`
	Brokers          map[BrokerName]BrokerStats `json:"brokers"            kprommap:"brokers"`
	Topics           map[TopicName]TopicStats   `json:"topics"             kprommap:"topics"`
	Cgrp             CgrpStats                  `json:"cgrp"               kprompnt:"cgrp"` //Consumer group metrics.
	Eos              EosStats                   `json:"eos"                kprompnt:"eos"`  //EOS / Idempotent producer state and metrics.
}

type BrokerName string

type CgrpStats struct {
	State           string `json:"state"            kpromlbl:"state"` //Local consumer group handler's state.
	Stateage        int    `json:"stateage"         kpromcol:"GaugeVec,Time elapsed since last state change (milliseconds)."`
	JoinState       string `json:"join_state"       kpromlbl:"join_state"` //Local consumer group handler's join state.
	RebalanceAge    int    `json:"rebalance_age"    kpromcol:"GaugeVec,Time elapsed since last rebalance (assign or revoke) (milliseconds)."`
	RebalanceCnt    int    `json:"rebalance_cnt"    kpromcol:"CounterVec,Total number of rebalances (assign or revoke)."`
	RebalanceReason string `json:"rebalance_reason" kpromlbl:"rebalance_reason"` //Last rebalance reason, or empty string.
	AssignmentSize  int    `json:"assignment_size"  kpromcol:"GaugeVec,Current assignment's partition count."`
}

// BrokerStats is per broker statistics.
type BrokerStats struct {
	Name           string                             `json:"name"             kpromlbl:"name"`     //Broker hostname, port and broker id
	Nodeid         int                                `json:"nodeid"           kpromlbl:"nodeid"`   //Broker id (-1 for bootstraps)
	Nodename       string                             `json:"nodename"         kpromlbl:"nodename"` //Broker hostname
	Source         string                             `json:"source"           kpromlbl:"source"`   //Broker source (learned, configured, internal, logical)
	State          string                             `json:"state"            kpromlbl:"state"`    //Broker state (INIT, DOWN, CONNECT, AUTH, APIVERSION_QUERY, AUTH_HANDSHAKE, UP, UPDATE)
	Stateage       int                                `json:"stateage"         kpromcol:"GaugeVec,Time since last broker state change (microseconds)"`
	OutbufCnt      int                                `json:"outbuf_cnt"       kpromcol:"GaugeVec,Number of requests awaiting transmission to broker"`
	OutbufMsgCnt   int                                `json:"outbuf_msg_cnt"   kpromcol:"GaugeVec,Number of messages awaiting transmission to broker"`
	WaitrespCnt    int                                `json:"waitresp_cnt"     kpromcol:"GaugeVec,Number of requests in-flight to broker awaiting response"`
	WaitrespMsgCnt int                                `json:"waitresp_msg_cnt" kpromcol:"GaugeVec,Number of messages in-flight to broker awaiting response"`
	Tx             int                                `json:"tx"               kpromcol:"CounterVec,Total number of requests sent"`
	Txbytes        int                                `json:"txbytes"          kpromcol:"CounterVec,Total number of bytes sent"`
	Txerrs         int                                `json:"txerrs"           kpromcol:"CounterVec,Total number of transmission errors"`
	Txretries      int                                `json:"txretries"        kpromcol:"CounterVec,Total number of request retries"`
	Txidle         int                                `json:"txidle"           kpromcol:"CounterVec,Microseconds since last socket send (or -1 if no sends yet for current connection)."`
	ReqTimeouts    int                                `json:"req_timeouts"     kpromcol:"CounterVec,Total number of requests timed out"`
	Rx             int                                `json:"rx"               kpromcol:"CounterVec,Total number of responses received"`
	Rxbytes        int                                `json:"rxbytes"          kpromcol:"CounterVec,Total number of bytes received"`
	Rxerrs         int                                `json:"rxerrs"           kpromcol:"CounterVec,Total number of receive errors"`
	Rxcorriderrs   int                                `json:"rxcorriderrs"     kpromcol:"CounterVec,Total number of unmatched correlation ids in response (typically for timed out requests)"`
	Rxpartial      int                                `json:"rxpartial"        kpromcol:"CounterVec,Total number of partial MessageSets received. The broker may return partial responses if the full MessageSet could not fit in the remaining Fetch response size."`
	Rxidle         int                                `json:"rxidle"           kpromcol:"CounterVec,Microseconds since last socket receive (or -1 if no receives yet for current connection)."`
	Req            map[RequestName]RequestsSent       `json:"req"` //Value is the number of requests sent.
	ZbufGrow       int                                `json:"zbuf_grow"        kpromcol:"CounterVec,Total number of decompression buffer size increases"`
	BufGrow        int                                `json:"buf_grow"         kpromcol:"CounterVec,Total number of buffer size increases (deprecated%2C unused)"`
	Wakeups        int                                `json:"wakeups"          kpromcol:"CounterVec,Broker thread poll loop wakeups"`
	Connects       int                                `json:"connects"         kpromcol:"CounterVec,Number of connection attempts%2C including successful and failed%2C and name resolution failures."`
	Disconnects    int                                `json:"disconnects"      kpromcol:"CounterVec,Number of disconnects (triggered by broker%2C network%2C load-balancer%2C etc.)."`
	IntLatency     WindowStats                        `json:"int_latency"      kprompnt:"int_latency"`    //Internal producer queue latency in microseconds.
	OutbufLatency  WindowStats                        `json:"outbuf_latency"   kprompnt:"outbuf_latency"` //Internal request queue latency in microseconds. This is the time between a request is enqueued on the transmit (outbuf) queue and the time the request is written to the TCP socket. Additional buffering and latency may be incurred by the TCP stack and network.
	Rtt            WindowStats                        `json:"rtt"              kprompnt:"rtt"`            //Broker latency / round-trip time in microseconds.
	Throttle       WindowStats                        `json:"throttle"         kprompnt:"throttle"`       //Broker throttling time in milliseconds.
	Toppars        map[TopicAndPartition]TopparsStats `json:"toppars"`                                    //Partitions handled by this broker handle.
}

type EosStats struct {
	IdempState    string `json:"idemp_state"     kpromlbl:"idemp_state"` //Current idempotent producer id state.
	IdempStateage int    `json:"idemp_stateage"  kpromcol:"GaugeVec,Time elapsed since last idemp_state change (milliseconds)."`
	TxnState      string `json:"txn_state"       kpromlbl:"txn_state"` //Current transactional producer state.
	TxnStateage   int    `json:"txn_stateage"    kpromcol:"GaugeVec,Time elapsed since last txn_state change (milliseconds)."`
	TxnMayEnq     bool   `json:"txn_may_enq"`                            //Transactional state allows enqueuing (producing) new messages.
	ProducerId    int    `json:"producer_id"     kpromlbl:"producer_id"` //The currently assigned Producer ID (or -1).
	ProducerEpoch int    `json:"producer_epoch"`                         //The current epoch (or -1).
	EpochCnt      int    `json:"epoch_cnt"       kpromcol:"GaugeVec,The number of Producer ID assignments since start."`
}

type PartitionId int

type PartitionStats struct {
	Partition         int    `json:"partition"           kpromlbl:"partition"` //Partition Id (-1 for internal UA/UnAssigned partition)
	Broker            int    `json:"broker"              kpromlbl:"broker"`    //The id of the broker that messages are currently being fetched from
	Leader            int    `json:"leader"              kpromlbl:"leader"`    //Current leader broker id
	Desired           bool   `json:"desired"`                                  //Partition is explicitly desired by application
	Unknown           bool   `json:"unknown"`                                  //Partition not seen in topic metadata from broker
	MsgqCnt           int    `json:"msgq_cnt"            kpromcol:"GaugeVec,Number of messages waiting to be produced in first-level queue"`
	MsgqBytes         int    `json:"msgq_bytes"          kpromcol:"GaugeVec,Number of bytes in msgq_cnt"`
	XmitMsgqCnt       int    `json:"xmit_msgq_cnt"       kpromcol:"GaugeVec,Number of messages ready to be produced in transmit queue"`
	XmitMsgqBytes     int    `json:"xmit_msgq_bytes"     kpromcol:"GaugeVec,Number of bytes in xmit_msgq"`
	FetchqCnt         int    `json:"fetchq_cnt"          kpromcol:"GaugeVec,Number of pre-fetched messages in fetch queue"`
	FetchqSize        int    `json:"fetchq_size"         kpromcol:"GaugeVec,Bytes in fetchq"`
	FetchState        string `json:"fetch_state"         kpromlbl:"fetch_state"` //Consumer fetch state for this partition (none, stopping, stopped, offset-query, offset-wait, active).
	QueryOffset       int    `json:"query_offset"        kpromcol:"GaugeVec,Current/Last logical offset query"`
	NextOffset        int    `json:"next_offset"         kpromcol:"GaugeVec,Next offset to fetch"`
	AppOffset         int    `json:"app_offset"          kpromcol:"GaugeVec,Offset of last message passed to application + 1"`
	StoredOffset      int    `json:"stored_offset"       kpromcol:"GaugeVec,Offset to be committed"`
	CommittedOffset   int    `json:"committed_offset"    kpromcol:"GaugeVec,Last committed offset"`
	EofOffset         int    `json:"eof_offset"          kpromcol:"GaugeVec,Last PARTITION_EOF signaled offset"`
	LoOffset          int    `json:"lo_offset"           kpromcol:"GaugeVec,Partition's low watermark offset on broker"`
	HiOffset          int    `json:"hi_offset"           kpromcol:"GaugeVec,Partition's high watermark offset on broker"`
	LsOffset          int    `json:"ls_offset"           kpromcol:"GaugeVec,Partition's last stable offset on broker%2C or same as hi_offset is broker version is less than 0.11.0.0."`
	ConsumerLag       int    `json:"consumer_lag"        kpromcol:"GaugeVec,Difference between (hi_offset or ls_offset) and committed_offset). hi_offset is used when isolation.level=read_uncommitted%2C otherwise ls_offset."`
	ConsumerLagStored int    `json:"consumer_lag_stored" kpromcol:"GaugeVec,Difference between (hi_offset or ls_offset) and stored_offset. See consumer_lag and stored_offset."`
	Txmsgs            int    `json:"txmsgs"              kpromcol:"CounterVec,Total number of messages transmitted (produced)"`
	Txbytes           int    `json:"txbytes"             kpromcol:"CounterVec,Total number of bytes transmitted for txmsgs"`
	Rxmsgs            int    `json:"rxmsgs"              kpromcol:"CounterVec,Total number of messages consumed%2C not including ignored messages (due to offset%2C etc)."`
	Rxbytes           int    `json:"rxbytes"             kpromcol:"CounterVec,Total number of bytes received for rxmsgs"`
	Msgs              int    `json:"msgs"                kpromcol:"CounterVec,Total number of messages received (consumer%2C same as rxmsgs)%2C or total number of messages produced (possibly not yet transmitted) (producer)."`
	RxVerDrops        int    `json:"rx_ver_drops"        kpromcol:"CounterVec,Dropped outdated messages"`
	MsgsInflight      int    `json:"msgs_inflight"       kpromcol:"GaugeVec,Current number of messages in-flight to/from broker"`
	NextAckSeq        int    `json:"next_ack_seq"        kpromcol:"GaugeVec,Next expected acked sequence (idempotent producer)"`
	NextErrSeq        int    `json:"next_err_seq"        kpromcol:"GaugeVec,Next expected errored sequence (idempotent producer)"`
	AckedMsgid        int    `json:"acked_msgid"` //Last acked internal message id (idempotent producer
}

type RequestName string
type RequestsSent int

type TopicName string
type TopicAndPartition string

type TopicStats struct {
	Topic       string                         `json:"topic"        kpromlbl:"topic"` //Topic name
	Age         int                            `json:"age"          kpromcol:"GaugeVec,Age of client's topic object (milliseconds)"`
	MetadataAge int                            `json:"metadata_age" kpromcol:"GaugeVec,Age of metadata from broker for this topic (milliseconds)"`
	Batchsize   WindowStats                    `json:"batchsize"    kprompnt:"batchsize"` //Batch sizes in bytes.
	Batchcnt    WindowStats                    `json:"batchcnt"     kprompnt:"batchcnt"`  //Batch message counts.
	Partitions  map[PartitionId]PartitionStats `json:"partitions"   kprommap:"partitions"`
}
type TopparsStats struct {
	Topic     string `json:"topic"`     //Topic name
	Partition int    `json:"partition"` //Partition id
}

// WindowStats has rolling window statistics. The values are in microseconds unless otherwise stated.
type WindowStats struct {
	Min        int `json:"min"        kpromcol:"GaugeVec,Smallest value"`
	Max        int `json:"max"        kpromcol:"GaugeVec,Largest value"`
	Avg        int `json:"avg"        kpromcol:"GaugeVec,Average value"`
	Sum        int `json:"sum"        kpromcol:"GaugeVec,Sum of values"`
	Cnt        int `json:"cnt"        kpromcol:"GaugeVec,Number of values sampled"`
	Stddev     int `json:"stddev"     kpromcol:"GaugeVec,Standard deviation (based on histogram)"`
	Hdrsize    int `json:"hdrsize"    kpromcol:"GaugeVec,Memory size of Hdr Histogram"`
	P50        int `json:"p50"        kpromcol:"GaugeVec,50th percentile"`
	P75        int `json:"p75"        kpromcol:"GaugeVec,75th percentile"`
	P90        int `json:"p90"        kpromcol:"GaugeVec,90th percentile"`
	P95        int `json:"p95"        kpromcol:"GaugeVec,95th percentile"`
	P99        int `json:"p99"        kpromcol:"GaugeVec,99th percentile"`
	P99_99     int `json:"p99_99"     kpromcol:"GaugeVec,99.99th percentile"`
	Outofrange int `json:"outofrange" kpromcol:"GaugeVec,Values skipped due to out of histogram range"`
}
