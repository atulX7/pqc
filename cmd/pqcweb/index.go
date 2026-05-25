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
      --bg: #f6f7f9;
      --panel: #ffffff;
      --text: #17202a;
      --muted: #5d6978;
      --line: #d9dee7;
      --accent: #126b5a;
      --accent-dark: #0d5145;
      --warn: #a33d16;
      --shadow: 0 10px 30px rgba(20, 32, 45, 0.08);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      background: var(--bg);
      color: var(--text);
    }
    header {
      background: #ffffff;
      border-bottom: 1px solid var(--line);
      padding: 18px 28px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
    }
    h1 { font-size: 21px; margin: 0; letter-spacing: 0; }
    main {
      max-width: 1180px;
      margin: 0 auto;
      padding: 24px;
      display: grid;
      grid-template-columns: 360px 1fr;
      gap: 22px;
    }
    section {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
      box-shadow: var(--shadow);
    }
    .form-panel { padding: 20px; }
    .results-panel { min-height: 680px; }
    h2 { font-size: 16px; margin: 0 0 16px; }
    label { display: block; font-size: 12px; font-weight: 700; color: #344252; margin: 14px 0 6px; }
    input, select {
      width: 100%;
      min-height: 38px;
      border: 1px solid var(--line);
      border-radius: 6px;
      padding: 8px 10px;
      color: var(--text);
      background: #fff;
      font: inherit;
    }
    .buttons { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; margin-top: 18px; }
    button {
      min-height: 40px;
      border: 0;
      border-radius: 6px;
      background: var(--accent);
      color: white;
      font-weight: 800;
      cursor: pointer;
    }
    button.secondary { background: #30435a; }
    button:disabled { opacity: 0.6; cursor: wait; }
    .summary {
      display: grid;
      grid-template-columns: repeat(4, minmax(120px, 1fr));
      gap: 12px;
      padding: 18px;
      border-bottom: 1px solid var(--line);
    }
    .metric {
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 12px;
      min-height: 84px;
      background: #fbfcfd;
    }
    .metric span { display: block; color: var(--muted); font-size: 12px; font-weight: 700; }
    .metric strong { display: block; margin-top: 8px; font-size: 24px; }
    .content { padding: 18px; }
    table { width: 100%; border-collapse: collapse; font-size: 13px; }
    th, td { text-align: left; padding: 10px; border-bottom: 1px solid var(--line); vertical-align: top; }
    th { color: #344252; font-size: 12px; }
    .empty { color: var(--muted); padding: 64px 20px; text-align: center; }
    .error { color: var(--warn); font-weight: 700; margin-top: 12px; min-height: 20px; }
    pre {
      white-space: pre-wrap;
      word-break: break-word;
      background: #111827;
      color: #e5e7eb;
      border-radius: 8px;
      padding: 14px;
      max-height: 320px;
      overflow: auto;
    }
    @media (max-width: 900px) {
      header { padding: 16px; }
      main { grid-template-columns: 1fr; padding: 16px; }
      .summary { grid-template-columns: repeat(2, minmax(120px, 1fr)); }
    }
  </style>
</head>
<body>
  <header>
    <h1>PQC Readiness</h1>
    <span id="status">Ready</span>
  </header>
  <main>
    <section class="form-panel">
      <h2>Assessment Inputs</h2>
      <form id="scan-form">
        <label for="app_name">Application name</label>
        <input id="app_name" name="app_name" value="Clinical API">

        <label for="repo">Repository zip</label>
        <input id="repo" name="repo" type="file" accept=".zip">

        <label for="sensitivity">Sensitivity flags</label>
        <input id="sensitivity" name="sensitivity" value="PHI,GxP">

        <label for="exposure">Exposure</label>
        <select id="exposure" name="exposure">
          <option value="internet_authenticated">Internet authenticated</option>
          <option value="internet_public">Internet public</option>
          <option value="partner">Partner</option>
          <option value="internal">Internal</option>
          <option value="offline">Offline</option>
        </select>

        <label for="business_criticality">Business criticality</label>
        <select id="business_criticality" name="business_criticality">
          <option value="high">High</option>
          <option value="mission_critical">Mission critical</option>
          <option value="medium">Medium</option>
          <option value="low">Low</option>
        </select>

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

        <div class="buttons">
          <button type="button" id="sample">Scan Sample</button>
          <button type="submit" class="secondary">Scan Zip</button>
        </div>
        <div class="error" id="error"></div>
      </form>
    </section>
    <section class="results-panel">
      <div id="results" class="empty">Run a scan to view crypto inventory and risk posture.</div>
    </section>
  </main>
  <script>
    const form = document.querySelector("#scan-form");
    const error = document.querySelector("#error");
    const statusText = document.querySelector("#status");
    const results = document.querySelector("#results");
    const sampleButton = document.querySelector("#sample");

    form.addEventListener("submit", event => {
      event.preventDefault();
      if (!document.querySelector("#repo").files.length) {
        showError("Choose a repository .zip file or run the sample scan.");
        return;
      }
      runScan("/api/scan/upload", new FormData(form));
    });

    sampleButton.addEventListener("click", () => {
      const data = new FormData(form);
      data.delete("repo");
      runScan("/api/scan/sample", data);
    });

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
      const rows = report.top_critical_assets.map(asset =>
        "<tr>" +
          "<td>" + escapeHTML(asset.algorithm) + "</td>" +
          "<td>" + escapeHTML(asset.risk_level) + "<br>" + asset.risk_score.toFixed(1) + "</td>" +
          "<td>" + escapeHTML(asset.file) + ":" + asset.line + "</td>" +
          "<td>" + escapeHTML(asset.recommendation) + "</td>" +
        "</tr>"
      ).join("");
      results.className = "";
      results.innerHTML =
        "<div class=\"summary\">" +
          "<div class=\"metric\"><span>Risk posture</span><strong>" + escapeHTML(report.risk_posture || "Low") + "</strong></div>" +
          "<div class=\"metric\"><span>Highest score</span><strong>" + Number(report.highest_quantum_risk_score).toFixed(1) + "</strong></div>" +
          "<div class=\"metric\"><span>Total findings</span><strong>" + report.total_findings + "</strong></div>" +
          "<div class=\"metric\"><span>Files scanned</span><strong>" + report.scan_summary.files_scanned + "</strong></div>" +
        "</div>" +
        "<div class=\"content\">" +
          "<h2>Top Critical Assets</h2>" +
          "<table>" +
            "<thead><tr><th>Algorithm</th><th>Risk</th><th>Location</th><th>Recommendation</th></tr></thead>" +
            "<tbody>" + (rows || "<tr><td colspan=\"4\">No findings detected.</td></tr>") + "</tbody>" +
          "</table>" +
          "<h2 style=\"margin-top:22px\">JSON Report</h2>" +
          "<pre>" + escapeHTML(JSON.stringify(report, null, 2)) + "</pre>" +
        "</div>";
    }

    function setBusy(isBusy) {
      statusText.textContent = isBusy ? "Scanning" : "Ready";
      document.querySelectorAll("button").forEach(button => button.disabled = isBusy);
    }

    function showError(message) {
      error.textContent = message;
    }

    function escapeHTML(value) {
      return String(value ?? "").replace(/[&<>"']/g, char => ({
        "&": "&amp;", "<": "&lt;", ">": "&gt;", "\"": "&quot;", "'": "&#39;"
      }[char]));
    }
  </script>
</body>
</html>`
