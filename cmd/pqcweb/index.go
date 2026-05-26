package main

import "net/http"

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

const indexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>PQC Readiness</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f4f6f8;
      --panel: #ffffff;
      --panel-soft: #f9fafb;
      --text: #17202a;
      --muted: #607080;
      --line: #d8dee8;
      --accent: #126b5a;
      --accent-dark: #0d5145;
      --blue: #2f5f9e;
      --warn: #a33d16;
      --critical-bg: #fff2ee;
      --critical-text: #8f2d11;
      --shadow: 0 12px 28px rgba(21, 35, 48, 0.08);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    header {
      position: sticky;
      top: 0;
      z-index: 5;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 18px;
      padding: 16px 28px;
      background: rgba(255, 255, 255, 0.96);
      border-bottom: 1px solid var(--line);
    }
    h1 {
      margin: 0;
      font-size: 21px;
      line-height: 1.2;
      letter-spacing: 0;
    }
    h2 {
      margin: 0;
      font-size: 15px;
      line-height: 1.3;
    }
    h3 {
      margin: 0;
      font-size: 13px;
      color: #344252;
      line-height: 1.3;
    }
    main {
      max-width: 1280px;
      margin: 0 auto;
      padding: 22px;
      display: grid;
      grid-template-columns: 390px minmax(0, 1fr);
      gap: 18px;
    }
    .panel {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
      box-shadow: var(--shadow);
    }
    .form-panel {
      align-self: start;
      overflow: hidden;
    }
    .panel-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      padding: 16px 18px;
      border-bottom: 1px solid var(--line);
      background: var(--panel-soft);
    }
    .status-pill {
      min-width: 78px;
      border: 1px solid var(--line);
      border-radius: 999px;
      padding: 6px 10px;
      background: #ffffff;
      color: #344252;
      font-size: 12px;
      font-weight: 800;
      text-align: center;
    }
    form { padding: 18px; }
    .tabs {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 8px;
      margin-bottom: 16px;
    }
    .tab-button {
      min-height: 38px;
      border: 1px solid var(--line);
      border-radius: 6px;
      background: #ffffff;
      color: #30435a;
      font-weight: 800;
      cursor: pointer;
    }
    .tab-button.active {
      background: #e9f3f0;
      border-color: #91beb3;
      color: var(--accent-dark);
    }
    .field-group {
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 14px;
      margin-bottom: 14px;
      background: #ffffff;
    }
    .field-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 10px;
    }
    label {
      display: block;
      margin: 12px 0 6px;
      color: #344252;
      font-size: 12px;
      font-weight: 800;
    }
    label:first-child { margin-top: 0; }
    input, select, textarea {
      width: 100%;
      border: 1px solid var(--line);
      border-radius: 6px;
      padding: 8px 10px;
      background: #ffffff;
      color: var(--text);
      font: inherit;
    }
    input, select { min-height: 38px; }
    textarea {
      min-height: 136px;
      resize: vertical;
      line-height: 1.4;
    }
    .hidden { display: none; }
    .actions {
      display: grid;
      grid-template-columns: 1fr;
      gap: 10px;
      margin-top: 16px;
    }
    button {
      min-height: 40px;
      border: 0;
      border-radius: 6px;
      padding: 0 14px;
      background: var(--accent);
      color: #ffffff;
      font-weight: 800;
      cursor: pointer;
    }
    button.secondary { background: #30435a; }
    button.ghost {
      background: #ffffff;
      color: #30435a;
      border: 1px solid var(--line);
    }
    button:disabled {
      cursor: wait;
      opacity: 0.62;
    }
    .error {
      min-height: 20px;
      margin-top: 12px;
      color: var(--warn);
      font-size: 13px;
      font-weight: 800;
    }
    .results-panel {
      min-height: 720px;
      overflow: hidden;
    }
    .empty-state {
      min-height: 720px;
      display: grid;
      place-items: center;
      padding: 24px;
      color: var(--muted);
      text-align: center;
    }
    .summary {
      display: grid;
      grid-template-columns: repeat(4, minmax(120px, 1fr));
      gap: 12px;
      padding: 16px;
      border-bottom: 1px solid var(--line);
      background: var(--panel-soft);
    }
    .metric {
      min-height: 88px;
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 12px;
      background: #ffffff;
    }
    .metric span {
      display: block;
      color: var(--muted);
      font-size: 12px;
      font-weight: 800;
    }
    .metric strong {
      display: block;
      margin-top: 8px;
      font-size: 24px;
      line-height: 1.1;
      letter-spacing: 0;
    }
    .metric.critical strong,
    .metric.very-high strong {
      color: var(--critical-text);
    }
    .content {
      padding: 16px;
      display: grid;
      gap: 16px;
    }
    .report-grid {
      display: grid;
      grid-template-columns: minmax(0, 1fr) 280px;
      gap: 16px;
    }
    .report-section {
      border: 1px solid var(--line);
      border-radius: 8px;
      overflow: hidden;
      background: #ffffff;
    }
    .report-section > .section-title {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      padding: 13px 14px;
      border-bottom: 1px solid var(--line);
      background: var(--panel-soft);
    }
    .section-title button {
      min-height: 32px;
      padding: 0 12px;
      font-size: 12px;
      flex: 0 0 auto;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      font-size: 13px;
    }
    th, td {
      padding: 10px;
      border-bottom: 1px solid var(--line);
      text-align: left;
      vertical-align: top;
    }
    th {
      color: #344252;
      background: #ffffff;
      font-size: 12px;
    }
    tr:last-child td { border-bottom: 0; }
    .risk-chip {
      display: inline-flex;
      align-items: center;
      min-height: 24px;
      border-radius: 999px;
      padding: 3px 9px;
      background: #eef3f8;
      color: #30435a;
      font-size: 12px;
      font-weight: 800;
      white-space: nowrap;
    }
    .risk-chip.critical,
    .risk-chip.very-high {
      background: var(--critical-bg);
      color: var(--critical-text);
    }
    .kv-list {
      display: grid;
      gap: 8px;
      padding: 12px;
    }
    .kv-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 10px;
      border-bottom: 1px solid #edf0f4;
      padding-bottom: 8px;
      font-size: 13px;
    }
    .kv-row:last-child {
      border-bottom: 0;
      padding-bottom: 0;
    }
    .kv-row strong { font-size: 13px; }
    .actions-list {
      margin: 0;
      padding: 12px 12px 12px 30px;
      font-size: 13px;
      line-height: 1.45;
    }
    pre {
      margin: 0;
      max-height: 360px;
      overflow: auto;
      white-space: pre-wrap;
      word-break: break-word;
      border-radius: 0 0 8px 8px;
      padding: 14px;
      background: #111827;
      color: #e5e7eb;
      font-size: 12px;
      line-height: 1.45;
    }
    .copy-note {
      min-height: 18px;
      padding: 0 14px 12px;
      color: var(--muted);
      font-size: 12px;
    }
    @media (max-width: 1040px) {
      main { grid-template-columns: 1fr; }
      .report-grid { grid-template-columns: 1fr; }
    }
    @media (max-width: 720px) {
      header { padding: 14px 16px; }
      main { padding: 14px; }
      .summary { grid-template-columns: repeat(2, minmax(120px, 1fr)); }
      .field-grid { grid-template-columns: 1fr; }
      .tabs { grid-template-columns: 1fr 1fr; }
      th:nth-child(4), td:nth-child(4) { display: none; }
    }
  </style>
</head>
<body>
  <header>
    <h1>PQC Readiness</h1>
    <span id="status" class="status-pill">Ready</span>
  </header>
  <main>
    <section class="panel form-panel">
      <div class="panel-header">
        <h2>Assessment</h2>
      </div>
      <form id="scan-form">
        <div class="tabs" role="tablist" aria-label="Scan source">
          <button type="button" class="tab-button active" data-mode="sample">Sample</button>
          <button type="button" class="tab-button" data-mode="zip">Zip</button>
          <button type="button" class="tab-button" data-mode="github">GitHub</button>
          <button type="button" class="tab-button" data-mode="domains">Domains</button>
        </div>

        <div class="field-group">
          <label for="app_name">Application name</label>
          <input id="app_name" name="app_name" value="Clinical API">

          <div id="source-zip" class="source-panel hidden">
            <label for="repo">Repository zip</label>
            <input id="repo" name="repo" type="file" accept=".zip">
          </div>

          <div id="source-github" class="source-panel hidden">
            <label for="repo_url">GitHub repo URL</label>
            <input id="repo_url" name="repo_url" placeholder="https://github.com/owner/repo">

            <label for="branch">Branch</label>
            <input id="branch" name="branch" placeholder="main">
          </div>

          <div id="source-domains" class="source-panel hidden">
            <label for="domains">TLS domains</label>
            <textarea id="domains" name="domains" placeholder="github.com&#10;google.com&#10;cloudflare.com&#10;microsoft.com&#10;amazon.com&#10;openai.com"></textarea>
          </div>
        </div>

        <div class="field-group">
          <h3>Business context</h3>
          <label for="sensitivity">Sensitivity flags</label>
          <input id="sensitivity" name="sensitivity" value="PHI,GxP">

          <div class="field-grid">
            <div>
              <label for="exposure">Exposure</label>
              <select id="exposure" name="exposure">
                <option value="internet_authenticated">Internet authenticated</option>
                <option value="internet_public">Internet public</option>
                <option value="partner">Partner</option>
                <option value="internal">Internal</option>
                <option value="offline">Offline</option>
              </select>
            </div>
            <div>
              <label for="business_criticality">Criticality</label>
              <select id="business_criticality" name="business_criticality">
                <option value="high">High</option>
                <option value="mission_critical">Mission critical</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>
          </div>

          <label for="secrecy_lifetime_years">Secrecy lifetime years</label>
          <input id="secrecy_lifetime_years" name="secrecy_lifetime_years" type="number" min="0" value="15">

          <label for="crypto_agility">Crypto agility</label>
          <select id="crypto_agility" name="crypto_agility">
            <option value="partially_configurable">Partially configurable</option>
            <option value="centralized_policy_driven">Centralized policy driven</option>
            <option value="configurable">Configurable</option>
            <option value="hardcoded">Hardcoded</option>
            <option value="vendor_controlled_or_legacy">Vendor controlled or legacy</option>
          </select>

          <label for="vendor_dependency">Vendor dependency</label>
          <select id="vendor_dependency" name="vendor_dependency">
            <option value="cloud_with_roadmap">Cloud with roadmap</option>
            <option value="none">None</option>
            <option value="saas_unknown_roadmap">SaaS unknown roadmap</option>
            <option value="legacy_no_roadmap">Legacy no roadmap</option>
          </select>

          <label for="migration_complexity">Migration complexity</label>
          <select id="migration_complexity" name="migration_complexity">
            <option value="multi_system">Multi-system</option>
            <option value="config_change">Config change</option>
            <option value="library_upgrade">Library upgrade</option>
            <option value="code_change">Code change</option>
            <option value="legacy_regulated_vendor_bound">Legacy regulated vendor-bound</option>
          </select>
        </div>

        <div class="actions">
          <button type="submit" id="run-scan">Run Scan</button>
        </div>
        <div class="error" id="error"></div>
      </form>
    </section>

    <section class="panel results-panel">
      <div id="results" class="empty-state">
        <div>
          <h2>Assessment report</h2>
        </div>
      </div>
    </section>
  </main>

  <script>
    const form = document.querySelector("#scan-form");
    const error = document.querySelector("#error");
    const statusText = document.querySelector("#status");
    const results = document.querySelector("#results");
    const runButton = document.querySelector("#run-scan");
    const modeButtons = Array.from(document.querySelectorAll(".tab-button"));
    let activeMode = "sample";
    let latestReportJSON = "";

    modeButtons.forEach(button => {
      button.addEventListener("click", () => setMode(button.dataset.mode));
    });

    form.addEventListener("submit", event => {
      event.preventDefault();
      const data = new FormData(form);
      data.delete("repo");
      let endpoint = "/api/scan/sample";

      if (activeMode === "zip") {
        const repo = document.querySelector("#repo");
        if (!repo.files.length) {
          showError("Choose a repository zip file.");
          return;
        }
        data.set("repo", repo.files[0]);
        endpoint = "/api/scan/upload";
      }
      if (activeMode === "github") {
        if (!String(data.get("repo_url") || "").trim()) {
          showError("Enter a public GitHub repository URL.");
          return;
        }
        endpoint = "/api/scan/git";
      }
      if (activeMode === "domains") {
        if (!String(data.get("domains") || "").trim()) {
          showError("Enter one or more TLS domains.");
          return;
        }
        endpoint = "/api/scan/domains";
      }

      runScan(endpoint, data);
    });

    function setMode(mode) {
      activeMode = mode;
      modeButtons.forEach(button => {
        button.classList.toggle("active", button.dataset.mode === mode);
      });
      document.querySelectorAll(".source-panel").forEach(panel => {
        panel.classList.add("hidden");
      });
      const panel = document.querySelector("#source-" + mode);
      if (panel) panel.classList.remove("hidden");
      runButton.textContent = mode === "sample" ? "Scan Sample" :
        mode === "zip" ? "Scan Zip" :
        mode === "github" ? "Scan GitHub" : "Scan Domains";
      showError("");
    }

    async function runScan(url, data) {
      setBusy(true);
      showError("");
      try {
        const response = await fetch(url, { method: "POST", body: data });
        const payload = await response.json();
        if (!response.ok) throw new Error(payload.error || "Scan failed");
        render(payload);
      } catch (err) {
        showError(err.message);
      } finally {
        setBusy(false);
      }
    }

    function render(report) {
      latestReportJSON = JSON.stringify(report, null, 2);
      const postureClass = riskClass(report.risk_posture);
      const summary = report.executive_summary || {};
      const topAssets = report.top_critical_assets || [];
      const inventory = report.full_crypto_inventory || [];
      const unreadable = (report.scan_summary && report.scan_summary.unreadable_files) || [];

      results.className = "";
      results.innerHTML =
        "<div class=\"summary\">" +
          metric("Risk posture", report.risk_posture || "Low", postureClass) +
          metric("Highest score", Number(report.highest_quantum_risk_score || 0).toFixed(1), "") +
          metric("Findings", report.total_findings || 0, "") +
          metric("Files scanned", safePath(report, "scan_summary.files_scanned", 0), "") +
        "</div>" +
        "<div class=\"content\">" +
          "<div class=\"report-grid\">" +
            section("Findings by algorithm", keyValueRows(report.findings_by_algorithm || {})) +
            section("Findings by usage", keyValueRows(report.findings_by_usage_type || {})) +
          "</div>" +
          section("Top assets", topAssetTable(topAssets)) +
          section("Inventory preview", inventoryTable(inventory.slice(0, 12))) +
          section("30 / 60 / 90 actions", actionsHTML(summary)) +
          unreadableSection(unreadable) +
          jsonSection() +
        "</div>";

      document.querySelector("#copy-json").addEventListener("click", copyJSONReport);
    }

    function metric(label, value, className) {
      return "<div class=\"metric " + className + "\"><span>" + escapeHTML(label) + "</span><strong>" + escapeHTML(value) + "</strong></div>";
    }

    function section(title, body) {
      return "<div class=\"report-section\"><div class=\"section-title\"><h2>" + escapeHTML(title) + "</h2></div>" + body + "</div>";
    }

    function keyValueRows(values) {
      const keys = Object.keys(values).sort();
      if (!keys.length) return "<div class=\"kv-list\"><div class=\"kv-row\"><span>No findings</span><strong>0</strong></div></div>";
      return "<div class=\"kv-list\">" + keys.map(key =>
        "<div class=\"kv-row\"><span>" + escapeHTML(key) + "</span><strong>" + escapeHTML(values[key]) + "</strong></div>"
      ).join("") + "</div>";
    }

    function topAssetTable(assets) {
      const rows = assets.map(asset =>
        "<tr>" +
          "<td>" + escapeHTML(asset.algorithm) + "</td>" +
          "<td><span class=\"risk-chip " + riskClass(asset.risk_level) + "\">" + escapeHTML(asset.risk_level) + " " + Number(asset.risk_score || 0).toFixed(1) + "</span></td>" +
          "<td>" + escapeHTML(asset.usage) + "</td>" +
          "<td>" + escapeHTML(asset.file) + ":" + escapeHTML(asset.line) + "</td>" +
        "</tr>"
      ).join("");
      return "<table><thead><tr><th>Algorithm</th><th>Risk</th><th>Usage</th><th>Location</th></tr></thead><tbody>" +
        (rows || "<tr><td colspan=\"4\">No findings detected.</td></tr>") +
        "</tbody></table>";
    }

    function inventoryTable(items) {
      const rows = items.map(item =>
        "<tr>" +
          "<td>" + escapeHTML(item.algorithm) + "</td>" +
          "<td>" + escapeHTML(item.crypto_usage_type) + "</td>" +
          "<td>" + escapeHTML(item.severity) + "</td>" +
          "<td>" + escapeHTML(item.confidence) + "</td>" +
          "<td>" + escapeHTML(item.file_path) + (item.line_number ? ":" + escapeHTML(item.line_number) : "") + "</td>" +
        "</tr>"
      ).join("");
      return "<table><thead><tr><th>Algorithm</th><th>Usage</th><th>Severity</th><th>Confidence</th><th>Location</th></tr></thead><tbody>" +
        (rows || "<tr><td colspan=\"5\">No inventory records.</td></tr>") +
        "</tbody></table>";
    }

    function actionsHTML(summary) {
      const groups = [
        ["30 days", summary.actions_30_days || []],
        ["60 days", summary.actions_60_days || []],
        ["90 days", summary.actions_90_days || []]
      ];
      return groups.map(group =>
        "<div class=\"section-title\"><h3>" + group[0] + "</h3></div><ol class=\"actions-list\">" +
          group[1].map(action => "<li>" + escapeHTML(action) + "</li>").join("") +
        "</ol>"
      ).join("");
    }

    function unreadableSection(items) {
      if (!items.length) return "";
      return section("Scan notices", "<div class=\"kv-list\">" + items.map(item =>
        "<div class=\"kv-row\"><span>" + escapeHTML(item) + "</span></div>"
      ).join("") + "</div>");
    }

    function jsonSection() {
      return "<div class=\"report-section\">" +
        "<div class=\"section-title\"><h2>JSON report</h2><button type=\"button\" id=\"copy-json\" class=\"ghost\">Copy JSON</button></div>" +
        "<pre id=\"json-report\">" + escapeHTML(latestReportJSON) + "</pre>" +
        "<div id=\"copy-note\" class=\"copy-note\"></div>" +
      "</div>";
    }

    async function copyJSONReport() {
      const note = document.querySelector("#copy-note");
      try {
        if (navigator.clipboard && window.isSecureContext) {
          await navigator.clipboard.writeText(latestReportJSON);
        } else {
          const temp = document.createElement("textarea");
          temp.value = latestReportJSON;
          temp.style.position = "fixed";
          temp.style.left = "-9999px";
          document.body.appendChild(temp);
          temp.focus();
          temp.select();
          document.execCommand("copy");
          temp.remove();
        }
        note.textContent = "JSON copied to clipboard.";
      } catch (err) {
        note.textContent = "Copy failed. Select the JSON text and copy manually.";
      }
    }

    function setBusy(isBusy) {
      statusText.textContent = isBusy ? "Scanning" : "Ready";
      document.querySelectorAll("button").forEach(button => button.disabled = isBusy);
    }

    function showError(message) {
      error.textContent = message;
    }

    function riskClass(value) {
      return String(value || "").toLowerCase().replace(/\s+/g, "-");
    }

    function safePath(object, path, fallback) {
      return path.split(".").reduce((value, key) => value && value[key], object) ?? fallback;
    }

    function escapeHTML(value) {
      return String(value ?? "").replace(/[&<>"']/g, char => ({
        "&": "&amp;", "<": "&lt;", ">": "&gt;", "\"": "&quot;", "'": "&#39;"
      }[char]));
    }
  </script>
</body>
</html>`
