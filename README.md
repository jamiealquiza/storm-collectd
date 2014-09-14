storm-collectd
==============

Apache Storm bolt metrics collectd plugin

### How it works

The plugin reports bolt metrics using the Storm API. It automatically finds all topologies, then fetches info for all bolts running under each topology.

The plugin honors the `COLLECTD_HOSTNAME` variable passed by collectd and will failover to detecting the system hostname if no value is passed.

### Example setup

Build / place binary in your plugin path. Example:
<pre>
go build storm-collectd.go
cp storm-collectd /opt/collectd/lib/collectd/
chown collectd:collectd /opt/collectd/lib/collectd/storm-collectd
</pre>

In collectd.conf:
<pre>
LoadPlugin exec
&lt;Plugin exec&gt;
  Exec "collectd" "/opt/collectd/lib/collectd/storm-collectd"
&lt;/Plugin&gt;
</pre>

### Example output

<pre>
./storm-collectd 
PUTVAL some-server/storm/gauge-some_bolt-capacity N:0.000
PUTVAL some-server/storm/gauge-some_bolt-acked N:0
PUTVAL some-server/storm/gauge-some_bolt-executeLatency N:0.000
PUTVAL some-server/storm/gauge-some_bolt-executed N:0
PUTVAL some-server/storm/gauge-some_bolt-executors N:1
PUTVAL some-server/storm/gauge-some_bolt-emitted N:0
PUTVAL some-server/storm/gauge-some_bolt-tasks N:1
PUTVAL some-server/storm/gauge-some_bolt-transferred N:0
PUTVAL some-server/storm/gauge-some_bolt-failed N:0
PUTVAL some-server/storm/gauge-some_bolt-processLatency N:0.000
PUTVAL some-server/storm/gauge-another_bolt-transferred N:0
PUTVAL some-server/storm/gauge-another_bolt-acked N:3974680
PUTVAL some-server/storm/gauge-another_bolt-capacity N:0.003
PUTVAL some-server/storm/gauge-another_bolt-tasks N:2
PUTVAL some-server/storm/gauge-another_bolt-failed N:0
PUTVAL some-server/storm/gauge-another_bolt-emitted N:0
PUTVAL some-server/storm/gauge-another_bolt-executed N:3974680
PUTVAL some-server/storm/gauge-another_bolt-processLatency N:0.034
PUTVAL some-server/storm/gauge-another_bolt-executors N:2
PUTVAL some-server/storm/gauge-another_bolt-executeLatency N:0.038
PUTVAL some-server/storm/gauge-even_more_bolts-tasks N:4
PUTVAL some-server/storm/gauge-even_more_bolts-emitted N:0
PUTVAL some-server/storm/gauge-even_more_bolts-acked N:0
PUTVAL some-server/storm/gauge-even_more_bolts-executed N:60
PUTVAL some-server/storm/gauge-even_more_bolts-transferred N:0
PUTVAL some-server/storm/gauge-even_more_bolts-capacity N:0.000
PUTVAL some-server/storm/gauge-even_more_bolts-processLatency N:0.000
PUTVAL some-server/storm/gauge-even_more_bolts-failed N:0
PUTVAL some-server/storm/gauge-even_more_bolts-executors N:4
PUTVAL some-server/storm/gauge-even_more_bolts-executeLatency N:0.000

...etc
</pre>
