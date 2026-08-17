package view

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1"/>
<title>PiMon · Infrastructure Monitoring</title>
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0f1419;--bg-elevated:#161b22;--bg-panel:#1a2029;--bg-hover:#20262f;
  --border:#262e38;--border-light:#2e3744;
  --text:#d6d9dd;--text-bright:#f4f5f7;--text-dim:#8b95a1;--text-faint:#5a6473;
  --purple:#8c5cf6;--purple-dim:#6938b8;--purple-bg:rgba(140,92,246,.12);
  --green:#4ce0a0;--green-bg:rgba(76,224,160,.12);
  --red:#e8485e;--red-bg:rgba(232,72,94,.12);
  --amber:#f5a623;--amber-bg:rgba(245,166,35,.12);
  --blue:#3ea6ff;--blue-bg:rgba(62,166,255,.12);
  --cyan:#3bc9db;--pink:#f06595;
  --mono:'SF Mono',ui-monospace,'Cascadia Mono',Consolas,monospace;
  --sans:-apple-system,BlinkMacSystemFont,'Segoe UI',system-ui,sans-serif;
  --radius:4px;
}
html{background:var(--bg);color:var(--text);font-family:var(--sans);font-size:13px;line-height:1.45;-webkit-font-smoothing:antialiased;}
body{min-height:100vh;display:flex;flex-direction:column;}
::-webkit-scrollbar{width:8px;height:8px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:var(--border-light);border-radius:4px}
::-webkit-scrollbar-thumb:hover{background:var(--text-faint)}
.topbar{background:var(--bg-elevated);border-bottom:1px solid var(--border);height:48px;display:flex;align-items:center;justify-content:space-between;padding:0 16px;position:sticky;top:0;z-index:100;}
.topbar-left{display:flex;align-items:center;gap:14px;}
.logo{display:flex;align-items:center;gap:8px;font-weight:600;font-size:14px;color:var(--text-bright);letter-spacing:-.01em;}
.logo .mark{width:22px;height:22px;background:var(--purple);border-radius:4px;display:flex;align-items:center;justify-content:center;font-size:12px;color:#fff;font-weight:700;}
.crumb{color:var(--text-faint);font-size:12px;}
.crumb-sep{color:var(--text-faint);margin:0 6px;}
.crumb-host{color:var(--text-dim);font-size:12px;font-family:var(--mono);}
.topbar-right{display:flex;align-items:center;gap:10px;}
.live-pill{display:flex;align-items:center;gap:6px;font-size:11px;color:var(--green);background:var(--green-bg);padding:4px 10px;border-radius:999px;font-weight:500;}
.live-dot{width:6px;height:6px;border-radius:50%;background:var(--green);animation:blink 1.6s ease-in-out infinite;}
@keyframes blink{0%,100%{opacity:1}50%{opacity:.35}}
.clock{font-family:var(--mono);font-size:11px;color:var(--text-faint);}
.subheader{background:var(--bg);border-bottom:1px solid var(--border);padding:14px 20px;display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:12px;}
.subheader-title{display:flex;align-items:center;gap:10px;}
.subheader-title h1{font-size:16px;font-weight:600;color:var(--text-bright);letter-spacing:-.01em;}
.tag{font-size:10.5px;font-family:var(--mono);background:var(--bg-panel);border:1px solid var(--border);color:var(--text-dim);padding:2px 7px;border-radius:3px;}
.tag-purple{color:var(--purple);border-color:var(--purple-dim);background:var(--purple-bg);}
.subheader-meta{display:flex;align-items:center;gap:16px;}
.meta-item{display:flex;flex-direction:column;gap:1px;}
.meta-label{font-size:9.5px;color:var(--text-faint);text-transform:uppercase;letter-spacing:.06em;}
.meta-val{font-size:12px;color:var(--text);font-family:var(--mono);}
.main{flex:1;padding:16px 20px;max-width:1480px;margin:0 auto;width:100%;display:flex;flex-direction:column;gap:14px;}
.row{display:grid;gap:12px;}
.row-4{grid-template-columns:repeat(4,1fr)}
.row-3{grid-template-columns:repeat(3,1fr)}
.row-2{grid-template-columns:repeat(2,1fr)}
.row-12-8{grid-template-columns:1.5fr 1fr}
@media(max-width:1100px){.row-4{grid-template-columns:repeat(2,1fr)}.row-3{grid-template-columns:repeat(2,1fr)}}
@media(max-width:760px){.row-4,.row-3,.row-2,.row-12-8{grid-template-columns:1fr}}
.panel{background:var(--bg-panel);border:1px solid var(--border);border-radius:var(--radius);overflow:hidden;}
.panel-head{padding:9px 12px;display:flex;align-items:center;justify-content:space-between;border-bottom:1px solid var(--border);}
.panel-title{font-size:11px;font-weight:600;color:var(--text-dim);text-transform:uppercase;letter-spacing:.05em;display:flex;align-items:center;gap:6px;}
.panel-title .ind{width:5px;height:5px;border-radius:50%;}
.panel-actions{display:flex;align-items:center;gap:8px;}
.panel-body{padding:12px;}
.panel-body-tight{padding:0;}
.metric-widget .panel-body{padding:14px 14px 10px;}
.metric-top{display:flex;align-items:baseline;justify-content:space-between;margin-bottom:2px;}
.metric-value{font-family:var(--mono);font-size:28px;font-weight:600;color:var(--text-bright);letter-spacing:-.02em;}
.metric-unit{font-size:13px;color:var(--text-faint);font-weight:400;margin-left:2px;}
.metric-delta{font-size:11px;font-family:var(--mono);padding:1px 6px;border-radius:3px;}
.delta-up{color:var(--red);background:var(--red-bg);}
.delta-down{color:var(--green);background:var(--green-bg);}
.delta-flat{color:var(--text-faint);background:var(--bg-elevated);}
.metric-label{font-size:11px;color:var(--text-faint);margin-bottom:8px;}
.spark-wrap{height:42px;margin:0 -4px;}
.spark-wrap canvas{display:block;width:100%;height:100%;}
.status-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(76px,1fr));gap:6px;}
.status-cell{background:var(--bg-elevated);border:1px solid var(--border);border-radius:3px;padding:7px 8px;text-align:center;position:relative;overflow:hidden;}
.status-cell .sc-label{font-size:9px;color:var(--text-faint);margin-bottom:3px;font-family:var(--mono);}
.status-cell .sc-val{font-size:14px;font-weight:600;font-family:var(--mono);line-height:1;}
.status-cell .sc-bar{position:absolute;bottom:0;left:0;height:2px;transition:width .5s ease;}
.dtable{width:100%;border-collapse:collapse;font-size:12px;}
.dtable thead th{text-align:left;font-size:10px;font-weight:600;text-transform:uppercase;letter-spacing:.05em;color:var(--text-faint);padding:7px 12px;border-bottom:1px solid var(--border);background:var(--bg-elevated);position:sticky;top:0;}
.dtable tbody td{padding:6px 12px;border-bottom:1px solid var(--border);font-family:var(--mono);color:var(--text);white-space:nowrap;}
.dtable tbody tr:last-child td{border-bottom:none;}
.dtable tbody tr:hover td{background:var(--bg-hover);}
.dtable .col-name{font-family:var(--sans);max-width:140px;overflow:hidden;text-overflow:ellipsis;color:var(--text-bright);}
.dtable-scroll{max-height:280px;overflow-y:auto;}
.pill{display:inline-flex;align-items:center;gap:4px;font-size:10px;font-weight:600;padding:2px 7px;border-radius:3px;font-family:var(--mono);text-transform:uppercase;letter-spacing:.03em;}
.pill-green{color:var(--green);background:var(--green-bg);}
.pill-red{color:var(--red);background:var(--red-bg);}
.pill-amber{color:var(--amber);background:var(--amber-bg);}
.pill-blue{color:var(--blue);background:var(--blue-bg);}
.pill-gray{color:var(--text-faint);background:var(--bg-elevated);}
.ubar-row{display:flex;align-items:center;gap:10px;padding:6px 0;}
.ubar-row:not(:last-child){border-bottom:1px solid var(--border);}
.ubar-label{font-size:11px;color:var(--text-dim);min-width:80px;font-family:var(--mono);}
.ubar-track{flex:1;height:5px;background:var(--bg-elevated);border-radius:3px;overflow:hidden;}
.ubar-fill{height:100%;border-radius:3px;transition:width .6s ease;}
.ubar-val{font-size:11px;font-family:var(--mono);color:var(--text);min-width:96px;text-align:right;}
.net-row{display:flex;align-items:center;gap:10px;padding:7px 0;border-bottom:1px solid var(--border);}
.net-row:last-child{border-bottom:none;}
.net-iface{font-family:var(--mono);font-size:12px;font-weight:600;color:var(--text-bright);min-width:64px;}
.net-stat{display:flex;align-items:center;gap:5px;font-size:11px;font-family:var(--mono);}
.net-stat .arrow{font-size:9px;}
.section-label{font-size:11px;font-weight:600;text-transform:uppercase;letter-spacing:.07em;color:var(--text-faint);display:flex;align-items:center;gap:8px;padding:2px 0;}
.section-label::after{content:'';flex:1;height:1px;background:var(--border);}
</style>
</head>
<body>
  <header class="topbar">
    <div class="topbar-left">
      <div class="logo"><span class="mark">P</span>PiMon</div>
      <div class="crumb">Infrastructure</div>
      <div class="crumb-sep">›</div>
      <div class="crumb-host" id="hostName">loading...</div>
    </div>
    <div class="topbar-right">
      <div class="live-pill"><span class="live-dot"></span>LIVE</div>
      <div class="clock" id="clock">--:--:--</div>
    </div>
  </header>

  <div class="subheader">
    <div class="subheader-title">
      <h1>Raspberry Pi Monitor</h1>
      <span class="tag tag-purple" id="piModel">--</span>
    </div>
    <div class="subheader-meta">
      <div class="meta-item"><div class="meta-label">Uptime</div><div class="meta-val" id="uptime">--</div></div>
      <div class="meta-item"><div class="meta-label">Kernel</div><div class="meta-val" id="kernel">--</div></div>
      <div class="meta-item"><div class="meta-label">Go</div><div class="meta-val" id="goVersion">--</div></div>
    </div>
  </div>

  <main class="main">
    <section class="row row-4" id="metricRow"></section>

    <section class="row row-12-8">
      <div class="panel">
        <div class="panel-head"><div class="panel-title"><span class="ind" style="background:#4ce0a0"></span>CPU usage</div></div>
        <div class="panel-body">
          <div class="section-label">Overall</div>
          <div class="ubar-row"><div class="ubar-label">System</div><div class="ubar-track"><div class="ubar-fill" id="cpuFill" style="width:0%;background:linear-gradient(90deg,#8c5cf6,#4ce0a0)"></div></div><div class="ubar-val" id="cpuPct">0%</div></div>
          <div class="section-label" style="margin-top:12px;">Core load</div>
          <div id="coreGrid" class="status-grid"></div>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head"><div class="panel-title"><span class="ind" style="background:#3ea6ff"></span>System</div></div>
        <div class="panel-body">
          <div id="statusGrid" class="status-grid"></div>
        </div>
      </div>
    </section>

    <section class="row row-2">
      <div class="panel">
        <div class="panel-head"><div class="panel-title"><span class="ind" style="background:#f5a623"></span>Memory</div></div>
        <div class="panel-body">
          <div id="memBars"></div>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head"><div class="panel-title"><span class="ind" style="background:#e8485e"></span>Storage</div></div>
        <div class="panel-body">
          <div id="diskTable"></div>
        </div>
      </div>
    </section>

    <section class="row row-2">
      <div class="panel">
        <div class="panel-head"><div class="panel-title"><span class="ind" style="background:#3bc9db"></span>Network</div></div>
        <div class="panel-body">
          <div id="netList"></div>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head"><div class="panel-title"><span class="ind" style="background:#f06595"></span>Processes</div></div>
        <div class="panel-body">
          <div class="dtable-scroll">
            <table class="dtable">
              <thead>
                <tr><th>PID</th><th>Name</th><th>CPU</th><th>Mem</th><th>User</th></tr>
              </thead>
              <tbody id="procTable"></tbody>
            </table>
          </div>
        </div>
      </div>
    </section>
  </main>

  <script>
    const metricDefs = [
      { key: 'cpu', label: 'CPU', unit: '%', color: '#8c5cf6' },
      { key: 'memory', label: 'MEM', unit: '%', color: '#4ce0a0' },
      { key: 'disk', label: 'DISK', unit: '%', color: '#f5a623' },
      { key: 'net', label: 'NET', unit: 'KB/s', color: '#3ea6ff' }
    ];

    let cpuHistory = [];
    let memHistory = [];

    function fmtNumber(v, digits = 1) {
      return Number(v || 0).toFixed(digits);
    }

    function pctColor(v) {
      if (v >= 80) return '#e8485e';
      if (v >= 60) return '#f5a623';
      return '#4ce0a0';
    }

    function renderMetricWidget(data) {
      const row = document.getElementById('metricRow');
      row.innerHTML = '';
      const cards = metricDefs.map(function(def) {
        let value = 0;
        let label = '0';
        switch (def.key) {
          case 'cpu': value = data.cpu && data.cpu.usage_percent !== undefined ? data.cpu.usage_percent : 0; label = fmtNumber(value, 1); break;
          case 'memory': value = data.memory && data.memory.used_percent !== undefined ? data.memory.used_percent : 0; label = fmtNumber(value, 1); break;
          case 'disk': value = data.disks && data.disks[0] && data.disks[0].used_percent !== undefined ? data.disks[0].used_percent : 0; label = fmtNumber(value, 1); break;
          case 'net': {
            const total = (data.network || []).reduce(function(sum, i) {
              return sum + (i.rx_rate_kbps || 0) + (i.tx_rate_kbps || 0);
            }, 0);
            value = total / 1024;
            label = fmtNumber(value, 1);
          } break;
        }
        const delta = def.key === 'net' ? 'LIVE' : value > 0 ? 'active' : 'idle';
        return '<div class="panel metric-widget">' +
          '<div class="panel-head">' +
          '  <div class="panel-title"><span class="ind" style="background:' + def.color + '"></span>' + def.label + '</div>' +
          '  <div class="pill ' + (value > 75 ? 'pill-red' : value > 45 ? 'pill-amber' : 'pill-green') + '">' + delta + '</div>' +
          '</div>' +
          '<div class="panel-body">' +
          '  <div class="metric-top">' +
          '    <div class="metric-value">' + label + '<span class="metric-unit">' + def.unit + '</span></div>' +
          '  </div>' +
          '  <div class="metric-label">' + (def.key === 'net' ? 'throughput' : 'current usage') + '</div>' +
          '  <div class="spark-wrap"><canvas data-key="' + def.key + '" width="250" height="42"></canvas></div>' +
          '</div>' +
          '</div>';
      });
      row.innerHTML = cards.join('');
      renderSparkline('cpu', cpuHistory);
      renderSparkline('memory', memHistory);
      renderSparkline('disk', data.disks ? data.disks.map(function(d) { return d.used_percent || 0; }) : []);
      renderSparkline('net', (data.network || []).map(function(i) { return (i.rx_rate_kbps || 0) + (i.tx_rate_kbps || 0); }));
    }

    function renderSparkline(key, arr) {
      const canvas = document.querySelector('canvas[data-key="' + key + '"]');
      if (!canvas) return;
      const ctx = canvas.getContext('2d');
      const values = Array.isArray(arr) && arr.length ? arr : [0, 0, 0, 0, 0, 0];
      const max = Math.max.apply(null, values.concat([1]));
      const min = Math.min.apply(null, values.concat([0]));
      const w = canvas.width;
      const h = canvas.height;
      ctx.clearRect(0, 0, w, h);
      ctx.beginPath();
      values.forEach(function(v, i) {
        const x = (i / Math.max(values.length - 1, 1)) * w;
        const y = h - ((v - min) / Math.max(max - min, 1)) * (h - 4) - 2;
        if (i === 0) ctx.moveTo(x, y); else ctx.lineTo(x, y);
      });
      ctx.strokeStyle = key === 'cpu' ? '#8c5cf6' : key === 'memory' ? '#4ce0a0' : key === 'disk' ? '#f5a623' : '#3ea6ff';
      ctx.lineWidth = 2;
      ctx.stroke();
    }

    function renderCoreGrid(data) {
      const container = document.getElementById('coreGrid');
      const cores = data.cpu && data.cpu.core_usages ? data.cpu.core_usages : [];
      container.innerHTML = cores.map(function(v, i) {
        const color = pctColor(v);
        return '<div class="status-cell"><div class="sc-label">cpu' + i + '</div><div class="sc-val" style="color:' + color + '">' + fmtNumber(v, 0) + '%</div><div class="sc-bar" style="width:' + Math.min(v, 100) + '%;background:' + color + ';"></div></div>';
      }).join('') || '<div class="status-cell"><div class="sc-label">cpu0</div><div class="sc-val">0%</div></div>';
      const cpuFill = document.getElementById('cpuFill');
      const cpuPct = document.getElementById('cpuPct');
      const pct = data.cpu && data.cpu.usage_percent !== undefined ? data.cpu.usage_percent : 0;
      cpuFill.style.width = Math.min(pct, 100) + '%';
      cpuPct.textContent = fmtNumber(pct, 1) + '%';
    }

    function renderStatusGrid(data) {
      const container = document.getElementById('statusGrid');
      const items = [
        ['Host', data.system && data.system.hostname ? data.system.hostname : '--'],
        ['OS', data.system && data.system.os ? data.system.os : '--'],
        ['Kernel', data.system && data.system.kernel ? data.system.kernel : '--'],
        ['CPU temp', fmtNumber(data.cpu && data.cpu.temp_celsius !== undefined ? data.cpu.temp_celsius : 0, 1) + '°C'],
        ['Freq', fmtNumber(data.cpu && data.cpu.freq_mhz !== undefined ? data.cpu.freq_mhz : 0, 0) + ' MHz'],
        ['Load', fmtNumber(data.cpu && data.cpu.load_avg_1 !== undefined ? data.cpu.load_avg_1 : 0, 2) + ' / ' + fmtNumber(data.cpu && data.cpu.load_avg_5 !== undefined ? data.cpu.load_avg_5 : 0, 2)],
        ['GPU', data.gpu && data.gpu.available ? fmtNumber(data.gpu.temp_celsius || 0, 1) + '°C' : 'n/a'],
        ['Throttled', data.cpu && data.cpu.throttled ? 'yes' : 'no']
      ];
      container.innerHTML = items.map(function(entry) {
        return '<div class="status-cell"><div class="sc-label">' + entry[0] + '</div><div class="sc-val">' + entry[1] + '</div></div>';
      }).join('');
    }

    function renderMem(data) {
      const total = data.memory && data.memory.total_mb !== undefined ? data.memory.total_mb : 0;
      const used = data.memory && data.memory.used_mb !== undefined ? data.memory.used_mb : 0;
      const pct = data.memory && data.memory.used_percent !== undefined ? data.memory.used_percent : 0;
      const lines = [
        ['Used', used, pct, '#4ce0a0'],
        ['Cached', data.memory && data.memory.cached_mb !== undefined ? data.memory.cached_mb : 0, Math.min(((data.memory && data.memory.cached_mb !== undefined ? data.memory.cached_mb : 0) / total) * 100, 100), '#3ea6ff'],
        ['Swap', data.memory && data.memory.swap_used_mb !== undefined ? data.memory.swap_used_mb : 0, data.memory && data.memory.swap_percent !== undefined ? data.memory.swap_percent : 0, '#8c5cf6']
      ];
      const html = lines.map(function(entry) {
        return '<div class="ubar-row"><div class="ubar-label">' + entry[0] + '</div><div class="ubar-track"><div class="ubar-fill" style="width:' + Math.min(entry[2] || 0, 100) + '%;background:' + entry[3] + ';"></div></div><div class="ubar-val">' + fmtNumber(entry[1], 1) + ' MB</div></div>';
      }).join('');
      document.getElementById('memBars').innerHTML = html + '<div class="ubar-row"><div class="ubar-label">Total</div><div class="ubar-track"><div class="ubar-fill" style="width:100%;background:#d6d9dd"></div></div><div class="ubar-val">' + fmtNumber(total, 0) + ' MB</div></div>';
    }

    function renderDisks(data) {
      const disks = (data.disks || []).slice(0, 4);
      document.getElementById('diskTable').innerHTML = disks.length ? disks.map(function(d) {
        return '<div class="ubar-row"><div class="ubar-label">' + (d.device || d.path) + '</div><div class="ubar-track"><div class="ubar-fill" style="width:' + Math.min(d.used_percent || 0, 100) + '%;background:' + pctColor(d.used_percent || 0) + ';"></div></div><div class="ubar-val">' + fmtNumber(d.used_percent || 0, 1) + '%</div></div>';
      }).join('') : '<div class="ubar-row"><div class="ubar-label">No disks</div><div class="ubar-val">--</div></div>';
    }

    function renderNetwork(data) {
      const net = data.network || [];
      const html = net.slice(0, 6).map(function(iface) {
        return '<div class="net-row"><div class="net-iface">' + iface.name + '</div><div class="net-stat"><span class="arrow">↓</span>' + fmtNumber(iface.rx_rate_kbps || 0, 1) + ' KB/s</div><div class="net-stat"><span class="arrow">↑</span>' + fmtNumber(iface.tx_rate_kbps || 0, 1) + ' KB/s</div></div>';
      }).join('');
      document.getElementById('netList').innerHTML = html || '<div class="net-row"><div class="net-iface">none</div></div>';
    }

    function renderProcesses(data) {
      const rows = (data.processes || []).slice(0, 15).map(function(p) {
        return '<tr><td>' + p.pid + '</td><td class="col-name">' + p.name + '</td><td>' + fmtNumber(p.cpu_percent || 0, 1) + '%</td><td>' + fmtNumber(p.mem_mb || 0, 1) + ' MB</td><td>' + (p.user || '--') + '</td></tr>';
      }).join('');
      document.getElementById('procTable').innerHTML = rows;
    }

    function updateClock() {
      document.getElementById('clock').textContent = new Date().toLocaleTimeString();
    }

    async function fetchSnapshot() {
      try {
        const res = await fetch('/api/snapshot');
        if (!res.ok) throw new Error('snapshot unavailable');
        const data = await res.json();
        document.getElementById('hostName').textContent = data.system && data.system.hostname ? data.system.hostname : 'unknown';
        document.getElementById('piModel').textContent = data.system && data.system.pi_model ? data.system.pi_model : 'unknown';
        document.getElementById('uptime').textContent = data.system && data.system.uptime ? data.system.uptime : '--';
        document.getElementById('kernel').textContent = data.system && data.system.kernel ? data.system.kernel : '--';
        document.getElementById('goVersion').textContent = data.system && data.system.go_version ? data.system.go_version : '--';

        renderMetricWidget(data);
        renderCoreGrid(data);
        renderStatusGrid(data);
        renderMem(data);
        renderDisks(data);
        renderNetwork(data);
        renderProcesses(data);

        cpuHistory = Array.isArray(cpuHistory) ? cpuHistory.concat([data.cpu && data.cpu.usage_percent !== undefined ? data.cpu.usage_percent : 0]).slice(-24) : [data.cpu && data.cpu.usage_percent !== undefined ? data.cpu.usage_percent : 0];
        memHistory = Array.isArray(memHistory) ? memHistory.concat([data.memory && data.memory.used_percent !== undefined ? data.memory.used_percent : 0]).slice(-24) : [data.memory && data.memory.used_percent !== undefined ? data.memory.used_percent : 0];
      } catch (e) {
        console.error(e);
      }
    }

    async function fetchHistory() {
      try {
        const res = await fetch('/api/history');
        if (!res.ok) return;
        const data = await res.json();
        cpuHistory = Array.isArray(data.cpu) ? data.cpu : cpuHistory;
        memHistory = Array.isArray(data.mem) ? data.mem : memHistory;
      } catch (e) {
        console.error(e);
      }
    }

    updateClock();
    fetchHistory();
    fetchSnapshot();
    setInterval(updateClock, 1000);
    setInterval(fetchSnapshot, 2000);
    setInterval(fetchHistory, 15000);
  </script>
</body>
</html>`
