package webapp

const dashboardPage = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Neferpitool Web Console</title>
  <style>
    @keyframes rise {
      from { opacity: 0; transform: translateY(10px); }
      to { opacity: 1; transform: translateY(0); }
    }
    :root {
      --bg: #0b1020;
      --panel: #131b2e;
      --panel-2: #0f172a;
      --line: #24314d;
      --text: #dbe6ff;
      --muted: #8ea1c8;
      --accent: #1dd3b0;
      --accent-2: #f59e0b;
      --danger: #ef4444;
      --ok: #22c55e;
      --radius: 14px;
      --shadow: 0 10px 30px rgba(0, 0, 0, .35);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      color: var(--text);
      font-family: "Space Grotesk", "Manrope", "Noto Sans", sans-serif;
      min-height: 100vh;
      background:
        radial-gradient(1200px 650px at 8% -4%, #1f2a46 0%, transparent 55%),
        radial-gradient(900px 600px at 105% 110%, #193550 0%, transparent 50%),
        var(--bg);
    }
    .wrap { max-width: 1500px; margin: 0 auto; padding: 18px; }
    .top {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 14px;
    }
    .brand {
      display: flex;
      align-items: center;
      gap: 12px;
    }
    .brand h1 { margin: 0; font-size: 25px; line-height: 1.05; letter-spacing: .3px; }
    .brand p { margin: 2px 0 0; color: var(--muted); font-size: 13px; }
    .actions { display: flex; gap: 8px; flex-wrap: wrap; }
    .layout {
      display: grid;
      grid-template-columns: 360px 1fr;
      gap: 14px;
    }
    .card {
      border: 1px solid var(--line);
      border-radius: var(--radius);
      background: linear-gradient(180deg, rgba(255,255,255,.02), rgba(255,255,255,.00)), var(--panel);
      box-shadow: var(--shadow);
      overflow: hidden;
      animation: rise .42s ease both;
    }
    .card h2 {
      margin: 0;
      padding: 12px 14px;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: .08em;
      border-bottom: 1px solid var(--line);
      background: var(--panel-2);
      color: #cde0ff;
    }
    .content { padding: 12px; }
    .stack { display: grid; gap: 10px; }
    .row { display: flex; gap: 8px; flex-wrap: wrap; }
    .row > * { flex: 1; min-width: 0; }
    .col-toggle {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      flex: 0 0 auto;
      border: 1px solid #2a3959;
      background: #0e1628;
      border-radius: 10px;
      padding: 6px 9px;
      font-size: 12px;
      color: #bcd2ff;
    }
    .col-toggle input {
      flex: 0 0 auto;
      width: auto;
      margin: 0;
      padding: 0;
    }
    .algo-filter {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      width: 100%;
      flex: 0 0 100%;
      border: 1px solid #2a3959;
      background: #0e1628;
      border-radius: 10px;
      padding: 6px 9px;
      font-size: 12px;
      color: #bcd2ff;
    }
    .algo-filter input {
      flex: 0 0 auto;
      width: auto;
      margin: 0;
      padding: 0;
    }
    .sort-btn {
      padding: 4px 8px;
      border-radius: 8px;
      border: 1px solid transparent;
      background: transparent;
      color: #bcd2ff;
      font-weight: 700;
    }
    .sort-btn.is-active {
      border-color: #2a3959;
      background: #14213b;
      color: #dbe6ff;
    }
    .multi-select {
      position: relative;
      width: 100%;
      max-width: 420px;
    }
    .multi-select-btn {
      width: 100%;
      text-align: left;
      display: inline-flex;
      align-items: center;
      justify-content: space-between;
      gap: 8px;
    }
    .multi-select-btn::after {
      content: "▾";
      color: #8ea1c8;
      font-size: 12px;
      transition: transform .15s ease;
    }
    .multi-select.is-open .multi-select-btn::after {
      transform: rotate(180deg);
    }
    .multi-select-menu {
      position: absolute;
      top: calc(100% + 6px);
      left: 0;
      width: 100%;
      max-height: 220px;
      overflow: auto;
      z-index: 20;
      border: 1px solid #2a3959;
      border-radius: 10px;
      background: #0a1528;
      box-shadow: 0 10px 30px rgba(0, 0, 0, .35);
      padding: 8px;
      display: none;
    }
    .multi-select.is-open .multi-select-menu {
      display: block;
    }
    .domains { max-height: 350px; overflow: auto; border-top: 1px dashed var(--line); padding-top: 8px; }
    .domain {
      display: flex; align-items: center; justify-content: space-between; gap: 8px;
      border-bottom: 1px dotted var(--line); padding: 8px 0;
    }
    .domain .name {
      max-width: 130px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
      font-weight: 700; font-size: 13px; color: #cde0ff;
      font-family: "IBM Plex Mono", "JetBrains Mono", Menlo, Consolas, monospace;
      cursor: pointer;
    }
    .domain .btns { display: flex; gap: 5px; }
    .domain .btns button { padding: 5px 7px; font-size: 12px; }
    input, textarea, button, select {
      color: var(--text);
      background: #0e1628;
      border: 1px solid #2a3959;
      border-radius: 10px;
      padding: 9px 10px;
      font: inherit;
    }
    input::placeholder, textarea::placeholder { color: #6f84ae; }
    textarea { min-height: 110px; resize: vertical; }
    button { cursor: pointer; transition: .15s ease; }
    button:hover { transform: translateY(-1px); filter: brightness(1.08); }
    .b-ok { border-color: #14532d; background: #0b2d1d; color: #9bf7c0; }
    .b-warn { border-color: #7c4a03; background: #3f2d09; color: #fcd68a; }
    .b-danger { border-color: #7f1d1d; background: #3f1212; color: #fecaca; }
    .b-main { border-color: #0f766e; background: #0e3834; color: #9ef3e7; }
    .muted { color: var(--muted); font-size: 13px; }
    .busy {
      border: 1px solid #7c4a03;
      background: #2a1f0b;
      color: #fcd68a;
      border-radius: 10px;
      padding: 8px 10px;
      font-size: 13px;
    }
    .progress-wrap {
      border: 1px solid #214266;
      background: #0a1528;
      border-radius: 12px;
      padding: 10px;
    }
    .progress-head {
      display: flex;
      justify-content: space-between;
      gap: 8px;
      font-size: 13px;
      color: var(--muted);
      margin-bottom: 6px;
    }
    .progress-track {
      width: 100%;
      height: 12px;
      border-radius: 999px;
      background: #17263e;
      overflow: hidden;
      border: 1px solid #273b58;
    }
    .progress-fill {
      height: 100%;
      width: 0%;
      background: linear-gradient(90deg, #17c7a7, #31e6c6);
      transition: width .25s ease;
    }
    .progress-meta {
      margin-top: 6px;
      font-size: 12px;
      color: #b2c5eb;
    }
    .chip {
      display: inline-block; padding: 4px 9px; border-radius: 999px; font-size: 12px;
      border: 1px solid var(--line); background: #0e1628;
    }
    .chip-ok { color: #9ef3e7; border-color: #0f766e; background: #083832; }
    .chip-stop { color: #ffd1d1; border-color: #7f1d1d; background: #351112; }
    .panel-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
    .logbox {
      min-height: 150px; max-height: 280px; overflow: auto;
      border: 1px solid var(--line); border-radius: 10px;
      background: #0a1222; padding: 10px; color: #cfe0ff;
      font-family: "IBM Plex Mono", "JetBrains Mono", Menlo, Consolas, monospace;
      font-size: 12px; white-space: pre-wrap;
    }
    table { width: 100%; border-collapse: collapse; }
    th, td {
      border-bottom: 1px solid #1f2d4b;
      text-align: left;
      vertical-align: top;
      padding: 8px 9px;
      font-size: 13px;
      word-break: break-word;
    }
    th {
      color: #bcd2ff; background: #0f172a;
      position: sticky; top: 0; z-index: 1;
    }
    .mono { font-family: "IBM Plex Mono", "JetBrains Mono", Menlo, Consolas, monospace; }
    .scroll-x { overflow: auto; }
    .pager {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 8px;
      margin: 8px 0 10px;
      flex-wrap: wrap;
    }
    .pager .nav { display: flex; gap: 6px; }
    @media (max-width: 1150px) {
      .layout { grid-template-columns: 1fr; }
      .panel-grid { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <div class="wrap">
    <div class="top">
      <div class="brand">
        <div>
          <h1>Neferpitool Web</h1>
          <p>Web control center for domains, typo monitoring, change checks and background jobs.</p>
        </div>
      </div>
      <div class="actions">
        <button type="button" id="bgStart" class="b-ok">Start Background</button>
        <button type="button" id="bgStop">Stop Background</button>
        <button type="button" id="bgRun" class="b-warn">Run One Cycle</button>
        <button type="button" id="refreshAll" class="b-main">Refresh</button>
      </div>
    </div>

    <div class="layout">
      <div class="stack">
        <div class="card">
          <h2>Domain Actions</h2>
          <div class="content stack">
            <form id="addDomainForm" class="stack">
              <input id="newDomain" placeholder="example.com" required>
              <button type="submit" id="addDomainBtn" class="b-main">Add Domain + Generate Typos</button>
            </form>
            <form id="importForm" class="stack">
              <input id="importDomain" placeholder="Main domain for import" required>
              <textarea id="importTypos" placeholder="One typodomain per line"></textarea>
              <button type="submit" id="importBtn" class="b-warn">Import Typodomain List</button>
            </form>
            <form id="presenceForm" class="stack">
              <textarea id="presenceDomains" placeholder="Presence check: one domain per line"></textarea>
              <button type="submit" id="presenceBtn">Check Domain Presence</button>
            </form>
            <div class="row">
              <button type="button" id="checkAllBtn">Check All Domains</button>
              <button type="button" id="expBtn">Expiring</button>
            </div>
            <div class="row">
              <button type="button" id="reliableBtn">Reliable Changes</button>
              <button type="button" id="statsBtn">Stats</button>
            </div>
          </div>
        </div>

        <div class="card">
          <h2>Domain List</h2>
          <div class="content">
            <div class="domains" id="domainsList"></div>
          </div>
        </div>

        <div class="card">
          <h2>System Status</h2>
          <div class="content stack">
            <div id="opStatus" class="busy" style="display:none"></div>
            <div id="opProgressBox" class="progress-wrap" style="display:none">
              <div class="progress-head">
                <span id="opProgressTitle">Operation</span>
                <span id="opProgressPct">0%</span>
              </div>
              <div class="progress-track"><div id="opProgressFill" class="progress-fill"></div></div>
              <div id="opProgressMeta" class="progress-meta"></div>
            </div>
            <div id="bgStatusText" class="muted">Background: n/a</div>
            <div id="statsText" class="muted">Stats: n/a</div>
            <div id="confText" class="muted">Config: n/a</div>
          </div>
        </div>
      </div>

      <div class="stack">
        <div class="card">
          <h2>Selected Domain</h2>
          <div class="content">
            <div class="row">
              <span class="chip" id="selectedDomainChip">none</span>
              <button type="button" id="checkSelectedBtn">Check Selected Domain</button>
            </div>
          </div>
        </div>

        <div class="panel-grid">
          <div class="card">
            <h2>Activity Log</h2>
            <div class="content"><div class="logbox" id="logs"></div></div>
          </div>
          <div class="card">
            <h2>Details</h2>
            <div class="content"><div class="logbox" id="details"></div></div>
          </div>
        </div>

        <div class="card">
          <h2>Typodomains</h2>
          <div class="content scroll-x">
            <div class="stack">
              <input id="typoFilter" placeholder="Global filter (name, status, algorithm, registrar)">
              <div class="row">
                <select id="typoStatusFilter">
                  <option value="">Status: all</option>
                </select>
                <input id="typoRegistrarFilter" placeholder="Registrar contains">
                <input id="typoNameFilter" placeholder="Name contains">
              </div>
              <div class="row">
                <input id="typoMainDomainFilter" placeholder="Main domain contains">
                <select id="typoHasErrorsFilter">
                  <option value="">Errors: any</option>
                  <option value="true">Errors only</option>
                  <option value="false">No errors only</option>
                </select>
                <input id="typoExpiryFromFilter" type="date" title="Expiry from">
                <input id="typoExpiryToFilter" type="date" title="Expiry to">
              </div>
              <div class="row">
                <div id="typoAlgorithmSelect" class="multi-select">
                  <button type="button" id="typoAlgorithmSelectBtn" class="multi-select-btn">Algorithms: all</button>
                  <div id="typoAlgorithmFilter" class="multi-select-menu"></div>
                </div>
              </div>
              <div id="typoColumns" class="row"></div>
            </div>
            <div class="pager">
              <div id="typoSummary" class="muted">No domain selected</div>
              <div class="nav">
                <button type="button" id="prevTyposPage">Prev</button>
                <button type="button" id="nextTyposPage">Next</button>
              </div>
            </div>
            <table>
              <thead>
                <tr id="typosHead"></tr>
              </thead>
              <tbody id="typosBody"></tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>

  <script>
    const TYPO_COLUMNS = [
      { key: 'name', label: 'Name', mono: true },
      { key: 'status', label: 'Status' },
      { key: 'algorithm', label: 'Algorithm' },
      { key: 'legit_domain', label: 'Main Domain' },
      { key: 'registrar', label: 'Registrar' },
      { key: 'expiry_date', label: 'Expiry' },
      { key: 'date_check', label: 'Date Check' },
      { key: 'unicode', label: 'Unicode' },
      { key: 'error_status', label: 'Error Status' },
      { key: 'error_whois', label: 'Error Whois' },
      { key: 'time_status', label: 'Time Status' },
      { key: 'time_whois', label: 'Time Whois' }
    ];
    const TYPO_DEFAULT_COLUMNS = ['name', 'status', 'algorithm', 'registrar', 'expiry_date', 'unicode'];

    const state = {
      selectedDomain: '',
      domains: [],
      currentTypos: [],
      typoFilter: '',
      typoStatusFilter: '',
      typoAlgorithmFilter: [],
      typoRegistrarFilter: '',
      typoNameFilter: '',
      typoMainDomainFilter: '',
      typoHasErrorsFilter: '',
      typoExpiryFromFilter: '',
      typoExpiryToFilter: '',
      typoAvailableStatuses: [],
      typoAvailableAlgorithms: [],
      typoSortBy: 'name',
      typoSortDir: 'asc',
      typoVisibleColumns: TYPO_DEFAULT_COLUMNS.slice(),
      typoPage: 1,
      typoSize: 10,
      typoPages: 0,
      typoTotal: 0,
      busyCount: 0,
      activeOperationId: '',
      operationPollTimer: null
    };

    function logLine(text) {
      const el = document.getElementById('logs');
      const stamp = new Date().toLocaleTimeString();
      el.textContent = '[' + stamp + '] ' + text + '\n' + el.textContent;
    }

    function showDetails(data) {
      document.getElementById('details').textContent = JSON.stringify(data, null, 2);
    }

    function setBusy(on, text) {
      if (on) {
        state.busyCount++;
      } else if (state.busyCount > 0) {
        state.busyCount--;
      }
      const busy = state.busyCount > 0;
      const op = document.getElementById('opStatus');
      op.style.display = busy ? 'block' : 'none';
      op.textContent = text || 'Working... this may take some time for large domains.';

      // Keep domain selection usable while long scans are running.
      document.querySelectorAll('[data-busy-lock="true"]').forEach(function(el) {
        if (el.id && el.id.indexOf('typo') === 0) return;
        el.disabled = busy;
      });
    }

    function renderOperationProgress(op) {
      const box = document.getElementById('opProgressBox');
      const title = document.getElementById('opProgressTitle');
      const pct = document.getElementById('opProgressPct');
      const fill = document.getElementById('opProgressFill');
      const meta = document.getElementById('opProgressMeta');

      if (!op) {
        box.style.display = 'none';
        return;
      }

      box.style.display = 'block';
      title.textContent = (op.type || 'operation') + ' - ' + (op.domain || '');
      const percent = typeof op.progress_percent === 'number' ? op.progress_percent : 0;
      pct.textContent = percent + '%';
      fill.style.width = String(Math.max(0, Math.min(100, percent))) + '%';
      const done = typeof op.done === 'number' ? op.done : 0;
      const total = typeof op.total === 'number' ? op.total : 0;
      const found = typeof op.typodomains_count === 'number' ? op.typodomains_count : 0;
      meta.textContent =
        'state=' + (op.state || 'running') +
        ' | scanned=' + done + '/' + total +
        ' | typodomains=' + found +
        (op.error ? (' | error=' + op.error) : '') +
        (op.message ? (' | ' + op.message) : '');
    }

    function stopOperationPolling() {
      if (state.operationPollTimer) {
        clearInterval(state.operationPollTimer);
        state.operationPollTimer = null;
      }
      state.activeOperationId = '';
    }

    async function pollOperation() {
      if (!state.activeOperationId) return;
      const op = await api('/api/ops/' + encodeURIComponent(state.activeOperationId));
      renderOperationProgress(op);
      if (op.state === 'done' || op.state === 'failed') {
        stopOperationPolling();
        if (op.state === 'done') {
          logLine('Operation completed: ' + (op.message || 'done'));
          await refreshOverview();
          if (op.domain) {
            await selectDomain(op.domain);
          } else if (state.selectedDomain) {
            await loadSelectedDomainPage();
          }
          await refreshStats();
          showDetails(op);
        } else {
          logLine('Operation failed: ' + (op.error || 'unknown error'));
          showDetails(op);
          alert(op.error || 'Operation failed');
        }
      }
    }

    async function startOperationPolling(opId) {
      stopOperationPolling();
      state.activeOperationId = opId;
      await pollOperation();
      if (!state.activeOperationId) return;
      state.operationPollTimer = setInterval(function() {
        pollOperation().catch(function(err) {
          logLine('ERROR polling operation: ' + (err.message || err));
        });
      }, 1000);
    }

    async function api(url, method, body, opts) {
      const options = opts || {};
      const ctrl = new AbortController();
      const timeoutMs = typeof options.timeoutMs === 'number' ? options.timeoutMs : 90000;
      const timeout = timeoutMs > 0 ? setTimeout(function() { ctrl.abort(); }, timeoutMs) : null;
      let res;
      try {
        res = await fetch(url, {
          method: method || 'GET',
          headers: { 'Content-Type': 'application/json' },
          body: body ? JSON.stringify(body) : undefined,
          signal: ctrl.signal
        });
      } finally {
        if (timeout) clearTimeout(timeout);
      }

      const text = await res.text();
      let parsed;
      try {
        parsed = text ? JSON.parse(text) : {};
      } catch (e) {
        throw new Error('Invalid server response');
      }

      if (!res.ok || !parsed.ok) {
        throw new Error((parsed && parsed.error) ? parsed.error : ('HTTP ' + res.status));
      }
      return parsed.data;
    }

    function makeButton(label, cls) {
      const b = document.createElement('button');
      b.type = 'button';
      b.textContent = label;
      if (cls) b.className = cls;
      return b;
    }

    function getVisibleTypoColumns() {
      const selected = state.typoVisibleColumns || [];
      return TYPO_COLUMNS.filter(function(col) {
        return selected.indexOf(col.key) !== -1;
      });
    }

    function isSortableTypoColumn(key) {
      return [
        'name',
        'status',
        'algorithm',
        'legit_domain',
        'registrar',
        'expiry_date',
        'date_check',
        'unicode',
        'error_status',
        'error_whois',
        'time_status',
        'time_whois'
      ].indexOf(key) !== -1;
    }

    function sortArrowForColumn(key) {
      if (state.typoSortBy !== key) return '↕';
      return state.typoSortDir === 'desc' ? '↓' : '↑';
    }

    function toggleSortByColumn(key) {
      if (!isSortableTypoColumn(key)) return;
      if (state.typoSortBy === key) {
        state.typoSortDir = state.typoSortDir === 'asc' ? 'desc' : 'asc';
      } else {
        state.typoSortBy = key;
        state.typoSortDir = 'asc';
      }
      state.typoPage = 1;
      renderTypoHead();
      loadSelectedDomainPage().catch(handleErr);
    }

    function renderTypoHead() {
      const head = document.getElementById('typosHead');
      head.innerHTML = '';
      getVisibleTypoColumns().forEach(function(col) {
        const th = document.createElement('th');
        if (isSortableTypoColumn(col.key)) {
          const btn = document.createElement('button');
          btn.type = 'button';
          btn.className = 'sort-btn' + (state.typoSortBy === col.key ? ' is-active' : '');
          btn.textContent = col.label + ' ' + sortArrowForColumn(col.key);
          btn.addEventListener('click', function() { toggleSortByColumn(col.key); });
          th.appendChild(btn);
        } else {
          th.textContent = col.label;
        }
        head.appendChild(th);
      });
      const actions = document.createElement('th');
      actions.textContent = 'Actions';
      head.appendChild(actions);
    }

    function renderTypoColumnsControls() {
      const root = document.getElementById('typoColumns');
      root.innerHTML = '';
      TYPO_COLUMNS.forEach(function(col) {
        const label = document.createElement('label');
        label.className = 'col-toggle';
        const cb = document.createElement('input');
        cb.type = 'checkbox';
        cb.checked = state.typoVisibleColumns.indexOf(col.key) !== -1;
        cb.addEventListener('change', function() {
          const current = state.typoVisibleColumns.slice();
          if (cb.checked) {
            if (current.indexOf(col.key) === -1) {
              current.push(col.key);
            }
          } else {
            const idx = current.indexOf(col.key);
            if (idx !== -1) {
              current.splice(idx, 1);
            }
          }
          if (current.length === 0) {
            cb.checked = true;
            return;
          }
          state.typoVisibleColumns = current;
          renderTypos(state.currentTypos);
        });
        const text = document.createElement('span');
        text.textContent = col.label;
        label.appendChild(cb);
        label.appendChild(text);
        root.appendChild(label);
      });
    }

    function renderStatusFilterOptions() {
      const select = document.getElementById('typoStatusFilter');
      const current = state.typoStatusFilter || '';
      select.innerHTML = '';

      const allOption = document.createElement('option');
      allOption.value = '';
      allOption.textContent = 'Status: all';
      select.appendChild(allOption);

      (state.typoAvailableStatuses || []).forEach(function(status) {
        const option = document.createElement('option');
        option.value = status;
        option.textContent = status;
        select.appendChild(option);
      });
      select.value = current;
      if (select.value !== current) {
        state.typoStatusFilter = '';
      }
    }

    function renderAlgorithmFilterOptions() {
      const root = document.getElementById('typoAlgorithmFilter');
      const button = document.getElementById('typoAlgorithmSelectBtn');
      root.innerHTML = '';
      const allowed = {};
      (state.typoAvailableAlgorithms || []).forEach(function(v) { allowed[v] = true; });
      state.typoAlgorithmFilter = (state.typoAlgorithmFilter || []).filter(function(v) { return !!allowed[v]; });
      const selected = state.typoAlgorithmFilter || [];

      if (!state.typoAvailableAlgorithms.length) {
        const empty = document.createElement('div');
        empty.className = 'muted';
        empty.textContent = 'Algorithms: n/a';
        root.appendChild(empty);
        button.textContent = 'Algorithms: n/a';
        return;
      }

      if (!selected.length) {
        button.textContent = 'Algorithms: all';
      } else if (selected.length === 1) {
        button.textContent = 'Algorithm: ' + selected[0];
      } else {
        button.textContent = 'Algorithms: ' + selected.length + ' selected';
      }

      state.typoAvailableAlgorithms.forEach(function(algorithm) {
        const label = document.createElement('label');
        label.className = 'algo-filter';
        const cb = document.createElement('input');
        cb.type = 'checkbox';
        cb.checked = selected.indexOf(algorithm) !== -1;
        cb.addEventListener('change', function() {
          const current = state.typoAlgorithmFilter.slice();
          if (cb.checked) {
            if (current.indexOf(algorithm) === -1) {
              current.push(algorithm);
            }
          } else {
            const idx = current.indexOf(algorithm);
            if (idx !== -1) {
              current.splice(idx, 1);
            }
          }
          state.typoAlgorithmFilter = current;
          state.typoPage = 1;
          loadSelectedDomainPage().catch(handleErr);
        });
        const text = document.createElement('span');
        text.textContent = algorithm;
        label.appendChild(cb);
        label.appendChild(text);
        root.appendChild(label);
      });
    }

    function closeAlgorithmDropdown() {
      const wrap = document.getElementById('typoAlgorithmSelect');
      wrap.classList.remove('is-open');
    }

    function toggleAlgorithmDropdown() {
      const wrap = document.getElementById('typoAlgorithmSelect');
      wrap.classList.toggle('is-open');
    }

    function parseDateForFilter(value) {
      const raw = String(value || '').trim();
      if (!raw) return null;
      const iso = new Date(raw);
      if (!Number.isNaN(iso.getTime())) return iso;
      const m = raw.match(/^(\d{2})\/(\d{2})\/(\d{4})$/);
      if (m) {
        const d = new Date(m[3] + '-' + m[2] + '-' + m[1] + 'T00:00:00');
        if (!Number.isNaN(d.getTime())) return d;
      }
      return null;
    }

    function compareLocalTypoRows(a, b, key) {
      function compareString(x, y) {
        const xs = String(x || '').toLowerCase().trim();
        const ys = String(y || '').toLowerCase().trim();
        if (xs < ys) return -1;
        if (xs > ys) return 1;
        return 0;
      }

      function compareDate(x, y) {
        const dx = parseDateForFilter(x);
        const dy = parseDateForFilter(y);
        if (dx && dy) {
          if (dx < dy) return -1;
          if (dx > dy) return 1;
          return 0;
        }
        if (dx && !dy) return -1;
        if (!dx && dy) return 1;
        return 0;
      }

      switch (key) {
        case 'status':
        case 'algorithm':
        case 'legit_domain':
        case 'registrar':
        case 'unicode':
        case 'error_status':
        case 'error_whois':
        case 'time_status':
        case 'time_whois':
          return compareString(a[key], b[key]);
        case 'expiry_date':
          return compareDate(a.expiry_date, b.expiry_date);
        case 'date_check':
          return compareDate(a.created_at_iso || a.date_check, b.created_at_iso || b.date_check);
        default:
          return compareString(a.name, b.name);
      }
    }

    function setSelectedDomain(name) {
      state.selectedDomain = name || '';
      document.getElementById('selectedDomainChip').textContent = state.selectedDomain || 'none';
      if (!state.selectedDomain) {
        state.typoPage = 1;
        state.typoPages = 0;
        state.typoTotal = 0;
        document.getElementById('typoSummary').textContent = 'No domain selected';
        renderTypos([]);
      }
    }

    function renderDomains() {
      const root = document.getElementById('domainsList');
      root.innerHTML = '';
      if (!state.domains.length) {
        const p = document.createElement('div');
        p.className = 'muted';
        p.textContent = 'No domains in database yet.';
        root.appendChild(p);
        return;
      }

      state.domains.forEach(function(domain) {
        const item = document.createElement('div');
        item.className = 'domain';

        const name = document.createElement('div');
        name.className = 'name';
        name.textContent = domain;
        name.title = domain;
        name.addEventListener('click', function() { selectDomain(domain).catch(handleErr); });

        const btns = document.createElement('div');
        btns.className = 'btns';

        const openBtn = makeButton('Open');
        openBtn.addEventListener('click', function() { selectDomain(domain).catch(handleErr); });

        const checkBtn = makeButton('Check');
        checkBtn.addEventListener('click', function() { checkDomain(domain).catch(handleErr); });

        const delBtn = makeButton('Del', 'b-danger');
        delBtn.addEventListener('click', function() { removeDomain(domain).catch(handleErr); });

        btns.appendChild(openBtn);
        btns.appendChild(checkBtn);
        btns.appendChild(delBtn);

        item.appendChild(name);
        item.appendChild(btns);
        root.appendChild(item);
      });
    }

    function renderTypos(list) {
      renderTypoHead();
      state.currentTypos = list || [];
      const body = document.getElementById('typosBody');
      body.innerHTML = '';
      const summary = document.getElementById('typoSummary');
      const total = state.typoTotal;
      const page = state.typoPage;
      const pages = state.typoPages;
      const q = state.typoFilter.trim();
      summary.textContent = q
        ? ('Filtered "' + q + '" | ' + total + ' total | page ' + page + '/' + (pages || 1))
        : (total + ' total | page ' + page + '/' + (pages || 1));
      document.getElementById('prevTyposPage').disabled = page <= 1 || state.busyCount > 0;
      document.getElementById('nextTyposPage').disabled = pages === 0 || page >= pages || state.busyCount > 0;

      if (!state.currentTypos.length) {
        const tr = document.createElement('tr');
        const td = document.createElement('td');
        td.colSpan = getVisibleTypoColumns().length + 1;
        td.className = 'muted';
        td.textContent = state.typoTotal === 0 ? 'No typodomains found.' : 'No typodomains on this page.';
        tr.appendChild(td);
        body.appendChild(tr);
        return;
      }

      state.currentTypos.forEach(function(td) {
        const tr = document.createElement('tr');
        function tdCell(text, cls) {
          const cell = document.createElement('td');
          if (cls) cell.className = cls;
          cell.textContent = text || '';
          return cell;
        }

        getVisibleTypoColumns().forEach(function(col) {
          const raw = td[col.key];
          tr.appendChild(tdCell(raw == null ? '' : String(raw), col.mono ? 'mono' : ''));
        });

        const actions = document.createElement('td');
        const bHistory = makeButton('History');
        bHistory.addEventListener('click', function() { showHistory(td.name).catch(handleErr); });
        const bUpdate = makeButton('Update');
        bUpdate.addEventListener('click', function() { updateTypo(td.name).catch(handleErr); });
        const bDelete = makeButton('Delete', 'b-danger');
        bDelete.addEventListener('click', function() { deleteTypo(td.name).catch(handleErr); });

        actions.appendChild(bHistory);
        actions.appendChild(document.createTextNode(' '));
        actions.appendChild(bUpdate);
        actions.appendChild(document.createTextNode(' '));
        actions.appendChild(bDelete);
        tr.appendChild(actions);

        body.appendChild(tr);
      });
    }

    function renderBackgroundStatus(st) {
      const el = document.getElementById('bgStatusText');
      const tagClass = st.running ? 'chip chip-ok' : 'chip chip-stop';
      const run = st.running ? 'running' : 'stopped';
      const hasLast = st.last_run && String(st.last_run).indexOf('0001-01-01') !== 0;
      const last = hasLast ? (' | last run: ' + new Date(st.last_run).toLocaleString()) : '';
      el.innerHTML = 'Background: <span class="' + tagClass + '">' + run + '</span>' +
        ' | crons=' + st.active_crons +
        ' | changes=' + st.last_changes +
        ' | scanErrs=' + st.last_scan_errors + last +
        (st.last_error ? (' | err=' + st.last_error) : '');
    }

    async function refreshOverview() {
      const data = await api('/api/overview');
      state.domains = data.domains || [];
      renderDomains();
      if (!state.selectedDomain && state.domains.length > 0) {
        setSelectedDomain(state.domains[0]);
      }
      renderBackgroundStatus(data.background_status);
      document.getElementById('statsText').textContent = 'Domains in DB: ' + data.domains_count;
    }

    async function refreshConfig() {
      const conf = await api('/api/config');
      document.getElementById('confText').textContent =
        'Sleep=' + conf.MINUTESLEEPBACKGROUNDMONITORING +
        'm | Expiration window=' + conf.EXPIRATIONTIME +
        'd | Cron rules=' + ((conf.REPORTFREQUENCY || []).length);
    }

    async function refreshStats() {
      const st = await api('/api/stats');
      document.getElementById('statsText').textContent =
        'Domains in DB: ' + state.domains.length +
        ' | Typodomains: ' + st.typodomains_total +
        ' | Scan errors: ' + st.typodomains_errors;
      return st;
    }

    async function loadSelectedDomainPage() {
      if (!state.selectedDomain) return;
      const params = new URLSearchParams();
      params.set('page', String(state.typoPage));
      params.set('size', String(state.typoSize));
      params.set('sort_by', state.typoSortBy);
      params.set('sort_dir', state.typoSortDir);
      if (state.typoFilter.trim()) {
        params.set('q', state.typoFilter.trim());
      }
      if (state.typoStatusFilter.trim()) {
        params.set('status', state.typoStatusFilter.trim());
      }
      if (state.typoAlgorithmFilter.length > 0) {
        params.set('algorithm', state.typoAlgorithmFilter.join(','));
      }
      if (state.typoRegistrarFilter.trim()) {
        params.set('registrar', state.typoRegistrarFilter.trim());
      }
      if (state.typoNameFilter.trim()) {
        params.set('name', state.typoNameFilter.trim());
      }
      if (state.typoMainDomainFilter.trim()) {
        params.set('legit_domain', state.typoMainDomainFilter.trim());
      }
      if (state.typoHasErrorsFilter) {
        params.set('has_errors', state.typoHasErrorsFilter);
      }
      if (state.typoExpiryFromFilter) {
        params.set('expiry_from', state.typoExpiryFromFilter);
      }
      if (state.typoExpiryToFilter) {
        params.set('expiry_to', state.typoExpiryToFilter);
      }
      try {
        const out = await api('/api/domain/' + encodeURIComponent(state.selectedDomain) + '/typos/page?' + params.toString());
        state.typoAvailableStatuses = out.available_statuses || [];
        state.typoAvailableAlgorithms = out.available_algorithms || [];
        renderStatusFilterOptions();
        renderAlgorithmFilterOptions();
        state.typoPage = out.page || 1;
        state.typoPages = out.pages || 0;
        state.typoTotal = out.total || 0;
        renderTypos(out.items || []);
      } catch (err) {
        // Fallback for old backend binaries that don't expose /typos/page yet.
        const all = await api('/api/domain/' + encodeURIComponent(state.selectedDomain) + '/typos');
        const q = state.typoFilter.trim().toLowerCase();
        const status = state.typoStatusFilter.trim().toLowerCase();
        const algorithms = (state.typoAlgorithmFilter || []).map(function(v) { return String(v || '').toLowerCase(); });
        const registrar = state.typoRegistrarFilter.trim().toLowerCase();
        const name = state.typoNameFilter.trim().toLowerCase();
        const legitDomain = state.typoMainDomainFilter.trim().toLowerCase();
        const hasErrors = state.typoHasErrorsFilter;
        const expiryFrom = state.typoExpiryFromFilter ? parseDateForFilter(state.typoExpiryFromFilter) : null;
        const expiryTo = state.typoExpiryToFilter ? parseDateForFilter(state.typoExpiryToFilter) : null;

        const statusSet = {};
        const algorithmSet = {};
        (all || []).forEach(function(td) {
          const s = String(td.status || '').trim();
          if (s) statusSet[s] = true;
          const a = String(td.algorithm || '').trim();
          if (a) algorithmSet[a] = true;
        });
        state.typoAvailableStatuses = Object.keys(statusSet).sort(function(a, b) { return a.toLowerCase().localeCompare(b.toLowerCase()); });
        state.typoAvailableAlgorithms = Object.keys(algorithmSet).sort(function(a, b) { return a.toLowerCase().localeCompare(b.toLowerCase()); });
        renderStatusFilterOptions();
        renderAlgorithmFilterOptions();

        const filtered = (all || []).filter(function(td) {
          const candidate = [
            td.name, td.status, td.algorithm, td.registrar, td.unicode, td.legit_domain
          ].join(' ').toLowerCase();
          if (q && candidate.indexOf(q) === -1) return false;
          if (status && String(td.status || '').toLowerCase() !== status) return false;
          if (algorithms.length > 0 && algorithms.indexOf(String(td.algorithm || '').toLowerCase()) === -1) return false;
          if (registrar && (td.registrar || '').toLowerCase().indexOf(registrar) === -1) return false;
          if (name && (td.name || '').toLowerCase().indexOf(name) === -1) return false;
          if (legitDomain && (td.legit_domain || '').toLowerCase().indexOf(legitDomain) === -1) return false;
          if (hasErrors === 'true') {
            if (!td.error_status && !td.error_whois) return false;
          }
          if (hasErrors === 'false') {
            if (td.error_status || td.error_whois) return false;
          }
          if (expiryFrom || expiryTo) {
            const expiry = parseDateForFilter(td.expiry_date);
            if (!expiry) return false;
            if (expiryFrom && expiry < expiryFrom) return false;
            if (expiryTo && expiry > expiryTo) return false;
          }
          return true;
        });
        filtered.sort(function(a, b) {
          const direction = state.typoSortDir === 'desc' ? -1 : 1;
          const key = state.typoSortBy || 'name';
          const cmp = compareLocalTypoRows(a, b, key);
          if (cmp !== 0) {
            return cmp * direction;
          }
          return compareLocalTypoRows(a, b, 'name') * direction;
        });
        state.typoTotal = filtered.length;
        state.typoPages = state.typoTotal === 0 ? 0 : Math.ceil(state.typoTotal / state.typoSize);
        if (state.typoPages > 0 && state.typoPage > state.typoPages) {
          state.typoPage = state.typoPages;
        }
        const start = Math.max(0, (state.typoPage - 1) * state.typoSize);
        const end = Math.min(filtered.length, start + state.typoSize);
        renderTypos(filtered.slice(start, end));
      }
      logLine('Loaded typos for ' + state.selectedDomain + ' | page ' + state.typoPage + '/' + (state.typoPages || 1));
    }

    async function selectDomain(domain) {
      setSelectedDomain(domain);
      state.typoPage = 1;
      await loadSelectedDomainPage();
    }

    async function addDomain(ev) {
      ev.preventDefault();
      const domain = document.getElementById('newDomain').value.trim();
      if (!domain) return;
      const out = await api('/api/domain/add/start', 'POST', { domain: domain });
      const opId = out.operation_id || '';
      if (!opId) {
        throw new Error('Cannot start operation');
      }
      logLine('Add domain started for ' + domain);
      document.getElementById('newDomain').value = '';
      showDetails(out);
      await refreshOverview();
      await startOperationPolling(opId);
    }

    async function importTypos(ev) {
      ev.preventDefault();
      const domain = document.getElementById('importDomain').value.trim();
      const typos = document.getElementById('importTypos').value;
      if (!domain || !typos.trim()) return;
      setBusy(true, 'Importing and scanning typodomains for ' + domain + '.');
      try {
        const out = await api('/api/domain/import', 'POST', { domain: domain, typos: typos, include_typodomains: false }, { timeoutMs: 0 });
        logLine('Imported typodomains for ' + domain + ' | typodomains=' + (out.typodomains_count || 0));
        showDetails(out);
        await refreshOverview();
        await selectDomain(domain);
      } finally {
        setBusy(false);
      }
    }

    async function checkPresence(ev) {
      ev.preventDefault();
      const value = document.getElementById('presenceDomains').value || '';
      const domains = value.split('\n').map(function(v) { return v.trim(); }).filter(Boolean);
      if (!domains.length) return;
      const out = await api('/api/presence', 'POST', { domains: domains }, { timeoutMs: 30000 });
      showDetails(out);
      logLine('Presence checked for ' + domains.length + ' domains');
    }

    async function removeDomain(domain) {
      if (!confirm('Remove domain "' + domain + '" and all its typodomains?')) return;
      await api('/api/domain/' + encodeURIComponent(domain) + '/remove', 'POST');
      logLine('Removed domain ' + domain);
      if (state.selectedDomain === domain) {
        setSelectedDomain('');
      }
      await refreshOverview();
    }

    async function checkDomain(domain) {
      setBusy(true, 'Checking domain ' + domain + ' for reliable changes...');
      try {
        const out = await api('/api/domain/' + encodeURIComponent(domain) + '/check', 'POST', null, { timeoutMs: 0 });
        logLine('Checked domain ' + domain + ' | changes=' + (out.changes || []).length);
        showDetails(out);
        if (state.selectedDomain === domain) {
          await loadSelectedDomainPage();
        }
      } finally {
        setBusy(false);
      }
    }

    async function checkAll() {
      setBusy(true, 'Checking all domains. This can be long on large datasets...');
      try {
        const out = await api('/api/check/all', 'POST', null, { timeoutMs: 0 });
        logLine('Checked all domains | changes=' + (out.changes || []).length);
        showDetails(out);
        if (state.selectedDomain) {
          await loadSelectedDomainPage();
        }
      } finally {
        setBusy(false);
      }
    }

    async function showHistory(name) {
      const out = await api('/api/typodomain/' + encodeURIComponent(name) + '/history');
      showDetails(out);
      logLine('Loaded history for ' + name + ' | entries=' + out.length);
    }

    async function updateTypo(name) {
      setBusy(true, 'Updating typodomain ' + name + '...');
      try {
        const out = await api('/api/typodomain/' + encodeURIComponent(name) + '/update', 'POST', null, { timeoutMs: 0 });
        logLine('Updated typodomain ' + name);
        showDetails(out);
        if (state.selectedDomain) {
          await loadSelectedDomainPage();
        }
      } finally {
        setBusy(false);
      }
    }

    async function deleteTypo(name) {
      if (!confirm('Delete typodomain history for "' + name + '"?')) return;
      setBusy(true, 'Deleting typodomain history for ' + name + '...');
      try {
        await api('/api/typodomain/' + encodeURIComponent(name) + '/delete', 'POST');
        logLine('Deleted typodomain ' + name);
        showDetails({ deleted: name });
        if (state.selectedDomain) {
          await loadSelectedDomainPage();
        }
      } finally {
        setBusy(false);
      }
    }

    async function showExpiration() {
      const out = await api('/api/expiration');
      showDetails(out);
      logLine('Loaded expiring typodomains | count=' + out.length);
    }

    async function showReliable() {
      const out = await api('/api/reliable');
      showDetails(out);
      logLine('Loaded reliable changes | count=' + out.length);
    }

    async function showStats() {
      const out = await refreshStats();
      showDetails(out);
      logLine('Loaded stats');
    }

    async function bgStart() {
      const st = await api('/api/background/start', 'POST');
      renderBackgroundStatus(st);
      logLine('Background started');
    }

    async function bgStop() {
      const st = await api('/api/background/stop', 'POST');
      renderBackgroundStatus(st);
      logLine('Background stopped');
    }

    async function bgRunOnce() {
      setBusy(true, 'Running one background monitoring cycle...');
      try {
        const st = await api('/api/background/run', 'POST', null, { timeoutMs: 0 });
        renderBackgroundStatus(st);
        logLine('Background cycle executed');
        if (state.selectedDomain) {
          await loadSelectedDomainPage();
        }
      } finally {
        setBusy(false);
      }
    }

    async function refreshAll() {
      await refreshOverview();
      await refreshConfig();
      await refreshStats();
      if (state.selectedDomain) {
        await loadSelectedDomainPage();
      }
    }

    function handleErr(err) {
      const msg = err && err.name === 'AbortError'
        ? 'Request timeout. Try again with a smaller scope.'
        : (err && err.message ? err.message : 'Unexpected error');
      logLine('ERROR: ' + msg);
      alert(msg);
    }

    document.getElementById('addDomainForm').addEventListener('submit', function(ev) { addDomain(ev).catch(handleErr); });
    document.getElementById('importForm').addEventListener('submit', function(ev) { importTypos(ev).catch(handleErr); });
    document.getElementById('presenceForm').addEventListener('submit', function(ev) { checkPresence(ev).catch(handleErr); });
    document.getElementById('checkAllBtn').addEventListener('click', function() { checkAll().catch(handleErr); });
    document.getElementById('expBtn').addEventListener('click', function() { showExpiration().catch(handleErr); });
    document.getElementById('reliableBtn').addEventListener('click', function() { showReliable().catch(handleErr); });
    document.getElementById('statsBtn').addEventListener('click', function() { showStats().catch(handleErr); });
    document.getElementById('bgStart').addEventListener('click', function() { bgStart().catch(handleErr); });
    document.getElementById('bgStop').addEventListener('click', function() { bgStop().catch(handleErr); });
    document.getElementById('bgRun').addEventListener('click', function() { bgRunOnce().catch(handleErr); });
    document.getElementById('refreshAll').addEventListener('click', function() { refreshAll().catch(handleErr); });
    document.getElementById('typoFilter').addEventListener('input', function(ev) {
      state.typoFilter = ev.target.value || '';
      state.typoPage = 1;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('typoStatusFilter').addEventListener('change', function(ev) {
      state.typoStatusFilter = ev.target.value || '';
      state.typoPage = 1;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('typoAlgorithmSelectBtn').addEventListener('click', function(ev) {
      ev.stopPropagation();
      toggleAlgorithmDropdown();
    });
    document.getElementById('typoRegistrarFilter').addEventListener('input', function(ev) {
      state.typoRegistrarFilter = ev.target.value || '';
      state.typoPage = 1;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('typoNameFilter').addEventListener('input', function(ev) {
      state.typoNameFilter = ev.target.value || '';
      state.typoPage = 1;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('typoMainDomainFilter').addEventListener('input', function(ev) {
      state.typoMainDomainFilter = ev.target.value || '';
      state.typoPage = 1;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('typoHasErrorsFilter').addEventListener('change', function(ev) {
      state.typoHasErrorsFilter = ev.target.value || '';
      state.typoPage = 1;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('typoExpiryFromFilter').addEventListener('change', function(ev) {
      state.typoExpiryFromFilter = ev.target.value || '';
      state.typoPage = 1;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('typoExpiryToFilter').addEventListener('change', function(ev) {
      state.typoExpiryToFilter = ev.target.value || '';
      state.typoPage = 1;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('prevTyposPage').addEventListener('click', function() {
      if (state.typoPage <= 1) return;
      state.typoPage--;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('nextTyposPage').addEventListener('click', function() {
      if (state.typoPages === 0 || state.typoPage >= state.typoPages) return;
      state.typoPage++;
      loadSelectedDomainPage().catch(handleErr);
    });
    document.getElementById('checkSelectedBtn').addEventListener('click', function() {
      if (!state.selectedDomain) {
        alert('Select a domain first.');
        return;
      }
      checkDomain(state.selectedDomain).catch(handleErr);
    });
    document.addEventListener('click', function(ev) {
      const wrap = document.getElementById('typoAlgorithmSelect');
      if (!wrap.contains(ev.target)) {
        closeAlgorithmDropdown();
      }
    });
    document.addEventListener('keydown', function(ev) {
      if (ev.key === 'Escape') {
        closeAlgorithmDropdown();
      }
    });

    [
      'newDomain',
      'addDomainBtn',
      'importDomain',
      'importTypos',
      'importBtn',
      'presenceDomains',
      'presenceBtn',
      'checkAllBtn',
      'expBtn',
      'reliableBtn',
      'statsBtn',
      'bgStart',
      'bgStop',
      'bgRun',
      'checkSelectedBtn'
    ].forEach(function(id) {
      const el = document.getElementById(id);
      if (el) {
        el.setAttribute('data-busy-lock', 'true');
      }
    });

    renderStatusFilterOptions();
    renderAlgorithmFilterOptions();
    renderTypoColumnsControls();
    renderTypoHead();

    refreshAll().catch(handleErr);
  </script>
</body>
</html>`
